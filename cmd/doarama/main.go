package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"

	"github.com/twpayne/go-doarama"
	"github.com/twpayne/go-doarama/doaramacache"
	"github.com/twpayne/go-doarama/doaramacli"
	"github.com/urfave/cli"

	_ "github.com/mattn/go-sqlite3"
)

func newCache(c *cli.Context, client *doarama.Client) (doaramacache.ActivityCreator, error) {
	dataSourceName := c.GlobalString("cache")
	if dataSourceName == "" {
		return client, nil
	}
	return doaramacache.NewSQLite3(dataSourceName, client)
}

func activityCreateOne(ctx context.Context, cache doaramacache.ActivityCreator, filename string, activityInfo *doarama.ActivityInfo) (*doarama.Activity, error) {
	gpsTrack, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer gpsTrack.Close()
	return cache.CreateActivityWithInfo(ctx, filepath.Base(filename), gpsTrack, activityInfo)
}

func activityCreate(c *cli.Context) error {
	ctx := context.Background()
	client, err := doaramacli.NewAuthenticatedDoaramaClient(c)
	if err != nil {
		return err
	}
	defer client.Close()
	activityType, err := doarama.DefaultActivityTypes.Find(doaramacli.ActivityType(c))
	if err != nil {
		return err
	}
	activityInfo := &doarama.ActivityInfo{
		TypeID: activityType.ID,
	}
	cache, err := newCache(c, client)
	if err != nil {
		return err
	}
	defer cache.Close()
	for _, arg := range c.Args() {
		a, err := activityCreateOne(ctx, cache, arg, activityInfo)
		if err != nil {
			log.Print(err)
			continue
		}
		fmt.Printf("ActivityId: %d\n", a.ID)
	}
	return nil
}

func activityDelete(c *cli.Context) error {
	ctx := context.Background()
	client, err := doaramacli.NewAuthenticatedDoaramaClient(c)
	if err != nil {
		return err
	}
	defer client.Close()
	var ids []int
	for _, arg := range c.Args() {
		id64, err := strconv.ParseInt(arg, 10, 0)
		if err != nil {
			return err
		}
		ids = append(ids, int(id64))
	}
	for _, id := range ids {
		a := client.Activity(id)
		if err := a.Delete(ctx); err != nil {
			log.Print(err)
			continue
		}
	}
	return nil
}

func create(c *cli.Context) error {
	ctx := context.Background()
	client, err := doaramacli.NewAuthenticatedDoaramaClient(c)
	if err != nil {
		return err
	}
	defer client.Close()
	activityType, err := doarama.DefaultActivityTypes.Find(doaramacli.ActivityType(c))
	if err != nil {
		return err
	}
	activityInfo := &doarama.ActivityInfo{
		TypeID: activityType.ID,
	}
	cache, err := newCache(c, client)
	if err != nil {
		return err
	}
	defer cache.Close()
	var as []*doarama.Activity
	for _, arg := range c.Args() {
		var a *doarama.Activity
		a, err = activityCreateOne(ctx, cache, arg, activityInfo)
		if err != nil {
			break
		}
		fmt.Printf("ActivityId: %d\n", a.ID)
		as = append(as, a)
	}
	if err != nil {
		for _, a := range as {
			a.Delete(ctx)
		}
		return err
	}
	if len(as) == 0 {
		return errors.New("no activities specified")
	}
	v, err := client.CreateVisualisation(ctx, as)
	if err != nil {
		return err
	}
	fmt.Printf("VisualisationKey: %s\n", v.Key)
	vuo := doaramacli.NewVisualisationURLOptions(c)
	fmt.Printf("VisualisationURL: %s\n", v.URL(vuo))
	return nil
}

type byName doarama.ActivityTypes

func (ats byName) Len() int           { return len(ats) }
func (ats byName) Less(i, j int) bool { return ats[i].Name < ats[j].Name }
func (ats byName) Swap(i, j int)      { ats[i], ats[j] = ats[j], ats[i] }

func queryActivityTypes(c *cli.Context) error {
	ctx := context.Background()
	client := doaramacli.NewDoaramaClient(c)
	defer client.Close()
	ats, err := client.ActivityTypes(ctx)
	if err != nil {
		return err
	}
	sort.Sort(byName(ats))
	for _, at := range ats {
		fmt.Printf("%s: %d\n", at.Name, at.ID)
	}
	return nil
}

func visualisationCreate(c *cli.Context) error {
	ctx := context.Background()
	client, err := doaramacli.NewAuthenticatedDoaramaClient(c)
	if err != nil {
		return err
	}
	defer client.Close()
	var as []*doarama.Activity
	for _, arg := range c.Args() {
		var id64 int64
		id64, err = strconv.ParseInt(arg, 10, 0)
		if err != nil {
			return err
		}
		a := client.Activity(int(id64))
		as = append(as, a)
	}
	v, err := client.CreateVisualisation(ctx, as)
	if err != nil {
		return err
	}
	fmt.Printf("VisualisationKey: %s\n", v.Key)
	return nil
}

func visualisationURL(c *cli.Context) error {
	client := doaramacli.NewDoaramaClient(c)
	vuo := doaramacli.NewVisualisationURLOptions(c)
	for _, arg := range c.Args() {
		v := client.Visualisation(arg)
		fmt.Printf("VisualisationURL: %s\n", v.URL(vuo))
	}
	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "doarama"
	app.Usage = "A command line interface to doarama.com"
	app.Flags = append([]cli.Flag{
		cli.StringFlag{
			Name:  "cache",
			Usage: "Path to cache",
			Value: path.Join(os.Getenv("HOME"), ".doaramacache.db"),
		},
	}, doaramacli.Flags...)
	app.Commands = []cli.Command{
		{
			Name:    "activity",
			Aliases: []string{"a"},
			Usage:   "Manages activities",
			Subcommands: []cli.Command{
				{
					Name:    "create",
					Aliases: []string{"c"},
					Usage:   "Creates an activity from one or more tracklogs",
					Action:  activityCreate,
					Flags:   []cli.Flag{doaramacli.ActivityTypeFlag},
				},
				{
					Name:    "delete",
					Aliases: []string{"d"},
					Usage:   "Deletes one or more activities by id",
					Action:  activityDelete,
				},
			},
		},
		{
			Name:    "create",
			Aliases: []string{"c"},
			Usage:   "Creates a visualisation URL from one or more tracklogs",
			Action:  create,
			Flags:   append([]cli.Flag{doaramacli.ActivityTypeFlag}, doaramacli.VisualisationFlags...),
		},
		{
			Name:    "query-activity-types",
			Aliases: []string{"qat"},
			Usage:   "Queries activity types",
			Action:  queryActivityTypes,
		},
		{
			Name:    "visualisation",
			Aliases: []string{"v"},
			Usage:   "Manages visualisations",
			Subcommands: []cli.Command{
				{
					Name:    "create",
					Aliases: []string{"c"},
					Usage:   "Creates a visualisation from a list of activities",
					Action:  visualisationCreate,
				},
				{
					Name:    "url",
					Aliases: []string{"u"},
					Usage:   "Creates a visualisation URL from a visualisation key",
					Action:  visualisationURL,
					Flags:   doaramacli.VisualisationFlags,
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
