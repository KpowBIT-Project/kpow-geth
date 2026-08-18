package params

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/params/types/genesisT"
)

var TKPOWGenesisHash = common.HexToHash("0xa5fe61c6f76e7fc8bf331ccc84c5e2ae286f02afbeb1e49a8cbd4e64dc9cf1a5")

func DefaultTKPOWGenesisBlock() *genesisT.Genesis {
	return &genesisT.Genesis{
		Config:     TKPOWChainConfig,
		Nonce:      hexutil.MustDecodeUint64("0x0"),
		ExtraData:  []byte(""),
		GasLimit:   hexutil.MustDecodeUint64("0x1c9c380"),
		Difficulty: hexutil.MustDecodeBig("0xE4E1C0"),
		Timestamp:  1786705293,
		Alloc: genesisT.GenesisAlloc{
			common.HexToAddress("0xD04079c9eF5fdF82110FF24B1F8aE11Aa06eD4E2"): genesisT.GenesisAccount{
				Balance: new(big.Int).Mul(
					big.NewInt(4_200_000),
					big.NewInt(1e18),
				),
			},
			common.HexToAddress("0xEa12141040580222dbDB432B2F6f23bD5E1d27E4"): genesisT.GenesisAccount{
				Balance: new(big.Int).Mul(
					big.NewInt(2_100_000),
					big.NewInt(1e18),
				),
			},
		},
	}
}
