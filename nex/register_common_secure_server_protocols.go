package nex

import (
	matchmake_extension "github.com/PretendoNetwork/nex-protocols-common-go/matchmake-extension"
	matchmaking "github.com/PretendoNetwork/nex-protocols-common-go/matchmaking"
	matchmaking_ext "github.com/PretendoNetwork/nex-protocols-common-go/matchmaking-ext"
	nattraversal "github.com/PretendoNetwork/nex-protocols-common-go/nat-traversal"
	"github.com/PretendoNetwork/nex-protocols-common-go/ranking"
	secureconnection "github.com/PretendoNetwork/nex-protocols-common-go/secure-connection"
	"github.com/PretendoNetwork/pokken-tournament/globals"
	pt_ranking "github.com/PretendoNetwork/pokken-tournament/nex/ranking"
)

func registerCommonSecureServerProtocols() {
	secure_protocol := secureconnection.NewCommonSecureConnectionProtocol(globals.SecureServer)
	secure_protocol.CreateReportDBRecord(func(pid, reportID uint32, reportData []byte) error {
		// Stubbed
		return nil
	})
	matchmake_extension.NewCommonMatchmakeExtensionProtocol(globals.SecureServer)
	matchmaking.NewCommonMatchMakingProtocol(globals.SecureServer)
	matchmaking_ext.NewCommonMatchMakingExtProtocol(globals.SecureServer)
	nattraversal.NewCommonNATTraversalProtocol(globals.SecureServer)
	ranking_protocol := ranking.NewCommonRankingProtocol(globals.SecureServer)
	ranking_protocol.GetRankingsAndCountByCategoryAndRankingOrderParam(pt_ranking.GetRankingsAndCountByCategoryAndRankingOrder)
	ranking_protocol.InsertRankingByPIDAndRankingScoreData(pt_ranking.InsertRankingByPIDAndRankingScoreData)
}
