package ethash

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params/vars"
)

const (
	// EMA smoothing window.
	kpowEMAWindow int64 = 8

	// Controller gain:
	//
	// 1 / 2 = 50%
	kpowGainDivisor int64 = 2

	// Maximum difficulty increase per block:
	//
	// 1 / 16 = 6.25%
	kpowMaxIncreaseDivisor int64 = 16

	// Maximum difficulty decrease per block:
	//
	// 1 / 8 = 12.5%
	kpowMaxDecreaseDivisor int64 = 8

	// Maximum timestamp delta considered by the controller:
	//
	// 6 * vars.DurationLimit
	// = 6 * 13
	// = 78 seconds
	kpowTimeDeltaClampMultiplier int64 = 6
)

var (
	kpowEMAWindowBig = big.NewInt(kpowEMAWindow)
	kpowEMAWeightBig = big.NewInt(kpowEMAWindow - 1)

	kpowGainDivisorBig        = big.NewInt(kpowGainDivisor)
	kpowMaxIncreaseDivisorBig = big.NewInt(kpowMaxIncreaseDivisor)
	kpowMaxDecreaseDivisorBig = big.NewInt(kpowMaxDecreaseDivisor)

	kpowTimeDeltaClampMultiplierBig = big.NewInt(kpowTimeDeltaClampMultiplier)
)

// calcKPOWDifficulty calculates the next KPOW/TKPOW PoW difficulty.
//
// Rules:
//   - Target block time: vars.DurationLimit (13 seconds).
//   - EMA-compatible smoothing window: 8.
//   - Controller gain: 50%.
//   - Maximum difficulty increase: +6.25% per block.
//   - Maximum difficulty decrease: -12.5% per block.
//   - Timestamp delta clamp: [1, 6 * target] = [1, 78] seconds.
//   - Minimum difficulty: vars.MinimumDifficulty.
//   - Difficulty bomb: none.
func calcKPOWDifficulty(time uint64, parent *types.Header) *big.Int {
	// Calculate elapsed time from the parent block.
	timeDelta := parent_time_delta(time, parent)

	// Header validation normally guarantees time > parent.Time.
	// Keep the calculator deterministic when called directly with
	// an invalid timestamp.
	if timeDelta.Sign() <= 0 {
		timeDelta.Set(big1)
	}

	// Clamp excessively large timestamp deltas.
	maxTimeDelta := new(big.Int).Mul(
		vars.DurationLimit,
		kpowTimeDeltaClampMultiplierBig,
	)
	if timeDelta.Cmp(maxTimeDelta) > 0 {
		timeDelta.Set(maxTimeDelta)
	}

	// Estimate the difficulty implied by the observed block time:
	//
	//     observed = parentDifficulty * target / actualTime
	observed := new(big.Int).Mul(
		parent.Difficulty,
		vars.DurationLimit,
	)
	observed.Div(observed, timeDelta)

	// Calculate the difference between the observed difficulty and
	// the current parent difficulty.
	diffError := new(big.Int).Sub(
		observed,
		parent.Difficulty,
	)

	// Apply 50% controller gain:
	//
	//     controlled = parentDifficulty + error / 2
	adjustment := new(big.Int).Div(
		diffError,
		kpowGainDivisorBig,
	)
	controlled := new(big.Int).Add(
		parent.Difficulty,
		adjustment,
	)

	// Apply window-8 smoothing:
	//
	//     next = (7 * parentDifficulty + controlled) / 8
	out := new(big.Int).Mul(
		parent.Difficulty,
		kpowEMAWeightBig,
	)
	out.Add(out, controlled)
	out.Div(out, kpowEMAWindowBig)

	// Maximum increase:
	//
	//     upper = parentDifficulty + parentDifficulty / 16
	//           = +6.25%
	maxIncrease := new(big.Int).Div(
		new(big.Int).Set(parent.Difficulty),
		kpowMaxIncreaseDivisorBig,
	)
	upperBound := new(big.Int).Add(
		parent.Difficulty,
		maxIncrease,
	)

	// Maximum decrease:
	//
	//     lower = parentDifficulty - parentDifficulty / 8
	//           = -12.5%
	maxDecrease := new(big.Int).Div(
		new(big.Int).Set(parent.Difficulty),
		kpowMaxDecreaseDivisorBig,
	)
	lowerBound := new(big.Int).Sub(
		parent.Difficulty,
		maxDecrease,
	)

	// Enforce per-block adjustment bounds.
	if out.Cmp(upperBound) > 0 {
		out.Set(upperBound)
	} else if out.Cmp(lowerBound) < 0 {
		out.Set(lowerBound)
	}

	// Enforce the standard Ethash minimum difficulty.
	out.Set(math.BigMax(out, vars.MinimumDifficulty))

	return out
}
