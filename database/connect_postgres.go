package database

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"

	"github.com/PretendoNetwork/pokken-tournament/globals"
)

func ConnectPostgres() {
	var err error

	globals.PostgresDB, err = sql.Open("postgres", os.Getenv("PN_POKKENTOURNAMENT_POSTGRES_URI"))
	if err != nil {
		globals.Logger.Critical(err.Error())
		panic(err)
	}

	globals.Logger.Success("Connected to Postgres!")

	initPostgres()
}
