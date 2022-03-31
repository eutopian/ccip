package ccip

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"sort"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/offramp"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/utils"
)

const ExecutionMaxInflightTimeSeconds = 180

var _ types.ReportingPluginFactory = &ExecutionReportingPluginFactory{}
var _ types.ReportingPlugin = &ExecutionReportingPlugin{}

// Message contains the data from a cross chain message
type Message struct {
	SequenceNumber *big.Int       `json:"sequenceNumber"`
	SourceChainId  *big.Int       `json:"sourceChainId"`
	Sender         common.Address `json:"sender"`
	Payload        struct {
		Tokens             []common.Address `json:"tokens"`
		Amounts            []*big.Int       `json:"amounts"`
		DestinationChainId *big.Int         `json:"destinationChainId"`
		Receiver           common.Address   `json:"receiver"`
		Executor           common.Address   `json:"executor"`
		Data               []uint8          `json:"data"`
		Options            []uint8          `json:"options"`
	} `json:"payload"`
}

type ExecutableMessage struct {
	Path    [][32]byte `json:"path"`
	Index   *big.Int   `json:"index"`
	Message Message    `json:"message"`
}

// ExecutionObservation Note there can be gaps in this range of sequence numbers,
// indicative of some messages being non-DON executed.
type ExecutionObservation struct {
	MinSeqNum utils.Big `json:"minSeqNum"`
	MaxSeqNum utils.Big `json:"maxSeqNum"`
}

func makeExecutionReportArgs() abi.Arguments {
	mustType := func(ts string, components []abi.ArgumentMarshaling) abi.Type {
		ty, _ := abi.NewType(ts, "", components)
		return ty
	}
	return []abi.Argument{
		{
			Name: "executableMessages",
			Type: mustType("tuple[]", []abi.ArgumentMarshaling{
				{
					Name: "Path",
					Type: "bytes32[]",
				},
				{
					Name: "Index",
					Type: "uint256",
				},
				{
					Name: "Message",
					Type: "tuple",
					Components: []abi.ArgumentMarshaling{
						{
							Name: "sequenceNumber",
							Type: "uint256",
						},
						{
							Name: "sourceChainId",
							Type: "uint256",
						},

						{
							Name: "sender",
							Type: "address",
						},
						{
							Name: "payload",
							Type: "tuple",
							Components: []abi.ArgumentMarshaling{
								{
									Name: "tokens",
									Type: "address[]",
								},
								{
									Name: "amounts",
									Type: "uint256[]",
								},
								{
									Name: "destinationChainId",
									Type: "uint256",
								},
								{
									Name: "receiver",
									Type: "address",
								},
								{
									Name: "executor",
									Type: "address",
								},
								{
									Name: "data",
									Type: "bytes",
								},
								{
									Name: "options",
									Type: "bytes",
								},
							},
						},
					},
				},
			}),
		},
	}
}

func EncodeExecutionReport(ems []ExecutableMessage) (types.Report, error) {
	report, err := makeExecutionReportArgs().PackValues([]interface{}{ems})
	if err != nil {
		return nil, err
	}
	return report, nil
}

func DecodeExecutionReport(report types.Report) ([]ExecutableMessage, error) {
	unpacked, err := makeExecutionReportArgs().Unpack(report)
	if err != nil {
		return nil, err
	}
	if len(unpacked) == 0 {
		return nil, errors.New("assumptionViolation: expected at least one element")
	}

	// Must be anonymous struct here
	msgs, ok := unpacked[0].([]struct {
		Path    [][32]uint8 `json:"Path"`
		Index   *big.Int    `json:"Index"`
		Message struct {
			SequenceNumber *big.Int       `json:"sequenceNumber"`
			SourceChainId  *big.Int       `json:"sourceChainId"`
			Sender         common.Address `json:"sender"`
			Payload        struct {
				Tokens             []common.Address `json:"tokens"`
				Amounts            []*big.Int       `json:"amounts"`
				DestinationChainId *big.Int         `json:"destinationChainId"`
				Receiver           common.Address   `json:"receiver"`
				Executor           common.Address   `json:"executor"`
				Data               []uint8          `json:"data"`
				Options            []uint8          `json:"options"`
			} `json:"payload"`
		} `json:"Message"`
	})
	if !ok {
		return nil, fmt.Errorf("got %T", unpacked[0])
	}
	if len(msgs) == 0 {
		return nil, errors.New("assumptionViolation: expected at least one element")
	}
	var ems []ExecutableMessage
	for _, emi := range msgs {
		ems = append(ems, ExecutableMessage{
			Path:    emi.Path,
			Index:   emi.Index,
			Message: emi.Message,
		})
	}
	return ems, nil
}

