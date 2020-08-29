package config

import "github.com/sirupsen/logrus"

type Common struct {
	LogLevel logrus.Level
}

type Crawler struct {
	Common
	Repositories               []string
	GHToken                    string
	MgoURI                     string
	UserExpirationDelay        int64
	UserExpirationMinFollowers int64
}

type Database struct {
	Host     string
	Port     int64
	SSL      bool
	Name     string
	User     string
	Password string
}

type Web struct {
	Common
	Port     int64
	Database Database
}
