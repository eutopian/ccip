package load

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/smartcontractkit/wasp"
	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"

	"github.com/smartcontractkit/chainlink/integration-tests/actions"
	"github.com/smartcontractkit/chainlink/integration-tests/testsetups"
)

type laneLoadCfg struct {
	schedule []*wasp.Segment
	lane     *actions.CCIPLane
}

type loadArgs struct {
	t             *testing.T
	lggr          zerolog.Logger
	ctx           context.Context
	ccipLoad      []*CCIPE2ELoad
	loadRunner    []*wasp.Generator
	LaneLoadCfg   chan laneLoadCfg
	RunnerWg      *errgroup.Group // to wait on individual load generators run
	LoadStarterWg *sync.WaitGroup // waits for all the runners to start
	TestCfg       *testsetups.CCIPTestConfig
	TestSetupArgs *testsetups.CCIPTestSetUpOutputs
	EstimatedEnd  time.Time
}

func (l *loadArgs) Setup() {
	transferAmounts := []*big.Int{big.NewInt(1)}
	lggr := l.lggr
	var setUpArgs *testsetups.CCIPTestSetUpOutputs
	if !l.TestCfg.ExistingDeployment {
		setUpArgs = testsetups.CCIPDefaultTestSetUp(l.TestCfg.Test, lggr, "load-ccip",
			map[string]interface{}{
				"replicas":   "6",
				"prometheus": "true",
			}, transferAmounts, 5, true, true, l.TestCfg)
	} else {
		setUpArgs = testsetups.CCIPExistingDeploymentTestSetUp(l.TestCfg.Test, lggr, transferAmounts, true, l.TestCfg)
	}
	l.TestSetupArgs = setUpArgs
	if len(setUpArgs.Lanes) == 0 {
		return
	}
	l.EstimatedEnd = time.Now().Add(l.TestCfg.TestDuration)
}

func (l *loadArgs) TriggerLoad(schedule ...*wasp.Segment) {
	l.Start()
	if len(schedule) == 0 {
		schedule = wasp.Plain(l.TestCfg.Load.LoadRPS, l.TestCfg.TestDuration)
	}
	for _, lane := range l.TestSetupArgs.Lanes {
		if lane.LaneDeployed {
			l.LaneLoadCfg <- laneLoadCfg{
				schedule: schedule,
				lane:     lane.ForwardLane,
			}
			if lane.ReverseLane != nil {
				l.LaneLoadCfg <- laneLoadCfg{
					schedule: schedule,
					lane:     lane.ReverseLane,
				}
			}
		}
	}
	l.TestSetupArgs.Reporter.SetDuration(l.TestCfg.TestDuration)
	l.TestSetupArgs.Reporter.SetRPS(l.TestCfg.Load.LoadRPS)
}

func (l *loadArgs) AddMoreLanesToRun() {
	require.Len(l.t, l.TestSetupArgs.Lanes, 1, "lane for first network pair should be deployed already")
	if len(l.TestSetupArgs.Lanes) == len(l.TestCfg.NetworkPairs) {
		l.lggr.Info().Msg("All lanes are already deployed, no need to add more lanes")
		return
	}
	transferAmounts := []*big.Int{big.NewInt(1)}
	// set the ticker duration based on number of network pairs and the total test duration
	noOfPair := int64(len(l.TestCfg.NetworkPairs))
	step := l.TestCfg.TestDuration.Nanoseconds() / noOfPair
	ticker := time.NewTicker(time.Duration(step))
	// Lane for the first network pair is already deployed
	netIndex := 1
	for {
		select {
		case <-ticker.C:
			n := l.TestCfg.NetworkPairs[netIndex]
			l.lggr.Info().
				Str("Network 1", n.NetworkA.Name).
				Str("Network 2", n.NetworkB.Name).
				Msg("Adding lanes for network pair")
			err := l.TestSetupArgs.AddLanesForNetworkPair(
				l.lggr, n.NetworkA, n.NetworkB,
				n.ChainClientA, n.ChainClientB,
				transferAmounts, 5, true,
				true, false)
			assert.NoError(l.t, err)
			// set the duration for the load to be the remaining time
			dur := l.EstimatedEnd.Sub(time.Now())
			// if the estimated end is already passed, set the duration to mentioned duration
			if l.EstimatedEnd.Before(time.Now()) {
				l.lggr.Warn().Msgf("Estimated end time is already passed, setting the load for additional %s", l.TestCfg.TestDuration)
				dur = l.TestCfg.TestDuration
			}
			l.LaneLoadCfg <- laneLoadCfg{
				schedule: wasp.Plain(l.TestCfg.Load.LoadRPS, dur),
				lane:     l.TestSetupArgs.Lanes[netIndex].ForwardLane,
			}
			if l.TestSetupArgs.Lanes[netIndex].ReverseLane != nil {
				l.LaneLoadCfg <- laneLoadCfg{
					schedule: wasp.Plain(l.TestCfg.Load.LoadRPS, dur),
					lane:     l.TestSetupArgs.Lanes[netIndex].ReverseLane,
				}
			}
			netIndex++
			if netIndex >= len(l.TestCfg.NetworkPairs) {
				ticker.Stop()
				return
			}
		}
	}
}

