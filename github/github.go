package gh

import (
	"log"
	"net/http"
	"time"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/v45/github"
)

type Client struct {
	github      *github.Client
	config      *config
	onCommented OnCommentedFunc
}

// config for github
type config struct {
	secret []byte
}

// Thread is info to designate github thread.
type Thread struct {
	ID     uint `gorm:"primarykey"`
	LinkID uint
	Owner  string
	Repo   string
	Num    int
}

// avoid table name conflict.
func (Thread) TableName() string {
	return "github_threads"
}

type OnCommentedFunc func(client *Client, thread *Thread, comment *github.IssueComment) error

func New(appID, installationID int64) *Client {
	log.Println("creating github client")
	defer log.Println("created github client")

	g, err := newGithubClient(appID, installationID)
	if err != nil {
		log.Fatal(err)
	}
	return &Client{
		github: g,
		config: &config{
			secret: []byte{},
		},
		onCommented: func(client *Client, thread *Thread, comment *github.IssueComment) error {
			log.Println("github.Client.onIssueCommented is not registered.")
			return nil
		},
	}
}

func newGithubClient(appID, installationID int64) (*github.Client, error) {
	itr, err := ghinstallation.NewKeyFromFile(http.DefaultTransport, appID, installationID, "private-key.pem")
	if err != nil {
		return nil, err
	}

	return github.NewClient(&http.Client{
		Transport: itr,
		Timeout:   5 * time.Second,
	}), nil
}

// We want to use slClient in On~Handlers.
// So, We shouldnt assign handlers in New func.

func (Client *Client) RegisterOnCommented(f OnCommentedFunc) {
	Client.onCommented = f
}
