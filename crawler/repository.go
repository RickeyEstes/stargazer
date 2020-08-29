package crawler

import (
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/richardlt/stargazer/database"
)

func execTaskRepositoryRoutine(pgClient *database.DB, mgoClient *databaseClient, ghClient *githubClient, mainRepo string) error {
	es, err := pgClient.GetAllWithStatus(database.StatusRequested)
	if err != nil {
		return err
	}

	// TODO working pool
	for _, e := range es {
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
			exists, err = mgoClient.existsRepositoryOrganizationStargazer(mainRepo, owner)
			if err != nil {
				return err
			}
		}
		if !exists {
			logrus.Debugf("execTaskRepositoryRoutine: no stargazer found on main repo for %s", e.Repository)
			continue
		}

		// Load stargazer for repo
		if err := execStargazerRoutine(mgoClient, ghClient, e.Repository); err != nil {
			return err
		}

		// Compute stats for repo
		e.Stats = database.Stats{
			Evolution: nil,
			PerDays:   nil,
			Last10:    nil,
			Top10:     nil,
		}
		e.Status = database.StatusGenerated
		if err := pgClient.Update(&e); err != nil {
			return err
		}
	}

	return nil
}
