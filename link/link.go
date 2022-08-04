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
	Gh *gh.Thread
	Sl *sl.Thread
}
