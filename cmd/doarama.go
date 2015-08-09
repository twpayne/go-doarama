package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"

	"github.com/codegangsta/cli"
	"github.com/twpayne/go-doarama"
)

func newDoaramaClient(c *cli.Context) *doarama.Client {
	return doarama.NewClient(c.GlobalString("apiurl"), c.GlobalString("apiname"), c.GlobalString("apikey"))
}

func newAuthenticatedDoaramaClient(c *cli.Context) (*doarama.Client, error) {
	client := newDoaramaClient(c)
	userId := c.GlobalString("userid")
	userKey := c.GlobalString("userkey")
	switch {
	case userId != "" && userKey == "":
		return client.Anonymous(userId), nil
	case userId == "" && userKey != "":
		return client.Delegate(userKey), nil
	default:
		return nil, errors.New("exactly one of -userid and -userkey must be specified")
	}
}

func activityCreateOne(client *doarama.Client, filename string) (*doarama.Activity, error) {
	gpsTrack, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer gpsTrack.Close()
	return client.CreateActivity(filepath.Base(filename), gpsTrack)
}

func activityCreate(c *cli.Context) error {
	client, err := newAuthenticatedDoaramaClient(c)
	if err != nil {
		return err
	}
	typeId := c.Int("typeid")
	for _, arg := range c.Args() {
		a, err := activityCreateOne(client, arg)
		if err != nil {
			log.Print(err)
			continue
		}
		fmt.Printf("ActivityId: %d\n", a.Id)
		if err := a.SetInfo(&doarama.ActivityInfo{
			TypeId: typeId,
		}); err != nil {
			log.Print(err)
			continue
		}
	}
	return nil
}

func activityDelete(c *cli.Context) error {
	client, err := newAuthenticatedDoaramaClient(c)
	if err != nil {
		return err
	}
	var ids []int
	for _, arg := range c.Args() {
		id64, err := strconv.ParseInt(arg, 10, 0)
		if err != nil {
			return err
		}
		ids = append(ids, int(id64))
	}
	for _, id := range ids {
		a := client.Activity(id, "WGS84")
		if err := a.Delete(); err != nil {
			log.Print(err)
			continue
		}
	}
	return nil
}

func queryActivityTypes(c *cli.Context) error {
	client := newDoaramaClient(c)
	ats, err := client.ActivityTypes()
	if err != nil {
		return err
	}
	var names []string
	for name, _ := range ats {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		fmt.Printf("%s: %d\n", name, ats[name])
	}
	return nil
}

func visualisationCreate(c *cli.Context) error {
	client, err := newAuthenticatedDoaramaClient(c)
	if err != nil {
		return err
	}
	var as []*doarama.Activity
	for _, arg := range c.Args() {
		id64, err := strconv.ParseInt(arg, 10, 0)
		if err != nil {
			return err
		}
		a := client.Activity(int(id64), "WGS84")
		as = append(as, a)
	}
	v, err := client.CreateVisualisation(as)
	if err != nil {
		return err
	}
	fmt.Printf("VisualisationKey: %s\n", v.Key)
	return nil
}

func visualisationDelete(c *cli.Context) error {
	client, err := newAuthenticatedDoaramaClient(c)
	if err != nil {
		return err
	}
	for _, arg := range c.Args() {
		v := client.Visualisation(arg)
		if err := v.Delete(); err != nil {
			log.Print(err)
			continue
		}
	}
	return nil
}

func visualisationURL(c *cli.Context) error {
	client := newDoaramaClient(c)
	for _, arg := range c.Args() {
		v := client.Visualisation(arg)
		var vuo doarama.VisualisationURLOptions
		if c.StringSlice("name") != nil {
			vuo.Names = c.StringSlice("name")
		}
		if c.StringSlice("avatar") != nil {
			vuo.Avatars = c.StringSlice("avatar")
		}
		if c.String("avatarbaseurl") != "" {
			vuo.AvatarBaseURL = c.String("avatarbaseurl")
		}
		if c.Bool("fixedaspect") {
			vuo.FixedAspect = c.Bool("fixedaspect")
		}
		if c.Bool("minimalview") {
			vuo.MinimalView = c.Bool("minimalview")
		}
		if c.String("dzml") != "" {
			vuo.DZML = c.String("dzml")
		}
		fmt.Printf("VisualisationURL: %s\n", v.URL(&vuo))
	}
	return nil
}

func logError(f func(*cli.Context) error) func(*cli.Context) {
	return func(c *cli.Context) {
		if err := f(c); err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "doarama"
	app.Usage = "A command line interface to doarama.com"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "apiurl",
			Value:  doarama.API_URL,
			Usage:  "Doarama API URL",
			EnvVar: "DOARAMA_API_URL",
		},
		cli.StringFlag{
			Name:   "apikey",
			Usage:  "Doarama API key",
			EnvVar: "DOARAMA_API_KEY",
		},
		cli.StringFlag{
			Name:   "apiname",
			Usage:  "Doarama API name",
			EnvVar: "DOARAMA_API_NAME",
		},
		cli.StringFlag{
			Name:   "userid",
			Usage:  "Doarama user ID",
			EnvVar: "DOARAMA_USER_ID",
		},
		cli.StringFlag{
			Name:   "userkey",
			Usage:  "Doarama user key",
			EnvVar: "DOARAMA_USER_KEY",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:    "activity",
			Aliases: []string{"a"},
			Usage:   "Activity",
			Subcommands: []cli.Command{
				{
					Name:    "create",
					Aliases: []string{"c"},
					Usage:   "Create activity",
					Action:  logError(activityCreate),
					Flags: []cli.Flag{
						cli.IntFlag{
							Name:  "typeid",
							Usage: "Type ID",
						},
					},
				},
				{
					Name:    "delete",
					Aliases: []string{"d"},
					Usage:   "Delete activity",
					Action:  logError(activityDelete),
				},
			},
		},
		{
			Name:    "query-activity-types",
			Aliases: []string{"qat"},
			Usage:   "Query activity types",
			Action:  logError(queryActivityTypes),
		},
		{
			Name:    "visualisation",
			Aliases: []string{"v"},
			Subcommands: []cli.Command{
				{
					Name:    "create",
					Aliases: []string{"c"},
					Usage:   "Create visualisation",
					Action:  logError(visualisationCreate),
				},
				{
					Name:    "delete",
					Aliases: []string{"d"},
					Usage:   "Delete visualisation",
					Action:  logError(visualisationDelete),
				},
				{
					Name:    "url",
					Aliases: []string{"u"},
					Usage:   "Visualisation URL",
					Action:  logError(visualisationURL),
					Flags: []cli.Flag{
						cli.StringSliceFlag{
							Name:  "name",
							Usage: "Name",
						},
						cli.StringSliceFlag{
							Name:  "avatar",
							Usage: "Avatar",
						},
						cli.StringFlag{
							Name:  "avatarbaseurl",
							Usage: "Avatar base URL",
						},
						cli.BoolTFlag{
							Name:  "fixedaspect",
							Usage: "Fixed aspect",
						},
						cli.BoolFlag{
							Name:  "minimalview",
							Usage: "Minimal view",
						},
						cli.StringFlag{
							Name:  "dzml",
							Usage: "DZML",
						},
					},
				},
			},
		},
	}
	app.Run(os.Args)
}
