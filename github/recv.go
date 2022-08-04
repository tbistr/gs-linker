package gh

import (
	"log"
	"net/http"

	"github.com/google/go-github/v45/github"
)

// HandleEvent returns handlerFunc for github webhook event requests.
func (client *Client) HandleEvent() func(http.ResponseWriter, *http.Request) {
	// Refer example.
	// https://github.com/google/go-github#webhooks
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
			// I couldnt find the doc about user-type.
			// https://docs.github.com/ja/rest/users/users
			if event.Comment.GetUser().GetType() == "User" {
				log.Printf("catch issue comment event from: %s\n", event.Comment.GetURL())
				// Get~ method avoids nil references.
				// (If the structure is nil, it returns a zero-value.)
				// memo: needs guard?
				thread := &Thread{
					// see objects.
					// https://docs.github.com/en/graphql/reference/objects
					Owner: event.GetRepo().GetOwner().GetLogin(), // GetName doesnt return user name.
					Repo:  event.GetRepo().GetName(),
					Num:   event.GetIssue().GetNumber(),
				}
				if err := client.onCommented(client, thread, event.GetComment()); err != nil {
					log.Println(err)
				}
			}
			return
			// case *github.PullRequestReviewCommentEvent:
			// pr comment is also issue comment.
			// https://docs.github.com/en/developers/webhooks-and-events/events/issue-event-types#commented
		}
	}
}
