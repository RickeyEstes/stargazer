package main

import (
	"os"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	cli "gopkg.in/urfave/cli.v1"

	"github.com/paper2code-bot/stargazer/config"
	"github.com/paper2code-bot/stargazer/crawler"
	"github.com/paper2code-bot/stargazer/web"
)

func main() {
	app := cli.NewApp()

	globalFlags := []cli.Flag{
		cli.StringFlag{
			Name:  "db-host",
			Value: "localhost",
			Usage: "Database URI",
		},
		cli.Int64Flag{
			Name:  "db-port",
			Value: 3306,
			Usage: "Database port",
		},
		cli.BoolFlag{
			Name:  "db-ssl",
			Usage: "Database ssl mode",
		},
		cli.StringFlag{
			Name:  "db-name",
			Value: "stargazer",
			Usage: "Database name",
		},
		cli.StringFlag{
			Name:  "db-user",
			Value: "stargazer",
			Usage: "Database user",
		},
		cli.StringFlag{
			Name:  "db-pass",
			Value: "stargazer",
			Usage: "Database password",
		},
		cli.StringFlag{
			Name:  "db-driver",
			Value: "mysql",
			Usage: "Database driver",
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
					Value: 30,
					Usage: "Set the delay for refresh users in seconds.",
				},
				cli.StringFlag{
					Name:  "main-repository",
					Value: "paper2code-bot/stargazer",
					Usage: "Set the path for main repository.",
				},
				cli.IntFlag{
					Name:  "main-repository-scan-delay",
					Value: 30,
					Usage: "Set the delay for main repository scanner in seconds.",
				},
				cli.IntFlag{
					Name:  "task-repository-scan-delay",
					Value: 30,
					Usage: "Set the delay for task repository scanner in seconds.",
				},
				cli.IntFlag{
					Name:  "task-repository-max-stargazer-pages",
					Value: 100,
					Usage: "Set the maximum stargazer pages to load for a repository.",
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
					MgoURI:                          c.String("mongo-uri"),
					GHToken:                         c.String("token"),
					UserExpirationDelay:             c.Int64("user-expiration-delay"),
					UserExpirationMinFollowers:      c.Int64("user-expiration-min-followers"),
					UserRefreshDelay:                c.Int64("user-refresh-delay"),
					MainRepository:                  c.String("main-repository"),
					MainRepositoryScanDelay:         c.Int64("main-repository-scan-delay"),
					TaskRepositoryScanDelay:         c.Int64("task-repository-scan-delay"),
					TaskRepositoryMaxStargazerPages: c.Int64("task-repository-max-stargazer-pages"),
					Database: config.Database{
						Host:     c.String("db-host"),
						Port:     c.Int64("db-port"),
						SSL:      c.Bool("db-ssl"),
						Name:     c.String("db-name"),
						User:     c.String("db-user"),
						Password: c.String("db-pass"),
						Driver:   c.String("db-driver"),
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
