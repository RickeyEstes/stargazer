package config

import "github.com/sirupsen/logrus"

type Common struct {
	LogLevel logrus.Level
}

type Crawler struct {
	Common
	GHToken                         string
	MgoURI                          string
	UserExpirationDelay             int64
	UserExpirationMinFollowers      int64
	UserRefreshDelay                int64
	MainRepository                  string
	MainRepositoryScanDelay         int64
	TaskRepositoryScanDelay         int64
	TaskRepositoryMaxStargazerPages int64
	Database                        Database
}

type Database struct {
	Host     string
	Port     int64
	SSL      bool
	Name     string
	User     string
	Password string
	Driver   string
}

type Web struct {
	Common
	Port     int64
	Database Database
}
