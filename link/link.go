package link

// TODO: Rewrite by bimap.
// TODO: Add testing.

import (
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
	Protocol string // tcp(db:3306)
	DbName   string // database
}

const DB_CON_RETRY = 20

// New creates Client.
func New(conf DbConfig) *Client {
	// root:passwd@tcp(db:3306)/database?charset=utf8&parseTime=True&loc=Local
	dsn := conf.UserName + ":" + conf.Pass + "@" + conf.Protocol + "/" + conf.DbName + "?charset=utf8&parseTime=True&loc=Local"

	log.Println("connectiong db...")
	var db *gorm.DB
	for i := 1; i <= DB_CON_RETRY; i++ {
		time.Sleep(2 * time.Second)
		var err error
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			if i == DB_CON_RETRY {
				log.Fatalf("cannnot open db(tried %dtimes): %v\n", i, err)
			}
		} else {
			break
		}
	}
	log.Println("connected db")

	log.Println("migrating db...")
	if err := db.AutoMigrate(&Link{}, &sl.Thread{}, &gh.Thread{}); err != nil {
		log.Fatalf("failed to migrate: %v\n", err)
	}
	log.Println("migrated db")

	return &Client{
		db:   db,
		conf: &conf,
	}
}
