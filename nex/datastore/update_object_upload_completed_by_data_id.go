package datastore

import (
	"database/sql"
	"errors"
	"github.com/PretendoNetwork/pokken-tournament/globals"
	"time"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
)

var updateUploadCompleteByIdStmt *sql.Stmt

func UpdateObjectUploadCompletedByDataID(dataID types.UInt64, uploadCompleted bool) *nex.Error {
	result, err := updateUploadCompleteByIdStmt.Exec(dataID, uploadCompleted, time.Now().UTC())
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

func initUpdateUploadCompleteByIdStmt() error {
	// TODO SMM has a check on update_password here, but it doesn't make sense to me. work it out and implement it
	stmt, err := globals.PostgresDB.Prepare(`UPDATE datastore.objects SET upload_completed = $2, update_date = $3 WHERE data_id = $1`)
	if err != nil {
		return err
	}

	updateUploadCompleteByIdStmt = stmt
	return nil
}
