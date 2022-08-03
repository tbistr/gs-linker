package gh

import (
	"context"

	"github.com/google/go-github/v45/github"
	"github.com/k0kubun/pp/v3"
)

func (client *Client) CreateIssueComment(ctx context.Context, body string) {
	comment := &github.IssueComment{
		Body: &body,
	}
	result, resp, err := client.github.Issues.CreateComment(
		ctx,
		"tbistr",
		"gs-linker",
		1,
		comment,
	)
	if err != nil {
		pp.Println(err)
	}
	pp.Println(result)
	pp.Println(resp)
}
