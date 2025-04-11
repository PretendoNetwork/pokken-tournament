package database

import "github.com/PretendoNetwork/pokken-tournament/globals"

func initPostgres() {
	var err error

	_, err = globals.PostgresDB.Exec(`CREATE TABLE IF NOT EXISTS common_data (
		unique_id serial PRIMARY KEY,
		owner_pid integer,
		common_data bytea
	)`)
	if err != nil {
		globals.Logger.Critical(err.Error())
		panic(err)
		return
	}

	globals.Logger.Success("Postgres tables created")
}
