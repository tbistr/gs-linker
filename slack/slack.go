package sl

import (
	"log"

	"github.com/slack-go/slack"
)

type Client struct {
	slack       *slack.Client
	config      *config
	onMentioned OnMentionedFunc

	// It is impossible to determine from the event whether the message is a reply or not.
	// If ThreadTS is "", go through, but the registrant should determine if the message is a reply to the bot.
	onMsgSent OnMsgSentFunc
}

// config for slack
type config struct {
	token         string
	signingSecret string
}

// Thread is info to designate slack thread.
type Thread struct {
	Channel string
	TS      string
}

type OnMentionedFunc func(client *Client, thread *Thread, text string) error
type OnMsgSentFunc func(client *Client, thread *Thread, text string) error

func New(token, signingSecret string) *Client {
	return &Client{
		slack: slack.New(token),
		config: &config{
			token:         token,
			signingSecret: signingSecret,
		},
		onMentioned: func(client *Client, thread *Thread, text string) error {
			// TODO: print thread info.
			log.Println("slack.Client.onMentioned is not registered.")
			return nil
		},
		onMsgSent: func(client *Client, thread *Thread, text string) error {
			log.Println("slack.Client.onMsgSent is not registered.")
			return nil
		},
	}
}

// We want to use ghClient in On~Handlers.
// So, We shouldnt assign handlers in New func.

func (Client *Client) RegisterOnMentioned(f OnMentionedFunc) {
	Client.onMentioned = f
}

func (Client *Client) RegisterOnMsgSent(f OnMsgSentFunc) {
	Client.onMsgSent = f
}
