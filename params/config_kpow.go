package params

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params/types/coregeth"
	"github.com/ethereum/go-ethereum/params/types/ctypes"
	"github.com/holiman/uint256"
)

var (
	// KPOWChainConfig is the chain parameters to run a node on the KPOW network (PoW).
	KPOWChainConfig = &coregeth.CoreGethChainConfig{
		NetworkID: KPOWChainID.Uint64(),
		ChainID:   new(big.Int).Set(KPOWChainID),
		Ethash:    new(ctypes.EthashConfig),

		// Homestead
		EIP2FBlock: big.NewInt(0),
		EIP7FBlock: big.NewInt(0),

		// Tangerine Whistle
		EIP150Block: big.NewInt(0),

		// Spurious Dragon
		EIP155Block: big.NewInt(0),

		// EIP158 eq
		EIP160FBlock: big.NewInt(0),
		EIP161FBlock: big.NewInt(0),
		EIP170FBlock: big.NewInt(0),

		// Byzantium eq
		EIP100FBlock: big.NewInt(0),
		EIP140FBlock: big.NewInt(0),
		EIP198FBlock: big.NewInt(0),
		EIP211FBlock: big.NewInt(0),
		EIP212FBlock: big.NewInt(0),
		EIP213FBlock: big.NewInt(0),
		EIP214FBlock: big.NewInt(0),
		EIP658FBlock: big.NewInt(0),

		// Constantinople eq, aka Agharta
		EIP145FBlock:  big.NewInt(0),
		EIP1014FBlock: big.NewInt(0),
		EIP1052FBlock: big.NewInt(0),

		// Do not enable EIP-1283 directly.
		// EIP-2200 below provides the final SSTORE net gas metering rules.
		EIP1283FBlock:   nil,
		PetersburgBlock: nil,

		// Istanbul eq, aka Phoenix
		// ECIP-1088
		EIP152FBlock:  big.NewInt(0),
		EIP1108FBlock: big.NewInt(0),
		EIP1344FBlock: big.NewInt(0),
		EIP1884FBlock: big.NewInt(0),
		EIP2028FBlock: big.NewInt(0),
		EIP2200FBlock: big.NewInt(0),

		// Berlin-compatible EVM rules
		EIP2565FBlock: big.NewInt(0), // ModExp gas repricing
		EIP2718FBlock: big.NewInt(0), // Typed transaction envelopes
		EIP2929FBlock: big.NewInt(0), // Cold/warm state access gas
		EIP2930FBlock: big.NewInt(0), // Access-list transactions

		// London EVM rules which don't require EIP-1559 fee market
		EIP3529FBlock: big.NewInt(0), // Reduce gas refunds
		EIP3541FBlock: big.NewInt(0), // Reject new contract code starting with 0xEF

		// No EIP-1559 / BASEFEE.
		// KPOW keeps the legacy PoW gas fee model.
		EIP1559FBlock: nil,
		EIP3198FBlock: nil,

		// No Ethereum difficulty-bomb-delay forks.
		// Difficulty bomb is disposed from genesis below.
		EIP3554FBlock: nil,

		// No Merge / PREVRANDAO.
		// KPOW remains Ethash PoW.
		EIP4399FBlock: nil,

		// Shanghai EVM features enabled from genesis.
		// These EVM changes can be used independently of PoS.
		EIP3651FBlock: big.NewInt(0), // Warm COINBASE
		EIP3855FBlock: big.NewInt(0), // PUSH0
		EIP3860FBlock: big.NewInt(0), // Limit and meter initcode

		// No Beacon Chain withdrawals.
		EIP4895FBlock: nil,

		// EIP-6049 only deprecates SELFDESTRUCT and does not itself
		// change SELFDESTRUCT execution semantics.
		EIP6049FBlock: nil,

		// ETC-specific transitions
		ECIP1099FBlock: nil, // No Etchash transition

		// Dispose difficulty bomb from genesis.
		DisposalBlock: big.NewInt(0),

		// No Ethereum Classic disinflationary monetary policy.
		ECIP1017FBlock:    nil,
		ECIP1017EraRounds: nil,

		// No need to delay difficulty bomb; it is disposed from genesis.
		ECIP1010PauseBlock: nil,
		ECIP1010Length:     nil,

		// No MESS artificial finality.
		ECBP1100FBlock: nil,

		RequireBlockHashes: map[uint64]common.Hash{},

		BlockRewardSchedule: ctypes.Uint64Uint256MapEncodesHex{
			0:          uint256.NewInt(2e18), // 2 KPOW
			4_000_000:  uint256.NewInt(1e18), // 1 KPOW
			8_400_000:  uint256.NewInt(5e17), // 0.5 KPOW
			13_000_000: uint256.NewInt(0),    // 0 KPOW
		},
	}
)
