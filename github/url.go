package gh

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"

	"github.com/google/go-github/v45/github"
)

// VerifyURL parses and verifies the URL and returns an error if the URL is incorrect.
// https://pkg.go.dev/github.com/google/go-github/v45@v45.0.0/github?utm_source=gopls#Issue
// > Note: As far as the GitHub API is concerned, every pull request is an issue, but not every issue is a pull request. Some endpoints, events, and webhooks may also return pull requests via this struct.
func (client *Client) VerifyURL(rawURL string) (*github.Issue, error) {
	owner, repo, num, err := parseURL(rawURL)
	if err != nil {
		return nil, err
	}
	issue, res, err := client.github.Issues.Get(context.Background(), owner, repo, num)
	if err != nil {
		return nil, err
	}

	log.Printf("get issue info: %+v\n", res)
	return issue, nil
}

func parseURL(rawURL string) (owner, repo string, num int, err error) {
	log.Printf("parsing URL: %#v\n", rawURL)
	maybeURL, err := url.Parse(rawURL)
	if err != nil {
		return "", "", 0, err
	}

	if maybeURL.Host != "github.com" {
		return "", "", 0, fmt.Errorf("%#v is not github url", maybeURL)
	}

	// assumes that []string{"", "owner", "repo", "pull", "0"} or []string{"", "owner", "repo", "issues", "0"}
	params := strings.Split(maybeURL.Path, "/")
	if len(params) != 5 {
		return "", "", 0, fmt.Errorf("%#v is not assumed url", maybeURL)
	}
	owner = params[1]
	repo = params[2]
	num, err = strconv.Atoi(params[4])
	if err != nil {
		return "", "", 0, err
	}

	return owner, repo, num, nil
}