//go:generate mockery --name OffRampLastReporter --output ./mocks/lastreporter --case=underscore
type OffRampLastReporter interface {
	GetLastReport(opts *bind.CallOpts) (offramp.CCIPRelayReport, error)
}

type ExecutionReportingPluginFactory struct {
	l            logger.Logger
	orm          ORM
	source, dest *big.Int
	lastReporter OffRampLastReporter
	executor     common.Address
	onRamp       common.Address
	offRamp      common.Address
}

func NewExecutionReportingPluginFactory(l logger.Logger, orm ORM, source, dest *big.Int, onRamp, offRamp common.Address, executor common.Address, lastReporter OffRampLastReporter) types.ReportingPluginFactory {
	return &ExecutionReportingPluginFactory{l: l, orm: orm, source: source, dest: dest, onRamp: onRamp, offRamp: offRamp, executor: executor, lastReporter: lastReporter}
}

func (rf *ExecutionReportingPluginFactory) NewReportingPlugin(config types.ReportingPluginConfig) (types.ReportingPlugin, types.ReportingPluginInfo, error) {
	return ExecutionReportingPlugin{rf.l, config.F, rf.orm, rf.source, rf.dest, rf.executor, rf.onRamp, rf.offRamp, rf.lastReporter}, types.ReportingPluginInfo{
		Name:          "CCIPExecution",
		UniqueReports: true,
		MaxQueryLen:   0, // We do not use the query phase.
		// TODO: https://app.shortcut.com/chainlinklabs/story/30171/define-report-plugin-limits
		MaxObservationLen: 100000, // TODO
		MaxReportLen:      100000, // TODO
	}, nil
}

type ExecutionReportingPlugin struct {
	l             logger.Logger
	F             int
	orm           ORM
	sourceChainId *big.Int
	destChainId   *big.Int
	executor      common.Address
	onRamp        common.Address
	offRamp       common.Address
	// We also use the offramp for defensive checks
	lastReporter OffRampLastReporter
}

func (r ExecutionReportingPlugin) Query(ctx context.Context, timestamp types.ReportTimestamp) (types.Query, error) {
	// We don't use a query for this reporting plugin, so we can just leave it empty here
	return types.Query{}, nil
}

func (r ExecutionReportingPlugin) Observation(ctx context.Context, timestamp types.ReportTimestamp, query types.Query) (types.Observation, error) {
	// We want to execute any messages which satisfy the following:
	// 1. Have the executor field set to the DONs message executor contract
	// 2. There exists a confirmed relay report containing its sequence number, i.e. it's status is RequestStatusRelayConfirmed
	relayedReqs, err := r.orm.Requests(r.sourceChainId, r.destChainId, r.onRamp, r.offRamp, nil, nil, RequestStatusRelayConfirmed, &r.executor, nil)
	if err != nil {
		return nil, err
	}
	// No request to process. Return an observation with MinSeqNum and MaxSeqNum equal to NoRequestsToProcess
	// which should not result in a new report being generated during the Report step.
	if len(relayedReqs) == 0 {
		b, jsonErr := json.Marshal(&ExecutionObservation{
			MinSeqNum: utils.Big(*NoRequestsToProcess),
			MaxSeqNum: utils.Big(*NoRequestsToProcess),
		})
		if jsonErr != nil {
			return nil, jsonErr
		}
		return b, nil
	}
	// Double-check the latest sequence number onchain is >= our max relayed seq num
	lr, err := r.lastReporter.GetLastReport(nil)
	if err != nil {
		return nil, err
	}
	if relayedReqs[len(relayedReqs)-1].SeqNum.ToInt().Cmp(lr.MaxSequenceNumber) > 0 {
		return nil, errors.Errorf("invariant violated, mismatch between relay_confirmed requests and last report")
	}
	b, err := json.Marshal(&ExecutionObservation{
		MinSeqNum: relayedReqs[0].SeqNum,
		MaxSeqNum: relayedReqs[len(relayedReqs)-1].SeqNum,
	})
	if err != nil {
		return nil, err
	}
	return b, nil
}

