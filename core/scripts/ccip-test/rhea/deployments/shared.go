package deployments

import (
	"fmt"
	"time"

	"github.com/smartcontractkit/chainlink/core/scripts/ccip-test/rhea"
)

const (
	BATCH_GAS_LIMIT                     = 5_000_000
	FEE_UPDATE_HEARTBEAT                = 24 * time.Hour
	FEE_UPDATE_DEVIATION_PPB            = 10e7 // 10%
	FEE_UPDATE_DEVIATION_PPB_FAST_CHAIN = 20e7 // 20%
	RELATIVE_BOOST_PER_WAIT_HOUR        = 600  // 10x every minute, fees artificially low (0.1 on testnet)
	INFLIGHT_CACHE_EXPIRY               = 3 * time.Minute
	ROOT_SNOOZE_TIME                    = 6 * time.Minute
)

func getBlockConfirmations(chain rhea.Chain) uint32 {
	var blockConfirmationPerChain = map[rhea.Chain]uint32{
		rhea.Goerli:         4,
		rhea.Sepolia:        4,
		rhea.OptimismGoerli: 4,
		rhea.AvaxFuji:       1,
		rhea.PolygonMumbai:  4,
		rhea.ArbitrumGoerli: 1,
		rhea.Quorum:         4,
	}

	if val, ok := blockConfirmationPerChain[chain]; ok {
		return val
	}
	panic(fmt.Sprintf("Block confirmation for %s not found", chain))
}

func getMaxGasPrice(chain rhea.Chain) uint64 {
	var maxGasPricePerChain = map[rhea.Chain]uint64{
		rhea.Goerli:         200e9,
		rhea.Sepolia:        200e9,
		rhea.OptimismGoerli: 200e9,
		rhea.AvaxFuji:       200e9,
		rhea.PolygonMumbai:  200e9,
		rhea.ArbitrumGoerli: 200e9,
		rhea.Quorum:         200e9,
	}

	if val, ok := maxGasPricePerChain[chain]; ok {
		return val
	}
	panic(fmt.Sprintf("Max gas price for %s not found", chain))
}
