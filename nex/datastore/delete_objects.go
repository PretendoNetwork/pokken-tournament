package datastore

import (
	"github.com/PretendoNetwork/pokken-tournament/globals"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"

	datastore "github.com/PretendoNetwork/nex-protocols-go/v2/datastore"
	datastore_types "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/types"
)

// slightly repurposed from
// https://github.com/PretendoNetwork/nex-protocols-common-go/tree/d7f8b585c13942350b1f60de6e2b97746ac8867b/datastore/delete_object.go
func DeleteObjects(err error, packet nex.PacketInterface, callID uint32, params types.List[datastore_types.DataStoreDeleteParam], transactional types.Bool) (*nex.RMCMessage, *nex.Error) {

	// TODO: properly implement transactional
	// https://discord.com/channels/408718485913468928/881852117550243860/1260329610495918211
	if globals.DatastoreCommon.GetObjectInfoByDataID == nil {
		globals.Logger.Warning("GetObjectInfoByDataID not defined")
		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "change_error")
	}

	if globals.DatastoreCommon.DeleteObjectByDataIDWithPassword == nil {
		globals.Logger.Warning("DeleteObjectByDataIDWithPassword not defined")
		return nil, nex.NewError(nex.ResultCodes.Core.NotImplemented, "change_error")
	}

	if err != nil {
		globals.Logger.Error(err.Error())
		return nil, nex.NewError(nex.ResultCodes.DataStore.Unknown, "change_error")
	}

	connection := packet.Sender()
	endpoint := connection.Endpoint()

	for i := range params {
		metaInfo, errCode := globals.DatastoreCommon.GetObjectInfoByDataID(params[i].DataID)
		if errCode != nil {
			if transactional {
				return nil, errCode
			}

			continue
		}

		errCode = globals.DatastoreCommon.VerifyObjectPermission(metaInfo.OwnerID, connection.PID(), metaInfo.DelPermission)
		if errCode != nil {
			if transactional {
				return nil, errCode
			}

			continue
		}

		errCode = globals.DatastoreCommon.DeleteObjectByDataIDWithPassword(params[i].DataID, params[i].UpdatePassword)
		if errCode != nil {
			if transactional {
				return nil, errCode
			}

			continue
		}
	}

	rmcResponse := nex.NewRMCSuccess(endpoint, nil)
	rmcResponse.ProtocolID = datastore.ProtocolID
	rmcResponse.MethodID = datastore.MethodDeleteObject
	rmcResponse.CallID = callID

	return rmcResponse, nil
}
