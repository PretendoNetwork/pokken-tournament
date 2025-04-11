package datastore

import (
	"database/sql"
	"errors"
	"github.com/PretendoNetwork/pokken-tournament/globals"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	datastoretypes "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/types"
)

var selectObjectByIdPasswordStmt *sql.Stmt

func GetObjectInfoByDataIDWithPassword(dataID types.UInt64, password types.UInt64) (datastoretypes.DataStoreMetaInfo, *nex.Error) {
	objects, err := getObjects(selectObjectByIdPasswordStmt, dataID, password)
	if errors.Is(err, sql.ErrNoRows) || len(objects) < 1 {
		// todo nex.ResultCodes.DataStore.InvalidPassword return?
		return datastoretypes.NewDataStoreMetaInfo(), nex.NewError(nex.ResultCodes.DataStore.NotFound, "Object not found or wrong password")
	} else if err != nil {
		return datastoretypes.NewDataStoreMetaInfo(), nex.NewError(nex.ResultCodes.DataStore.SystemFileError, err.Error())
	}

	return objects[0], nil
}

func initSelectObjectByIdPasswordStmt() error {
	stmt, err := globals.PostgresDB.Prepare(selectObject + ` WHERE data_id = $1 AND access_password = $2 LIMIT 1`)
	if err != nil {
		return err
	}

	selectObjectByIdPasswordStmt = stmt
	return nil
}
