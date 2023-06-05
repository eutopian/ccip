package legacygasstation

import (
	"context"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	geth_types "github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/evm_2_evm_offramp"
	forwarder_wrapper "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/forwarder"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/legacygasstation/types"
	"github.com/smartcontractkit/chainlink/v2/core/services/pg"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

// forwarderInterface is a narrow interface for forwarder contract gethwrapper
type forwarderInterface interface {
	Address() common.Address
	ParseLog(log geth_types.Log) (generated.AbigenLog, error)
}

// ccipOffRampInterface is a narrow interface for CCIP offramp interface gethwrapper
type ccipOffRampInterface interface {
	Address() common.Address
	ParseLog(log geth_types.Log) (generated.AbigenLog, error)
}

// Sidecar is responsible for listening to on-chain events and
// applying necessary status updates for legacy gasless txs
type Sidecar struct {
	orm         ORM
	lp          logpoller.LogPoller
	lggr        logger.Logger
	forwarder   forwarderInterface
	ccipOffRamp ccipOffRampInterface
	cfg         Config
	// ccipChainSelector is used to query legacy_gasless_txs database
	// the txs in database use ccip chain selectors instead of EVM chain IDs
	ccipChainSelector uint64
	lookbackBlocks    uint32
}

func NewSidecar(
	lggr logger.Logger,
	lp logpoller.LogPoller,
	forwarderInterface forwarderInterface,
	ccipOffRampInterface ccipOffRampInterface,
	cfg Config,
	ccipChainSelector uint64,
	lookbackBlocks uint32,
	orm ORM,
) (*Sidecar, error) {
	err := lp.RegisterFilter(logpoller.Filter{
		Name: logpoller.FilterName("Legacy Gas Station Sidecar", forwarderInterface.Address(), ccipOffRampInterface.Address()),
		EventSigs: []common.Hash{
			forwarder_wrapper.ForwarderForwardSucceeded{}.Topic(),
			evm_2_evm_offramp.EVM2EVMOffRampExecutionStateChanged{}.Topic(),
		},
		Addresses: []common.Address{
			forwarderInterface.Address(),
			ccipOffRampInterface.Address(),
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "register filter")
	}

	return &Sidecar{
		lp:                lp,
		lggr:              lggr,
		forwarder:         forwarderInterface,
		ccipOffRamp:       ccipOffRampInterface,
		cfg:               cfg,
		ccipChainSelector: ccipChainSelector,
		lookbackBlocks:    lookbackBlocks,
		orm:               orm,
	}, nil
}

func (sc *Sidecar) Run(ctx context.Context) error {
	latestBlock, err := sc.lp.LatestBlock(pg.WithParentCtx(ctx))
	if err != nil {
		return errors.Wrap(err, "latest block")
	}

	// we only care about logs that have passed finality depth
	toBlock := latestBlock - int64(sc.cfg.EvmFinalityDepth())
	if toBlock < 0 {
		return errors.Errorf("negative toBlock: %d", toBlock)
	}

	fromBlock := toBlock - int64(sc.lookbackBlocks)
	if fromBlock < 0 {
		fromBlock = 0
	}

	// Handle txs with status == Submitted
	// Following transitions are possible:
	// Submitted -> Finalized: same-chain transfer was finalized
	// Submitted -> SourceFinalized: cross-chain transfer was finalized on source chain
	// Submitted -> Failure: same-chain or cross-chain transfer failed
	err = sc.handleSubmittedTxs(ctx, fromBlock, toBlock)
	if err != nil {
		return errors.Wrap(err, "handle submitted transactions")
	}

	// Handle txs with status == SourceFinalized
	// Following transitions are possible:
	// SourceFinalized -> Finalized: cross-chain transfer was finalized on destination chain
	// SourceFinalized -> Failure: TODO: figure out failure scenarios for CCIP DON
	err = sc.handleSourceFinalizedTxs(ctx, fromBlock, toBlock)
	if err != nil {
		return errors.Wrap(err, "handle source finalized transactions")
	}
	return nil

}

func (sc *Sidecar) handleSubmittedTxs(ctx context.Context, fromBlock, toBlock int64) error {
	submittedTxs, err := sc.orm.SelectBySourceChainIDAndStatus(sc.ccipChainSelector, types.Submitted, pg.WithParentCtx(ctx))
	if err != nil {
		return errors.Wrap(err, "find by status")
	}

	txSenderAddresses := txSenderAddresses(submittedTxs)

	logs, err := sc.lp.IndexedLogsByBlockRange(
		fromBlock,
		toBlock,
		forwarder_wrapper.ForwarderForwardSucceeded{}.Topic(),
		sc.forwarder.Address(),
		1, // From address is the first indexed field
		txSenderAddresses,
		pg.WithParentCtx(ctx))
	if err != nil {
		return errors.Wrap(err, "log poller indexed logs")
	}

	fsLogs, err := sc.unmarshalForwardSucceededLogs(logs)
	if err != nil {
		return errors.Wrap(err, "unmarshal forward succeeded logs")
	}

	finalizedLogs, err := finalizedLogs(fsLogs, submittedTxs)
	if err != nil {
		return errors.Wrap(err, "finalized logs")
	}

	for _, log := range finalizedLogs {
		err = sc.orm.UpdateLegacyGaslessTx(log)
		if err != nil {
			return errors.Wrap(err, "update legacy gasless tx")
		}
	}

	failedTxs, err := sc.orm.SelectEthTxsBySourceChainIDAndState(sc.ccipChainSelector, txmgr.EthTxFatalError, pg.WithParentCtx(ctx))
	if err != nil {
		return errors.Wrap(err, "failed transactions")
	}

	return sc.updateFailedTxs(submittedTxs, failedTxs)
}

func (sc *Sidecar) updateFailedTxs(submittedTxs []types.LegacyGaslessTx, failedTxs []txmgr.DbEthTx) error {
	failedTxMap := make(map[string]txmgr.DbEthTx)
	for _, ethTx := range failedTxs {
		failedTxMap[strconv.FormatInt(ethTx.ID, 10)] = ethTx
	}

	for _, tx := range submittedTxs {
		if ethTx, ok := failedTxMap[tx.ID]; ok {
			tx.Status = types.Failure
			tx.FailureReason = &ethTx.Error.String
			err := sc.orm.UpdateLegacyGaslessTx(tx)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (sc *Sidecar) handleSourceFinalizedTxs(ctx context.Context, fromBlock, toBlock int64) error {
	txs, err := sc.orm.SelectByDestChainIDAndStatus(sc.ccipChainSelector, types.SourceFinalized, pg.WithParentCtx(ctx))
	if err != nil {
		return errors.Wrap(err, "find by status")
	}

	ccipMessageIDs := ccipMessageIDs(txs)

	logs, err := sc.lp.IndexedLogsByBlockRange(
		fromBlock,
		toBlock,
		evm_2_evm_offramp.EVM2EVMOffRampExecutionStateChanged{}.Topic(),
		sc.ccipOffRamp.Address(),
		2, //CCIP Message ID is the second indexed field
		ccipMessageIDs,
		pg.WithParentCtx(ctx))
	if err != nil {
		return errors.Wrap(err, "log poller indexed logs")
	}

	escLogs, err := sc.unmarshalExecutionStateChanged(logs)
	if err != nil {
		return errors.Wrap(err, "unmarshal execution state changed logs")
	}

	finalizedLogs, err := destinationFinalizedLogs(escLogs, txs)
	if err != nil {
		return errors.Wrap(err, "finalized logs")
	}

	for _, log := range finalizedLogs {
		err = sc.orm.UpdateLegacyGaslessTx(log)
		if err != nil {
			return err
		}
	}
	return nil
}

func (sc *Sidecar) unmarshalForwardSucceededLogs(logs []logpoller.Log) (unmarshalledLogs []*forwarder_wrapper.ForwarderForwardSucceeded, err error) {
	for _, log := range logs {
		rawLog := log.ToGethLog()
		if log.EventSig != (forwarder_wrapper.ForwarderForwardSucceeded{}).Topic() {
			err = errors.Errorf("unexpected event signature: %x", log.EventSig)
			return
		}
		unpacked, err2 := sc.forwarder.ParseLog(rawLog)
		if err2 != nil {
			// should never happen
			err = errors.Wrap(err2, "unmarshal ForwarderForwardSucceeded failed")
			return
		}
		fs, ok := unpacked.(*forwarder_wrapper.ForwarderForwardSucceeded)
		if !ok {
			// should never happen
			err = errors.New("cast to ForwarderForwardSucceeded")
			return
		}
		unmarshalledLogs = append(unmarshalledLogs, fs)
	}
	return
}

func (sc *Sidecar) unmarshalExecutionStateChanged(logs []logpoller.Log) (unmarshalledLogs []*evm_2_evm_offramp.EVM2EVMOffRampExecutionStateChanged, err error) {
	for _, log := range logs {
		rawLog := log.ToGethLog()
		if log.EventSig != (evm_2_evm_offramp.EVM2EVMOffRampExecutionStateChanged{}).Topic() {
			err = errors.Errorf("unexpected event signature: %x", log.EventSig)
			return
		}
		unpacked, err2 := sc.ccipOffRamp.ParseLog(rawLog)
		if err2 != nil {
			// should never happen
			err = errors.Wrap(err2, "unmarshal RampExecutionStateChanged failed")
			return
		}
		fs, ok := unpacked.(*evm_2_evm_offramp.EVM2EVMOffRampExecutionStateChanged)
		if !ok {
			// should never happen
			err = errors.New("cast to RampExecutionStateChanged")
			return
		}
		unmarshalledLogs = append(unmarshalledLogs, fs)
	}
	return
}

func txSenderAddresses(txs []types.LegacyGaslessTx) (txSenderAddresses []common.Hash) {
	senderAddressesMap := make(map[common.Address]struct{})
	for _, tx := range txs {
		senderAddressesMap[tx.From] = struct{}{}
	}
	for sender := range senderAddressesMap {
		txSenderAddresses = append(txSenderAddresses, sender.Hash())
	}
	return
}

func ccipMessageIDs(txs []types.LegacyGaslessTx) (ccipMessageIDs []common.Hash) {
	for _, tx := range txs {
		if tx.CCIPMessageID == nil {
			continue
		}
		ccipMessageIDs = append(ccipMessageIDs, *tx.CCIPMessageID)
	}
	return
}

func finalizedLogs(
	logs []*forwarder_wrapper.ForwarderForwardSucceeded, submittedTxs []types.LegacyGaslessTx) (finalizedLogs []types.LegacyGaslessTx, err error) {
	keyToLogs := make(map[string]*forwarder_wrapper.ForwarderForwardSucceeded)
	for _, log := range logs {
		gt := types.LegacyGaslessTx{
			Forwarder: log.Raw.Address,
			From:      log.From,
			Nonce:     utils.NewBig(log.Nonce),
		}
		key, err2 := gt.Key()
		if err2 != nil {
			// should not happen
			err = err2
			return
		}
		keyToLogs[*key] = log
	}

	for _, tx := range submittedTxs {
		key, err2 := tx.Key()
		if err2 != nil {
			// should not happen
			err = err2
			return
		}
		log, exists := keyToLogs[*key]
		if exists {
			if tx.SourceChainID == tx.DestinationChainID {
				tx.Status = types.Finalized
			} else {
				tx.Status = types.SourceFinalized
				ccipMessageID := common.Hash(log.ReturnValue)
				tx.CCIPMessageID = &ccipMessageID
			}
			finalizedLogs = append(finalizedLogs, tx)
		}
	}
	return
}

func destinationFinalizedLogs(
	logs []*evm_2_evm_offramp.EVM2EVMOffRampExecutionStateChanged, sourceFinalizedTxs []types.LegacyGaslessTx) (finalizedLogs []types.LegacyGaslessTx, err error) {
	ccipMessageIDToLogs := make(map[common.Hash]*evm_2_evm_offramp.EVM2EVMOffRampExecutionStateChanged)
	for _, log := range logs {
		ccipMessageIDToLogs[common.Hash(log.MessageId)] = log
	}

	for _, tx := range sourceFinalizedTxs {
		if tx.CCIPMessageID == nil {
			// should not happen
			err = errors.New("empty CCIP message ID")
			return
		}
		_, exists := ccipMessageIDToLogs[*tx.CCIPMessageID]
		if exists {
			tx.Status = types.Finalized
			finalizedLogs = append(finalizedLogs, tx)
		}
	}
	return
}
