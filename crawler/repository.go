package crawler

import (
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/paper2code-bot/stargazer/database"
)

func execMainRepositoryRoutine(dbClient *databaseClient, ghClient *githubClient, repo string) error {
	logrus.Infof("execMainRepositoryRoutine: get main repository %s from Github", repo)

	o, err := ghClient.getRepository(repo)
	if err != nil {
		return err
	}

	githubStargazersCount := int(o["stargazers_count"].(float64))

	logrus.Infof("stargazer routine: get repository %s from database", repo)
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

		logrus.Infof("stargazer routine: create repository %s in database", repo)
		if err := dbClient.insertRepository(r); err != nil {
			return err
		}
	}

	databaseStargazersCount, err := dbClient.countStargazers(r.ID)
	if err != nil {
		return err
	}

	logrus.Infof("stargazer routine: found %d stargazers from GH for repo %s and %d in database", githubStargazersCount, repo, databaseStargazersCount)
	change := !repoExists || githubStargazersCount != databaseStargazersCount

	// if counts are different then reload all stargazers
	if change {
		if repoExists {
			logrus.Infof("stargazer routine: update repository %s in database", r.Path)
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

func execTaskRepositoryRoutine(pgClient *database.DB, mgoClient *databaseClient, ghClient *githubClient, mainRepo string, maxStargazerPageToScan int64, mainRepoPath string) error {
	es, err := pgClient.GetAllWithStatus(database.StatusRequested)
	if err != nil {
		return err
	}

	// TODO working pool
	for _, e := range es {
		logrus.Infof("execTaskRepositoryRoutine: starting scan for repository %s", e.Repository)
		rs := strings.Split(e.Repository, "/")
		if len(rs) != 2 {
			logrus.Warnf("execTaskRepositoryRoutine: invalid repository path %s", e.Repository)
			continue
		}
		owner := rs[0]

		// Check that the repository owner starred the main repository
		// For organization repository, check that one stargazer of the main repository is in the organization
		exists, err := mgoClient.existsRepositoryStargazer(mainRepo, owner)
		if err != nil {
			return err
		}
		if !exists {
			logrus.Debugf("execTaskRepositoryRoutine: no stargazer found on main repo for %s", e.Repository)
			continue
		}
		logrus.Debugf("execTaskRepositoryRoutine: starting compute stats for repo for %s", e.Repository)

		// Load stargazer for repo
		o, err := ghClient.getRepository(e.Repository)
		if err != nil {
			logrus.Warnf("execTaskRepositoryRoutine: repository not found on GH %s", e.Repository)
			continue
		}
		if err := loadStargazerForRepo(mgoClient, ghClient, o, maxStargazerPageToScan, mainRepoPath); err != nil {
			return err
		}

		// Compute evolution stats
		msPage, err := mgoClient.getRepoStarCountPerDaysAndPage(e.Repository)
		if err != nil {
			return err
		}
		if len(msPage) == 0 {
			continue
		}
		previousPage := int64(1)
		count := int64(0)
		e.Stats.Evolution = nil
		for i := range msPage {
			if previousPage < msPage[i].Page {
				pageGap := msPage[i].Page - previousPage
				if pageGap > 1 {
					count += (pageGap - 1) * 100
				}
			}
			previousPage = msPage[i].Page
			count += msPage[i].Count
			e.Stats.Evolution = append(e.Stats.Evolution, database.Measure{Date: msPage[i].Date, Count: count})
		}

		// Compute count per days stats
		ms, err := mgoClient.getRepoStarCountPerDays(e.Repository)
		if err != nil {
			return err
		}
		if len(ms) == 0 {
			continue
		}
		e.Stats.PerDays = nil
		for i := 30; i > 0; i-- {
			if len(ms) >= i {
				m := ms[len(ms)-i]
				e.Stats.PerDays = append(e.Stats.PerDays, database.Measure{Date: m.Date, Count: m.Count})
			}
		}

		e.Status = database.StatusGenerated
		if err := pgClient.Update(&e); err != nil {
			return err
		}
	}

	return nil
}
