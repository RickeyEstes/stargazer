package main

import (
	"context"
	"os"
	"time"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	cli "gopkg.in/urfave/cli.v1"
)

const (
	baseURL = "https://api.github.com"
)

func main() {
	app := cli.NewApp()

	app.Commands = []cli.Command{{
		Name: "start",
		Flags: []cli.Flag{
			cli.StringSliceFlag{
				Name:   "repository",
				EnvVar: "STARGAZER_REPOSITORY",
				Usage:  "owner/repository",
			},
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
			cli.StringFlag{
				Name:  "log-level",
				Value: "warning",
				Usage: "[panic fatal error warning info debug]",
			},
			cli.IntFlag{
				Name:  "user-expiration-delay",
				Usage: "Set expiration delay for users in hours (0 means no expiration).",
			},
			cli.IntFlag{
				Name:  "user-expiration-min-followers",
				Usage: "Set the min count of followers needed for a user to expire.",
			},
		},
		Action: start,
	}}

	if err := app.Run(os.Args); err != nil {
		logrus.Errorf("%+v", err)
	}
}

func start(c *cli.Context) error {
	level, err := logrus.ParseLevel(c.String("log-level"))
	if err != nil {
		return errors.Wrap(err, "invalid given log level")
	}
	logrus.SetLevel(level)

	// init database
	client, err := mongo.NewClient(c.String("mongo-uri"))
	if err != nil {
		return errors.WithStack(err)
	}
	if err := client.Connect(context.Background()); err != nil {
		return errors.WithStack(err)
	}

	db := client.Database("stargazer")
	dbClient := &databaseClient{db}
	if err := dbClient.init(); err != nil {
		return err
	}

	ghClient := &githubClient{c.String("token")}

	rs := c.StringSlice("repository")

	for i := range rs {
		go func(repo string) {
			for {
				if err := execStargazerRoutine(dbClient, ghClient, repo); err != nil {
					logrus.Errorf("%+v", err)
				}

				logrus.Infof("main: waiting 30s for next repo %s update", repo)
				time.Sleep(30 * time.Second)
			}
		}(rs[i])
	}

	go func() {
		for {
			if err := execUserRoutine(dbClient, ghClient, c.Int("user-expiration-delay"),
				c.Int("user-expiration-min-followers")); err != nil {
				logrus.Errorf("%+v", err)
			}

			logrus.Info("main: waiting 30s for next user update")
			time.Sleep(30 * time.Second)
		}
	}()

	date := time.Now()
	for {
		logrus.Infof("main: running since %s", time.Since(date).String())
		time.Sleep(time.Minute)
	}
}
