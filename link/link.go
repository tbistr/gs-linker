package link

import (
	"log"

	gh "github.com/tbistr/gs-linker/github"
	"github.com/tbistr/gs-linker/link/db"
	sl "github.com/tbistr/gs-linker/slack"
)

// Client control DB to keep links.
type Client struct {
	db *db.Client
}

// Link is bidirectional map between slack thread and github issue.
type Link struct {
	ID uint
	Gh *gh.Thread
	Sl *sl.Thread
}

// New creates Client.
func New(conf *db.Config) *Client {
	log.Println("connectiong db...")
	dbClient, err := db.New(conf)
	if err != nil {
		log.Fatalf("failed to init db: %v\n", err)
	}
	log.Println("connected db")

	log.Println("initializing db...")
	if err := dbClient.TouchTables(); err != nil {
		log.Fatalf("failed to initialize: %v\n", err)
	}
	log.Println("initialized db")

	return &Client{
		db: dbClient,
	}
}
