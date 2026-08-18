package params

import (
	"math/big"

	"github.com/ethereum/go-ethereum/params/types/ctypes"
)

var (
	KPOWChainID  = big.NewInt(112016)
	TKPOWChainID = big.NewInt(122016)
)

func IsKPOW(config ctypes.ChainConfigurator) bool {
	if config == nil {
		return false
	}

	chainID := config.GetChainID()
	if chainID == nil {
		return false
	}

	return chainID.Cmp(KPOWChainID) == 0 ||
		chainID.Cmp(TKPOWChainID) == 0
}
