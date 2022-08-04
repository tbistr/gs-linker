package gh

import (
	"context"
	"log"

	"github.com/google/go-github/v45/github"
)

// TODO: send as slack user name and image?

func (client *Client) CreateComment(ctx context.Context, thread *Thread, body string) {
	comment := &github.IssueComment{
		Body: &body,
	}
	res, _, err := client.github.Issues.CreateComment(
		ctx,
		thread.Owner,
		thread.Repo,
		thread.Num,
		comment,
	)
	if err != nil {
		log.Printf("cant create comment: %v\n", err)
	}
	log.Printf("comment is created: %s\n", res.GetHTMLURL())
}
