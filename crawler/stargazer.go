package crawler

import (
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
			logrus.Info("stargazer routine: update repository %s in database", r.Path)
			r.Data = o
			if err := dbClient.updateRepository(r); err != nil {
				return err
			}
		}

		logrus.Info("stargazer routine: load stargazers from Github")

		os, err := ghClient.getRepositoryStargazer(r.Path)
		if err != nil {
			return err
		}

		logrus.Infof("stargazer routine: delete all stargazers for repository %s in database", r.Path)
		if err := dbClient.deleteStargazers(r.ID); err != nil {
			return err
		}

		ss := make([]stargazer, len(os))
		for i := range os {
			ss[i].RepositoryID = r.ID
			ss[i].RepositoryPath = r.Path
			ss[i].Data = os[i]
		}

		logrus.Infof("stargazer routine: insert %d stargazers for repository %s in database", len(ss), r.Path)
		if err := dbClient.insertStargazers(ss); err != nil {
			return err
		}
	}

	return nil
}
