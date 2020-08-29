package web

import (
	"html/template"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/richardlt/stargazer/database"
	"github.com/sirupsen/logrus"
)

func (s *Server) notFoundPageHandler(t *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := t.ExecuteTemplate(w, "not-found", map[string]interface{}{}); err != nil {
			logrus.Errorf("%+v", errors.WithStack(err))
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
}

func (s *Server) repositoryPageHandler(t *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		organization := vars["organization"]
		repository := vars["repository"]

		if organization == "" || repository == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		repoPath := organization + "/" + repository

		e, err := s.db.Get(repoPath)
		if err != nil && errors.Cause(err) != gorm.ErrRecordNotFound {
			logrus.Errorf("%+v", errors.WithStack(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if e == nil {
			e = &database.Entry{
				Repository: repoPath,
				Status:     database.StatusRequested,
			}
			if err := s.db.Create(e); err != nil {
				logrus.Errorf("%+v", errors.WithStack(err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			logrus.Debugf("New entry created for repository: %s", repoPath)
		} else {
			e.LastRequestedAt = time.Now()
			if err := s.db.Update(e); err != nil {
				logrus.Errorf("%+v", errors.WithStack(err))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			logrus.Debugf("Entry updated for repository: %s", repoPath)
		}

		if err := t.ExecuteTemplate(w, "repository", *e); err != nil {
			logrus.Errorf("%+v", errors.WithStack(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
