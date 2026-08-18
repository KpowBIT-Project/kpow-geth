package ethash

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/params/types/coregeth"
	"github.com/ethereum/go-ethereum/params/types/ctypes"
	"github.com/ethereum/go-ethereum/params/vars"
)

const kpowTestParentDifficulty int64 = 16_000_000

func newKPOWDifficultyTestParent() *types.Header {
	return &types.Header{
		Number:     big.NewInt(100),
		Time:       1_000,
		Difficulty: big.NewInt(kpowTestParentDifficulty),
	}
}

func newKPOWDifficultyTestConfig(chainID *big.Int) *coregeth.CoreGethChainConfig {
	return &coregeth.CoreGethChainConfig{
		ChainID: new(big.Int).Set(chainID),
		Ethash:  new(ctypes.EthashConfig),
	}
}

func checkKPOWDifficulty(t *testing.T, got *big.Int, want int64) {
	t.Helper()

	if got == nil {
		t.Fatalf("difficulty: got nil, want %d", want)
	}
	if got.Cmp(big.NewInt(want)) != 0 {
		t.Errorf("difficulty: got %v, want %d", got, want)
	}
}

func TestCalcKPOWDifficulty(t *testing.T) {
	tests := []struct {
		name      string
		timeDelta uint64
		want      int64
	}{
		{
			name:      "target block time",
			timeDelta: 13,
			want:      16_000_000,
		},
		{
			// observed = 16M * 13 / 12 = 17,333,333
			// gain 50% => controlled = 16,666,666
			// window 8 => next = 16,083,333
			name:      "slightly fast block",
			timeDelta: 12,
			want:      16_083_333,
		},
		{
			// observed = 20.8M
			// gain 50% => controlled = 18.4M
			// window 8 => next = 16.3M
			name:      "fast block",
			timeDelta: 10,
			want:      16_300_000,
		},
		{
			// Fast, but still below the +6.25% clamp.
			name:      "fast block below increase clamp",
			timeDelta: 7,
			want:      16_857_142,
		},
		{
			// Raw output exceeds +6.25%.
			name:      "maximum increase clamp",
			timeDelta: 6,
			want:      17_000_000,
		},
		{
			name:      "extremely fast block",
			timeDelta: 1,
			want:      17_000_000,
		},
		{
			name:      "slightly slow block",
			timeDelta: 14,
			want:      15_928_571,
		},
		{
			// observed = 10.4M
			// gain 50% => controlled = 13.2M
			// window 8 => next = 15.65M
			name:      "slow block",
			timeDelta: 20,
			want:      15_650_000,
		},
		{
			name:      "double target block time",
			timeDelta: 26,
			want:      15_500_000,
		},
		{
			name:      "long block",
			timeDelta: 50,
			want:      15_260_000,
		},
		{
			// Maximum timestamp delta used by the controller is 78 seconds.
			name:      "timestamp clamp boundary",
			timeDelta: 78,
			want:      15_166_666,
		},
		{
			// 1000 seconds is treated exactly like 78 seconds.
			name:      "timestamp above clamp",
			timeDelta: 1_000,
			want:      15_166_666,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parent := newKPOWDifficultyTestParent()

			got := calcKPOWDifficulty(
				parent.Time+tt.timeDelta,
				parent,
			)

			checkKPOWDifficulty(t, got, tt.want)
		})
	}
}

func TestCalcKPOWDifficultyTargetKeepsDifficulty(t *testing.T) {
	parent := newKPOWDifficultyTestParent()

	got := calcKPOWDifficulty(
		parent.Time+vars.DurationLimit.Uint64(),
		parent,
	)

	if got.Cmp(parent.Difficulty) != 0 {
		t.Errorf(
			"target block time changed difficulty: got %v, want %v",
			got,
			parent.Difficulty,
		)
	}
}

func TestCalcKPOWDifficultyTimestampClamp(t *testing.T) {
	parent := newKPOWDifficultyTestParent()

	maxTimeDelta := new(big.Int).Mul(
		vars.DurationLimit,
		big.NewInt(kpowTimeDeltaClampMultiplier),
	)

	atClamp := calcKPOWDifficulty(
		parent.Time+maxTimeDelta.Uint64(),
		parent,
	)

	beyondClamp := calcKPOWDifficulty(
		parent.Time+10_000,
		parent,
	)

	if atClamp.Cmp(beyondClamp) != 0 {
		t.Errorf(
			"timestamp clamp mismatch: at clamp %v, beyond clamp %v",
			atClamp,
			beyondClamp,
		)
	}
}

