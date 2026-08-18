package mutations

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params/types/ctypes"
	"github.com/holiman/uint256"
)

// getKPOWRewards calculates mining rewards for KPOW/TKPOW.
//
// KPOW uncle reward rules:
//   - Winning miner: base block reward.
//   - Winning miner: +1/32 block reward for each included uncle.
//   - Uncle miner:    1/32 block reward.
//   - Uncle depth does not affect the reward.
//
// Uncle validation, maximum uncle count and allowed uncle depth
// are handled separately by the block validator.
func getKPOWRewards(
	config ctypes.ChainConfigurator,
	header *types.Header,
	uncles []*types.Header,
) (*uint256.Int, []*uint256.Int) {
	blockReward := ctypes.EthashBlockReward(config, header.Number)

	// Base reward for the winning miner.
	reward := new(uint256.Int).Set(blockReward)

	// Rewards for uncle miners.
	uncleRewards := make([]*uint256.Int, len(uncles))

	// Uncle miner and winning miner receive the same
	// 1/32 reward for each included uncle.
	perUncleReward := new(uint256.Int).Div(
		new(uint256.Int).Set(blockReward),
		big32,
	)

	for i := range uncles {
		uncleRewards[i] = new(uint256.Int).Set(perUncleReward)
		reward.Add(reward, perUncleReward)
	}

	return reward, uncleRewards
}
