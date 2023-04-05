package load

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog"

	"github.com/smartcontractkit/chainlink-testing-framework/loadgen"
	"github.com/stretchr/testify/require"
	"go.uber.org/atomic"

	"github.com/smartcontractkit/chainlink/integration-tests/testreporters"
	"github.com/smartcontractkit/chainlink/integration-tests/testsetups"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/evm_2_evm_onramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/testhelpers"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type CCIPE2ELoad struct {
	t                     *testing.T
	Source                *actions.SourceCCIPModule // all source contracts
	Destination           *actions.DestCCIPModule   // all destination contracts
	NoOfReq               int64                     // no of Request fired - required for balance assertion at the end
	totalGEFee            *big.Int
	BalanceStats          BalanceStats  // balance assertion details
	CurrentMsgSerialNo    *atomic.Int64 // current msg serial number in the load sequence
	InitialSourceBlockNum uint64
	InitialDestBlockNum   uint64        // blocknumber before the first message is fired in the load sequence
	CallTimeOut           time.Duration // max time to wait for various on-chain events
	TickerDuration        time.Duration // poll frequency while waiting for on-chain events
	callStatsMu           *sync.Mutex
	reports               *testreporters.CCIPLaneStats
	seqNumCommittedMu     *sync.Mutex
	seqNumCommitted       map[uint64]uint64 // key : seqNumber in the ReportAccepted event, value : blocknumber for corresponding event
	msg                   router.ClientEVM2AnyMessage
	msgMu                 *sync.Mutex
}
type BalanceStats struct {
	SourceBalanceReq        map[string]*big.Int
	SourceBalanceAssertions []testhelpers.BalanceAssertion
	DestBalanceReq          map[string]*big.Int
	DestBalanceAssertions   []testhelpers.BalanceAssertion
}

func NewCCIPLoad(t *testing.T, source *actions.SourceCCIPModule, dest *actions.DestCCIPModule, timeout time.Duration, noOfReq int64, reporter *testreporters.CCIPLaneStats) *CCIPE2ELoad {
	return &CCIPE2ELoad{
		t:                  t,
		Source:             source,
		Destination:        dest,
		CurrentMsgSerialNo: atomic.NewInt64(1),
		TickerDuration:     time.Second,
		CallTimeOut:        timeout,
		NoOfReq:            noOfReq,
		reports:            reporter,
		callStatsMu:        &sync.Mutex{},
		seqNumCommittedMu:  &sync.Mutex{},
		seqNumCommitted:    make(map[uint64]uint64),
		msgMu:              &sync.Mutex{},
	}
}

