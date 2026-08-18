package mutations

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/params/types/coregeth"
	"github.com/ethereum/go-ethereum/params/types/ctypes"
	"github.com/holiman/uint256"
)

const (
	// Synthetic KPOW rewards used only by these tests.
	//
	// These values intentionally do not depend on the production
	// KPOW BlockRewardSchedule heights.
	kpowTestInitialReward uint64 = 2e18
	kpowTestSecondReward  uint64 = 1e18
	kpowTestThirdReward   uint64 = 5e17

	// Synthetic transition heights used only to exercise the
	// reward schedule selection logic.
	kpowTestFirstRewardTransition  uint64 = 100
	kpowTestSecondRewardTransition uint64 = 200
	kpowTestRewardEnd              uint64 = 300

	kpowTestInitialUncleReward = kpowTestInitialReward / 32
	kpowTestSecondUncleReward  = kpowTestSecondReward / 32
	kpowTestThirdUncleReward   = kpowTestThirdReward / 32
)

func newKPOWTestConfig(chainID *big.Int) *coregeth.CoreGethChainConfig {
	return &coregeth.CoreGethChainConfig{
		ChainID: new(big.Int).Set(chainID),
		Ethash:  new(ctypes.EthashConfig),

		BlockRewardSchedule: ctypes.Uint64Uint256MapEncodesHex{
			0:                              uint256.NewInt(kpowTestInitialReward),
			kpowTestFirstRewardTransition:  uint256.NewInt(kpowTestSecondReward),
			kpowTestSecondRewardTransition: uint256.NewInt(kpowTestThirdReward),
			kpowTestRewardEnd:              uint256.NewInt(0),
		},
	}
}

func checkKPOWReward(
	t *testing.T,
	name string,
	got *uint256.Int,
	want uint64,
) {
	t.Helper()

	if got == nil {
		t.Fatalf("%s: got nil, want %d", name, want)
	}

	if got.Cmp(uint256.NewInt(want)) != 0 {
		t.Errorf(
			"%s: got %v, want %d",
			name,
			got,
			want,
		)
	}
}

func TestKPOWBlockRewardSchedule(t *testing.T) {
	config := newKPOWTestConfig(params.KPOWChainID)

	tests := []struct {
		name  string
		block uint64
		want  uint64
	}{
		{
			name:  "genesis",
			block: 0,
			want:  kpowTestInitialReward,
		},
		{
			name:  "before first transition",
			block: kpowTestFirstRewardTransition - 1,
			want:  kpowTestInitialReward,
		},
		{
			name:  "at first transition",
			block: kpowTestFirstRewardTransition,
			want:  kpowTestSecondReward,
		},
		{
			name:  "before second transition",
			block: kpowTestSecondRewardTransition - 1,
			want:  kpowTestSecondReward,
		},
		{
			name:  "at second transition",
			block: kpowTestSecondRewardTransition,
			want:  kpowTestThirdReward,
		},
		{
			name:  "before reward end",
			block: kpowTestRewardEnd - 1,
			want:  kpowTestThirdReward,
		},
		{
			name:  "at reward end",
			block: kpowTestRewardEnd,
			want:  0,
		},
		{
			name:  "after reward end",
			block: kpowTestRewardEnd + 1,
			want:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ctypes.EthashBlockReward(
				config,
				new(big.Int).SetUint64(tt.block),
			)

			checkKPOWReward(
				t,
				"block reward",
				got,
				tt.want,
			)
		})
	}
}

