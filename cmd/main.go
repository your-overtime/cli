package main

import (
	"errors"
	"fmt"
	"os"
	"time"

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
		os.Exit(1)
	}
}

func createState() error {
	loadConfig()
	if config != nil && len(config.Token) > 0 && len(config.Host) > 0 {
		c := client.Init(config.Host, config.Token)
		otc = &c
		return nil
	}
	return errors.New("No valid config found")
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
			return nil
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
							return conf.InitConf()
						},
					},
				},
			},
			{
				Name:    "token",
				Aliases: []string{"t"},
				Usage:   "handles tokens",
				Before: func(c *cli.Context) error {
					err := createState()
					if err == nil {
						c := client.Init(config.Host, config.Token)
						otc = &c
						return nil
					}
					os.Exit(1)
					return errors.New("No conf loaded")
				},
				Subcommands: []*cli.Command{
					{
						Name:    "create",
						Aliases: []string{"c"},
						Usage:   "creates a token",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "name",
								Required: true,
							},
						},
						Action: func(c *cli.Context) error {
							return otc.CreateTokenViaCli(c.String("name"))
						},
					},
					{
						Name:    "get",
						Aliases: []string{"g"},
						Usage:   "get tokens",
						Action: func(c *cli.Context) error {
							return otc.GetTokens()
						},
					},
					{
						Name:    "delete",
						Aliases: []string{"c"},
						Usage:   "deletes a token",
						Flags: []cli.Flag{
							&cli.UintFlag{
								Name:     "id",
								Required: true,
							},
						},
						Action: func(c *cli.Context) error {
							return otc.DeleteToken(c.Uint("id"))
						},
					},
				},
			},
			{
				Name:    "activity",
				Aliases: []string{"a"},
				Usage:   "handles activities",
				Before: func(c *cli.Context) error {
					err := createState()
					if err == nil {
						c := client.Init(config.Host, config.Token)
						otc = &c
						return nil
					}
					os.Exit(1)
					return errors.New("No conf loaded")
				},
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
								Value:       cli.NewTimestamp(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Now().Location())),
								DefaultText: "now -1 day",
								Layout:      "2006-01-02",
							},
							&cli.TimestampFlag{
								Name:        "end",
								Value:       cli.NewTimestamp(time.Now()),
								DefaultText: "now",
								Layout:      "2006-01-02",
							},
							&cli.BoolFlag{
								Name: "json",
							},
						},
						Action: func(c *cli.Context) error {
							return otc.GetActivities(*c.Timestamp("start"), *c.Timestamp("end"), c.Bool("json"))
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
								Required:    true,
							},
							&cli.TimestampFlag{
								Name:        "end",
								DefaultText: "now",
								Layout:      "2006-01-02 15:04",
								Required:    true,
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
					{
						Name:    "import",
						Aliases: []string{"i"},
						Usage:   "imports activities from kimai",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "csv",
								Required: true,
							},
						},
						Action: func(c *cli.Context) error {
							return otc.ImportKimai(c.Path("csv"))
						},
					},
				},
			},
			{
				Name:    "employee",
				Aliases: []string{"e"},
				Usage:   "handles employees",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "adminToken",
						Aliases:  []string{"at"},
						Required: true,
					},
				},
				Before: func(c *cli.Context) error {
					err := createState()
					if err == nil {
						c := client.Init(config.Host, config.Token)
						otc = &c
						return nil
					}
					os.Exit(1)
					return errors.New("No conf loaded")
				},
				Subcommands: []*cli.Command{
					{
						Name:    "add",
						Aliases: []string{"a"},
						Usage:   "creates a new employee",
						Flags: []cli.Flag{
							&cli.UintFlag{
								Name:     "WeekWorkingTimeInMinutes",
								Aliases:  []string{"wwtim"},
								Required: true,
							},
							&cli.StringFlag{
								Name:     "name",
								Aliases:  []string{"n"},
								Required: true,
							},
							&cli.StringFlag{
								Name:     "surname",
								Aliases:  []string{"s"},
								Required: true,
							},
							&cli.StringFlag{
								Name:     "login",
								Aliases:  []string{"l"},
								Required: true,
							},
							&cli.StringFlag{
								Name:     "password",
								Aliases:  []string{"p"},
								Required: true,
							},
						},
						Action: func(c *cli.Context) error {
							return otc.AddEmployee(c.String("name"), c.String("surname"), c.String("login"), c.String("password"), c.Uint("WeekWorkingTimeInMinutes"), c.String("adminToken"))
						},
					},
					{
						Name:    "delete",
						Aliases: []string{"d"},
						Usage:   "deletes a employee",
						Flags: []cli.Flag{&cli.StringFlag{
							Name:     "login",
							Aliases:  []string{"l"},
							Required: true,
						}},
						Action: func(c *cli.Context) error {
							return otc.DeleteEmployee(c.String("login"), c.String("adminToken"))
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
