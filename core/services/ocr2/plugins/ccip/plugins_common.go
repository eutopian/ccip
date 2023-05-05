package ccip

import (
	"context"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/commit_store"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/evm_2_evm_offramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/evm_2_evm_onramp"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	ccipconfig "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/config"
)

var ErrCommitStoreIsDown = errors.New("commitStore is down")

func LoadOnRamp(onRampAddress common.Address, client client.Client) (*evm_2_evm_onramp.EVM2EVMOnRamp, error) {
	err := ccipconfig.VerifyTypeAndVersion(onRampAddress, client, ccipconfig.EVM2EVMOnRamp)
	if err != nil {
		return nil, errors.Wrap(err, "Invalid onRamp contract")
	}
	return evm_2_evm_onramp.NewEVM2EVMOnRamp(onRampAddress, client)
}

func LoadOffRamp(offRampAddress common.Address, client client.Client) (*evm_2_evm_offramp.EVM2EVMOffRamp, error) {
	err := ccipconfig.VerifyTypeAndVersion(offRampAddress, client, ccipconfig.EVM2EVMOffRamp)
	if err != nil {
		return nil, errors.Wrap(err, "Invalid offRamp contract")
	}
	return evm_2_evm_offramp.NewEVM2EVMOffRamp(offRampAddress, client)
}

func LoadCommitStore(commitStoreAddress common.Address, client client.Client) (*commit_store.CommitStore, error) {
	err := ccipconfig.VerifyTypeAndVersion(commitStoreAddress, client, ccipconfig.CommitStore)
	if err != nil {
		return nil, errors.Wrap(err, "Invalid commitStore contract")
	}
	return commit_store.NewCommitStore(commitStoreAddress, client)
}

func contiguousReqs(lggr logger.Logger, min, max uint64, seqNrs []uint64) bool {
	if int(max-min+1) != len(seqNrs) {
		return false
	}
	for i, j := min, 0; i <= max && j < len(seqNrs); i, j = i+1, j+1 {
		if seqNrs[j] != i {
			lggr.Errorw("unexpected gap in seq nums", "seq", i)
			return false
		}
	}
	return true
}

// Extracts the hashed leaves from a given set of logs
func leavesFromIntervals(
	lggr logger.Logger,
	seqParser func(logpoller.Log) (uint64, error),
	interval commit_store.CommitStoreInterval,
	hasher LeafHasherInterface[[32]byte],
	logs []logpoller.Log,
) ([][32]byte, error) {
	var seqNrs []uint64
	for _, log := range logs {
		seqNr, err2 := seqParser(log)
		if err2 != nil {
			return nil, err2
		}
		seqNrs = append(seqNrs, seqNr)
	}
	if !contiguousReqs(lggr, interval.Min, interval.Max, seqNrs) {
		return nil, errors.Errorf("do not have full range [%v, %v] have %v", interval.Min, interval.Max, seqNrs)
	}
	var leaves [][32]byte
	for _, log := range logs {
		hash, err2 := hasher.HashLeaf(log.ToGethLog())
		if err2 != nil {
			return nil, err2
		}
		leaves = append(leaves, hash)
	}

	return leaves, nil
}

// Checks whether the commit store is down by doing an onchain check for
// Paused and AFN status
func isCommitStoreDownNow(ctx context.Context, lggr logger.Logger, commitStore commit_store.CommitStoreInterface) bool {
	paused, err := commitStore.Paused(&bind.CallOpts{Context: ctx})
	if err != nil {
		// Air on side of caution by halting if we cannot read the state?
		lggr.Errorw("Unable to read CommitStore paused", "err", err)
		return true
	}
	healthy, err := commitStore.IsAFNHealthy(&bind.CallOpts{Context: ctx})
	if err != nil {
		lggr.Errorw("Unable to read CommitStore AFN state", "err", err)
		return true
	}
	return paused || !healthy
}
