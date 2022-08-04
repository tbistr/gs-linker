package gh

import (
	"log"
	"net/http"

	"github.com/google/go-github/v45/github"
)

// HandleEvent returns handlerFunc for github webhook event requests.
func (client *Client) HandleEvent() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		payload, err := github.ValidatePayload(r, client.config.secret)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		event, err := github.ParseWebHook(github.WebHookType(r), payload)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		switch event := event.(type) {
		case *github.IssueCommentEvent:
			// Get~ method avoids nil references.
			// (If the structure is nil, it returns a zero-value.)
			// memo: needs guard?
			owner := event.GetRepo().GetOwner().GetName()
			repo := event.GetRepo().GetName()
			num := event.GetIssue().GetNumber()
			thread := &Thread{
				SubType: ISSUE,
				Owner:   owner,
				Repo:    repo,
				Num:     num,
			}
			if err := client.onIssueCommented(client, thread, event.GetComment()); err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		case *github.PullRequestReviewCommentEvent:
			owner := event.GetRepo().GetOwner().GetName()
			repo := event.GetRepo().GetName()
			num := event.GetPullRequest().GetNumber()
			thread := &Thread{
				SubType: PR,
				Owner:   owner,
				Repo:    repo,
				Num:     num,
			}
			if err := client.onPrCommented(client, thread, event.GetComment()); err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}
	}
}
