package datastore

import (
	"database/sql"
	"github.com/PretendoNetwork/pokken-tournament/globals"
	"time"

	"github.com/PretendoNetwork/nex-go/v2/types"
	datastoreconstants "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/constants"
	datastoretypes "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/types"
	"github.com/lib/pq"
)

func initDatabase() error {
	inits := []func() error{
		initTables,
		initInsertObjectStmt,                         // initialize_object_by_prepare_post_param.go
		initSelectObjectByIdPasswordStmt,             // get_object_info_by_data_id_with_password.go
		initSelectObjectByOwnerPersistenceNoPassStmt, // get_object_info_by_persistence_target.go
		initSelectObjectByOwnerPersistenceStmt,       // get_object_info_by_persistence_target_with_password.go
		initSelectObjectByIdStmt,                     // get_object_info_by_data_id.go
		initSelectObjectsBySearchParamStmt,           // get_object_infos_by_data_store_search_param.go
		initSelectOwnerByIdStmt,                      // get_object_owner_by_data_id.go
		initSelectSizeByIdStmt,                       // get_object_size_by_data_id.go
		initUpdateUploadCompleteByIdStmt,             // update_object_upload_completed_by_data_id.go
		initUpdateMetaBinaryByIdPasswordStmt,         // update_object_meta_binary_by_data_id_with_password.go
		initUpdatePeriodByIdPasswordStmt,             // update_object_period_by_data_id_with_password.go
		initUpdateDataTypeByIdPasswordStmt,           // update_object_data_type_by_data_id_with_password.go
		initUpdateDeletedByIdStmt,                    // delete_object_by_data_id.go
		initUpdateDeletedByIdPasswordStmt,            // delete_object_by_data_id_with_password.go
	}

	for _, init := range inits {
		err := init()
		if err != nil {
			return err
		}
	}

	return nil
}

func initTables() error {
	_, err := globals.PostgresDB.Exec(`CREATE SCHEMA IF NOT EXISTS datastore`)
	if err != nil {
		return err
	}

	globals.Logger.Success("datastore Postgres schema created")

	_, err = globals.PostgresDB.Exec(`CREATE SEQUENCE IF NOT EXISTS datastore.object_data_id_seq
		INCREMENT 1
		MINVALUE 1
		MAXVALUE 281474976710656
		START 1
		CACHE 1`, // * Honestly I don't know what CACHE does but I saw it recommended so here it is
	)
	if err != nil {
		return err
	}

	_, err = globals.PostgresDB.Exec(`CREATE TABLE IF NOT EXISTS datastore.objects (
		data_id bigint NOT NULL DEFAULT nextval('datastore.object_data_id_seq') PRIMARY KEY,
		upload_completed boolean NOT NULL DEFAULT FALSE,
		deleted boolean NOT NULL DEFAULT FALSE,
		owner bigint,
		size int,
		name text,
		data_type int,
		meta_binary bytea,
		permission int,
		permission_recipients int[],
		delete_permission int,
		delete_permission_recipients int[],
		flag int,
		period int,
		refer_data_id bigint,
		tags text[],
		persistence_slot_id int,
		extra_data text[],
		access_password bigint NOT NULL DEFAULT 0,
		update_password bigint NOT NULL DEFAULT 0,
		creation_date timestamp,
		update_date timestamp
	)`)
	if err != nil {
		return err
	}

	globals.Logger.Success("Postgres tables created")
	return nil
}

const selectObject = `
	SELECT
		data_id,
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
		creation_date,
		update_date
	FROM datastore.objects`

// Helper to unpack things selected with (selectObject + ` WHERE ....`)
func getObjects(stmt *sql.Stmt, args ...any) ([]datastoretypes.DataStoreMetaInfo, error) {
	rows, err := stmt.Query(args...)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	// surely we know the length of the result set at this point?
	var results []datastoretypes.DataStoreMetaInfo

	for rows.Next() {
		result := datastoretypes.NewDataStoreMetaInfo()

		var createdTime time.Time
		var updatedTime time.Time
		var tagArray []string
		var persistentSlotId uint16

		err := rows.Scan(
			&result.DataID,
			&result.OwnerID,
			&result.Size,
			&result.Name,
			&result.DataType,
			&result.MetaBinary,
			&result.Permission.Permission,
			pq.Array(&result.Permission.RecipientIDs),
			&result.DelPermission.Permission,
			pq.Array(&result.DelPermission.RecipientIDs),
			&result.Flag,
			&result.Period,
			&result.ReferDataID,
			pq.Array(&tagArray),
			&persistentSlotId,
			&createdTime,
			&updatedTime,
		)

		if err != nil {
			return nil, err
			//globals.Logger.Error(err.Error())
			//continue
		}

		// TODO: is this a good implementation of this?

		// default to never expire (persistent slots)
		result.ExpireTime = types.NewDateTime(0x9C3F3F7EFB) // 9999-12-31 23:59:59
		// this isn't vanilla behavior (vanilla is 9999-12-31 00:00:00)
		// but it should be close enough, and i dont see direct checks
		// for either of the times in ghidra (don't quote me please)

		// only set it otherwise if object doesn't get persisted
		if persistentSlotId == datastoreconstants.InvalidPersistenceSlotID {
			// this is definitely a bit of a line of code
			// basically, make a time from the updated time, and add the object's period value as days
			dateTime := types.NewDateTime(0)
			result.ExpireTime = dateTime.FromTimestamp(updatedTime.AddDate(0, 0, int(result.Period)))
		}

		// add the tags from tagArray (created because we can't really use NEX types in postgres) to the nex array
		result.Tags = make(types.List[types.String], 0, len(tagArray))
		for i := range tagArray {
			result.Tags = append(result.Tags, types.String(tagArray[i]))
		}

		// I'm not sure how this API is meant to be used but this works
		result.CreatedTime = result.CreatedTime.FromTimestamp(createdTime)
		result.UpdatedTime = result.UpdatedTime.FromTimestamp(updatedTime)
		result.ReferredTime = result.ReferredTime.FromTimestamp(createdTime)
		// note from another dev: referred time does seem to equal created(/updated?)
		// time in packet dumps, so looks good to me!

		results = append(results, result)
	}

	return results, rows.Err()
}
