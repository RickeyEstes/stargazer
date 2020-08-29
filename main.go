package main

import (
	"os"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	cli "gopkg.in/urfave/cli.v1"

	"github.com/richardlt/stargazer/config"
	"github.com/richardlt/stargazer/crawler"
	"github.com/richardlt/stargazer/web"
)

func main() {
	app := cli.NewApp()

	globalFlags := []cli.Flag{
		cli.StringFlag{
			Name:  "db-host",
			Value: "localhost",
			Usage: "Postgres database URI",
		},
		cli.Int64Flag{
			Name:  "db-port",
			Value: 5432,
			Usage: "Postgres database port",
		},
		cli.BoolFlag{
			Name:  "db-ssl",
			Usage: "Postgres database ssl mode",
		},
		cli.StringFlag{
			Name:  "db-name",
			Value: "stargazer",
			Usage: "Postgres database name",
		},
		cli.StringFlag{
			Name:  "db-user",
			Value: "stargazer",
			Usage: "Postgres database user",
		},
		cli.StringFlag{
			Name:  "db-pass",
			Value: "stargazer",
			Usage: "Postgres database password",
		},
		cli.StringFlag{
			Name:  "log-level",
			Value: "info",
			Usage: "[panic fatal error warning info debug]",
		},
	}

	app.Commands = []cli.Command{
		{
			Name: "crawler",
			Flags: append(globalFlags,
				cli.StringFlag{
					Name:  "token",
					Value: "secret",
					Usage: "Github api token",
				},
				cli.StringFlag{
					Name:  "mongo-uri",
					Value: "mongodb://localhost:27017",
					Usage: "Mongo database URI",
				},
				cli.IntFlag{
					Name:  "user-expiration-delay",
					Value: 3600,
					Usage: "Set expiration delay for users in seconds (0 means no expiration).",
				},
				cli.IntFlag{
					Name:  "user-expiration-min-followers",
					Value: 1000,
					Usage: "Set the min count of followers needed for a user to expire.",
				},
				cli.IntFlag{
					Name:  "user-refresh-delay",
					Value: 60,
					Usage: "Set the delay for refresh users in seconds.",
				},
				cli.StringFlag{
					Name:  "main-repository",
					Value: "richardlt/stargazer",
					Usage: "Set the path for main repository.",
				},
				cli.IntFlag{
					Name:  "main-repository-scan-delay",
					Value: 60,
					Usage: "Set the delay for main repository scanner in seconds.",
				},
				cli.IntFlag{
					Name:  "task-repository-scan-delay",
					Value: 60,
					Usage: "Set the delay for task repository scanner in seconds.",
				},
			),
			Action: func(c *cli.Context) error {
				level, err := logrus.ParseLevel(c.String("log-level"))
				if err != nil {
					return errors.Wrap(err, "invalid given log level")
				}

				return crawler.Start(config.Crawler{
					Common: config.Common{
						LogLevel: level,
					},
					MgoURI:                     c.String("mongo-uri"),
					GHToken:                    c.String("token"),
					UserExpirationDelay:        c.Int64("user-expiration-delay"),
					UserExpirationMinFollowers: c.Int64("user-expiration-min-followers"),
					UserRefreshDelay:           c.Int64("user-refresh-delay"),
					MainRepository:             c.String("main-repository"),
					MainRepositoryScanDelay:    c.Int64("main-repository-scan-delay"),
					TaskRepositoryScanDelay:    c.Int64("task-repository-scan-delay"),
					Database: config.Database{
						Host:     c.String("db-host"),
						Port:     c.Int64("db-port"),
						SSL:      c.Bool("db-ssl"),
						Name:     c.String("db-name"),
						User:     c.String("db-user"),
						Password: c.String("db-pass"),
					},
				})
			},
		},
		{
			Name: "web",
			Flags: append(globalFlags,
				cli.Int64Flag{
					Name:  "port",
					Value: 8080,
					Usage: "Stargazer webserver port",
				},
			),
			Action: func(c *cli.Context) error {
				level, err := logrus.ParseLevel(c.String("log-level"))
				if err != nil {
					return errors.WithStack(err)
				}

				return web.Start(config.Web{
					Common: config.Common{
						LogLevel: level,
					},
					Port: c.Int64("port"),
					Database: config.Database{
						Host:     c.String("db-host"),
						Port:     c.Int64("db-port"),
						SSL:      c.Bool("db-ssl"),
						Name:     c.String("db-name"),
						User:     c.String("db-user"),
						Password: c.String("db-pass"),
					},
				})
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Errorf("%+v", err)
	}
}