// Start polls the LaneLoadCfg channel for new lanes and starts the load runner.
// LaneLoadCfg channel should receive a lane whenever the deployment is complete.
func (l *loadArgs) Start() {
	l.LoadStarterWg.Add(1)
	go func() {
		defer l.LoadStarterWg.Done()
		loadCount := 0
		for {
			select {
			case cfg := <-l.LaneLoadCfg:
				loadCount++
				lane := cfg.lane
				l.lggr.Info().
					Str("Source Network", lane.SourceNetworkName).
					Str("Destination Network", lane.DestNetworkName).
					Msg("Starting load for lane")

				schedule := cfg.schedule
				ccipLoad := NewCCIPLoad(l.TestCfg.Test, lane, l.TestCfg.PhaseTimeout, 100000, lane.Reports)
				ccipLoad.BeforeAllCall(testsetups.DataOnlyTransfer)
				var namespace string
				if lane.TestEnv != nil && lane.TestEnv.K8Env != nil && lane.TestEnv.K8Env.Cfg != nil {
					namespace = lane.TestEnv.K8Env.Cfg.Namespace
				}

				loadRunner, err := wasp.NewGenerator(&wasp.Config{
					T:           l.TestCfg.Test,
					Schedule:    schedule,
					LoadType:    wasp.RPSScheduleType,
					CallTimeout: l.TestCfg.Load.LoadTimeOut,
					Gun:         ccipLoad,
					Logger:      zerolog.Logger{},
					SharedData:  l.TestCfg.MsgType,
					LokiConfig:  wasp.NewEnvLokiConfig(),
					Labels: map[string]string{
						"test_group":   "load",
						"cluster":      "sdlc",
						"namespace":    namespace,
						"test_id":      "ccip",
						"source_chain": lane.SourceNetworkName,
						"dest_chain":   lane.DestNetworkName,
					},
				})
				require.NoError(l.TestCfg.Test, err, "initiating loadgen for lane %s --> %s",
					lane.SourceNetworkName, lane.DestNetworkName)
				loadRunner.Run(false)
				l.ccipLoad = append(l.ccipLoad, ccipLoad)
				l.loadRunner = append(l.loadRunner, loadRunner)
				l.RunnerWg.Go(func() error {
					_, failed := loadRunner.Wait()
					if failed {
						return fmt.Errorf("load run is failed")
					}
					if len(loadRunner.Errors()) > 0 {
						return fmt.Errorf("error in load sequence call %v", loadRunner.Errors())
					}
					return nil
				})
				if loadCount == len(l.TestCfg.NetworkPairs)*2 {
					l.lggr.Info().Msg("load is running for all lanes now")
					return
				}
			}
		}
	}()
}

func (l *loadArgs) Wait() {
	// wait for load runner to start on all lanes
	l.LoadStarterWg.Wait()
	l.lggr.Info().Msg("Waiting for load to finish on all lanes")
	// wait for load runner to finish
	err := l.RunnerWg.Wait()
	require.NoError(l.t, err, "load run is failed")
}

func (l *loadArgs) TearDown() {
	if l.TestSetupArgs.TearDown != nil {
		for i := range l.ccipLoad {
			l.ccipLoad[i].ReportAcceptedLog()
		}
		l.TestSetupArgs.TearDown()
	}
}

func NewLoadArgs(t *testing.T, lggr zerolog.Logger, parent context.Context) *loadArgs {
	wg, ctx := errgroup.WithContext(parent)
	return &loadArgs{
		t:             t,
		lggr:          lggr,
		RunnerWg:      wg,
		ctx:           ctx,
		TestCfg:       testsetups.NewCCIPTestConfig(t, lggr, testsetups.Load),
		LaneLoadCfg:   make(chan laneLoadCfg),
		LoadStarterWg: &sync.WaitGroup{},
	}
}
