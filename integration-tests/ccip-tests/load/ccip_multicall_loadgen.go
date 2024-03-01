package load

import (
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/prometheus/common/model"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"

	chain_selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/wasp"

	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"
	"github.com/smartcontractkit/chainlink-testing-framework/logging"

	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/contracts"
	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/testreporters"
	"github.com/smartcontractkit/chainlink/integration-tests/ccip-tests/testsetups"
)

// CCIPMultiCallLoadGenerator represents a load generator for the CCIP lanes originating from same network
// The purpose of this load generator is to group ccip-send calls for the CCIP lanes originating from same network
// This is to avoid the scenario of hitting rpc rate limit for the same network if the load generator is sending
// too many ccip-send calls to the same network hitting the rpc rate limit
type CCIPMultiCallLoadGenerator struct {
	t                       *testing.T
	logger                  zerolog.Logger
	client                  blockchain.EVMClient
	E2ELoads                map[string]*CCIPE2ELoad
	MultiCall               string
	ChunkSize               *int
	NoOfRequestsPerUnitTime int64
	labels                  model.LabelSet
	loki                    *wasp.LokiClient
	responses               chan map[string]MultiCallReturnValues
	Done                    chan struct{}
}

type MultiCallReturnValues struct {
	Msgs  map[int][]contracts.CCIPMsgData
	Stats map[int][]*testreporters.RequestStat
}

func NewMultiCallLoadGenerator(testCfg *testsetups.CCIPTestConfig, lanes []*actions.CCIPLane, noOfRequestsPerUnitTime int64, labels map[string]string) (*CCIPMultiCallLoadGenerator, error) {
	// check if all lanes are from same network
	source := lanes[0].SourceChain.GetChainID()
	multiCall := lanes[0].SrcNetworkLaneCfg.Multicall
	if multiCall == "" {
		return nil, fmt.Errorf("multicall address cannot be empty")
	}
	for i := 1; i < len(lanes); i++ {
		if source.String() != lanes[i].SourceChain.GetChainID().String() {
			return nil, fmt.Errorf("all lanes should be from same network; expected %s, got %s", source, lanes[i].SourceChain.GetChainID())
		}
	}
	client := lanes[0].SourceChain
	lggr := logging.GetTestLogger(testCfg.Test).With().Str("Source Network", client.GetNetworkName()).Logger()
	ls := wasp.LabelsMapToModel(labels)
	if err := ls.Validate(); err != nil {
		return nil, err
	}
	lokiConfig := testCfg.EnvInput.Logging.Loki
	loki, err := wasp.NewLokiClient(wasp.NewLokiConfig(lokiConfig.Endpoint, lokiConfig.TenantId, lokiConfig.BasicAuth, lokiConfig.BearerToken))
	if err != nil {
		return nil, err
	}
	m := &CCIPMultiCallLoadGenerator{
		t:                       testCfg.Test,
		client:                  client,
		MultiCall:               multiCall,
		logger:                  lggr,
		NoOfRequestsPerUnitTime: noOfRequestsPerUnitTime,
		ChunkSize:               testCfg.TestGroupInput.ChunkSizeInMulticall,
		E2ELoads:                make(map[string]*CCIPE2ELoad),
		labels:                  ls,
		loki:                    loki,
		responses:               make(chan map[string]MultiCallReturnValues),
		Done:                    make(chan struct{}),
	}
	for _, lane := range lanes {
		ccipLoad := NewCCIPLoad(testCfg.Test, lane, testCfg.TestGroupInput.PhaseTimeout.Duration(), 100000)
		// for multicall load generator, we don't want to send max data intermittently, it might
		// cause oversized data for multicall
		ccipLoad.SendMaxDataIntermittently = false
		ccipLoad.BeforeAllCall(testCfg.TestGroupInput.MsgType, big.NewInt(*testCfg.TestGroupInput.DestGasLimit))
		m.E2ELoads[fmt.Sprintf("%s-%s", lane.SourceNetworkName, lane.DestNetworkName)] = ccipLoad
	}

	m.StartLokiStream()
	return m, nil
}

