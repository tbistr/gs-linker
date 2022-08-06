package gh

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
)

// VerifyURL parses and verifies the URL and returns an error if the URL is incorrect.
// https://pkg.go.dev/github.com/google/go-github/v45@v45.0.0/github?utm_source=gopls#Issue
// > Note: As far as the GitHub API is concerned, every pull request is an issue, but not every issue is a pull request. Some endpoints, events, and webhooks may also return pull requests via this struct.
func (client *Client) VerifyURL(rawURL string) (owner, repo string, num int, err error) {
	log.Printf("verifing URL: %#v\n", rawURL)
	owner, repo, num, err = parseURL(rawURL)
	if err != nil {
		return "", "", 0, err
	}
	// confirmation of existence
	if _, _, err := client.github.Issues.Get(context.Background(), owner, repo, num); err != nil {
		return "", "", 0, err
	}

	log.Printf("confirmed the issue (or pr) is exist. owner: %#v repo: %#v num: %#v\n", owner, repo, num)
	return owner, repo, num, nil
}

func parseURL(rawURL string) (owner, repo string, num int, err error) {
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
