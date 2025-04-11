package datastore

import (
	"github.com/PretendoNetwork/nex-go/v2"
	datastoretypes "github.com/PretendoNetwork/nex-protocols-go/v2/datastore/types"
)

func InitializeObjectRatingWithSlot(dataID uint64, param datastoretypes.DataStoreRatingInitParamWithSlot) *nex.Error {
	return nex.NewError(nex.ResultCodes.Core.NotImplemented, "Ratings not yet implemented")
}
