package crawler

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/richardlt/stargazer/config"
)

func Start(cfg config.Crawler) error {
	logrus.SetLevel(cfg.LogLevel)

	// init database
	client, err := mongo.NewClient(options.Client().ApplyURI(cfg.MgoURI))
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

	ghClient := &githubClient{cfg.GHToken}

	rs := cfg.Repositories

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
			if err := execUserRoutine(dbClient, ghClient, cfg.UserExpirationDelay,
				cfg.UserExpirationMinFollowers); err != nil {
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
