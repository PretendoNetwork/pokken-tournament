package datastore

import (
	"database/sql"
	"errors"
	"github.com/PretendoNetwork/pokken-tournament/globals"
	"slices"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	datastoreconstants "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/constants"
	datastoretypes "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/types"
	"github.com/lib/pq"
)

var selectObjectsBySearchParamStmt *sql.Stmt

func GetObjectInfosByDataStoreSearchParam(param datastoretypes.DataStoreSearchParam, pid types.PID) ([]datastoretypes.DataStoreMetaInfo, uint32, *nex.Error) {
	if param.CreatedAfter == 0x9C3F3F7EFB { // 9999-12-31 23:59:59
		param.CreatedAfter = 0x4420000 // year 1, month 1, day 1 00:00:00 (go zero date)
	}

	if param.UpdatedAfter == 0x9C3F3F7EFB {
		param.UpdatedAfter = 0x4420000
	}

	var dataTypes []int32
	var idArray []int64
	var tagArray []string

	// use one search id array to pass to postgres and make our lives Significantly Easier:tm:
	if uint8(param.SearchTarget) == uint8(datastoreconstants.SearchTypeFriend) {
		pids := globals.GetUserFriendPIDs(uint32(pid))

		// this is guessed behavior, it probably is just filtered to friends only with param.OwnerIDs ignored on official servers but no evidence for now
		// if we arent trying to filter then copy over the pids to idArray
		if len(param.OwnerIDs) == 0 {
			idArray = make([]int64, 0, len(pids))
			for i := range pids {
				idArray = append(idArray, int64(pids[i]))
			}
		} else {
			// otherwise if we are then filter to pids present in param.OwnerIDs
			idArray = make([]int64, 0)
			for i := range pids {
				if slices.Contains(param.OwnerIDs, types.NewPID(uint64(pids[i]))) {
					idArray = append(idArray, int64(pids[i]))
				}
			}
		}
	} else {
		idArray = make([]int64, 0, len(param.OwnerIDs))
		for i := range param.OwnerIDs {
			idArray = append(idArray, int64(param.OwnerIDs[i]))
		}
	}

	if uint8(param.SearchTarget) == uint8(datastoreconstants.SearchTypeOwnAll) {
		idArray = append(idArray, int64(pid))
	}

	// handle param.DataType vs param.DataTypes
	if len(param.DataTypes) != 0 {
		dataTypes = make([]int32, 0, len(param.DataTypes))
		for i := range param.DataTypes {
			dataTypes = append(dataTypes, int32(param.DataTypes[i]))
		}
	} else {
		dataTypes = append(dataTypes, int32(param.DataType))
	}

	// convert param.Tags into a proper string array
	tagArray = make([]string, 0, len(param.Tags))
	for i := range param.Tags {
		// types.String.String() escapes which is Not Good:tm:
		tagArray = append(tagArray, string(param.Tags[i]))
	}

	objects, err := getObjects(selectObjectsBySearchParamStmt,
		pq.Int64Array(idArray),
		pq.Int32Array(dataTypes),
		pq.FormatTimestamp(param.CreatedAfter.Standard()),
		pq.FormatTimestamp(param.CreatedBefore.Standard()),
		pq.FormatTimestamp(param.UpdatedAfter.Standard()),
		pq.FormatTimestamp(param.UpdatedBefore.Standard()),
		int64(param.ReferDataID),
		pq.StringArray(tagArray),
		int64(param.ResultRange.Length))

	if errors.Is(err, sql.ErrNoRows) {
		return []datastoretypes.DataStoreMetaInfo{}, 0, nil
	} else if err != nil {
		return []datastoretypes.DataStoreMetaInfo{}, 0, nex.NewError(nex.ResultCodes.DataStore.SystemFileError, err.Error())
	}

	return objects, uint32(len(objects)), nil
}

func initSelectObjectsBySearchParamStmt() error {
	stmt, err := globals.PostgresDB.Prepare(
		selectObject + ` WHERE (owner = ANY($1) OR cardinality($1) = 0)
		AND data_type = ANY($2)
		AND creation_date BETWEEN $3 AND $4
		AND update_date BETWEEN $5 AND $6
		AND refer_data_id = $7
		AND tags @> $8
		AND deleted IS FALSE
		ORDER BY data_id DESC
		LIMIT $9`)
	if err != nil {
		return err
	}

	selectObjectsBySearchParamStmt = stmt
	return nil
}
