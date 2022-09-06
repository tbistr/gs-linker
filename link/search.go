package link

import (
	"log"

	gh "github.com/tbistr/gs-linker/github"
	sl "github.com/tbistr/gs-linker/slack"
)

// SearchByG searchs links by github thread.
// Returns nil if not found.
func (Client *Client) SearchByG(g *gh.Thread) *sl.Thread {
	found := Client.db.JoinByG(g).GetS()
	if found == nil {
		log.Printf("link not found. searched by: %+v\n", g)
	} else {
		log.Printf("search by github thread: %+v find slack thread: %+v\n", g, found)
	}
	return found
}

// SearchByS searchs links by slack thread.
// Returns nil if not found.
func (Client *Client) SearchByS(s *sl.Thread) *gh.Thread {
	found := Client.db.JoinByS(s).GetG()
	if found == nil {
		log.Printf("link not found. searched by: %+v", s)
	} else {
		log.Printf("search by slack thread: %+v find github thread: %+v\n", s, found)
	}
	return found
}
