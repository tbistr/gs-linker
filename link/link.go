package link

// TODO: Rewrite by bimap.
// TODO: Add testing.

import (
	"fmt"
	"log"
	"time"

	gh "github.com/tbistr/gs-linker/github"
	sl "github.com/tbistr/gs-linker/slack"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Client control DB to keep links.
type Client struct {
	conf *DbConfig
	db   *gorm.DB
}

// Link is bidirectional map between slack thread and github issue.
type Link struct {
	ID uint `gorm:"primarykey"`
	// "OnDelete:CASCADE" means that deleting this Link will cascade delete gh.Thread and sl.Thread from DB.
	Gh *gh.Thread `gorm:"constraint:OnDelete:CASCADE;"`
	Sl *sl.Thread `gorm:"constraint:OnDelete:CASCADE;"`
}

// DbConfig is config and credential for DB.
// Assumes mysql.
type DbConfig struct {
	UserName string // root
	Pass     string // passwd
	Addr     string // db
	DbName   string // database
}

// dsn makes dsn string.
func (conf *DbConfig) dsn() string {
	// root:passwd@tcp(db:3306)/database?charset=utf8&parseTime=True&loc=Local
	return conf.UserName + ":" + conf.Pass + "@tcp(" + conf.Addr + ":3306)/" + conf.DbName + "?charset=utf8&parseTime=True&loc=Local"

}

const DB_CON_RETRY = 20

// initDB connects to DB server.
func initDB(dsn string) (*gorm.DB, error) {
	var db *gorm.DB
	var err error
	for i := 0; i < DB_CON_RETRY; i++ {
		time.Sleep(2 * time.Second)
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			return db, nil
		}
	}
	return nil, fmt.Errorf("cannnot open db(tried %dtimes): %w", DB_CON_RETRY, err)
}

// New creates Client.
func New(conf *DbConfig) *Client {
	log.Println("connectiong db...")
	db, err := initDB(conf.dsn())
	if err != nil {
		log.Fatalf("failed to init db: %v\n", err)
	}
	log.Println("connected db")

	log.Println("migrating db...")
	if err := db.AutoMigrate(&Link{}, &sl.Thread{}, &gh.Thread{}); err != nil {
		log.Fatalf("failed to migrate: %v\n", err)
	}
	log.Println("migrated db")

	return &Client{
		db:   db,
		conf: conf,
	}
}
