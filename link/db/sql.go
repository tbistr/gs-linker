package db

import (
	gh "github.com/tbistr/gs-linker/github"
	sl "github.com/tbistr/gs-linker/slack"
)

func (client *Client) Create(gh *gh.Thread, sl *sl.Thread) (*LinkTable, error) {
	return nil, nil
}

func (client *Client) DeleteByS(sl *sl.Thread) error {
	return nil
}

func (client *Client) JoinByG(g *gh.Thread) *LinkTable {
	return nil
}

func (client *Client) JoinByS(s *sl.Thread) *LinkTable {
	return nil
}
