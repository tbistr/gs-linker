package gh

import (
	"context"

	"github.com/google/go-github/v45/github"
	"github.com/k0kubun/pp/v3"
)

// TODO: send as slack user name and image?

func (client *Client) CreateIssueComment(ctx context.Context, thread *Thread, body string) {
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

func (client *Client) CreatePrComment(ctx context.Context, thread *Thread, body string) {
	comment := &github.IssueComment{
		Body: &body,
	}
	// https://dev.classmethod.jp/articles/get-and-post-comment-on-pull-request-with-github-api-v3/
	// https://docs.github.com/en/rest/issues/comments
	// > Every pull request is an issue, but not every issue is a pull request.
	// TODO: unite Issue~func and Pr~func?
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
