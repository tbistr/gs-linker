package gs

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

// HandleSlackEvent returns handlerFunc for github webhook event request.
func (gs *Linker) HandleSlackEvent() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		sv, err := slack.NewSecretsVerifier(r.Header, gs.sConfig.signingSecret)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if _, err := sv.Write(body); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err := sv.Ensure(); err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		event, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionNoVerifyToken())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if event.Type == slackevents.URLVerification {
			var r *slackevents.ChallengeResponse
			err := json.Unmarshal([]byte(body), &r)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "text")
			w.Write([]byte(r.Challenge))
		}

		if event.Type == slackevents.CallbackEvent {
			switch event := event.InnerEvent.Data.(type) {
			case *slackevents.AppMentionEvent:
				if err := processMentioned(event); err != nil {
					return
				}
			}
		}
	}
}

// processMentioned reacts to mentioned message.
func (gs *Linker) processMentioned(event *slackevents.AppMentionEvent) error {
	threadTS := event.ThreadTimeStamp
	for _, maybeURL := range strings.Split(event.Text, " ") {
		maybeURL, err := url.Parse(maybeURL)
		if err != nil {
			continue
		}

		owner, repo, num, ok := isIssueURL(maybeURL)
		if !ok {
		}
		owner, repo, num, ok = isPrURL(maybeURL)
		if !ok {
		}
	}

	// _, _, _, err := gs.slack.SendMessage(
	// 	event.Channel,
	// 	slack.MsgOptionText("hoge huga", false),
	// 	slack.MsgOptionTS(event.TimeStamp),
	// 	slack.MsgOptionUsername("username change test"),
	// )
	// if err != nil {
	// 	pp.Println(err)
	// }
	return nil
}

func isIssueURL(maybeURL *url.URL) (string, string, int, bool) {
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

func isPrURL(maybeURL *url.URL) (string, string, int, bool) {
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
