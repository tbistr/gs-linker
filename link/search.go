package link

import (
	"fmt"
	"log"
	"reflect"

	gh "github.com/tbistr/gs-linker/github"
	sl "github.com/tbistr/gs-linker/slack"
)

// SearchByG searchs links by github thread.
// Returns non-nil error if not found.
func (links *Links) SearchByG(g *gh.Thread) (*sl.Thread, error) {
	for _, l := range links.links {
		if reflect.DeepEqual(g, l.Gh) {
			log.Printf("search by github thread: %+v find slack thread: %+v\n", g, l.Sl)
			// should not return just l.Sl.
			// it is uncontroled out of func.
			// TODO: is there something good way?
			return &sl.Thread{
				Channel: l.Sl.Channel,
				TS:      l.Sl.TS,
			}, nil
		}
	}

	return nil, fmt.Errorf("link not found. searched by: %+v", g)
}

// SearchByS searchs links by slack thread.
// Returns non-nil error if not found.
func (links *Links) SearchByS(s *sl.Thread) (*gh.Thread, error) {
	for _, l := range links.links {
		if reflect.DeepEqual(s, l.Sl) {
			log.Printf("search by slack thread: %+v find github thread: %+v\n", s, l.Gh)
			return &gh.Thread{
				SubType: l.Gh.SubType,
				Owner:   l.Gh.Owner,
				Repo:    l.Gh.Repo,
				Num:     l.Gh.Num,
			}, nil
		}
	}

	return nil, fmt.Errorf("link not found. searched by: %+v", s)
}
