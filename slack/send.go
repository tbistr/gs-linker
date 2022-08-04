package sl

import "github.com/slack-go/slack"

func (client *Client) SendMsg(thread *Thread, text string) error {
	_, _, err := client.slack.PostMessage(
		thread.Channel,
		slack.MsgOptionTS(thread.TS),
		slack.MsgOptionText(text, false),
	)
	return err
}
