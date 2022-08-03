package link

// Links consists of link
type Links []*Link

// Link is bidirectional map between slack thread and github issue
type Link struct {
	G              string
	SlackTimestamp string
}

// Sub subscribes a link
//
// does not check if the link is already registered
func (links *Links) Sub(g, s string) {
	*links = append(*links, &Link{G: g, SlackTimestamp: s})
}

// UnSub unsbscribe specified link
//
// assume that a link is unsbscribed from slack
func (links *Links) UnSub(s string) {
	for i := range *links {
		if (*links)[i].SlackTimestamp == s {
			(*links)[i] = (*links)[len(*links)-1]
			(*links)[len(*links)-1] = nil
			*links = (*links)[:len(*links)-1]
			return
		}
	}
}