func (r ExecutionReportingPlugin) Report(ctx context.Context, timestamp types.ReportTimestamp, query types.Query, observations []types.AttributedObservation) (bool, types.Report, error) {
	var nonEmptyObservations []ExecutionObservation
	for _, ao := range observations {
		var ob ExecutionObservation
		err := json.Unmarshal(ao.Observation, &ob)
		if err != nil {
			r.l.Errorw("Received unmarshallable observation", "err", err, "observation", string(ao.Observation))
			continue
		}
		minSeqNum := ob.MinSeqNum.ToInt()
		if minSeqNum.Sign() < 0 {
			if minSeqNum.Cmp(NoRequestsToProcess) == 0 {
				r.l.Tracew("Discarded empty observation %+v", ao)
			} else {
				r.l.Warnf("Discarded invalid observation %+v", ao)
			}
			continue
		}
		nonEmptyObservations = append(nonEmptyObservations, ob)
	}
	// Need at least F+1 observations
	if len(nonEmptyObservations) <= r.F {
		r.l.Tracew("Non-empty observations <= F, need at least F+1 to continue")
		return false, nil, nil
	}
	// We have at least F+1 valid observations
	// Extract the min and max
	sort.Slice(nonEmptyObservations, func(i, j int) bool {
		return nonEmptyObservations[i].MinSeqNum.ToInt().Cmp(nonEmptyObservations[j].MinSeqNum.ToInt()) < 0
	})
	// r.F < len(nonEmptyObservations) because of the check above and therefore this is safe
	min := *nonEmptyObservations[r.F].MinSeqNum.ToInt()
	sort.Slice(nonEmptyObservations, func(i, j int) bool {
		return nonEmptyObservations[i].MaxSeqNum.ToInt().Cmp(nonEmptyObservations[j].MaxSeqNum.ToInt()) < 0
	})
	// We use a conservative maximum. If we pick a value that some honest oracles might not
	// have seen they’ll end up not agreeing on a report, stalling the protocol.
	max := *nonEmptyObservations[r.F].MaxSeqNum.ToInt()
	if max.Cmp(&min) < 0 {
		return false, nil, errors.New("max seq num smaller than min")
	}
	reqs, err := r.orm.Requests(r.sourceChainId, r.destChainId, r.onRamp, r.offRamp, &min, &max, RequestStatusRelayConfirmed, &r.executor, nil)
	if err != nil {
		return false, nil, err
	}
	// Cannot construct a report for which we haven't seen all the messages.
	if len(reqs) == 0 {
		return false, nil, errors.Errorf("do not have all the messages in report, have zero messages, report has min %v max %v", min, max)
	}
	lr, err := r.lastReporter.GetLastReport(nil)
	if err != nil {
		return false, nil, err
	}
	if reqs[len(reqs)-1].SeqNum.ToInt().Cmp(lr.MaxSequenceNumber) > 0 {
		return false, nil, errors.Errorf("invariant violated, mismatch between relay_confirmed requests (max %v) and last report (max %v)", reqs[len(reqs)-1].SeqNum, lr.MaxSequenceNumber)
	}
	report, err := r.buildReport(reqs)
	if err != nil {
		return false, nil, err
	}
	return true, report, nil
}