func TestGetKPOWRewards(t *testing.T) {
	config := newKPOWTestConfig(params.KPOWChainID)

	tests := []struct {
		name        string
		block       uint64
		uncleDepths []int64
		wantWinner  uint64
		wantUncle   uint64
	}{
		{
			name:       "initial reward without uncles",
			block:      50,
			wantWinner: kpowTestInitialReward,
		},
		{
			name:        "initial reward with one uncle",
			block:       50,
			uncleDepths: []int64{1},
			wantWinner: kpowTestInitialReward +
				kpowTestInitialUncleReward,
			wantUncle: kpowTestInitialUncleReward,
		},
		{
			name:        "initial reward with two uncles",
			block:       50,
			uncleDepths: []int64{1, 6},
			wantWinner: kpowTestInitialReward +
				2*kpowTestInitialUncleReward,
			wantUncle: kpowTestInitialUncleReward,
		},
		{
			name:        "second reward with one uncle",
			block:       kpowTestFirstRewardTransition,
			uncleDepths: []int64{1},
			wantWinner: kpowTestSecondReward +
				kpowTestSecondUncleReward,
			wantUncle: kpowTestSecondUncleReward,
		},
		{
			name:        "second reward with two uncles",
			block:       kpowTestFirstRewardTransition,
			uncleDepths: []int64{1, 6},
			wantWinner: kpowTestSecondReward +
				2*kpowTestSecondUncleReward,
			wantUncle: kpowTestSecondUncleReward,
		},
		{
			name:        "third reward with one uncle",
			block:       kpowTestSecondRewardTransition,
			uncleDepths: []int64{1},
			wantWinner: kpowTestThirdReward +
				kpowTestThirdUncleReward,
			wantUncle: kpowTestThirdUncleReward,
		},
		{
			name:        "third reward with two uncles",
			block:       kpowTestSecondRewardTransition,
			uncleDepths: []int64{1, 6},
			wantWinner: kpowTestThirdReward +
				2*kpowTestThirdUncleReward,
			wantUncle: kpowTestThirdUncleReward,
		},
		{
			name:        "zero reward with one uncle",
			block:       kpowTestRewardEnd,
			uncleDepths: []int64{1},
			wantWinner:  0,
			wantUncle:   0,
		},
		{
			name:        "zero reward with two uncles",
			block:       kpowTestRewardEnd,
			uncleDepths: []int64{1, 6},
			wantWinner:  0,
			wantUncle:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			header := &types.Header{
				Number: new(big.Int).SetUint64(tt.block),
			}

			uncles := make(
				[]*types.Header,
				len(tt.uncleDepths),
			)

			for i, depth := range tt.uncleDepths {
				uncles[i] = &types.Header{
					Number: new(big.Int).Sub(
						header.Number,
						big.NewInt(depth),
					),
				}
			}

			winnerReward, uncleRewards := getKPOWRewards(
				config,
				header,
				uncles,
			)

			checkKPOWReward(
				t,
				"winner reward",
				winnerReward,
				tt.wantWinner,
			)

			if len(uncleRewards) != len(uncles) {
				t.Fatalf(
					"uncle reward count: got %d, want %d",
					len(uncleRewards),
					len(uncles),
				)
			}

			for _, reward := range uncleRewards {
				checkKPOWReward(
					t,
					"uncle reward",
					reward,
					tt.wantUncle,
				)
			}
		})
	}
}

func TestGetKPOWRewardsIgnoresUncleDepth(t *testing.T) {
	config := newKPOWTestConfig(params.KPOWChainID)

	header := &types.Header{
		Number: big.NewInt(50),
	}

	uncles := []*types.Header{
		{
			Number: big.NewInt(49), // depth 1
		},
		{
			Number: big.NewInt(44), // depth 6
		},
	}

	_, uncleRewards := getKPOWRewards(
		config,
		header,
		uncles,
	)

	if len(uncleRewards) != len(uncles) {
		t.Fatalf(
			"uncle reward count: got %d, want %d",
			len(uncleRewards),
			len(uncles),
		)
	}

	checkKPOWReward(
		t,
		"depth-1 uncle reward",
		uncleRewards[0],
		kpowTestInitialUncleReward,
	)

	checkKPOWReward(
		t,
		"depth-6 uncle reward",
		uncleRewards[1],
		kpowTestInitialUncleReward,
	)
}

func TestGetKPOWRewardsDoesNotMutateBlockReward(t *testing.T) {
	config := newKPOWTestConfig(params.KPOWChainID)

	header := &types.Header{
		Number: big.NewInt(50),
	}

	uncles := []*types.Header{
		{
			Number: big.NewInt(49),
		},
		{
			Number: big.NewInt(44),
		},
	}

	before := new(uint256.Int).Set(
		config.BlockRewardSchedule[0],
	)

	getKPOWRewards(
		config,
		header,
		uncles,
	)

	if got := config.BlockRewardSchedule[0]; got.Cmp(before) != 0 {
		t.Errorf(
			"configured block reward mutated: got %v, want %v",
			got,
			before,
		)
	}
}

