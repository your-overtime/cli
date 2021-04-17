package main

import (
	"fmt"
	"os"

	"git.goasum.de/jasper/overtime-cli/cmd"
	"git.goasum.de/jasper/overtime-cli/internal/client"
	"git.goasum.de/jasper/overtime-cli/internal/conf"
	"github.com/pkg/errors"
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
		Version:              "0.0.1",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "debug",
				Value: false,
				Usage: "enables debug logging",
			},
		},
		Action: func(c *cli.Context) error {
			setLogger(c.Bool("debug"))
			createState()
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
						Name:    "overview",
						Aliases: []string{"o"},
						Usage:   "shows current overview",
						Action: func(c *cli.Context) error {
							setLogger(c.Bool("debug"))
							createState()
							if otc != nil {
								return otc.CalcCurrentOverview()
							}
							return errors.New("No conf loaded")
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
