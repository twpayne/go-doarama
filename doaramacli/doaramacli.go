// Package doaramacli provides integration between
// github.com/twpayne/go-doarama and github.com/urfave/cli.
package doaramacli

import (
	"errors"

	"github.com/twpayne/go-doarama"
	"github.com/urfave/cli"
)

// Flags specify connection and authentication options.
var Flags = []cli.Flag{
	cli.StringFlag{
		Name:   "apiurl",
		Value:  doarama.DefaultAPIURL,
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

// ActivityTypeFlag specifies the activity type.
var ActivityTypeFlag = cli.StringFlag{
	Name:  "activitytype",
	Usage: "activity type",
}

// VisualisationFlags specify visualisation options.
var VisualisationFlags = []cli.Flag{
	cli.StringSliceFlag{
		Name:  "name",
		Usage: "name",
	},
	cli.StringSliceFlag{
		Name:  "avatar",
		Usage: "avatar",
	},
	cli.StringFlag{
		Name:  "avatarbaseurl",
		Usage: "avatar base URL",
	},
	cli.BoolFlag{
		Name:  "fixedaspect",
		Usage: "fixed aspect",
	},
	cli.BoolFlag{
		Name:  "minimalview",
		Usage: "minimal view",
	},
	cli.StringFlag{
		Name:  "dzml",
		Usage: "DZML",
	},
}

// ActivityType returns the activity type from c.
func ActivityType(c *cli.Context) string {
	return c.String("activitytype")
}

// BaseDoaramaOptions returns the doarama.Options from c.
func BaseDoaramaOptions(c *cli.Context) []doarama.ClientOption {
	return []doarama.ClientOption{
		doarama.APIURL(c.GlobalString("apiurl")),
		doarama.APIName(c.GlobalString("apiname")),
		doarama.APIKey(c.GlobalString("apikey")),
	}
}

// NewDoaramaClient returns a new doarama.Client constructed from c.
func NewDoaramaClient(c *cli.Context) *doarama.Client {
	options := BaseDoaramaOptions(c)
	return doarama.NewClient(options...)
}

// NewAuthenticatedDoaramaOptions returns the doaram.Options for an
// authenticated doarama.Client from c.
func NewAuthenticatedDoaramaOptions(c *cli.Context) ([]doarama.ClientOption, error) {
	options := BaseDoaramaOptions(c)
	userID := c.GlobalString("userid")
	userKey := c.GlobalString("userkey")
	switch {
	case userID != "" && userKey == "":
		options = append(options, doarama.Anonymous(userID))
	case userID == "" && userKey != "":
		options = append(options, doarama.Delegate(userKey))
	default:
		return nil, errors.New("exactly one of -userid and -userkey must be specified")
	}
	return options, nil
}

// NewAuthenticatedDoaramaClient returns a new authenticated doarama.Client
// from c.
func NewAuthenticatedDoaramaClient(c *cli.Context) (*doarama.Client, error) {
	options, err := NewAuthenticatedDoaramaOptions(c)
	if err != nil {
		return nil, err
	}
	return doarama.NewClient(options...), nil
}

// NewVisualisationURLOptions returns a new doarama.VisualisationURLOptions
// from c.
func NewVisualisationURLOptions(c *cli.Context) *doarama.VisualisationURLOptions {
	vuo := &doarama.VisualisationURLOptions{}
	if name := c.StringSlice("name"); name != nil {
		vuo.Names = name
	}
	if avatar := c.StringSlice("avatar"); avatar != nil {
		vuo.Avatars = avatar
	}
	if avatarBaseURL := c.String("avatarbaseurl"); avatarBaseURL != "" {
		vuo.AvatarBaseURL = avatarBaseURL
	}
	if fixedAspect := c.Bool("fixedaspect"); fixedAspect {
		vuo.FixedAspect = fixedAspect
	}
	if minimalView := c.Bool("minimalview"); minimalView {
		vuo.MinimalView = minimalView
	}
	if dzml := c.String("dzml"); dzml != "" {
		vuo.DZML = dzml
	}
	return vuo
}
