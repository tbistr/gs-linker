package gs

import (
	"github.com/slack-go/slack"
)

// Linker links github and slack
type Linker struct {
	slack *slack.Client
	// github  *github.Client
	gConfig *GithubConfig
	sConfig *SlackConfig
}

// GithubConfig is config for github
type GithubConfig struct {
	secret []byte
}

// SlackConfig is config for slack
type SlackConfig struct {
	token         string
	signingSecret string
}

// New returns new client
func New(gc *GithubConfig, sc *SlackConfig) *Linker {
	s := slack.New(sc.token)
	return &Linker{
		slack:   s,
		gConfig: gc,
		sConfig: sc,
	}
}
