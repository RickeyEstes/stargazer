package crawler

import (
	"github.com/sirupsen/logrus"
)

func loadStargazerForRepo(dbClient *databaseClient, ghClient *githubClient, repoO object, maxStargazerPageToScan int64, mainRepoPath string) error {
	repo := repoO["full_name"].(string)

	githubStargazersCount := int(repoO["stargazers_count"].(float64))

	logrus.Infof("stargazer routine: get repository %s from database", repo)
	r, err := dbClient.getRepository(repo)
	if err != nil {
		return err
	}

	repoExists := r != nil

	if !repoExists {
		r = &repository{
			Path: repo,
			Data: repoO,
		}

		logrus.Infof("stargazer routine: create repository %s in database", repo)
		if err := dbClient.insertRepository(r); err != nil {
			return err
		}
	}

	databaseStargazersCount := int(r.Data["stargazers_count"].(float64))

	logrus.Infof("stargazer routine: found %d stargazers from GH for repo %s and %d in database", githubStargazersCount, repo, databaseStargazersCount)
	change := !repoExists || githubStargazersCount != databaseStargazersCount

	// if counts are different then reload all stargazers
	if change {
		if repoExists {
			logrus.Infof("stargazer routine: update repository %s in database", r.Path)
			r.Data = repoO
			if err := dbClient.updateRepository(r); err != nil {
				return err
			}
		}

		expectedPageCount := int64((githubStargazersCount / 100) + 1)
		if expectedPageCount > 400 { // GH limit on page count is 400
			expectedPageCount = 400
		}
		logrus.Infof("stargazer routine: load stargazers for repo %s from Github from %d pages expected", r.Path, expectedPageCount)
		getPage := func(path string, page, expectedPageCount int64) ([]stargazer, error) {
			logrus.Infof("stargazer routine: load stargazers page %d for repo %s from Github", page, r.Path)
			os, err := ghClient.getRepositoryStargazerPage(r.Path, page)
			if err != nil {
				return nil, err
			}
			ss := make([]stargazer, len(os))
			for i := range os {
				ss[i].RepositoryID = r.ID
				ss[i].RepositoryPath = r.Path
				ss[i].Page = page
				ss[i].LastPage = page == expectedPageCount
				ss[i].Data = os[i]
			}
			return ss, nil
		}

		ss := make([]stargazer, 0, maxStargazerPageToScan*100)

		if expectedPageCount > maxStargazerPageToScan && r.Path != mainRepoPath {
			for i := int64(1); i <= maxStargazerPageToScan; i++ {
				var page int64
				if i == 1 {
					page = 1
				} else if i == maxStargazerPageToScan {
					page = expectedPageCount
				} else {
					page = i * (expectedPageCount / maxStargazerPageToScan)
				}
				os, err := getPage(r.Path, page, expectedPageCount)
				if err != nil {
					return err
				}
				ss = append(ss, os...)
			}
		} else {
			for i := int64(1); i <= expectedPageCount; i++ {
				os, err := getPage(r.Path, i, expectedPageCount)
				if err != nil {
					return err
				}
				ss = append(ss, os...)
			}
		}

		logrus.Infof("stargazer routine: delete all stargazers for repository %s in database", r.Path)
		if err := dbClient.deleteStargazers(r.ID); err != nil {
			return err
		}

		logrus.Infof("stargazer routine: insert %d stargazers for repository %s in database", len(ss), r.Path)
		if err := dbClient.insertStargazers(ss); err != nil {
			return err
		}
	}

	return nil
}
