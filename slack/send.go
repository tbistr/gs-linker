package sl

import (
	"log"

	"github.com/slack-go/slack"
)

func (client *Client) SendMsg(thread *Thread, text string) error {
	log.Printf("Sending msg to slack. channel: %s, thread ts: %s, text: %s\n", thread.Channel, thread.TS, text)
	_, _, err := client.slack.PostMessage(
		thread.Channel,
		slack.MsgOptionTS(thread.TS),
		slack.MsgOptionText(text, false),
	)
	return err
}
