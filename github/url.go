package gh

import (
	"net/url"
	"strconv"
	"strings"
)

func IsIssue(maybeURL *url.URL) (string, string, int, bool) {
	if maybeURL.Host != "github.com" {
		return "", "", 0, false
	}

	// assume that []string{"", "owner", "repo", "pull", "0"} or []string{"", "owner", "repo", "issues", "0"}
	params := strings.Split(maybeURL.Path, "/")
	if len(params) != 5 {
		return "", "", 0, false
	}
	if params[3] != "issues" {
		return "", "", 0, false
	}

	issueNum, err := strconv.Atoi(params[4])
	if err != nil {
		return "", "", 0, false
	}

	return params[1], params[2], issueNum, true
}

func IsPR(maybeURL *url.URL) (string, string, int, bool) {
	if maybeURL.Host != "github.com" {
		return "", "", 0, false
	}

	// assume that []string{"", "owner", "repo", "pull", "0"} or []string{"", "owner", "repo", "issues", "0"}
	params := strings.Split(maybeURL.Path, "/")
	if len(params) != 5 {
		return "", "", 0, false
	}
	if params[3] != "issues" {
		return "", "", 0, false
	}

	issueNum, err := strconv.Atoi(params[4])
	if err != nil {
		return "", "", 0, false
	}

	return params[1], params[2], issueNum, true
}
