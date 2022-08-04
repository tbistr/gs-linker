package gh

import (
	"log"
	"net/http"
	"time"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/google/go-github/v45/github"
	"github.com/k0kubun/pp"
)

type Client struct {
	github           *github.Client
	config           *config
	onIssueCommented OnIssueCommentedFunc
	onPrCommented    OnPrCommentedFunc
}

// config for github
type config struct {
	secret []byte
}

// Thread is info to designate github thread.
type Thread struct {
	SubType SubType
	Owner   string
	Repo    string
	Num     int
}

// SubType is issue or pull request.
type SubType string

const (
	ISSUE = SubType("issue")
	PR    = SubType("pull request")
)

type OnIssueCommentedFunc func(client *Client, owner string, repo string, num int, comment *github.IssueComment) error
type OnPrCommentedFunc func(client *Client, owner string, repo string, num int, comment *github.PullRequestComment) error

func New(appID, installationID int64) *Client {
	g, err := newGithubClient(appID, installationID)
	if err != nil {
		panic(pp.Sprintln(err))
	}
	return &Client{
		github: g,
		config: &config{
			secret: []byte{},
		},
		onIssueCommented: func(client *Client, owner, repo string, num int, comment *github.IssueComment) error {
			log.Println("github.Client.onIssueCommented is not registered.")
			return nil
		},
		onPrCommented: func(client *Client, owner, repo string, num int, comment *github.PullRequestComment) error {
			log.Println("github.Client.onPrCommented is not registered.")
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

func (Client *Client) RegisterOnIssueCommented(f OnIssueCommentedFunc) {
	Client.onIssueCommented = f
}

func (Client *Client) RegisterOnPrCommented(f OnPrCommentedFunc) {
	Client.onPrCommented = f
}