// BeforeAllCall funds subscription, approves the token transfer amount.
// Needs to be called before load sequence is started.
// Needs to approve and fund for the entire sequence.
func (c *CCIPE2ELoad) BeforeAllCall(msgType string) {
	sourceCCIP := c.Source
	destCCIP := c.Destination
	var tokenAndAmounts []router.ClientEVMTokenAmount
	for i, token := range sourceCCIP.Common.BridgeTokens {
		tokenAndAmounts = append(tokenAndAmounts, router.ClientEVMTokenAmount{
			Token: common.HexToAddress(token.Address()), Amount: c.Source.TransferAmount[i],
		})
	}

	err := sourceCCIP.Common.ChainClient.WaitForEvents()
	require.NoError(c.t, err, "Failed to wait for events")

	// save the current block numbers to use in various filter log requests
	currentBlockOnSource, err := sourceCCIP.Common.ChainClient.LatestBlockNumber(context.Background())
	require.NoError(c.t, err, "failed to fetch latest source block num")
	currentBlockOnDest, err := destCCIP.Common.ChainClient.LatestBlockNumber(context.Background())
	require.NoError(c.t, err, "failed to fetch latest dest block num")
	c.InitialDestBlockNum = currentBlockOnDest
	c.InitialSourceBlockNum = currentBlockOnSource
	// collect the balance requirement to verify balances after transfer
	sourceBalances, err := testhelpers.GetBalances(sourceCCIP.CollectBalanceRequirements(c.t))
	require.NoError(c.t, err, "fetching source balance")
	destBalances, err := testhelpers.GetBalances(destCCIP.CollectBalanceRequirements(c.t))
	require.NoError(c.t, err, "fetching dest balance")
	c.BalanceStats = BalanceStats{
		SourceBalanceReq: sourceBalances,
		DestBalanceReq:   destBalances,
	}
	extraArgsV1, err := testhelpers.GetEVMExtraArgsV1(big.NewInt(100_000), false)
	require.NoError(c.t, err, "Failed encoding the options field")

	receiver, err := utils.ABIEncode(`[{"type":"address"}]`, destCCIP.ReceiverDapp.EthAddress)
	require.NoError(c.t, err, "Failed encoding the receiver address")
	c.msg = router.ClientEVM2AnyMessage{
		Receiver:  receiver,
		ExtraArgs: extraArgsV1,
		FeeToken:  common.HexToAddress(sourceCCIP.Common.FeeToken.Address()),
		Data:      []byte("message with Id 1"),
	}
	if msgType == testsetups.TokenTransfer {
		c.msg.TokenAmounts = tokenAndAmounts
	}

	sourceCCIP.Common.ChainClient.ParallelTransactions(false)
	destCCIP.Common.ChainClient.ParallelTransactions(false)
}

func (c *CCIPE2ELoad) AfterAllCall() {
	c.BalanceStats.DestBalanceAssertions = c.Destination.BalanceAssertions(
		c.t,
		c.BalanceStats.DestBalanceReq,
		c.Source.TransferAmount,
		c.NoOfReq,
	)
	c.BalanceStats.SourceBalanceAssertions = c.Source.BalanceAssertions(c.t, c.BalanceStats.SourceBalanceReq, c.NoOfReq, c.totalGEFee)
	actions.AssertBalances(c.t, c.BalanceStats.DestBalanceAssertions)
	actions.AssertBalances(c.t, c.BalanceStats.SourceBalanceAssertions)
}

