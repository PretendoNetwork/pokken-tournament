package ranking

import (
	"fmt"

	"github.com/PretendoNetwork/nex-go/v2/types"
	ranking_types "github.com/PretendoNetwork/nex-protocols-go/v2/ranking/types"
)

// Stubbed
func InsertRankingByPIDAndRankingScoreData(pid types.PID, rankingScoreData ranking_types.RankingScoreData, uniqueID types.UInt64) error {
	fmt.Println(rankingScoreData)
	fmt.Println(uniqueID)
	fmt.Println(pid)

	return nil
}
