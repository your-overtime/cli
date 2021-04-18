package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"git.goasum.de/jasper/overtime-cli/cmd"
	"git.goasum.de/jasper/overtime-cli/internal/client"
	"git.goasum.de/jasper/overtime-cli/internal/conf"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

var (
	config *conf.Config
	otc    *client.Client
	err    error
)

func setLogger(debug bool) {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	log.SetReportCaller(true)
	if debug {
		log.SetLevel(log.DebugLevel)
	}
}

func loadConfig() {
	config, err = conf.LoadConfig()
	if err != nil {
		log.Debug(err)
		fmt.Println("Please run \"conf init\"")
	}
}

func createState() {
	loadConfig()
	if len(config.Token) > 0 && len(config.Host) > 0 {
		c := client.Init(config.Host, config.Token)
		otc = &c
	}
}

func main() {
	app := &cli.App{
		EnableBashCompletion: true,
		Name:                 "Overtime CLI",
		Usage:                "Controll your working time",
		Version:              "1.0.0",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "debug",
				Aliases: []string{"d"},
				Value:   false,
				Usage:   "enables debug logging",
			},
		},
		Before: func(c *cli.Context) error {
			setLogger(c.Bool("debug"))
			createState()
			if len(config.Token) > 0 && len(config.Host) > 0 {
				c := client.Init(config.Host, config.Token)
				otc = &c
				return nil
			}
			return errors.New("No conf loaded")
		},
		Commands: []*cli.Command{
			{
				Name:    "conf",
				Aliases: []string{"c"},
				Usage:   "handles config values",
				Subcommands: []*cli.Command{
					{
						Name:    "init",
						Aliases: []string{"i"},
						Usage:   "create or replace the config files",
						Action: func(ctx *cli.Context) error {
							return cmd.InitConf()
						},
					},
				},
			},
			{
				Name:    "activity",
				Aliases: []string{"a"},
				Usage:   "handles activities",
				Subcommands: []*cli.Command{
					{
						Name:    "start",
						Aliases: []string{"s"},
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:  "description",
								Value: "",
							},
						},
						Usage: "starts new activity",
						Action: func(c *cli.Context) error {
							desc := c.String("description")
							if len(desc) == 0 {
								desc = config.DefaultActivityDesc
							}
							return otc.StartActivity(desc)
						},
					},
					{
						Name:    "end",
						Aliases: []string{"e"},
						Usage:   "end currently running activity",
						Action: func(c *cli.Context) error {
							return otc.StopActivity()
						},
					},
					{
						Name:    "overview",
						Aliases: []string{"o"},
						Usage:   "shows current overview",
						Action: func(c *cli.Context) error {
							return otc.CalcCurrentOverview()
						},
					},
					{
						Name:    "activities",
						Aliases: []string{"a"},
						Usage:   "fetch activities between start and end",
						Flags: []cli.Flag{
							&cli.TimestampFlag{
								Name:        "start",
								Value:       cli.NewTimestamp(time.Now().AddDate(0, 0, -1)),
								DefaultText: "now -1 day",
								Layout:      "2006-01-02",
							},
							&cli.TimestampFlag{
								Name:        "end",
								Value:       cli.NewTimestamp(time.Now()),
								DefaultText: "now",
								Layout:      "2006-01-02",
							},
						},
						Action: func(c *cli.Context) error {
							return otc.GetActivities(*c.Timestamp("start"), *c.Timestamp("end"))
						},
					},
					{
						Name:    "create",
						Aliases: []string{"c"},
						Usage:   "creates a activity",
						Flags: []cli.Flag{
							&cli.TimestampFlag{
								Name:        "start",
								DefaultText: "now -1 day",
								Layout:      "2006-01-02 15:04",
							},
							&cli.TimestampFlag{
								Name:        "end",
								DefaultText: "now",
								Layout:      "2006-01-02 15:04",
							}, &cli.StringFlag{
								Name:  "description",
								Value: "",
							},
						},
						Action: func(c *cli.Context) error {
							desc := c.String("description")
							if len(desc) == 0 {
								desc = config.DefaultActivityDesc
							}
							return otc.AddActivity(desc, c.Timestamp("start"), c.Timestamp("end"))
						},
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
