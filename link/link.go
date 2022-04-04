package link

// Links consists of link
type Links struct {
	l []*link
}

// link is bidirectional map between slack thread and github issue
type link struct {
	g string
	s string
}

// Sub subscribes a link
//
// does not check if the link is already registered
func (links *Links) Sub(g, s string) {
	links.l = append(links.l, &link{g: g, s: s})
}

// UnSub unsbscribe specified link
//
// assume that a link is unsbscribed from slack
func (links *Links) UnSub(s string) {
	for i := range links.l {
		if links.l[i].s == s {
			links.l[i] = links.l[len(links.l)-1]
			links.l[len(links.l)-1] = nil
			links.l = links.l[:len(links.l)-1]
			return
		}
	}
}
