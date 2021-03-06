package database

import (
	"fmt"
	// "strings"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/pkg/errors"

	"github.com/paper2code-bot/stargazer/config"
)

func New(cfg config.Database) (*DB, error) {

	var dsn string
	switch cfg.Driver {
	// case "postgres":
	// 	opts = []string{
	// 		"host=" + cfg.Host,
	// 		fmt.Sprintf("port=%d", cfg.Port),
	// 		"dbname=" + cfg.Name,
	// 		"user=" + cfg.User,
	// 		"password=" + cfg.Password,
	// 	}
	// 	if !cfg.SSL {
	// 		opts = append(opts, "sslmode=disable")
	// 	}
	case "mysql":
		dsn = fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4,utf8&parseTime=True", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
	case "sqlite":
		fallthrough
	default:
		dsn = fmt.Sprintf("%s.db", cfg.Name)
	}

	db, err := gorm.Open(cfg.Driver, dsn) // strings.Join(opts, " "))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	res := db.AutoMigrate(&Entry{})
	if res.Error != nil {
		return nil, errors.WithStack(res.Error)
	}

	return &DB{
		db: db,
	}, nil
}

type DB struct {
	db *gorm.DB
}

func (d *DB) Close() {
	d.db.Close()
}

func (d *DB) Get(repo string) (*Entry, error) {
	var e Entry
	res := d.db.First(&e, "repository = ?", repo)
	if res.Error != nil {
		return nil, errors.WithStack(res.Error)
	}
	return &e, nil
}

func (d *DB) GetAllWithStatus(status Status) ([]Entry, error) {
	var es []Entry
	res := d.db.Find(&es, "status = ?", status)
	if res.Error != nil {
		return nil, errors.WithStack(res.Error)
	}
	return es, nil
}

func (d *DB) Create(e *Entry) error {
	res := d.db.Create(e)
	return errors.WithStack(res.Error)
}

func (d *DB) Update(e *Entry) error {
	res := d.db.Save(e)
	e.UpdatedAt = time.Now()
	return errors.WithStack(res.Error)
}

func (d *DB) Delete(repo string) error {
	res := d.db.Exec("DELETE FROM entries WHERE repository = ?", repo)
	return errors.WithStack(res.Error)
}
