package web

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/richardlt/stargazer/database"
)

type Server struct {
	router *mux.Router
	db     *database.DB
}

func (s *Server) initRouter() error {
	r := mux.NewRouter()

	ht, err := template.New("not-found").Parse(notFoundTemplate)
	if err != nil {
		return errors.WithStack(err)
	}
	rt, err := template.New("repository").Parse(repositoryTemplate)
	if err != nil {
		return errors.WithStack(err)
	}

	r.HandleFunc("/{organization}/{repository}", s.repositoryPageHandler(rt))
	r.NotFoundHandler = s.notFoundPageHandler(ht)

	s.router = r

	return nil
}

func (s *Server) Close() {
	s.db.Close()
}

func (s *Server) Start(port int64) error {
	logrus.Infof("Starting webserver at :%d", port)
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      s.router,
	}
	return errors.WithStack(srv.ListenAndServe())
}
