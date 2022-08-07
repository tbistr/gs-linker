package sl

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

// HandleEvent returns handlerFunc for slack webhook event requests.
func (client *Client) HandleEvent() func(http.ResponseWriter, *http.Request) {
	// Refer exemple.
	// https://github.com/slack-go/slack/blob/master/examples/eventsapi/events.go
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Verify credentials.
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
		// Verify credentials.

		event, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionNoVerifyToken())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// For initial setup.
		// Check if the URL is accessible.
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
			// TODO: queuing events and avoid dupulicate fire.
			// TODO: extract functions.
			switch event := event.InnerEvent.Data.(type) {
			case *slackevents.AppMentionEvent:
				// get gslinker's user id at the first event.
				// if client.config.mentionedText == "" {
				// 	log.Printf("first mentioned event fired. mentionedText reloaded: %v\n", client.config.mentionedText)
				// }
				// TODO: consider if Mentioned as single msg. (event.ThreadTimeStamp=="")
				log.Printf("catch mentioned event: \n%#v\n", event)
				c, rawURL, err := parseCommand(event.Text)
				if err != nil {
					// TODO: show help?
					// TODO: add response?
					log.Println(err)
					return
				}

				thread := &Thread{Channel: event.Channel, TS: event.ThreadTimeStamp}
				switch c {
				case subscribe:
					log.Println("fire subscribe command")
					// HandleEvent must returns in 3s.
					// https://api.slack.com/apis/connections/events-api#the-events-api__responding-to-events
					// > Your app should respond to the event request with an HTTP 2xx within three seconds. If it does not, we'll consider the event delivery attempt failed. After a failure, we'll retry three times, backing off exponentially.
					go client.handleSub(client, thread, rawURL)
				case unsubscribe:
					log.Println("fire unsubscribe command")
					go client.handleUnsub(client, thread)
				case summary:
					log.Println("fire summary command")
					go client.handleSummary(client, thread)
				}
				return
			case *slackevents.MessageEvent:
				// It is impossible to determine from the event whether the message is a reply or not.
				if event.ThreadTimeStamp == "" {
					log.Printf("dispose message event(it is not from thread): \n%#v\n", event.BotID)
					return
				}
				// It also fires `AppMentionEvent`
				// Ignore mentioned msg that have already been handled.
				if strings.Contains(event.Text, client.config.mentionedText) {
					log.Printf("dispose message event(it is mentioned): \n%#v\n", event)
					return
				}

				log.Printf("catch message event: \n%#v\n", event)
				go client.onMsgSent(client, &Thread{Channel: event.Channel, TS: event.ThreadTimeStamp}, event.Text)
				return
			}
		}
	}
}
