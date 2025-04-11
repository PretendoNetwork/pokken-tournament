package datastore

import (
	"database/sql"
	"github.com/PretendoNetwork/pokken-tournament/globals"
	"time"

	"github.com/PretendoNetwork/nex-go/v2"
	types "github.com/PretendoNetwork/nex-go/v2/types"
	datastoreconstants "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/constants"
	datastoretypes "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/types"
	"github.com/lib/pq"
)

var insertObjectStmt *sql.Stmt

func InitializeObjectByPreparePostParam(ownerPID types.PID, param datastoretypes.DataStorePreparePostParam) (uint64, *nex.Error) {
	if param.PersistenceInitParam.DeleteLastObject && uint16(param.PersistenceInitParam.PersistenceSlotID) != datastoreconstants.InvalidPersistenceSlotID {
		persistenceTarget := datastoretypes.NewDataStorePersistenceTarget()
		persistenceTarget.OwnerID = ownerPID
		persistenceTarget.PersistenceSlotID = param.PersistenceInitParam.PersistenceSlotID

		meta, err := GetObjectInfoByPersistenceTarget(persistenceTarget)
		if err == nil {
			err = globals.DatastoreCommon.VerifyObjectPermission(ownerPID, ownerPID, meta.DelPermission)
			if err != nil {
				return 0, nex.NewError(nex.ResultCodes.Core.AccessDenied, err.Error())
			}

			err = globals.DatastoreCommon.DeleteObjectByDataID(meta.DataID)
			if err != nil {
				return 0, nex.NewError(nex.ResultCodes.DataStore.SystemFileError, err.Error())
			}
		}
	}

	var dataID uint64
	now := time.Now().UTC()

	err := insertObjectStmt.QueryRow(
		ownerPID,
		param.Size,
		param.Name,
		param.DataType,
		param.MetaBinary,
		param.Permission.Permission,
		pq.Array(param.Permission.RecipientIDs),
		param.DelPermission.Permission,
		pq.Array(param.DelPermission.RecipientIDs),
		param.Flag,
		param.Period,
		param.ReferDataID,
		pq.Array(param.Tags),
		param.PersistenceInitParam.PersistenceSlotID,
		pq.Array(param.ExtraData),
		now,
		now,
	).Scan(&dataID)

	if err != nil {
		return 0, nex.NewError(nex.ResultCodes.DataStore.SystemFileError, err.Error())
	}

	return dataID, nil
}

func initInsertObjectStmt() error {
	stmt, err := globals.PostgresDB.Prepare(`INSERT INTO datastore.objects 
	(
		owner,
		size,
		name,
		data_type,
		meta_binary,
		permission,
		permission_recipients,
		delete_permission,
		delete_permission_recipients,
		flag,
		period,
		refer_data_id,
		tags,
		persistence_slot_id,
		extra_data,
		creation_date,
		update_date
	)
	VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17
	) RETURNING data_id`)
	if err != nil {
		return err
	}

	insertObjectStmt = stmt
	return nil
}