func TestGetRewardsUsesKPOWRewards(t *testing.T) {
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
			config := newKPOWTestConfig(tt.chainID)

			header := &types.Header{
				Number: big.NewInt(50),
			}

			uncles := []*types.Header{
				{
					Number: big.NewInt(49),
				},
			}

			winnerReward, uncleRewards := GetRewards(
				config,
				header,
				uncles,
			)

			checkKPOWReward(
				t,
				"winner reward",
				winnerReward,
				kpowTestInitialReward+
					kpowTestInitialUncleReward,
			)

			if len(uncleRewards) != 1 {
				t.Fatalf(
					"uncle reward count: got %d, want 1",
					len(uncleRewards),
				)
			}

			checkKPOWReward(
				t,
				"uncle reward",
				uncleRewards[0],
				kpowTestInitialUncleReward,
			)
		})
	}
}

func TestGetRewardsKPOWPrecedesECIP1017(t *testing.T) {
	config := newKPOWTestConfig(params.KPOWChainID)

	// Intentionally enable ECIP-1017.
	//
	// KPOW must still use its own reward path before
	// the generic ECIP-1017 reward path.
	config.ECIP1017FBlock = big.NewInt(0)
	config.ECIP1017EraRounds = new(big.Int).SetUint64(
		kpowTestFirstRewardTransition,
	)

	header := &types.Header{
		Number: new(big.Int).SetUint64(
			kpowTestFirstRewardTransition + 1,
		),
	}

	uncles := []*types.Header{
		{
			Number: new(big.Int).SetUint64(
				kpowTestFirstRewardTransition,
			),
		},
	}

	winnerReward, uncleRewards := GetRewards(
		config,
		header,
		uncles,
	)

	checkKPOWReward(
		t,
		"winner reward",
		winnerReward,
		kpowTestSecondReward+
			kpowTestSecondUncleReward,
	)

	if len(uncleRewards) != 1 {
		t.Fatalf(
			"uncle reward count: got %d, want 1",
			len(uncleRewards),
		)
	}

	checkKPOWReward(
		t,
		"uncle reward",
		uncleRewards[0],
		kpowTestSecondUncleReward,
	)
}

func TestGetRewardsNonKPOWKeepsLegacyUncleReward(t *testing.T) {
	config := newKPOWTestConfig(big.NewInt(1))

	header := &types.Header{
		Number: big.NewInt(50),
	}

	uncles := []*types.Header{
		{
			Number: big.NewInt(49), // depth 1
		},
	}

	winnerReward, uncleRewards := GetRewards(
		config,
		header,
		uncles,
	)

	// Legacy winner inclusion reward remains R/32.
	checkKPOWReward(
		t,
		"winner reward",
		winnerReward,
		kpowTestInitialReward+
			kpowTestInitialUncleReward,
	)

	if len(uncleRewards) != 1 {
		t.Fatalf(
			"uncle reward count: got %d, want 1",
			len(uncleRewards),
		)
	}

	// Legacy depth-1 uncle reward:
	//
	//     R * (8 - 1) / 8
	//   = 7R / 8
	checkKPOWReward(
		t,
		"legacy uncle reward",
		uncleRewards[0],
		7*kpowTestInitialReward/8,
	)
}

func TestAccumulateKPOWRewards(t *testing.T) {
	config := newKPOWTestConfig(params.KPOWChainID)

	db := rawdb.NewMemoryDatabase()
	defer db.Close()

	stateDB, err := state.New(
		common.Hash{},
		state.NewDatabase(db),
		nil,
	)
	if err != nil {
		t.Fatalf(
			"could not open statedb: %v",
			err,
		)
	}

	header := &types.Header{
		Number: big.NewInt(50),
		Coinbase: common.HexToAddress(
			"0000000000000000000000000000000000000001",
		),
	}

	uncles := []*types.Header{
		{
			Number: big.NewInt(49),
			Coinbase: common.HexToAddress(
				"0000000000000000000000000000000000000002",
			),
		},
		{
			Number: big.NewInt(44),
			Coinbase: common.HexToAddress(
				"0000000000000000000000000000000000000003",
			),
		},
	}

	AccumulateRewards(
		config,
		stateDB,
		header,
		uncles,
	)

	checkKPOWReward(
		t,
		"winner balance",
		stateDB.GetBalance(header.Coinbase),
		kpowTestInitialReward+
			2*kpowTestInitialUncleReward,
	)

	checkKPOWReward(
		t,
		"uncle 1 balance",
		stateDB.GetBalance(uncles[0].Coinbase),
		kpowTestInitialUncleReward,
	)

	checkKPOWReward(
		t,
		"uncle 2 balance",
		stateDB.GetBalance(uncles[1].Coinbase),
		kpowTestInitialUncleReward,
	)
}
