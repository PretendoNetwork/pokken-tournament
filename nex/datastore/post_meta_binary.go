package datastore

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/PretendoNetwork/nex-go"
	datastore "github.com/PretendoNetwork/nex-protocols-go/datastore"
	"github.com/PretendoNetwork/nex-protocols-go/datastore/types"
	"github.com/PretendoNetwork/pokken-tournament/globals"
)

// Basic implementation using local filesystem, will eventually switch to S3
func PostMetaBinary(err error, client *nex.Client, callID uint32, param *types.DataStorePreparePostParam) uint32 {
	fmt.Println(param)

	os.WriteFile(fmt.Sprintf("./data/%d/%d", param.ReferDataID, client.PID()), param.MetaBinary, 0777)

	metadata := make([]string, 0)

	metadata = append(metadata, fmt.Sprintf("%d", param.Size))
	metadata = append(metadata, param.Name)
	metadata = append(metadata, fmt.Sprintf("%d", param.DataType))
	recipients := make([]string, 0)
	for _, recipient := range param.Permission.RecipientIDs {
		recipients = append(recipients, fmt.Sprintf("%d", recipient))
	}
	metadata = append(metadata, fmt.Sprintf("%d;%s", param.Permission.Permission, strings.Join(recipients, ",")))
	delrecipients := make([]string, 0)
	for _, delrecipient := range param.DelPermission.RecipientIDs {
		delrecipients = append(delrecipients, fmt.Sprintf("%d", delrecipient))
	}
	metadata = append(metadata, fmt.Sprintf("%d;%s", param.DelPermission.Permission, strings.Join(delrecipients, ",")))
	metadata = append(metadata, time.Now().Format(time.DateTime))
	metadata = append(metadata, time.Now().Format(time.DateTime))
	metadata = append(metadata, fmt.Sprintf("%d", param.Period))
	metadata = append(metadata, fmt.Sprintf("%d", param.Flag))

	os.WriteFile(fmt.Sprintf("./data/%d/%d.meta", param.ReferDataID, client.PID()), []byte(strings.Join(metadata, "\n")), 0777)

	rmcResponseBody := nex.NewStreamOut(globals.SecureServer)

	rmcResponseBody.WriteUInt64LE(uint64(param.ReferDataID))

	rmcResponse := nex.NewRMCResponse(datastore.ProtocolID, callID)
	rmcResponse.SetSuccess(datastore.MethodPostMetaBinary, rmcResponseBody.Bytes())

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
