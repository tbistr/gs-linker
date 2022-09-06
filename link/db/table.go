package db

import (
	"database/sql"
)

const (
	LINKS_TABLE_NAME = "links"
	GH_TABLE_NAME    = "slack_threads"
	SL_TABLE_NAME    = "github_threads"
)

type LinkTable struct {
}

type GhTable struct {
}

type SlTable struct {
}

func (client *Client) TouchTables(db *sql.DB) error {
	// link
	if _, err := db.Exec(
		"CREATE TABLE IF NOT EXISTS `links`" +
			"(id BIGINT UNSIGNED AUTO_INCREMENT,PRIMARY KEY(id))",
	); err != nil {
		return err
	}

	// github
	if _, err := db.Exec(
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
	if _, err := db.Exec(
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