func TestCalcKPOWDifficultyInvalidTimestamp(t *testing.T) {
	tests := []struct {
		name string
		time uint64
	}{
		{
			name: "same timestamp as parent",
			time: 1_000,
		},
		{
			name: "timestamp before parent",
			time: 999,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parent := newKPOWDifficultyTestParent()

			got := calcKPOWDifficulty(
				tt.time,
				parent,
			)

			// Non-positive delta is treated as 1 second.
			// The result is then capped at +6.25%.
			checkKPOWDifficulty(t, got, 17_000_000)
		})
	}
}

func TestCalcKPOWDifficultyIncreaseLimit(t *testing.T) {
	parent := newKPOWDifficultyTestParent()

	maxIncrease := new(big.Int).Div(
		new(big.Int).Set(parent.Difficulty),
		big.NewInt(kpowMaxIncreaseDivisor),
	)
	upperBound := new(big.Int).Add(
		new(big.Int).Set(parent.Difficulty),
		maxIncrease,
	)

	for timeDelta := uint64(1); timeDelta <= 12; timeDelta++ {
		got := calcKPOWDifficulty(
			parent.Time+timeDelta,
			parent,
		)

		if got.Cmp(upperBound) > 0 {
			t.Errorf(
				"delta %d: difficulty exceeded +6.25%% limit: got %v, max %v",
				timeDelta,
				got,
				upperBound,
			)
		}
	}
}

func TestCalcKPOWDifficultyDecreaseLimit(t *testing.T) {
	parent := newKPOWDifficultyTestParent()

	maxDecrease := new(big.Int).Div(
		new(big.Int).Set(parent.Difficulty),
		big.NewInt(kpowMaxDecreaseDivisor),
	)
	lowerBound := new(big.Int).Sub(
		new(big.Int).Set(parent.Difficulty),
		maxDecrease,
	)

	tests := []uint64{
		14,
		20,
		26,
		50,
		78,
		100,
		1_000,
		10_000,
	}

	for _, timeDelta := range tests {
		got := calcKPOWDifficulty(
			parent.Time+timeDelta,
			parent,
		)

		if got.Cmp(lowerBound) < 0 {
			t.Errorf(
				"delta %d: difficulty exceeded -12.5%% limit: got %v, min %v",
				timeDelta,
				got,
				lowerBound,
			)
		}
	}
}

func TestCalcKPOWDifficultyMinimumDifficulty(t *testing.T) {
	parent := &types.Header{
		Number:     big.NewInt(100),
		Time:       1_000,
		Difficulty: new(big.Int).Set(vars.MinimumDifficulty),
	}

	got := calcKPOWDifficulty(
		parent.Time+78,
		parent,
	)

	if got.Cmp(vars.MinimumDifficulty) != 0 {
		t.Errorf(
			"minimum difficulty: got %v, want %v",
			got,
			vars.MinimumDifficulty,
		)
	}
}

func TestCalcKPOWDifficultyDoesNotMutateParent(t *testing.T) {
	parent := newKPOWDifficultyTestParent()

	beforeNumber := new(big.Int).Set(parent.Number)
	beforeDifficulty := new(big.Int).Set(parent.Difficulty)
	beforeTime := parent.Time

	_ = calcKPOWDifficulty(
		parent.Time+10,
		parent,
	)

	if parent.Number.Cmp(beforeNumber) != 0 {
		t.Errorf(
			"parent number mutated: got %v, want %v",
			parent.Number,
			beforeNumber,
		)
	}
	if parent.Difficulty.Cmp(beforeDifficulty) != 0 {
		t.Errorf(
			"parent difficulty mutated: got %v, want %v",
			parent.Difficulty,
			beforeDifficulty,
		)
	}
	if parent.Time != beforeTime {
		t.Errorf(
			"parent timestamp mutated: got %d, want %d",
			parent.Time,
			beforeTime,
		)
	}
}

func TestCalcDifficultyUsesKPOWDifficulty(t *testing.T) {
	tests := []struct {
		name    string
		chainID *big.Int
	}{
		{
			name:    "KPOW",
			chainID: params.KPOWChainID,
		},
		{
			name:    "TKPOW",
			chainID: params.TKPOWChainID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := newKPOWDifficultyTestConfig(tt.chainID)
			parent := newKPOWDifficultyTestParent()
			time := parent.Time + 10

			got := CalcDifficulty(config, time, parent)
			want := calcKPOWDifficulty(time, parent)

			if got.Cmp(want) != 0 {
				t.Errorf(
					"difficulty: got %v, want %v",
					got,
					want,
				)
			}
		})
	}
}
