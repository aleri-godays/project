package main

import (
	"github.com/aleri-godays/project/internal"
	"github.com/aleri-godays/project/internal/config"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	version     = "dev-snapshot"
	serviceName = "project"
	dbPath      string
)

func main() {
	cliApp := cli.NewApp()
	cliApp.Name = serviceName
	cliApp.Version = version

	cliApp.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "log-level,l",
			Usage:  "log level: TRACE, DEBUG, INFO, WARN, ERROR",
			EnvVar: "LOG_LEVEL",
			Value:  "DEBUG",
		},
		cli.StringFlag{
			Name:        "db",
			Usage:       "db file path",
			EnvVar:      "DB",
			Value:       ".",
			Destination: &dbPath,
		},
		cli.IntFlag{
			Name:   "port",
			Usage:  "port to server",
			EnvVar: "PORT",
			Value:  5010,
		},
		cli.StringFlag{
			Name:   "jwt-secret",
			Usage:  "jwt-secret",
			EnvVar: "JWT_SECRET",
			Value:  "aslkdhasljkdhasjdh",
		},
	}
	cliApp.Commands = []cli.Command{
		{
			Name:  "run",
			Usage: "server",
			Action: func(c *cli.Context) error {
				conf := getConfig(c)
				app := internal.NewApp(conf)
				app.Run()
				return nil
			},
		},
	}

	err := cliApp.Run(os.Args)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("could not initialize app")
	}
}

func getConfig(c *cli.Context) *config.Config {
	conf := &config.Config{
		Version:     version,
		ServiceName: serviceName,
		LogLevel:    c.GlobalString("log-level"),
		DbPath:      c.GlobalString("db"),
		HTTPPort:    c.GlobalInt("port"),
		JWTSecret:   c.GlobalString("jwt-secret"),
	}

	return conf
}
