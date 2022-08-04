package sl

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

// HandleEvent returns handlerFunc for slack webhook event requests.
func (client *Client) HandleEvent() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		sv, err := slack.NewSecretsVerifier(r.Header, client.config.signingSecret)
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
				// TODO: consider if Mentioned as single msg. (event.ThreadTimeStamp=="")
				if err := client.onMentioned(client, &Thread{Channel: event.Channel, TS: event.ThreadTimeStamp}, event.Text); err != nil {
					log.Println(err)
					w.WriteHeader(http.StatusInternalServerError)
				}
				return
			case *slackevents.MessageEvent:
				// It is impossible to determine from the event whether the message is a reply or not.
				// If ThreadTS is "", go through, but the registrant should determine if the message is a reply to the bot.
				if event.ThreadTimeStamp == "" {
					return
				}
				// If mentioned, handlerFunc is called twice.
				// TODO: ignore mentioned msg that have already been handled.
				if err := client.onMsgSent(client, &Thread{Channel: event.Channel, TS: event.ThreadTimeStamp}, event.Text); err != nil {
					log.Println(err)
					w.WriteHeader(http.StatusInternalServerError)
				}
				return
			}
		}
	}
}