func (c *CCIPE2ELoad) Call(_ *loadgen.Generator) loadgen.CallResult {
	var res loadgen.CallResult
	sourceCCIP := c.Source
	msgSerialNo := c.CurrentMsgSerialNo.Load()
	c.CurrentMsgSerialNo.Inc()

	lggr := zerolog.New(zerolog.NewConsoleWriter(zerolog.ConsoleTestWriter(c.t))).
		With().Timestamp().Logger().
		With().Int("msg Number", int(msgSerialNo)).
		Str("Lane",
			fmt.Sprintf("%d-->%d", c.Source.Common.ChainClient.GetChainID().Int64(),
				c.Destination.Common.ChainClient.GetChainID().Int64())).
		Logger()
	// form the message for transfer
	msgStr := fmt.Sprintf("message with Id %d", msgSerialNo)
	c.msgMu.Lock()
	msg := c.msg
	c.msgMu.Unlock()
	msg.Data = []byte(msgStr)

	feeToken := sourceCCIP.Common.FeeToken.EthAddress
	// initiate the transfer
	lggr.Debug().Str("triggeredAt", time.Now().GoString()).Msg("triggering transfer")
	var sendTx *types.Transaction
	var err error

	// initiate the transfer
	// if the token address is 0x0 it will use Native as fee token and the fee amount should be mentioned in bind.TransactOpts's value
	startTime := time.Now()
	if feeToken != common.HexToAddress("0x0") {
		sendTx, err = sourceCCIP.Common.Router.CCIPSend(sourceCCIP.DestinationChainId, msg, nil)
	} else {
		fee, err := sourceCCIP.Common.Router.GetFee(sourceCCIP.DestinationChainId, msg)
		if err != nil {
			res.Error = err.Error()
			return res
		}
		sendTx, err = sourceCCIP.Common.Router.CCIPSend(sourceCCIP.DestinationChainId, msg, fee)
	}

	if err != nil {
		c.reports.UpdatePhaseStats(msgSerialNo, 0, testreporters.TX, time.Since(startTime), testreporters.Failure)
		lggr.Err(err).Msg("ccip-send tx error for msg ID")
		res.Error = fmt.Sprintf("ccip-send tx error %+v for msg ID %d", err, msgSerialNo)
		res.Data = c.reports.GetPhaseStatsForRequest(msgSerialNo)
		return res
	}
	c.reports.UpdatePhaseStats(msgSerialNo, 0, testreporters.TX, time.Since(startTime), testreporters.Success)

	// wait for
	// - CCIPSendRequested Event log to be generated,
	ticker := time.NewTicker(c.TickerDuration)
	defer ticker.Stop()
	sentMsg, err := c.waitForCCIPSendRequested(lggr, ticker, c.InitialSourceBlockNum, msgSerialNo, sendTx.Hash().Hex(), time.Now())
	if err != nil {
		lggr.Err(err).Msg("CCIPSendRequested event error")
		res.Error = err.Error()
		res.Data = c.reports.GetPhaseStatsForRequest(msgSerialNo)
		return res
	}
	commitStartTime := time.Now()
	seqNum := sentMsg.SequenceNumber
	messageID := sentMsg.MessageId
	if bytes.Compare(sentMsg.Data, []byte(msgStr)) != 0 {
		res.Error = fmt.Sprintf("the message byte didnot match expected %s received %s msg ID %d", msgStr, string(sentMsg.Data), msgSerialNo)
		res.Data = c.reports.GetPhaseStatsForRequest(msgSerialNo)
		return res
	}

	// wait for
	// - CommitStore to increase the seq number,
	err = c.waitForSeqNumberIncrease(lggr, ticker, seqNum, msgSerialNo, commitStartTime)
	if err != nil {
		lggr.Err(err).Msgf("waiting for seq num increase for msg ID %d", msgSerialNo)
		res.Error = err.Error()
		res.Data = c.reports.GetPhaseStatsForRequest(msgSerialNo)
		return res
	}
	// wait for ReportAccepted event
	err = c.waitForReportAccepted(lggr, ticker, msgSerialNo, seqNum, c.InitialDestBlockNum, commitStartTime)
	if err != nil {
		lggr.Err(err).Msgf("waiting for ReportAcceptedEvent for msg ID %d", msgSerialNo)
		res.Error = err.Error()
		res.Data = c.reports.GetPhaseStatsForRequest(msgSerialNo)
		return res
	}

	// wait for ExecutionStateChanged event
	c.seqNumCommittedMu.Lock()
	currentBlockOnDest := c.seqNumCommitted[seqNum] - 2
	c.seqNumCommittedMu.Unlock()
	err = c.waitForExecStateChange(lggr, ticker, []uint64{seqNum}, [][32]byte{messageID}, currentBlockOnDest, msgSerialNo, time.Now())
	if err != nil {
		lggr.Err(err).Msgf("waiting for ExecutionStateChangedEvent for msg ID %d", msgSerialNo)
		res.Error = err.Error()
		res.Data = c.reports.GetPhaseStatsForRequest(msgSerialNo)
		return res
	}

	res.Data = c.reports.GetPhaseStatsForRequest(msgSerialNo)
	return res
}

func (c *CCIPE2ELoad) updateSeqNumCommitted(seqNum []uint64, blockNum uint64) {
	c.seqNumCommittedMu.Lock()
	defer c.seqNumCommittedMu.Unlock()
	for _, num := range seqNum {
		if _, ok := c.seqNumCommitted[num]; ok {
			return
		}
		c.seqNumCommitted[num] = blockNum
	}
}

