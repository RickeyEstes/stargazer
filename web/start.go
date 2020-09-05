package web

import (
	"github.com/sirupsen/logrus"

	"github.com/paper2code-bot/stargazer/config"
	"github.com/paper2code-bot/stargazer/database"
)

func Start(cfg config.Web) error {
	logrus.SetLevel(cfg.LogLevel)

	db, err := database.New(cfg.Database)
	if err != nil {
		return err
	}

	s := &Server{db: db}
	if err := s.initRouter(); err != nil {
		return err
	}
	defer s.Close()

	return s.Start(cfg.Port)
}
