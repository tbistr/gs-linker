package gs

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/k0kubun/pp"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

// HandleSlackEvent returns handlerFunc for github webhook event request.
func (gs *Linker) HandleSlackEvent() func(http.ResponseWriter, *http.Request) {
	token := "xoxb-hogehuga"
	signingSecret := "0000"

	client := slack.New(token)
	return func(w http.ResponseWriter, r *http.Request) {
		pp.Println("handle func!")
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		sv, err := slack.NewSecretsVerifier(r.Header, signingSecret)
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
		eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionNoVerifyToken())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if eventsAPIEvent.Type == slackevents.URLVerification {
			var r *slackevents.ChallengeResponse
			err := json.Unmarshal([]byte(body), &r)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "text")
			w.Write([]byte(r.Challenge))
		}

		if eventsAPIEvent.Type == slackevents.CallbackEvent {
			innerEvent := eventsAPIEvent.InnerEvent
			switch ev := innerEvent.Data.(type) {
			case *slackevents.AppMentionEvent:
				pp.Println(innerEvent)
				a, b, c, err := client.SendMessage(
					ev.Channel,
					slack.MsgOptionText("hoge huga", false),
					slack.MsgOptionTS(ev.TimeStamp),
					slack.MsgOptionUsername("username change test"),
				)
				pp.Println(a)
				pp.Println(b)
				pp.Println(c)
				if err != nil {
					pp.Println(err)
				}
			}
		}
	}

}
