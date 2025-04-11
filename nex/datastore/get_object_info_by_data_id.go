package datastore

import (
	"database/sql"
	"errors"
	"github.com/PretendoNetwork/pokken-tournament/globals"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	datastoretypes "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/types"
)

var selectObjectByIdStmt *sql.Stmt

func GetObjectInfoByDataID(dataID types.UInt64) (datastoretypes.DataStoreMetaInfo, *nex.Error) {
	objects, err := getObjects(selectObjectByIdStmt, dataID)
	if errors.Is(err, sql.ErrNoRows) || len(objects) < 1 {
		return datastoretypes.NewDataStoreMetaInfo(), nex.NewError(nex.ResultCodes.DataStore.NotFound, "Object not found")
	} else if err != nil {
		return datastoretypes.NewDataStoreMetaInfo(), nex.NewError(nex.ResultCodes.DataStore.SystemFileError, err.Error())
	}

	return objects[0], nil
}

func initSelectObjectByIdStmt() error {
	stmt, err := globals.PostgresDB.Prepare(selectObject + ` WHERE data_id = $1 LIMIT 1`)
	if err != nil {
		return err
	}

	selectObjectByIdStmt = stmt
	return nil
}