func (m *CCIPMultiCallLoadGenerator) Stop() error {
	m.Done <- struct{}{}
	tokenMap := make(map[string]struct{})
	var tokens []*contracts.ERC20Token
	for _, e2eLoad := range m.E2ELoads {
		for i := range e2eLoad.Lane.Source.TransferAmount {
			if _, ok := tokenMap[e2eLoad.Lane.Source.Common.BridgeTokens[i].Address()]; !ok {
				tokens = append(tokens, e2eLoad.Lane.Source.Common.BridgeTokens[i])
			}
		}
	}
	if len(tokens) > 0 {
		return contracts.TransferTokens(m.client, common.HexToAddress(m.MultiCall), tokens)
	}
	return nil
}

func (m *CCIPMultiCallLoadGenerator) StartLokiStream() {
	go func() {
		for {
			select {
			case <-m.Done:
				m.logger.Info().Msg("stopping loki client from multi call load generator")
				m.loki.Stop()
				return
			case rValues := <-m.responses:
				m.HandleLokiLogs(rValues)
			}
		}
	}()
}

func (m *CCIPMultiCallLoadGenerator) HandleLokiLogs(rValues map[string]MultiCallReturnValues) {
	for dest, rValue := range rValues {
		labels := m.labels.Merge(model.LabelSet{
			"dest_chain":     model.LabelValue(dest),
			"test_data_type": "responses",
			"go_test_name":   model.LabelValue(m.t.Name()),
		})
		for _, allstats := range rValue.Stats {
			for _, stat := range allstats {
				err := m.loki.HandleStruct(labels, time.Now().UTC(), stat.StatusByPhase)
				if err != nil {
					m.logger.Error().Err(err).Msg("error while handling loki logs")
				}
			}
		}
	}
}

func (m *CCIPMultiCallLoadGenerator) Call(_ *wasp.Generator) *wasp.Response {
	res := &wasp.Response{}
	msgBatches, returnValuesByDest, err := m.MergeCalls()
	if err != nil {
		res.Error = err.Error()
		res.Failed = true
		return res
	}
	defer func() {
		m.responses <- returnValuesByDest
	}()
	source := m.client.GetNetworkName()
	m.logger.Info().Int("Number of msg Batches", len(msgBatches)).Msgf("Sending ccip-send calls for network %s", source)
	startTime := time.Now().UTC()

	// for now we are using all ccip-sends with native
	for counter, batch := range msgBatches {
		m.logger.Debug().Interface("msgs in one tx", batch).Msgf("Sending %d ccip-send calls", len(batch))
		sendTx, err := contracts.MultiCallCCIP(m.client, m.MultiCall, batch, true)
		if err != nil {
			res.Error = err.Error()
			res.Failed = true
			return res
		}

		lggr := m.logger.With().Str("Msg Tx", sendTx.Hash().String()).Logger()
		txConfirmationTime := time.Now().UTC()
		rcpt, err1 := bind.WaitMined(context.Background(), m.client.DeployBackend(), sendTx)
		if err1 == nil {
			hdr, err1 := m.client.HeaderByNumber(context.Background(), rcpt.BlockNumber)
			if err1 == nil {
				txConfirmationTime = hdr.Timestamp
			}
		}
		var gasUsed uint64
		if rcpt != nil {
			gasUsed = rcpt.GasUsed
		}
		validateGrp := errgroup.Group{}
		for dest, mcReturnValues := range returnValuesByDest {
			// locate the batch in returnValuesByDest for the destination
			allStats, statsexists := mcReturnValues.Stats[counter]
			allMsgs, msgsexists := mcReturnValues.Msgs[counter]

			if statsexists != msgsexists {
				res.Error = fmt.Sprintf("msgs and stats does not match for batch num %d for lane %s-->%s", counter, source, dest)
				res.Failed = true
				return res
			}
			if !statsexists {
				continue
			}

			if len(allStats) != len(allMsgs) {
				res.Error = fmt.Sprintf("number of stats %d and msgs %d should be same", len(allStats), len(allMsgs))
				res.Failed = true
				return res
			}
			for i, stat := range allStats {
				msg := allMsgs[i]
				stat.UpdateState(lggr, 0, testreporters.TX, startTime.Sub(txConfirmationTime), testreporters.Success,
					testreporters.TransactionStats{
						Fee:                msg.Fee.String(),
						GasUsed:            gasUsed,
						TxHash:             sendTx.Hash().Hex(),
						NoOfTokensSent:     len(msg.Msg.TokenAmounts),
						MessageBytesLength: len(msg.Msg.Data),
					})
			}

			key := fmt.Sprintf("%s-%s", allStats[0].SourceNetwork, allStats[0].DestNetwork)
			c, ok := m.E2ELoads[key]
			if !ok {
				res.Error = fmt.Sprintf("load for %s not found", key)
				res.Failed = true
				return res
			}

			lggr = lggr.With().Str("Source Network", c.Lane.SourceChain.GetNetworkName()).Str("Dest Network", c.Lane.DestChain.GetNetworkName()).Logger()
			stats := allStats
			txConfirmationTime := txConfirmationTime
			sendTx := sendTx
			lggr := lggr
			validateGrp.Go(func() error {
				return c.Validate(lggr, sendTx, txConfirmationTime, stats)
			})
		}
		err = validateGrp.Wait()
		if err != nil {
			res.Error = err.Error()
			res.Failed = true
			return res
		}
	}
	return res
}

