package ranking

import (
	"fmt"

	ranking_types "github.com/PretendoNetwork/nex-protocols-go/ranking/types"
)

// Stubbed
func InsertRankingByPIDAndRankingScoreData(pid uint32, rankingScoreData *ranking_types.RankingScoreData, uniqueID uint64) error {
	fmt.Println(rankingScoreData)
	fmt.Println(uniqueID)
	fmt.Println(pid)

	return nil
}
