package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Config is config and credential for DB.
// Assumes mysql.
type Config struct {
	UserName string // root
	Pass     string // passwd
	Addr     string // db
	DbName   string // database
}

type Client struct {
	conf *Config
	db   *sql.DB
}

// dsn makes dsn string.
func (conf *Config) dsn() string {
	// root:passwd@tcp(db:3306)/database?charset=utf8&parseTime=True&loc=Local
	return conf.UserName + ":" + conf.Pass + "@tcp(" + conf.Addr + ":3306)/" + conf.DbName + "?charset=utf8&parseTime=True&loc=Local"
}

const DB_CON_RETRY = 20

// InitDB connects to DB server.
func InitDB(conf *Config) (*Client, error) {
	var db *sql.DB
	var err error
	for i := 0; i < DB_CON_RETRY; i++ {
		time.Sleep(2 * time.Second)
		db, err = sql.Open("mysql", conf.dsn())
		if err == nil {
			err = db.Ping()
			if err != nil {
				continue
			}
			return &Client{conf: conf, db: db}, nil
		}
	}
	return nil, fmt.Errorf("cannnot open db(tried %dtimes): %w", DB_CON_RETRY, err)
}