func (m *CCIPMultiCallLoadGenerator) MergeCalls() (map[int][]contracts.CCIPMsgData, map[string]MultiCallReturnValues, error) {
	ccipMsgs := make(map[int][]contracts.CCIPMsgData)
	statDetails := make(map[string]MultiCallReturnValues)
	batchNum := 1
	chunksize := 2
	if m.ChunkSize != nil {
		chunksize = pointer.GetInt(m.ChunkSize)
	}
	for _, e2eLoad := range m.E2ELoads {
		destChainSelector, err := chain_selectors.SelectorFromChainId(e2eLoad.Lane.Source.DestinationChainId)
		if err != nil {
			return ccipMsgs, statDetails, err
		}

		allFee := big.NewInt(0)
		allStatsForDest := make(map[int][]*testreporters.RequestStat)
		allMsgsForDest := make(map[int][]contracts.CCIPMsgData)
		for i := int64(0); i < m.NoOfRequestsPerUnitTime; i++ {
			msg, stats := e2eLoad.CCIPMsg()
			msg.FeeToken = common.Address{}
			fee, err := e2eLoad.Lane.Source.Common.Router.GetFee(destChainSelector, msg)
			if err != nil {
				return ccipMsgs, statDetails, err
			}
			// transfer fee to the multicall address
			if msg.FeeToken != (common.Address{}) {
				allFee = new(big.Int).Add(allFee, fee)
			}
			msgData := contracts.CCIPMsgData{
				RouterAddr:    e2eLoad.Lane.Source.Common.Router.EthAddress,
				ChainSelector: destChainSelector,
				Msg:           msg,
				Fee:           fee,
			}

			// if length of the batch exceeds 10 create another batch
			if _, exists := ccipMsgs[batchNum]; exists && len(ccipMsgs[batchNum]) > chunksize {
				batchNum++
				ccipMsgs[batchNum] = []contracts.CCIPMsgData{msgData}
				allStatsForDest[batchNum] = []*testreporters.RequestStat{stats}
				allMsgsForDest[batchNum] = []contracts.CCIPMsgData{msgData}
			} else {
				ccipMsgs[batchNum] = append(ccipMsgs[batchNum], msgData)
				allStatsForDest[batchNum] = append(allStatsForDest[batchNum], stats)
				allMsgsForDest[batchNum] = append(allMsgsForDest[batchNum], msgData)
			}
		}
		statDetails[e2eLoad.Lane.DestNetworkName] = MultiCallReturnValues{
			Stats: allStatsForDest,
			Msgs:  allMsgsForDest,
		}
		// transfer fee to the multicall address
		if allFee.Cmp(big.NewInt(0)) > 0 {
			if err := e2eLoad.Lane.Source.Common.FeeToken.Transfer(e2eLoad.Lane.Source.Common.MulticallContract.Hex(), allFee); err != nil {
				return ccipMsgs, statDetails, err
			}
		}
	}
	return ccipMsgs, statDetails, nil
}
