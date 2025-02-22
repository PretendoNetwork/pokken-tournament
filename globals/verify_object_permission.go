package globals

import (
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	datastore_types "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/types"
)

func VerifyObjectPermission(ownerPID, accessorPID types.PID, permission datastore_types.DataStorePermission) uint32 {
	if permission.Permission > 3 {
		return nex.ResultCodes.DataStore.InvalidArgument
	}

	// * Allow anyone
	if permission.Permission == 0 {
		return 0
	}

	// * Allow friends
	// TODO - Implement this
	if permission.Permission == 1 {
		return nex.ResultCodes.DataStore.PermissionDenied
	}

	// * Allow people in permission.RecipientIDs
	if permission.Permission == 2 {
		if !permission.RecipientIDs.Contains(accessorPID) {
			return nex.ResultCodes.DataStore.PermissionDenied
		}
	}

	// * Allow only the owner
	if permission.Permission == 3 {
		if !ownerPID.Equals(accessorPID) {
			return nex.ResultCodes.DataStore.PermissionDenied
		}
	}

	return 0
}
