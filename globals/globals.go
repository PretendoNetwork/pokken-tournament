package globals

import (
	"database/sql"

	pb_account "github.com/PretendoNetwork/grpc-go/account"
	pb_friends "github.com/PretendoNetwork/grpc-go/friends"
	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/plogger-go"
	"github.com/minio/minio-go/v7"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	datastore "github.com/PretendoNetwork/nex-protocols-common-go/v2/datastore"
)

var Logger *plogger.Logger
var KerberosPassword = "password" // * Default password
var AuthenticationServer *nex.PRUDPServer
var AuthenticationEndpoint *nex.PRUDPEndPoint
var SecureServer *nex.PRUDPServer
var SecureEndpoint *nex.PRUDPEndPoint
var MinIOClient *minio.Client
var Presigner *S3Presigner
var GRPCAccountClientConnection *grpc.ClientConn
var GRPCAccountClient pb_account.AccountClient
var GRPCAccountCommonMetadata metadata.MD
var GRPCFriendsClientConnection *grpc.ClientConn
var GRPCFriendsClient pb_friends.FriendsClient
var GRPCFriendsCommonMetadata metadata.MD
var DatastoreCommon *datastore.CommonProtocol

var PostgresDB *sql.DB
