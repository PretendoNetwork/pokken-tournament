package nex

import (
	"os"

	"github.com/PretendoNetwork/nex-go/v2/types"
	common_datastore "github.com/PretendoNetwork/nex-protocols-common-go/v2/datastore"
	common_globals "github.com/PretendoNetwork/nex-protocols-common-go/v2/globals"
	common_match_making "github.com/PretendoNetwork/nex-protocols-common-go/v2/match-making"
	common_match_making_ext "github.com/PretendoNetwork/nex-protocols-common-go/v2/match-making-ext"
	common_matchmake_extension "github.com/PretendoNetwork/nex-protocols-common-go/v2/matchmake-extension"
	common_nat_traversal "github.com/PretendoNetwork/nex-protocols-common-go/v2/nat-traversal"
	common_ranking "github.com/PretendoNetwork/nex-protocols-common-go/v2/ranking"
	common_secure "github.com/PretendoNetwork/nex-protocols-common-go/v2/secure-connection"
	datastore "github.com/PretendoNetwork/nex-protocols-go/v2/datastore"
	match_making "github.com/PretendoNetwork/nex-protocols-go/v2/match-making"
	match_making_ext "github.com/PretendoNetwork/nex-protocols-go/v2/match-making-ext"
	matchmake_extension "github.com/PretendoNetwork/nex-protocols-go/v2/matchmake-extension"
	nat_traversal "github.com/PretendoNetwork/nex-protocols-go/v2/nat-traversal"
	ranking "github.com/PretendoNetwork/nex-protocols-go/v2/ranking"
	secure "github.com/PretendoNetwork/nex-protocols-go/v2/secure-connection"
	"github.com/PretendoNetwork/pokken-tournament/database"
	"github.com/PretendoNetwork/pokken-tournament/globals"
	local_ranking "github.com/PretendoNetwork/pokken-tournament/nex/ranking"
)

func registerCommonSecureServerProtocols() {
	secureProtocol := secure.NewProtocol()
	globals.SecureEndpoint.RegisterServiceProtocol(secureProtocol)
	common_secure_protocol := common_secure.NewCommonProtocol(secureProtocol)
	common_secure_protocol.CreateReportDBRecord = func(pid types.PID, reportID types.UInt32, reportData types.QBuffer) error {
		// Stubbed
		return nil
	}

	natTraversalProtocol := nat_traversal.NewProtocol()
	globals.SecureEndpoint.RegisterServiceProtocol(natTraversalProtocol)
	common_nat_traversal.NewCommonProtocol(natTraversalProtocol)

	matchmakingManager := common_globals.NewMatchmakingManager(globals.SecureEndpoint, database.Postgres)

	matchMakingProtocol := match_making.NewProtocol()
	globals.SecureEndpoint.RegisterServiceProtocol(matchMakingProtocol)
	commonMatchMakingProtocol := common_match_making.NewCommonProtocol(matchMakingProtocol)
	commonMatchMakingProtocol.SetManager(matchmakingManager)

	matchMakingExtProtocol := match_making_ext.NewProtocol()
	globals.SecureEndpoint.RegisterServiceProtocol(matchMakingExtProtocol)
	commonMatchMakingExtProtocol := common_match_making_ext.NewCommonProtocol(matchMakingExtProtocol)
	commonMatchMakingExtProtocol.SetManager(matchmakingManager)

	matchmakeExtensionProtocol := matchmake_extension.NewProtocol()
	globals.SecureEndpoint.RegisterServiceProtocol(matchmakeExtensionProtocol)
	commonMatchmakeExtensionProtocol := common_matchmake_extension.NewCommonProtocol(matchmakeExtensionProtocol)
	commonMatchmakeExtensionProtocol.SetManager(matchmakingManager)

	rankingProtocol := ranking.NewProtocol()
	globals.SecureEndpoint.RegisterServiceProtocol(rankingProtocol)
	commonRankingProtocol := common_ranking.NewCommonProtocol(rankingProtocol)
	commonRankingProtocol.GetRankingsAndCountByCategoryAndRankingOrderParam = local_ranking.GetRankingsAndCountByCategoryAndRankingOrder
	commonRankingProtocol.InsertRankingByPIDAndRankingScoreData = local_ranking.InsertRankingByPIDAndRankingScoreData

	datastoreProtocol := datastore.NewProtocol()
	globals.SecureEndpoint.RegisterServiceProtocol(datastoreProtocol)
	commonDatastoreProtocol := common_datastore.NewCommonProtocol(datastoreProtocol)
	commonDatastoreProtocol.S3Bucket = os.Getenv("PN_POKKENTOURNAMENT_S3_BUCKET")
	commonDatastoreProtocol.SetMinIOClient(globals.MinIOClient)
}
