package link

import (
	"fmt"
	"log"

	gh "github.com/tbistr/gs-linker/github"
	sl "github.com/tbistr/gs-linker/slack"
)

// Sub subscribes a link.
//
// Returns error if already subscribed.
func (Client *Client) Sub(g *gh.Thread, s *sl.Thread) error {
	if Client.SearchByS(s) != nil {
		return fmt.Errorf("already subscribed by this thread: %+v", s)
	}

	Client.db.Create(g, s)
	// TODO: meaningless logs
	// ex. create link: &{0xc000210180 0xc00020a040}
	// log.Printf("create link: %+v\n", l)
	return nil
}

// UnSub unsbscribes specified link.
//
// Assumes that the link is unsbscribed only from slack.
// Returns error if yet subscribed.
func (Client *Client) UnSub(s *sl.Thread) error {
	g := Client.SearchByS(s)
	if g == nil {
		return fmt.Errorf("not subscribed: %+v", s)
	}

	Client.db.DeleteByS(s)
	log.Printf("remove link: %+v\n", Link{Gh: g, Sl: s})
	return nil
}
