package nex

import (
	datastore "github.com/PretendoNetwork/nex-protocols-go/datastore"
	"github.com/PretendoNetwork/pokken-tournament/globals"
	pt_datastore "github.com/PretendoNetwork/pokken-tournament/nex/datastore"
)

func registerSecureServerNEXProtocols() {
	datastoreProtocol := datastore.NewProtocol(globals.SecureServer)
	datastoreProtocol.GetMeta(pt_datastore.GetMeta)
	datastoreProtocol.PostMetaBinary(pt_datastore.PostMetaBinary)
}