// For each message in the given range of sequence numbers (with potential holes):
// 1. Lookup the report associated with that sequence number
// 2. Generate a merkle proof that the message was in that report
// 3. Encode those proofs and messages into a report for the executor contract
// TODO: We may want to combine these queries for performance, hold off
// until we decide whether we move forward with batch proving.
func (r ExecutionReportingPlugin) buildReport(reqs []*Request) ([]byte, error) {
	var executable []ExecutableMessage
	for _, req := range reqs {
		// Look up all the messages that are in the same report
		// as this one (even externally executed ones), generate a Proof and double-check the root checks out.
		rep, err2 := r.orm.RelayReport(req.SeqNum.ToInt())
		if err2 != nil {
			r.l.Errorw("Could not find relay report for request", "err", err2, "seq num", req.SeqNum.String())
			continue
		}
		allReqsInReport, err3 := r.orm.Requests(r.sourceChainId, r.destChainId, req.OnRamp, req.OffRamp, rep.MinSeqNum.ToInt(), rep.MaxSeqNum.ToInt(), "", nil, nil)
		if err3 != nil {
			continue
		}
		var leaves [][]byte
		for _, reqInReport := range allReqsInReport {
			leaves = append(leaves, reqInReport.Raw)
		}
		index := big.NewInt(0).Sub(req.SeqNum.ToInt(), rep.MinSeqNum.ToInt())
		root, proof := GenerateMerkleProof(32, leaves, int(index.Int64()))
		if !bytes.Equal(root[:], rep.Root[:]) {
			continue
		}
		executable = append(executable, ExecutableMessage{
			Path:    proof.PathForExecute(),
			Message: req.ToMessage(),
			Index:   proof.Index(),
		})
	}

	report, err := EncodeExecutionReport(executable)
	if err != nil {
		return nil, err
	}

	return report, nil
}

func (r ExecutionReportingPlugin) ShouldAcceptFinalizedReport(ctx context.Context, timestamp types.ReportTimestamp, report types.Report) (bool, error) {
	ems, err := DecodeExecutionReport(report)
	if err != nil {
		return false, nil
	}

	var seqNums []*big.Int
	for i := range ems {
		seqNums = append(seqNums, ems[i].Message.SequenceNumber)
	}

	// If the first message is executed already, this execution report is stale, and we do not accept it.
	stale, err := r.isStale(seqNums[0])
	if err != nil {
		return !stale, err
	}
	if stale {
		return false, err
	}
	// Any timed out requests should be set back to RequestStatusExecutionPending so their execution can be retried in a subsequent report.
	if err = r.orm.ResetExpiredRequests(r.sourceChainId, r.destChainId, r.onRamp, r.offRamp, ExecutionMaxInflightTimeSeconds, RequestStatusExecutionPending, RequestStatusRelayConfirmed); err != nil {
		// Ok to continue here, we'll try to reset them again on the next round.
		r.l.Errorw("Unable to reset expired requests", "err", err)
	}

	if err := r.orm.UpdateRequestSetStatus(r.sourceChainId, r.destChainId, r.onRamp, r.offRamp, seqNums, RequestStatusExecutionPending); err != nil {
		return false, err
	}
	return true, nil
}

func (r ExecutionReportingPlugin) ShouldTransmitAcceptedReport(ctx context.Context, timestamp types.ReportTimestamp, report types.Report) (bool, error) {
	parsedReport, err := DecodeExecutionReport(report)
	if err != nil {
		return false, nil
	}
	// If report is not stale we transmit.
	// When the executeTransmitter enqueues the tx for bptxm,
	// we mark it as execution_sent, removing it from the set of inflight messages.
	stale, err := r.isStale(parsedReport[0].Message.SequenceNumber)
	return !stale, err
}

func (r ExecutionReportingPlugin) isStale(min *big.Int) (bool, error) {
	// If the first message is executed already, this execution report is stale.
	req, err := r.orm.Requests(r.sourceChainId, r.destChainId, r.onRamp, r.offRamp, min, min, "", nil, nil)
	if err != nil {
		// if we can't find the request, assume transient db issue
		// and wait until the next OCR2 round (don't submit)
		return true, err
	}
	if len(req) != 1 {
		// If we don't have the request at all, this likely means we never had the request to begin with
		// (say our eth subscription is down) and we want to let other oracles continue the protocol.
		return false, errors.New("could not find first message in execution report")
	}
	return req[0].Status == RequestStatusExecutionConfirmed, nil
}

func (r ExecutionReportingPlugin) Start() error {
	return nil
}

func (r ExecutionReportingPlugin) Close() error {
	return nil
}
