package params

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/params/types/genesisT"
)

var KPOWGenesisHash = common.HexToHash("0x0f61e878f4be493cce0290e81a52364a3e15cda23a0bde4717b5554d76ef9e38")

func DefaultKPOWGenesisBlock() *genesisT.Genesis {
	return &genesisT.Genesis{
		Config:     KPOWChainConfig,
		Nonce:      hexutil.MustDecodeUint64("0x0"),
		ExtraData:  []byte(""),
		GasLimit:   hexutil.MustDecodeUint64("0x1c9c380"),
		Difficulty: hexutil.MustDecodeBig("0x2CB417800"),
		Timestamp:  1786777406,
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
