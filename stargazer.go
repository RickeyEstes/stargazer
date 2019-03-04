package main

import (
	"time"

	"github.com/sirupsen/logrus"
)

func execStargazerRoutine(dbClient *databaseClient, ghClient *githubClient, repo string) error {
	logrus.Info("stargazer routine: start")

	logrus.Info("stargazer routine: get repository from Github")
	o, err := ghClient.getRepository(repo)
	if err != nil {
		return err
	}

	logrus.Info("stargazer routine: get repository from database")
	r, err := dbClient.getRepository(repo)
	if err != nil {
		return err
	}

	repoExists := r != nil

	if !repoExists {
		r = &repository{
			Path: repo,
			Data: o,
		}

		logrus.Info("stargazer routine: create repository in database")
		if err := dbClient.insertRepository(r); err != nil {
			return err
		}
	}

	// compare stargazers count
	githubStargazersCount := int(o["stargazers_count"].(float64))

	databaseStargazersCount, err := dbClient.countStargazers(r.ID)
	if err != nil {
		return err
	}

	logrus.Infof("stargazer routine: found %d stargazers from Github and %d in database", githubStargazersCount, databaseStargazersCount)
	change := !repoExists || githubStargazersCount != databaseStargazersCount

	// if counts are different then reload all stargazers
	if change {
		if repoExists {
			logrus.Info("stargazer routine: update repository in database")
			r.Data = o
			if err := dbClient.updateRepository(r); err != nil {
				return err
			}
		}

		logrus.Info("stargazer routine: load stargazers from Github")
		expectedPageCount := (githubStargazersCount / 100) + 1
		for page := 1; page <= expectedPageCount; page++ {
			logrus.Infof("stargazer routine: load stargazers page %d/%d from Github", page, expectedPageCount)
			os, err := ghClient.getRepositoryStargazer(r.Path, page)
			if err != nil {
				return err
			}

			logrus.Infof("stargazer routine: delete all stargazers for page %d in database", page)
			if err := dbClient.deleteStargazers(r.ID, page); err != nil {
				return err
			}

			ss := make([]stargazer, len(os))
			for i := range os {
				ss[i].RepositoryID = r.ID
				ss[i].RepositoryPath = r.Path
				ss[i].Page = page
				ss[i].Data = os[i]
			}

			logrus.Infof("stargazer routine: insert stargazers for page %d in database", page)
			if err := dbClient.insertStargazers(ss); err != nil {
				return err
			}

			logrus.Info("stargazer routine: wait 50ms")
			time.Sleep(50 * time.Millisecond)
		}
	}

	return nil
}
