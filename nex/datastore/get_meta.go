package datastore

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PretendoNetwork/nex-go"
	datastore "github.com/PretendoNetwork/nex-protocols-go/datastore"
	"github.com/PretendoNetwork/nex-protocols-go/datastore/types"
	"github.com/PretendoNetwork/pokken-tournament/globals"
)

// Basic implementation using local filesystem, will eventually switch to S3
func GetMeta(err error, client *nex.Client, callID uint32, param *types.DataStoreGetMetaParam) uint32 {
	fmt.Println(param)

	rmcResponse := nex.NewRMCResponse(datastore.ProtocolID, callID)

	if _, err := os.Stat(fmt.Sprintf("./data/%d/%d", param.DataID, client.PID())); err == nil {
		metaFile, _ := os.ReadFile(fmt.Sprintf("./data/%d/%d", param.DataID, client.PID()))
		metadata, _ := os.ReadFile(fmt.Sprintf("./data/%d/%d.meta", param.DataID, client.PID()))
		metadataStr := strings.Split(string(metadata), "\n")

		metaInfo := types.NewDataStoreMetaInfo()
		metaInfo.DataID = param.DataID
		metaInfo.OwnerID = client.PID()
		size, _ := strconv.ParseUint(metadataStr[0], 10, 32)
		metaInfo.Size = uint32(size)
		metaInfo.Name = metadataStr[1]
		dataType, _ := strconv.ParseUint(metadataStr[2], 10, 16)
		metaInfo.DataType = uint16(dataType)
		metaInfo.MetaBinary = metaFile
		metaInfo.Permission = types.NewDataStorePermission()
		perm := strings.Split(metadataStr[3], ";")
		permissionLvl, _ := strconv.ParseUint(perm[0], 10, 8)
		recipients := make([]uint32, 0)
		for _, recipient := range strings.Split(perm[1], ",") {
			r, _ := strconv.ParseUint(recipient, 10, 32)
			recipients = append(recipients, uint32(r))
		}
		metaInfo.Permission.Permission = uint8(permissionLvl)
		metaInfo.Permission.RecipientIDs = recipients
		metaInfo.DelPermission = types.NewDataStorePermission()
		perm = strings.Split(metadataStr[4], ";")
		permissionLvl, _ = strconv.ParseUint(perm[0], 10, 8)
		recipients = make([]uint32, 0)
		for _, recipient := range strings.Split(perm[1], ",") {
			r, _ := strconv.ParseUint(recipient, 10, 32)
			recipients = append(recipients, uint32(r))
		}
		metaInfo.DelPermission.Permission = uint8(permissionLvl)
		metaInfo.DelPermission.RecipientIDs = recipients
		t1, _ := time.Parse(time.DateTime, metadataStr[5])
		t2, _ := time.Parse(time.DateTime, metadataStr[6])
		metaInfo.CreatedTime = nex.NewDateTime(nex.NewDateTime(0).FromTimestamp(t1))
		metaInfo.UpdatedTime = nex.NewDateTime(nex.NewDateTime(0).FromTimestamp(t2))
		period, _ := strconv.ParseUint(metadataStr[7], 10, 16)
		metaInfo.Period = uint16(period)
		flag, _ := strconv.ParseUint(metadataStr[7], 10, 32)
		metaInfo.Flag = uint32(flag)
		metaInfo.ExpireTime = nex.NewDateTime(nex.NewDateTime(0).FromTimestamp(time.Now().Add(30 * 24 * time.Hour)))
		metaInfo.Ratings = []*types.DataStoreRatingInfoWithSlot{}
		metaInfo.ReferredTime = nex.NewDateTime(0)

		rmcResponseBody := nex.NewStreamOut(globals.SecureServer)
		rmcResponseBody.WriteStructure(metaInfo)
		rmcResponse.SetSuccess(datastore.MethodGetMeta, rmcResponseBody.Bytes())
	} else if errors.Is(err, os.ErrNotExist) {
		rmcResponse.SetError(nex.Errors.DataStore.NotFound)
	} else {
		panic(err)
	}

	rmcResponseBytes := rmcResponse.Bytes()

	var responsePacket nex.PacketInterface

	responsePacket, _ = nex.NewPacketV1(client, nil)
	responsePacket.SetVersion(1)

	responsePacket.SetSource(0xA1)
	responsePacket.SetDestination(0xAF)
	responsePacket.SetType(nex.DataPacket)
	responsePacket.SetPayload(rmcResponseBytes)

	responsePacket.AddFlag(nex.FlagNeedsAck)
	responsePacket.AddFlag(nex.FlagReliable)

	globals.SecureServer.Send(responsePacket)

	return 0

}
