package link

import (
	"fmt"
	"log"
	"reflect"

	gh "github.com/tbistr/gs-linker/github"
	sl "github.com/tbistr/gs-linker/slack"
)

// Sub subscribes a link.
//
// Returns error if already subscribed.
func (links *Links) Sub(g *gh.Thread, s *sl.Thread) error {
	links.mu.Lock()
	defer links.mu.Unlock()

	for _, l := range links.links {
		if reflect.DeepEqual(s, l.Sl) {
			return fmt.Errorf("already subscribed by slack thread: %+v", l)
		}
	}

	l := &Link{Gh: g, Sl: s}
	links.links = append(links.links, l)
	log.Printf("link is created: %+v\n", l)
	return nil
}

// UnSub unsbscribes specified link.
//
// Assumes that the link is unsbscribed only from slack.
// Returns error if yet subscribed.
func (links *Links) UnSub(s *sl.Thread) error {
	links.mu.Lock()
	defer links.mu.Unlock()

	// TODO: linked list can reduce computation.
	for i, l := range links.links {
		if reflect.DeepEqual(s, l.Sl) {
			// should not use `res = l`, because I dont know whether overwittern by del-sequence or not.
			res := links.links[len(links.links)-1]
			// https://zenn.dev/mattn/articles/31dfed3c89956d
			links.links[i] = links.links[len(links.links)-1]
			links.links[len(links.links)-1] = nil
			links.links = links.links[:len(links.links)-1]

			log.Printf("link is removed: %+v\n", res)
			return nil
		}
	}

	return fmt.Errorf("not subscribed: %+v", s)
}
