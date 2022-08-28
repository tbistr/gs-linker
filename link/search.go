package link

import (
	"log"

	gh "github.com/tbistr/gs-linker/github"
	sl "github.com/tbistr/gs-linker/slack"
)

// SearchByG searchs links by github thread.
// Returns nil if not found.
func (Client *Client) SearchByG(g *gh.Thread) *sl.Thread {
	var gFound gh.Thread
	Client.db.First(&gFound, *g)
	if gFound.ID == 0 {
		log.Printf("link not found. searched by: %+v\n", g)
		return nil
	} else {
		var l Link
		Client.db.Joins("Gh").Joins("Sl").First(&l, gFound.LinkID)
		log.Printf("search by github thread: %+v find slack thread: %+v\n", g, l.Sl)
		return l.Sl
	}
}

// SearchByS searchs links by slack thread.
// Returns nil if not found.
func (Client *Client) SearchByS(s *sl.Thread) *gh.Thread {
	var sFound sl.Thread
	Client.db.First(&sFound, *s)
	if sFound.ID == 0 {
		log.Printf("link not found. searched by: %+v", s)
		return nil
	} else {
		var l Link
		Client.db.Joins("Gh").Joins("Sl").First(&l, sFound.LinkID)
		log.Printf("search by slack thread: %+v find github thread: %+v\n", s, l.Gh)
		return l.Gh
	}
}
