package gs

import (
	"log"
	"net/http"

	"github.com/google/go-github/v43/github"
)

// HandleGithubEvent returns handlerFunc for github webhook event request.
func (linker *Linker) HandleGithubEvent() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		payload, err := github.ValidatePayload(r, linker.gConfig.secret)
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
			if err := processIssueComment(event); err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	}
}

// processIssueComment processes IssueCommentEvent
func processIssueComment(event *github.IssueCommentEvent) error {
	return nil
}

// func newGithubClient(installationID int64) (*github.Client, error) {
// 	tr := http.DefaultTransport
// 	itr, err := ghinstallation.NewKeyFromFile(tr, appID, installationID, "private-key.pem")
// 	if err != nil {
// 		return nil, err
// 	}

// 	return github.NewClient(&http.Client{
// 		Transport: itr,
// 		Timeout:   5 * time.Second,
// 	}), nil
// }
