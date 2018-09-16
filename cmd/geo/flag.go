package main

import (
	"github.com/major1201/goutils"
	"github.com/urfave/cli"
)

func getApp() *cli.App {
	app := cli.NewApp()
	app.Name = "geo"
	app.HelpName = app.Name
	app.Usage = "geo ip cli tool"
	app.ArgsUsage = "[hostname/ip[s]]"
	app.Version = AppVer
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "help, h",
			Usage: "show help",
		},
		cli.VersionFlag,
		cli.BoolFlag{
			Name:  "detail, d",
			Usage: "display geo info in detail",
		},
		cli.StringFlag{
			Name:   "language, l",
			Usage:  "specify output language, should be in [en, de, es, fr, ja, pt-BR, ru, zh-CN], default: en",
			Value:  "en",
			EnvVar: "GEO_LANG",
		},
		cli.StringFlag{
			Name:   "mmdb-file, m",
			Usage:  "MaxMind GeoIP database file location(required)",
			EnvVar: "GEO_MMDBFILE",
		},
		cli.BoolFlag{
			Name:  "json",
			Usage: "output in json format",
		},
	}
	app.Action = func(c *cli.Context) error {
		if c.Bool("help") {
			cli.ShowAppHelpAndExit(c, 0)
		}
		verifyFlags(c)
		runApp(c)
		return nil
	}
	app.HideHelp = true
	return app
}

func verifyFlags(c *cli.Context) {
	// mmdb-file must be specified
	if !c.IsSet("mmdb-file") {
		logger.Error("mmdb-file is required")
		cli.ShowAppHelpAndExit(c, 1)
	}

	// language should be in limited values
	if c.IsSet("language") {
		language := c.String("language")
		if !goutils.Contains(language, "en", "de", "es", "fr", "ja", "pt-BR", "ru", "zh-CN") {
			logger.Errorf("invalid language value: %v", language)
			cli.ShowAppHelpAndExit(c, 1)
		}
	}
}
