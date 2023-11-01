package database

import (
	"database/sql"

	"github.com/PretendoNetwork/pokken-tournament/globals"
	_ "github.com/mattn/go-sqlite3"
)

// Still WIP, unused
func InitDatabase() {
	var err error
	globals.RankingDatabase, err = sql.Open("sqlite3", "./ranking.db")
	if err != nil {
		panic(err)
	}

	globals.RankingDatabase.Exec(`
	CREATE TABLE IF NOT EXISTS commondata (
		pid MEDIUMINT NOT NULL DEFAULT 0,
		uniqueId BIGINT NOT NULL DEFAULT 0,
		data BLOB NOT NULL DEFAULT ''
	);
	`)

	globals.RankingDatabase.Exec(`
	CREATE TABLE IF NOT EXISTS scoredata (
		pid MEDIUMINT NOT NULL DEFAULT 0,
		uniqueId BIGINT NOT NULL DEFAULT 0,
		category MEDIUMINT NOT NULL DEFAULT 0,
		score MEDIUMINT NOT NULL DEFAULT 0,
		groups BLOB NOT NULL DEFAULT '',
		param BIGINT NOT NULL DEFAULT 0,
	);
	`)
}
