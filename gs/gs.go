package gs

import (
	"github.com/google/go-github/v43/github"
	"github.com/slack-go/slack"
)

// Linker links github and slack
type Linker struct {
	slack   *slack.Client
	github  *github.Client
	gConfig *githubConfig
	sConfig *slackConfig
}

// githubConfig is config for github
type githubConfig struct {
	secret []byte
}

// slackConfig is config for slack
type slackConfig struct {
}

// New returns new client
func New() *Linker {
	return &Linker{}
}
