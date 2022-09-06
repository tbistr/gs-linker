package db

import (
	gh "github.com/tbistr/gs-linker/github"
	sl "github.com/tbistr/gs-linker/slack"
)

const (
	LINKS_TABLE_NAME = "links"
	GH_TABLE_NAME    = "slack_threads"
	SL_TABLE_NAME    = "github_threads"
)

type LinkTable struct {
	ID uint
	gh *ghTable
	sl *slTable
}

type ghTable struct {
	ID     uint
	LinkID uint
	Owner  string
	Repo   string
	Num    int
}

type slTable struct {
	ID      uint
	LinkID  uint
	Channel string
	TS      string
}

func (client *Client) TouchTables() error {
	// link
	if _, err := client.db.Exec(
		"CREATE TABLE IF NOT EXISTS `links`" +
			"(id BIGINT UNSIGNED AUTO_INCREMENT,PRIMARY KEY(id))",
	); err != nil {
		return err
	}

	// github
	if _, err := client.db.Exec(
		"CREATE TABLE IF NOT EXISTS `github_threads`" +
			"(`id` bigint unsigned AUTO_INCREMENT," +
			"`link_id` bigint unsigned," +
			"`owner` longtext," +
			"`repo` longtext," +
			"`num` bigint," +
			"PRIMARY KEY (`id`)," +
			"CONSTRAINT `fk_links_gh` FOREIGN KEY (`link_id`) REFERENCES `links`(`id`) ON DELETE CASCADE)",
	); err != nil {
		return err
	}

	// slack
	if _, err := client.db.Exec(
		"CREATE TABLE IF NOT EXISTS `slack_threads`(" +
			"`id` BIGINT UNSIGNED AUTO_INCREMENT," +
			"`link_id` BIGINT UNSIGNED," +
			"`channel` LONGTEXT," +
			"`ts` LONGTEXT," +
			"PRIMARY KEY(`id`)," +
			"CONSTRAINT `fk_links_sl` FOREIGN KEY(`link_id`) REFERENCES `links`(`id`) ON DELETE CASCADE",
	); err != nil {
		return err
	}

	return nil
}

func (link *LinkTable) GetG() *gh.Thread {
	if link == nil {
		return nil
	} else {
		return &gh.Thread{
			Owner: link.gh.Owner,
			Repo:  link.gh.Repo,
			Num:   link.gh.Num,
		}
	}
}

func (link *LinkTable) GetS() *sl.Thread {
	if link == nil {
		return nil
	} else {
		return &sl.Thread{
			Channel: link.sl.Channel,
			TS:      link.sl.TS,
		}
	}
}