func (c *CCIPE2ELoad) waitForExecStateChange(log zerolog.Logger, ticker *time.Ticker, seqNums []uint64, messageID [][32]byte, currentBlockOnDest uint64, msgSerialNo int64, timeNow time.Time) error {
	log.Info().Msgf("waiting for ExecutionStateChanged for seqNums %v", seqNums)
	ctx, cancel := context.WithTimeout(context.Background(), c.CallTimeOut)
	defer cancel()
	for {
		select {
		case <-ticker.C:
			iterator, err := c.Destination.OffRamp.FilterExecutionStateChanged(seqNums, messageID, currentBlockOnDest)
			if err != nil {
				for _, seqNum := range seqNums {
					c.reports.UpdatePhaseStats(msgSerialNo, seqNum, testreporters.ExecStateChanged, time.Since(timeNow), testreporters.Failure)
				}
				return fmt.Errorf("filtering event ExecutionStateChanged returned error %+v msg ID %d and seqNum %v", err, msgSerialNo, seqNums)
			}
			for iterator.Next() {
				switch ccip.MessageExecutionState(iterator.Event.State) {
				case ccip.Success:
					for _, seqNum := range seqNums {
						c.reports.UpdatePhaseStats(msgSerialNo, seqNum, testreporters.ExecStateChanged, time.Since(timeNow), testreporters.Success)
					}
					return nil
				case ccip.Failure:
					for _, seqNum := range seqNums {
						c.reports.UpdatePhaseStats(msgSerialNo, seqNum, testreporters.ExecStateChanged, time.Since(timeNow), testreporters.Failure)
					}
					return fmt.Errorf("ExecutionStateChanged event returned failure for seq num %v msg ID %d", seqNums, msgSerialNo)
				}
			}
		case <-ctx.Done():
			for _, seqNum := range seqNums {
				c.reports.UpdatePhaseStats(msgSerialNo, seqNum, testreporters.ExecStateChanged, time.Since(timeNow), testreporters.Failure)
			}
			return fmt.Errorf("ExecutionStateChanged event not found for seq num %v msg ID %d", seqNums, msgSerialNo)
		}
	}
}

func (c *CCIPE2ELoad) waitForSeqNumberIncrease(log zerolog.Logger, ticker *time.Ticker, seqNum uint64, msgSerialNo int64, timeNow time.Time) error {
	log.Info().Msgf("waiting for seq number %d to get increased", int(seqNum))
	ctx, cancel := context.WithTimeout(context.Background(), c.CallTimeOut)
	defer cancel()
	for {
		select {
		case <-ticker.C:
			seqNumberAfter, err := c.Destination.CommitStore.GetNextSeqNumber()
			if err != nil {
				c.reports.UpdatePhaseStats(msgSerialNo, seqNum, testreporters.Commit, time.Since(timeNow), testreporters.Failure)
				return fmt.Errorf("error %+v in GetNextExpectedSeqNumber by commitStore for msg ID %d", err, msgSerialNo)
			}
			if seqNumberAfter > seqNum {
				return nil
			}
		case <-ctx.Done():
			c.reports.UpdatePhaseStats(msgSerialNo, seqNum, testreporters.Commit, time.Since(timeNow), testreporters.Failure)
			return fmt.Errorf("sequence number is not increased for seq num %d msg ID %d", seqNum, msgSerialNo)
		}
	}
}

