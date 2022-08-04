package link

// TODO: Rewrite by bimap.
// TODO: Add testing.

import (
	"sync"

	gh "github.com/tbistr/gs-linker/github"
	sl "github.com/tbistr/gs-linker/slack"
)

// Links consists of Link.
type Links struct {
	mu    sync.RWMutex
	links []*Link
}

// Link is bidirectional map between slack thread and github issue.
type Link struct {
	Gh gh.Thread
	Sl sl.Thread
}

// Sub subscribes a link.
//
// Returns subscribed *Link, nil if already subscribed.
func (links *Links) Sub(g gh.Thread, s sl.Thread) *Link {
	links.mu.Lock()
	defer links.mu.Unlock()

	for _, l := range links.links {
		if s.Channel == l.Sl.Channel && s.TS == l.Sl.TS {
			return nil
		}
	}

	l := &Link{Gh: g, Sl: s}
	links.links = append(links.links, l)
	return l
}

// UnSub unsbscribes specified link.
//
// Assumes that the link is unsbscribed only from slack.
// Returns ubsubscribed *Link, or nil if already ubsubscribed.
func (links *Links) UnSub(s sl.Thread) *Link {
	links.mu.Lock()
	defer links.mu.Unlock()

	// TODO: linked list can reduce computation.
	for i, l := range links.links {
		if s.Channel == l.Sl.Channel && s.TS == l.Sl.TS {
			// should not use `res = l`, because I dont know whether overwittern by del-sequence or not.
			res := links.links[len(links.links)-1]
			// https://zenn.dev/mattn/articles/31dfed3c89956d
			links.links[i] = links.links[len(links.links)-1]
			links.links[len(links.links)-1] = nil
			links.links = links.links[:len(links.links)-1]
			return res
		}
	}

	return nil
}
