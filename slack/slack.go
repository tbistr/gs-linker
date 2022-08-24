package sl

import (
	"log"

	"github.com/slack-go/slack"
)

type Client struct {
	slack  *slack.Client
	config *config

	handleSub     HandleSubFunc
	handleUnsub   HandleUnsubFunc
	handleSummary HandleSummaryFunc

	onMsgSent OnMsgSentFunc
}

// config for slack
type config struct {
	token         string
	signingSecret string
	// userID is mentioned id string.
	// assumes `<@userID> this is mentioned message`.
	mentionedText string
}

// Thread is info to designate slack thread.
type Thread struct {
	ID      uint `gorm:"primarykey"`
	LinkID  uint
	Channel string
	TS      string
}

// avoid table name conflict.
func (Thread) TableName() string {
	return "slack_threads"
}

// command means how it is processed in OnMentioned.
type command string

// TODO: add abbreviation
const (
	// subscribe thread
	subscribe command = "subscribe"
	// unsubscribe thread
	unsubscribe = "unsubscribe"
	// post a thread summary
	summary = "summary"
)

type HandleSubFunc func(client *Client, thread *Thread, rawURL string)
type HandleUnsubFunc func(client *Client, thread *Thread)
type HandleSummaryFunc func(client *Client, thread *Thread)

// It is impossible to determine from the event whether the message is a reply or not.
// If ThreadTS is "", go through, but the registrant should determine if the message is a reply to the bot.
type OnMsgSentFunc func(client *Client, thread *Thread, text string)

func New(token, signingSecret, botUserID string) *Client {
	log.Println("creating github client")
	defer log.Println("created github client")
	return &Client{
		slack: slack.New(token),
		config: &config{
			token:         token,
			signingSecret: signingSecret,
			mentionedText: "<@" + botUserID + ">",
		},
		handleSub: func(client *Client, thread *Thread, rawURL string) {
			log.Println("slack.Client.handleSub is not registered")
			log.Printf("thread info: %#v\n", thread)
		},
		handleUnsub: func(client *Client, thread *Thread) {
			log.Println("slack.Client.handleUnsub is not registered")
			log.Printf("thread info: %#v\n", thread)
		},
		handleSummary: func(client *Client, thread *Thread) {
			log.Println("slack.Client.handleSummary is not registered")
			log.Printf("thread info: %#v\n", thread)
		},
		onMsgSent: func(client *Client, thread *Thread, text string) {
			log.Println("slack.Client.onMsgSent is not registered")
			log.Printf("thread info: %#v\n", thread)
		},
	}
}

// We want to use ghClient in the handlers.
// So, We should not assign handlers in New func.

func (Client *Client) RegisterHandleSub(f HandleSubFunc) {
	Client.handleSub = f
}

func (Client *Client) RegisterHandleUnsub(f HandleUnsubFunc) {
	Client.handleUnsub = f
}

func (Client *Client) RegisterHandleSummary(f HandleSummaryFunc) {
	Client.handleSummary = f
}

func (Client *Client) RegisterOnMsgSent(f OnMsgSentFunc) {
	Client.onMsgSent = f
}