func (c *CCIPE2ELoad) waitForReportAccepted(log zerolog.Logger, ticker *time.Ticker, msgSerialNo int64, seqNum uint64, currentBlockOnDest uint64, timeNow time.Time) error {
	log.Info().Int("seq Number", int(seqNum)).Msg("waiting for ReportAccepted")
	ctx, cancel := context.WithTimeout(context.Background(), c.CallTimeOut)
	defer cancel()
	for {
		select {
		case <-ticker.C:
			c.seqNumCommittedMu.Lock()
			_, ok := c.seqNumCommitted[seqNum]
			c.seqNumCommittedMu.Unlock()
			// skip calling FilterReportAccepted if the seqNum is present in the map
			if ok {
				c.reports.UpdatePhaseStats(msgSerialNo, seqNum, testreporters.Commit, time.Since(timeNow), testreporters.Success)
				return nil
			}
			it, err := c.Destination.CommitStore.FilterReportAccepted(currentBlockOnDest)
			if err != nil {
				c.reports.UpdatePhaseStats(msgSerialNo, seqNum, testreporters.Commit, time.Since(timeNow), testreporters.Failure)
				return fmt.Errorf("error %+v in filtering by ReportAccepted event for seq num %d", err, seqNum)
			}
			for it.Next() {
				in := it.Event.Report.Interval
				seqNums := make([]uint64, in.Max-in.Min+1)
				var i uint64
				for range seqNums {
					seqNums[i] = in.Min + i
					i++
				}
				// update SeqNumCommitted map for all seqNums in the emitted ReportAccepted event
				c.updateSeqNumCommitted(seqNums, it.Event.Raw.BlockNumber)
				if in.Max >= seqNum && in.Min <= seqNum {
					c.reports.UpdatePhaseStats(msgSerialNo, seqNum, testreporters.Commit, time.Since(timeNow), testreporters.Success)
					return nil
				}
			}
		case <-ctx.Done():
			c.reports.UpdatePhaseStats(msgSerialNo, seqNum, testreporters.Commit, time.Since(timeNow), testreporters.Failure)
			return fmt.Errorf("ReportAccepted is not found for seq num %d", seqNum)
		}
	}
}

func (c *CCIPE2ELoad) waitForCCIPSendRequested(
	log zerolog.Logger,
	ticker *time.Ticker,
	currentBlockOnSource uint64,
	msgSerialNo int64,
	txHash string,
	timeNow time.Time,
) (evm_2_evm_onramp.InternalEVM2EVMMessage, error) {
	var sentmsg evm_2_evm_onramp.InternalEVM2EVMMessage
	log.Info().Msg("waiting for CCIPSendRequested")
	ctx, cancel := context.WithTimeout(context.Background(), c.CallTimeOut)
	defer cancel()
	for {
		select {
		case <-ticker.C:
			iterator, err := c.Source.OnRamp.FilterCCIPSendRequested(currentBlockOnSource)
			if err != nil {
				c.reports.UpdatePhaseStats(msgSerialNo, 0, testreporters.CCIPSendRe, time.Since(timeNow), testreporters.Failure)
				return sentmsg, fmt.Errorf("error %+v in filtering CCIPSendRequested event for msg ID %d tx %s", err, msgSerialNo, txHash)
			}
			for iterator.Next() {
				if iterator.Event.Raw.TxHash.Hex() == txHash {
					sentmsg = iterator.Event.Message
					c.reports.UpdatePhaseStats(msgSerialNo, sentmsg.SequenceNumber, testreporters.CCIPSendRe, time.Since(timeNow), testreporters.Success)
					return sentmsg, nil
				}
			}
		case <-ctx.Done():
			c.reports.UpdatePhaseStats(msgSerialNo, 0, testreporters.CCIPSendRe, time.Since(timeNow), testreporters.Failure)
			latest, _ := c.Source.Common.ChainClient.LatestBlockNumber(context.Background())
			return sentmsg, fmt.Errorf("CCIPSendRequested event is not found for msg ID %d tx %s startblock %d latestblock %d",
				msgSerialNo, txHash, currentBlockOnSource, latest)
		}
	}
}

func (c *CCIPE2ELoad) ReportAcceptedLog(log zerolog.Logger) {
	log.Info().Msg("Commit Report stats")
	it, err := c.Destination.CommitStore.FilterReportAccepted(c.InitialDestBlockNum)
	require.NoError(c.t, err, "report committed result")
	i := 1
	event := log.Info()
	for it.Next() {
		event.Interface(fmt.Sprintf("%d Report Intervals", i), it.Event.Report.Interval)
		i++
	}
	event.Msgf("CommitStore-Reports Accepted")
}
