package datastore

import (
	"database/sql"
	"errors"
	"github.com/PretendoNetwork/pokken-tournament/globals"
	"time"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
)

var updateMetaBinaryByIdPasswordStmt *sql.Stmt

func UpdateObjectMetaBinaryByDataIDWithPassword(dataID types.UInt64, metaBinary types.QBuffer, password types.UInt64) *nex.Error {
	result, err := updateMetaBinaryByIdPasswordStmt.Exec(dataID, password, metaBinary, time.Now().UTC())
	if errors.Is(err, sql.ErrNoRows) {
		return nex.NewError(nex.ResultCodes.DataStore.NotFound, "Object not found or wrong password")
	} else if err != nil {
		return nex.NewError(nex.ResultCodes.DataStore.SystemFileError, err.Error())
	}

	rows, err := result.RowsAffected()
	if err != nil && rows < 1 {
		return nex.NewError(nex.ResultCodes.DataStore.NotFound, "Object not found or wrong password")
	}

	return nil
}

func initUpdateMetaBinaryByIdPasswordStmt() error {
	stmt, err := globals.PostgresDB.Prepare(`UPDATE datastore.objects SET meta_binary = $3, update_date = $4 WHERE data_id = $1 AND update_password = $2`)
	if err != nil {
		return err
	}

	updateMetaBinaryByIdPasswordStmt = stmt
	return nil
}
