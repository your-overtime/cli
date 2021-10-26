package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/your-overtime/cli/internal/client"
	"github.com/your-overtime/cli/internal/conf"
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

func fixLocation(t *time.Time) *time.Time {
	if t != nil {
		newT := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, time.Local)
		return &newT
	}
	return nil
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
							return client.InitConf()
						},
					},
				},
			},
			{
				Name:  "export",
				Usage: "export all data since given start (default the last year)",
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
				Flags: []cli.Flag{
					&cli.TimestampFlag{
						Name:        "since",
						Aliases:     []string{"s"},
						DefaultText: "now -1 Year",
						Layout:      "2006-01-02",
						Required:    false,
					},
					&cli.StringFlag{
						Name:     "output",
						Required: true,
						Aliases:  []string{"o"},
					},
				},
				Action: func(c *cli.Context) error {
					s := fixLocation(c.Timestamp("since"))
					return otc.Export(s, c.String("output"))
				},
			},
			{
				Name:  "import",
				Usage: "import all data",
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
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "input",
						Required: true,
						Aliases:  []string{"i"},
					},
				},
				Action: func(c *cli.Context) error {
					return otc.Import(c.String("input"))
				},
			},
			{
				Name:  "app",
				Usage: "handles app setup",
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
						Name:    "link",
						Aliases: []string{"l"},
						Usage:   "Links a app",
						Action: func(ctx *cli.Context) error {
							return otc.LinkApp()
						},
					},
				},
			},
			{
				Name:    "account",
				Aliases: []string{"ac"},
				Usage:   "update user information",
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
						Name:    "update",
						Aliases: []string{"u"},
						Usage:   "Update account values",
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:    "name",
								Aliases: []string{"n"},
								Usage:   "change name",
							},
							&cli.BoolFlag{
								Name:    "surname",
								Aliases: []string{"s"},
								Usage:   "change surname",
							},
							&cli.BoolFlag{
								Name:    "login",
								Aliases: []string{"l"},
								Usage:   "change login",
							},
							&cli.BoolFlag{
								Name:    "password",
								Aliases: []string{"p"},
								Usage:   "change password",
							},
							&cli.BoolFlag{
								Name:  "wwt",
								Usage: "change week working time",
							},
							&cli.BoolFlag{
								Name:  "nwwd",
								Usage: "change number of week working days",
							},
						},
						Action: func(c *cli.Context) error {
							return otc.ChangeAccount(c.Bool("name"), c.Bool("surname"), c.Bool("login"), c.Bool("password"), c.Bool("wwt"), c.Bool("nwwd"), c.Args().First())
						},
					},
					{
						Name:    "get",
						Aliases: []string{"g"},
						Usage:   "get account values",
						Action: func(ctx *cli.Context) error {
							return otc.GetAccount()
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
						Aliases: []string{"d"},
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
								Name: "description",
							},
						},
						Usage: "starts new activity",
						Before: func(c *cli.Context) error {
							if !c.IsSet("description") {
								if c.NArg() > 0 {
									c.Set("description", strings.Join(c.Args().Slice(), " "))
									return nil
								}
								if len(config.DefaultActivityDesc) > 0 {
									c.Set("description", config.DefaultActivityDesc)
									return nil
								}
								var desc string
								err := survey.AskOne(&survey.Input{
									Message: "Description",
								}, &desc)
								if err != nil {
									return err
								}
								if len(desc) > 0 {
									c.Set("description", desc)
									return nil
								}
							}
							return errors.New("description missing")
						},
						Action: func(c *cli.Context) error {

							return otc.StartActivity(c.String("description"))
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
						Name:    "get",
						Aliases: []string{"g"},
						Usage:   "fetch activities between start and end",
						Flags: []cli.Flag{
							&cli.TimestampFlag{
								Name:        "start",
								Aliases:     []string{"s"},
								Value:       cli.NewTimestamp(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Now().Location())),
								DefaultText: "now -1 day",
								Layout:      "2006-01-02",
							},
							&cli.TimestampFlag{
								Name:        "end",
								Aliases:     []string{"e"},
								Value:       cli.NewTimestamp(time.Now()),
								DefaultText: "now",
								Layout:      "2006-01-02",
							},
							&cli.BoolFlag{
								Name: "json",
							},
						},
						Action: func(c *cli.Context) error {
							return otc.GetActivities(*fixLocation(c.Timestamp("start")), *fixLocation(c.Timestamp("end")), c.Bool("json"))
						},
					},
					{
						Name:    "create",
						Aliases: []string{"c"},
						Usage:   "creates a activity",
						Flags: []cli.Flag{
							&cli.TimestampFlag{
								Name:        "start",
								Aliases:     []string{"s"},
								DefaultText: "now -1 day",
								Layout:      "2006-01-02 15:04",
								Required:    true,
							},
							&cli.TimestampFlag{
								Name:        "end",
								Aliases:     []string{"e"},
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
							return otc.AddActivity(desc, fixLocation(c.Timestamp("start")), fixLocation(c.Timestamp("end")))
						},
					},
					{
						Name:    "update",
						Aliases: []string{"u"},
						Usage:   "updates a activity",
						Flags: []cli.Flag{
							&cli.TimestampFlag{
								Name:        "start",
								Aliases:     []string{"s"},
								DefaultText: "now -1 day",
								Layout:      "2006-01-02 15:04",
							},
							&cli.TimestampFlag{
								Name:        "end",
								Aliases:     []string{"e"},
								DefaultText: "now",
								Layout:      "2006-01-02 15:04",
							},

							&cli.StringFlag{
								Name:  "description",
								Value: "",
							},
							&cli.UintFlag{
								Name:     "id",
								Required: true,
							},
						},
						Action: func(c *cli.Context) error {
							return otc.UpdateActivity(c.String("desc"), fixLocation(c.Timestamp("start")), fixLocation(c.Timestamp("end")), c.Uint("id"))
						},
					},
					{
						Name:    "delete",
						Aliases: []string{"d"},
						Usage:   "deletes a activity",
						Flags: []cli.Flag{
							&cli.UintFlag{
								Name:     "id",
								Required: true,
							},
						},
						Action: func(c *cli.Context) error {
							return otc.DeleteActivity(c.Uint("id"))
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
				Name:    "holidays",
				Aliases: []string{"h"},
				Usage:   "handles holidays",
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
						Name:    "get",
						Aliases: []string{"g"},
						Usage:   "fetch holidays between start and end",
						Flags: []cli.Flag{
							&cli.TimestampFlag{
								Name:        "start",
								Aliases:     []string{"s"},
								Value:       cli.NewTimestamp(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Now().Location())),
								DefaultText: "now -1 day",
								Layout:      "2006-01-02",
							},
							&cli.TimestampFlag{
								Name:        "end",
								Aliases:     []string{"e"},
								Value:       cli.NewTimestamp(time.Now()),
								DefaultText: "now",
								Layout:      "2006-01-02",
							},
							&cli.BoolFlag{
								Name: "json",
							},
						},
						Action: func(c *cli.Context) error {
							return otc.GetHolidays(*fixLocation(c.Timestamp("start")), *fixLocation(c.Timestamp("end")), c.Bool("json"))
						},
					},
					{
						Name:    "create",
						Aliases: []string{"c"},
						Usage:   "creates a holiday",
						Flags: []cli.Flag{
							&cli.TimestampFlag{
								Name:     "start",
								Layout:   "2006-01-02",
								Required: true,
								Aliases:  []string{"s"},
							},
							&cli.TimestampFlag{
								Name:    "end",
								Layout:  "2006-01-02",
								Aliases: []string{"e"},
							},
							&cli.StringFlag{
								Name:     "description",
								Required: true,
								Aliases:  []string{"d"},
							},
							&cli.BoolFlag{
								Name:    "legalholiday",
								Aliases: []string{"l"},
								Value:   false,
							},
							&cli.BoolFlag{
								Name:    "sick",
								Aliases: []string{"si"},
								Value:   false,
							},
						},
						Action: func(c *cli.Context) error {
							e := fixLocation(c.Timestamp("end"))
							s := fixLocation(c.Timestamp("start"))
							if e == nil {
								ce := time.Date(s.Year(), s.Month(), s.Day(), 23, 59, 59, 59, s.Location())
								e = &ce
							}
							return otc.AddHoliday(c.String("description"), *s, *e, c.Bool("legalholiday"), c.Bool("sick"))
						},
					},
					{
						Name:    "update",
						Aliases: []string{"u"},
						Usage:   "updates a holiday",
						Flags: []cli.Flag{
							&cli.TimestampFlag{
								Name:        "start",
								Aliases:     []string{"s"},
								DefaultText: "now -1 day",
								Layout:      "2006-01-02",
							},
							&cli.TimestampFlag{
								Name:        "end",
								Aliases:     []string{"e"},
								DefaultText: "now",
								Layout:      "2006-01-02",
							}, &cli.StringFlag{
								Name:    "description",
								Aliases: []string{"d"},
							},
							&cli.UintFlag{
								Name:     "id",
								Required: true,
							},
							&cli.BoolFlag{
								Name:    "legalholiday",
								Aliases: []string{"l"},
								Value:   false,
							},
							&cli.BoolFlag{
								Name:    "sick",
								Aliases: []string{"si"},
								Value:   false,
							},
							&cli.BoolFlag{
								Name:    "free",
								Aliases: []string{"f"},
								Value:   false,
							},
						},
						Action: func(c *cli.Context) error {
							e := fixLocation(c.Timestamp("end"))
							s := fixLocation(c.Timestamp("start"))
							if e == nil {
								ce := time.Date(s.Year(), s.Month(), s.Day(), 23, 59, 59, 59, s.Location())
								e = &ce
							}
							return otc.UpdateHoliday(c.String("description"), fixLocation(c.Timestamp("start")), fixLocation(c.Timestamp("end")), c.Uint("id"), c.Bool("legalholiday"), c.Bool("sick"), c.Bool("free"))
						},
					},
					{
						Name:    "delete",
						Aliases: []string{"d"},
						Usage:   "deletes a holiday",
						Flags: []cli.Flag{
							&cli.UintFlag{
								Name:     "id",
								Required: true,
							},
						},
						Action: func(c *cli.Context) error {
							return otc.DeleteHoliday(c.Uint("id"))
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
								Name:     "weekWorkingTimeInMinutes",
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
							&cli.UintFlag{
								Name:     "numberOfWeekWorkingDays",
								Aliases:  []string{"nwwd"},
								Required: true,
							},
						},
						Action: func(c *cli.Context) error {
							return otc.AddEmployee(c.String("name"), c.String("surname"), c.String("login"), c.String("password"), c.Uint("weekWorkingTimeInMinutes"), c.Uint("numberOfWeekWorkingDays"), c.String("adminToken"))
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
