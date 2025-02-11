package ranking

import (
	"github.com/PretendoNetwork/nex-go/v2/types"
	ranking_types "github.com/PretendoNetwork/nex-protocols-go/v2/ranking/types"
)

// Stubbed
func GetRankingsAndCountByCategoryAndRankingOrder(category types.UInt32, rankingOrderParam ranking_types.RankingOrderParam) (types.List[ranking_types.RankingRankData], uint32, error) {
	return types.List[ranking_types.RankingRankData]{}, 0, nil
}
