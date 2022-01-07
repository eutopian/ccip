package ccip_test

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/onsi/gomega"
	uuid "github.com/satori/go.uuid"

	"github.com/smartcontractkit/chainlink/core/chains/evm"
	"github.com/smartcontractkit/chainlink/core/chains/evm/addressparser"
	evmtypes "github.com/smartcontractkit/chainlink/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/afn_contract"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/lock_unlock_pool"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/message_executor"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/simple_message_receiver"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/single_token_offramp"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/single_token_onramp"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/single_token_receiver"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/single_token_sender"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	"github.com/smartcontractkit/chainlink/core/services/ccip"
	"github.com/smartcontractkit/chainlink/core/services/ccip/abihelpers"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/services/headtracker"
	httypes "github.com/smartcontractkit/chainlink/core/services/headtracker/types"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ocr2key"
	"github.com/smartcontractkit/chainlink/core/services/log"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/shutdown"
	"github.com/smartcontractkit/chainlink/core/utils"
	"github.com/smartcontractkit/libocr/commontypes"
	ocrnetworking "github.com/smartcontractkit/libocr/networking"
	confighelper2 "github.com/smartcontractkit/libocr/offchainreporting2/confighelper"
	ocrtypes2 "github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"
)

func setupChain(t *testing.T) (*backends.SimulatedBackend, *bind.TransactOpts) {
	key, err := crypto.GenerateKey()
	require.NoError(t, err)
	user, err := bind.NewKeyedTransactorWithChainID(key, big.NewInt(1337))
	chain := backends.NewSimulatedBackend(core.GenesisAlloc{
		user.From: {Balance: big.NewInt(0).Mul(big.NewInt(100), big.NewInt(1000000000000000000))}},
		ethconfig.Defaults.Miner.GasCeil)
	return chain, user
}

type CCIPContracts struct {
	sourceUser, destUser           *bind.TransactOpts
	sourceChain, destChain         *backends.SimulatedBackend
	sourcePool, destPool           *lock_unlock_pool.LockUnlockPool
	onRamp                         *single_token_onramp.SingleTokenOnRamp
	sourceLinkToken, destLinkToken *link_token_interface.LinkToken
	offRamp                        *single_token_offramp.SingleTokenOffRamp
	messageReceiver                *simple_message_receiver.SimpleMessageReceiver
	eoaTokenSender                 *single_token_sender.EOASingleTokenSender
	eoaTokenReceiver               *single_token_receiver.EOASingleTokenReceiver
	executor                       *message_executor.MessageExecutor
}

func setupCCIPContracts(t *testing.T) CCIPContracts {
	sourceChain, sourceUser := setupChain(t)
	destChain, destUser := setupChain(t)

	// Deploy link token and pool on source chain
	sourceLinkTokenAddress, _, _, err := link_token_interface.DeployLinkToken(sourceUser, sourceChain)
	require.NoError(t, err)
	sourceChain.Commit()
	sourceLinkToken, err := link_token_interface.NewLinkToken(sourceLinkTokenAddress, sourceChain)
	require.NoError(t, err)
	sourcePoolAddress, _, _, err := lock_unlock_pool.DeployLockUnlockPool(sourceUser, sourceChain, sourceLinkTokenAddress)
	require.NoError(t, err)
	sourceChain.Commit()
	sourcePool, err := lock_unlock_pool.NewLockUnlockPool(sourcePoolAddress, sourceChain)
	require.NoError(t, err)

	// Deploy link token and pool on destination chain
	destLinkTokenAddress, _, _, err := link_token_interface.DeployLinkToken(destUser, destChain)
	require.NoError(t, err)
	destChain.Commit()
	destLinkToken, err := link_token_interface.NewLinkToken(destLinkTokenAddress, destChain)
	require.NoError(t, err)
	destPoolAddress, _, _, err := lock_unlock_pool.DeployLockUnlockPool(destUser, destChain, destLinkTokenAddress)
	require.NoError(t, err)
	destChain.Commit()
	destPool, err := lock_unlock_pool.NewLockUnlockPool(destPoolAddress, destChain)
	require.NoError(t, err)
	destChain.Commit()

	// Float the offramp pool with 1M juels
	// Dest user is the owner of the dest pool, so he can store
	o, err := destPool.Owner(nil)
	require.NoError(t, err)
	require.Equal(t, destUser.From.String(), o.String())
	b, err := destLinkToken.BalanceOf(nil, destUser.From)
	require.NoError(t, err)
	t.Log("balance", b)
	_, err = destLinkToken.Approve(destUser, destPoolAddress, big.NewInt(1000000))
	require.NoError(t, err)
	destChain.Commit()
	_, err = destPool.LockOrBurn(destUser, destUser.From, big.NewInt(1000000))
	require.NoError(t, err)

	afnSourceAddress, _, _, err := afn_contract.DeployAFNContract(
		sourceUser,
		sourceChain,
		[]common.Address{sourceUser.From},
		[]*big.Int{big.NewInt(1)},
		big.NewInt(1),
		big.NewInt(1),
	)
	require.NoError(t, err)
	sourceChain.Commit()

	// Deploy onramp source chain
	onRampAddress, _, _, err := single_token_onramp.DeploySingleTokenOnRamp(
		sourceUser,             // users
		sourceChain,            // backend
		sourceChainID,          // source chain id
		sourceLinkTokenAddress, // token
		sourcePoolAddress,      // pool
		destChainID,            // remoteChainId
		destLinkTokenAddress,   // remoteToken
		[]common.Address{},     // allow list
		false,                  // enableAllowList
		big.NewInt(1),          // token bucket rate
		big.NewInt(1000),       // token bucket capacity,
		afnSourceAddress,       // AFN
		big.NewInt(86400),      //maxTimeWithoutAFNSignal 86400 seconds = one day
	)
	require.NoError(t, err)
	// We do this so onRamp.Address() works
	onRamp, err := single_token_onramp.NewSingleTokenOnRamp(onRampAddress, sourceChain)
	require.NoError(t, err)
	_, err = sourcePool.SetOnRamp(sourceUser, onRampAddress, true)
	require.NoError(t, err)

	afnDestAddress, _, _, err := afn_contract.DeployAFNContract(
		destUser,
		destChain,
		[]common.Address{destUser.From},
		[]*big.Int{big.NewInt(1)},
		big.NewInt(1),
		big.NewInt(1),
	)
	require.NoError(t, err)
	destChain.Commit()

	// Deploy offramp dest chain
	offRampAddress, _, _, err := single_token_offramp.DeploySingleTokenOffRamp(
		destUser,
		destChain,
		sourceChainID,
		destChainID,
		destLinkTokenAddress,
		destPoolAddress,
		big.NewInt(1),     // token bucket rate
		big.NewInt(1000),  // token bucket capacity,
		afnDestAddress,    // AFN
		big.NewInt(86400), //maxTimeWithoutAFNSignal 86400 seconds = one day
		big.NewInt(0),     // execution delay in seconds
	)
	require.NoError(t, err)
	offRamp, err := single_token_offramp.NewSingleTokenOffRamp(offRampAddress, destChain)
	require.NoError(t, err)
	// Set the pool to be the offramp
	_, err = destPool.SetOffRamp(destUser, offRampAddress, true)
	require.NoError(t, err)

	// Deploy offramp contract token receiver
	messageReceiverAddress, _, _, err := simple_message_receiver.DeploySimpleMessageReceiver(destUser, destChain)
	require.NoError(t, err)
	messageReceiver, err := simple_message_receiver.NewSimpleMessageReceiver(messageReceiverAddress, destChain)
	require.NoError(t, err)
	// Deploy offramp EOA token receiver
	eoaTokenReceiverAddress, _, _, err := single_token_receiver.DeployEOASingleTokenReceiver(destUser, destChain, offRampAddress)
	require.NoError(t, err)
	eoaTokenReceiver, err := single_token_receiver.NewEOASingleTokenReceiver(eoaTokenReceiverAddress, destChain)
	require.NoError(t, err)
	// Deploy onramp EOA token sender
	eoaTokenSenderAddress, _, _, err := single_token_sender.DeployEOASingleTokenSender(sourceUser, sourceChain, onRampAddress, eoaTokenReceiverAddress)
	require.NoError(t, err)
	eoaTokenSender, err := single_token_sender.NewEOASingleTokenSender(eoaTokenSenderAddress, sourceChain)
	require.NoError(t, err)

	// Deploy the message executor ocr2 contract
	executorAddress, _, _, err := message_executor.DeployMessageExecutor(destUser, destChain, offRampAddress)
	require.NoError(t, err)
	executor, err := message_executor.NewMessageExecutor(executorAddress, destChain)
	require.NoError(t, err)

	sourceChain.Commit()
	destChain.Commit()

	return CCIPContracts{
		sourceUser:       sourceUser,
		destUser:         destUser,
		sourceChain:      sourceChain,
		destChain:        destChain,
		sourcePool:       sourcePool,
		destPool:         destPool,
		onRamp:           onRamp,
		sourceLinkToken:  sourceLinkToken,
		destLinkToken:    destLinkToken,
		offRamp:          offRamp,
		messageReceiver:  messageReceiver,
		eoaTokenReceiver: eoaTokenReceiver,
		eoaTokenSender:   eoaTokenSender,
		executor:         executor,
	}
}

var (
	sourceChainID = big.NewInt(1000)
	destChainID   = big.NewInt(2000)
)

type EthKeyStoreSim struct {
	keystore.Eth
}

func (ks EthKeyStoreSim) SignTx(address common.Address, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {
	if chainID.String() == "1000" || chainID.String() == "2000" {
		// A terrible hack, just for the multichain test
		// Needs to actually use the sim
		return ks.Eth.SignTx(address, tx, big.NewInt(1337))
	}
	return ks.Eth.SignTx(address, tx, chainID)
}

var _ keystore.Eth = EthKeyStoreSim{}

func setupNodeCCIP(t *testing.T, owner *bind.TransactOpts, port int64, dbName string, sourceChain *backends.SimulatedBackend, destChain *backends.SimulatedBackend) (chainlink.Application, string, common.Address, ocr2key.KeyBundle, *configtest.TestGeneralConfig, func()) {
	// Do not want to load fixtures as they contain a dummy chainID.
	config, db := heavyweight.FullTestDB(t, fmt.Sprintf("%s%d", dbName, port), true, false)
	config.Overrides.FeatureOffchainReporting = null.BoolFrom(false)
	config.Overrides.FeatureOffchainReporting2 = null.BoolFrom(true)
	config.Overrides.GlobalGasEstimatorMode = null.NewString("FixedPrice", true)
	config.Overrides.DefaultChainID = nil
	config.Overrides.P2PListenPort = null.NewInt(port, true)
	config.Overrides.P2PNetworkingStack = ocrnetworking.NetworkingStackV2
	// Disables ocr spec validation so we can have fast polling for the test.
	config.Overrides.Dev = null.BoolFrom(true)

	var lggr = logger.TestLogger(t)
	eventBroadcaster := pg.NewEventBroadcaster(config.DatabaseURL(), 0, 0, lggr, uuid.NewV1())

	// We fake different chainIDs using the wrapped sim cltest.SimulatedBackend
	chainORM := evm.NewORM(db)
	_, err := chainORM.CreateChain(*utils.NewBig(sourceChainID), evmtypes.ChainCfg{})
	require.NoError(t, err)
	_, err = chainORM.CreateChain(*utils.NewBig(destChainID), evmtypes.ChainCfg{})
	require.NoError(t, err)
	sourceClient := cltest.NewSimulatedBackendClient(t, sourceChain, sourceChainID)
	destClient := cltest.NewSimulatedBackendClient(t, destChain, destChainID)

	keyStore := keystore.New(db, utils.FastScryptParams, lggr, config)
	simEthKeyStore := EthKeyStoreSim{Eth: keyStore.Eth()}

	// Create our chainset manually so we can have custom eth clients
	// (the wrapped sims faking different chainIDs)
	chainSet, err := evm.LoadChainSet(evm.ChainSetOpts{
		ORM:              chainORM,
		Config:           config,
		Logger:           lggr,
		DB:               db,
		KeyStore:         simEthKeyStore,
		EventBroadcaster: eventBroadcaster,
		GenEthClient: func(c evmtypes.Chain) eth.Client {
			if c.ID.String() == sourceChainID.String() {
				return sourceClient
			} else if c.ID.String() == destChainID.String() {
				return destClient
			}
			t.Fatalf("invalid chain ID %v", c.ID.String())
			return nil
		},
		GenHeadTracker: func(c evmtypes.Chain) httypes.HeadTracker {
			if c.ID.String() == sourceChainID.String() {
				return headtracker.NewHeadTracker(lggr, sourceClient, evmtest.NewChainScopedConfig(t, config), headtracker.NewORM(db, *sourceChainID), ht)
			} else if c.ID.String() == destChainID.String() {
				return headtracker.NewHeadTracker(lggr, destClient, evmtest.NewChainScopedConfig(t, config), headtracker.NewORM(db, *destChainID), ht)
			}
			t.Fatalf("invalid chain ID %v", c.ID.String())
			return nil
		},
		GenLogBroadcaster: func(c evmtypes.Chain) log.Broadcaster {
			if c.ID.String() == sourceChainID.String() {
				t.Log("Generating log broadcaster source")
				return log.NewBroadcaster(log.NewORM(db, lggr, config, *sourceChainID), sourceClient,
					evmtest.NewChainScopedConfig(t, config), lggr, nil)
			} else if c.ID.String() == destChainID.String() {
				return log.NewBroadcaster(log.NewORM(db, lggr, config, *destChainID), destClient,
					evmtest.NewChainScopedConfig(t, config), lggr, nil)
			}
			t.Fatalf("invalid chain ID %v", c.ID.String())
			return nil
		},
		GenTxManager: func(c evmtypes.Chain) bulletprooftxmanager.TxManager {
			if c.ID.String() == sourceChainID.String() {
				return bulletprooftxmanager.NewBulletproofTxManager(db, sourceClient, evmtest.NewChainScopedConfig(t, config), simEthKeyStore, eventBroadcaster, lggr)
			} else if c.ID.String() == destChainID.String() {
				return bulletprooftxmanager.NewBulletproofTxManager(db, destClient, evmtest.NewChainScopedConfig(t, config), simEthKeyStore, eventBroadcaster, lggr)
			}
			t.Fatalf("invalid chain ID %v", c.ID.String())
			return nil
		},
	})
	if err != nil {
		lggr.Fatal(err)
	}
	sig := shutdown.NewSignal()
	app, err := chainlink.NewApplication(chainlink.ApplicationOpts{
		Config:                   config,
		EventBroadcaster:         eventBroadcaster,
		ShutdownSignal:           sig,
		SqlxDB:                   db,
		KeyStore:                 keyStore,
		ChainSet:                 chainSet,
		Logger:                   lggr,
		ExternalInitiatorManager: nil,
	})
	require.NoError(t, err)
	require.NoError(t, app.GetKeyStore().Unlock("password"))
	_, err = app.GetKeyStore().P2P().Create()
	require.NoError(t, err)

	p2pIDs, err := app.GetKeyStore().P2P().GetAll()
	require.NoError(t, err)
	require.Len(t, p2pIDs, 1)
	peerID := p2pIDs[0].PeerID()

	config.Overrides.P2PPeerID = peerID
	config.Overrides.P2PListenPort = null.NewInt(port, true)
	p2paddresses := []string{
		fmt.Sprintf("127.0.0.1:%d", port),
	}
	config.Overrides.P2PV2ListenAddresses = p2paddresses
	config.Overrides.P2PV2AnnounceAddresses = p2paddresses

	_, err = app.GetKeyStore().Eth().Create(destChainID)
	require.NoError(t, err)
	sendingKeys, err := app.GetKeyStore().Eth().SendingKeys()
	require.NoError(t, err)
	require.Len(t, sendingKeys, 1)
	transmitter := sendingKeys[0].Address.Address()
	s, err := app.GetKeyStore().Eth().GetState(sendingKeys[0].ID())
	require.NoError(t, err)
	lggr.Debug(fmt.Sprintf("Transmitter address %s chainID %s", transmitter, s.EVMChainID.String()))

	// Fund the relayTransmitter address with some ETH
	n, err := destChain.NonceAt(context.Background(), owner.From, nil)
	require.NoError(t, err)

	tx := types.NewTransaction(n, transmitter, big.NewInt(1000000000000000000), 21000, big.NewInt(1000000000), nil)
	signedTx, err := owner.Signer(owner.From, tx)
	require.NoError(t, err)
	err = destChain.SendTransaction(context.Background(), signedTx)
	require.NoError(t, err)
	destChain.Commit()

	kb, err := app.GetKeyStore().OCR2().Create(chaintype.EVM)
	require.NoError(t, err)
	return app, peerID.Raw(), transmitter, kb, config, func() {
		app.Stop()
	}
}

func TestIntegration_CCIP(t *testing.T) {
	ccipContracts := setupCCIPContracts(t)
	// Oracles need ETH on the destination chain
	bootstrapNodePort := int64(19599)
	appBootstrap, bootstrapPeerID, _, _, _, _ := setupNodeCCIP(t, ccipContracts.destUser, bootstrapNodePort, "bootstrap_ccip", ccipContracts.sourceChain, ccipContracts.destChain)
	var (
		oracles      []confighelper2.OracleIdentityExtra
		transmitters []common.Address
		kbs          []ocr2key.KeyBundle
		apps         []chainlink.Application
	)
	// Set up the minimum 4 oracles all funded with destination ETH
	for i := int64(0); i < 4; i++ {
		app, peerID, transmitter, kb, cfg, _ := setupNodeCCIP(t, ccipContracts.destUser, bootstrapNodePort+1+i, fmt.Sprintf("oracle_ccip%d", i), ccipContracts.sourceChain, ccipContracts.destChain)
		// Supply the bootstrap IP and port as a V2 peer address
		cfg.Overrides.P2PV2Bootstrappers = []commontypes.BootstrapperLocator{
			{PeerID: bootstrapPeerID, Addrs: []string{
				fmt.Sprintf("127.0.0.1:%d", bootstrapNodePort),
			}},
		}
		kbs = append(kbs, kb)
		apps = append(apps, app)
		transmitters = append(transmitters, transmitter)
		oracles = append(oracles, confighelper2.OracleIdentityExtra{
			OracleIdentity: confighelper2.OracleIdentity{
				OnchainPublicKey:  []byte(kb.OnChainPublicKey()),
				TransmitAccount:   ocrtypes2.Account(transmitter.String()),
				OffchainPublicKey: kb.OffchainPublicKey(),
				PeerID:            peerID,
			},
			ConfigEncryptionPublicKey: kb.ConfigEncryptionPublicKey(),
		})
	}

	reportingPluginConfig, err := ccip.OffchainConfig{
		SourceIncomingConfirmations: 0,
		DestIncomingConfirmations:   1,
	}.Encode()
	require.NoError(t, err)

	setupOnchainConfig(t, ccipContracts, oracles, reportingPluginConfig)

	err = appBootstrap.Start()
	require.NoError(t, err)
	defer appBootstrap.Stop()

	// Add the bootstrap job
	chainSet := appBootstrap.GetChainSet()
	require.NotNil(t, chainSet)
	ocrJob, err := ccip.ValidatedCCIPBootstrapSpec(fmt.Sprintf(`
type               = "ccip-bootstrap"
schemaVersion      = 1
evmChainID         = "%s"
name               = "boot"
contractAddress    = "%s"
isBootstrapPeer    = true
contractConfigConfirmations = 1
contractConfigTrackerPollInterval = "1s"
`, destChainID, ccipContracts.offRamp.Address()))
	require.NoError(t, err)
	err = appBootstrap.AddJobV2(context.Background(), &ocrJob)
	require.NoError(t, err)

	// For each oracle add a relayer and job
	for i := 0; i < 4; i++ {
		err = apps[i].Start()
		require.NoError(t, err)
		defer apps[i].Stop()
		// Wait for peer wrapper to start
		time.Sleep(1 * time.Second)
		ccipJob, err := ccip.ValidatedCCIPSpec(fmt.Sprintf(`
type               = "ccip-relay"
schemaVersion      = 1
name               = "ccip-job-%d"
onRampAddress = "%s"
offRampAddress = "%s"
sourceEvmChainID   = "%s"
destEvmChainID     = "%s"
keyBundleID        = "%s"
transmitterAddress = "%s"
contractConfigConfirmations = 1
contractConfigTrackerPollInterval = "1s"
`, i, ccipContracts.onRamp.Address(), ccipContracts.offRamp.Address(), sourceChainID, destChainID, kbs[i].ID(), transmitters[i]))
		require.NoError(t, err)
		err = apps[i].AddJobV2(context.Background(), &ccipJob)
		require.NoError(t, err)
		// Add executor job
		ccipExecutionJob, err := ccip.ValidatedCCIPSpec(fmt.Sprintf(`
type               = "ccip-execution"
schemaVersion      = 1
name               = "ccip-executor-job-%d"
onRampAddress = "%s"
offRampAddress = "%s"
executorAddress = "%s"
sourceEvmChainID   = "%s"
destEvmChainID     = "%s"
keyBundleID        = "%s"
transmitterAddress = "%s"
contractConfigConfirmations = 1
contractConfigTrackerPollInterval = "1s"
`, i, ccipContracts.onRamp.Address(), ccipContracts.offRamp.Address(), ccipContracts.executor.Address(), sourceChainID, destChainID, kbs[i].ID(), transmitters[i]))
		require.NoError(t, err)
		err = apps[i].AddJobV2(context.Background(), &ccipExecutionJob)
		require.NoError(t, err)
	}
	// Send a request.
	// Jobs are booting but that is ok, the log broadcaster
	// will backfill this request log.
	ccipContracts.sourceUser.GasLimit = 500000
	_, err = ccipContracts.sourceLinkToken.Approve(ccipContracts.sourceUser, ccipContracts.sourcePool.Address(), big.NewInt(100))
	ccipContracts.sourceChain.Commit()
	msg := single_token_onramp.CCIPMessagePayload{
		Receiver: ccipContracts.messageReceiver.Address(),
		Data:     []byte("hello xchain world"),
		Tokens:   []common.Address{ccipContracts.sourceLinkToken.Address()},
		Amounts:  []*big.Int{big.NewInt(100)},
		Options:  nil,
	}
	tx, err := ccipContracts.onRamp.RequestCrossChainSend(ccipContracts.sourceUser, msg)
	require.NoError(t, err)
	ccipContracts.sourceChain.Commit()
	rec, err := ccipContracts.sourceChain.TransactionReceipt(context.Background(), tx.Hash())
	require.NoError(t, err)
	require.Equal(t, uint64(1), rec.Status)

	reportingPluginConfig, err = ccip.OffchainConfig{
		SourceIncomingConfirmations: 1,
		DestIncomingConfirmations:   0,
	}.Encode()
	require.NoError(t, err)

	setupOnchainConfig(t, ccipContracts, oracles, reportingPluginConfig)

	// Request should appear on all nodes eventually
	for i := 0; i < 4; i++ {
		var reqs []*ccip.Request
		ccipReqORM := ccip.NewORM(apps[i].GetSqlxDB())
		gomega.NewGomegaWithT(t).Eventually(func() bool {
			ccipContracts.sourceChain.Commit()
			reqs, err = ccipReqORM.Requests(sourceChainID, destChainID, big.NewInt(0), nil, ccip.RequestStatusUnstarted, nil, nil)
			return len(reqs) == 1
		}, 5*time.Second, 1*time.Second).Should(gomega.BeTrue())
	}

	// Once all nodes have the request, the reporting plugin should run to generate and submit a report onchain.
	// So we should eventually see a successful offramp submission.
	// Note that since we only send blocks here, it's likely that all the nodes will enter the transmission
	// phase before someone has submitted, so 1 report will succeed and 3 will revert.
	var report single_token_offramp.CCIPRelayReport
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		report, err = ccipContracts.offRamp.GetLastReport(nil)
		require.NoError(t, err)
		ccipContracts.destChain.Commit()
		return report.MinSequenceNumber.String() == "1" && report.MaxSequenceNumber.String() == "1"
	}, 10*time.Second, 1*time.Second).Should(gomega.BeTrue())

	// We should see the request in a fulfilled state on all nodes
	// after the offramp submission. There should be no
	// remaining valid requests.
	for i := 0; i < 4; i++ {
		gomega.NewGomegaWithT(t).Eventually(func() bool {
			ccipReqORM := ccip.NewORM(apps[i].GetSqlxDB())
			ccipContracts.destChain.Commit()
			reqs, err := ccipReqORM.Requests(sourceChainID, destChainID, report.MinSequenceNumber, report.MaxSequenceNumber, ccip.RequestStatusRelayConfirmed, nil, nil)
			require.NoError(t, err)
			valid, err := ccipReqORM.Requests(sourceChainID, destChainID, report.MinSequenceNumber, nil, ccip.RequestStatusUnstarted, nil, nil)
			require.NoError(t, err)
			return len(reqs) == 1 && len(valid) == 0
		}, 10*time.Second, 1*time.Second).Should(gomega.BeTrue())
	}

	// Now the merkle root is across.
	// Let's try to execute a request as an external party.
	// The raw log in the merkle root should be the abi-encoded version of the CCIPMessage
	ccipReqORM := ccip.NewORM(apps[0].GetSqlxDB())
	reqs, err := ccipReqORM.Requests(sourceChainID, destChainID, report.MinSequenceNumber, report.MaxSequenceNumber, "", nil, nil)
	require.NoError(t, err)
	root, proof := ccip.GenerateMerkleProof(32, [][]byte{reqs[0].Raw}, 0)
	// Root should match the report root
	require.True(t, bytes.Equal(root[:], report.MerkleRoot[:]))

	// Proof should verify.
	genRoot := ccip.GenerateMerkleRoot(reqs[0].Raw, proof)
	require.True(t, bytes.Equal(root[:], genRoot[:]))
	exists, err := ccipContracts.offRamp.GetMerkleRoot(nil, report.MerkleRoot)
	require.NoError(t, err)
	require.True(t, exists.Int64() > 0)

	h, err := utils.Keccak256(append([]byte{0x00}, reqs[0].Raw...))
	var leaf [32]byte
	copy(leaf[:], h)
	onchainRoot, err := ccipContracts.offRamp.GenerateMerkleRoot(nil, proof.PathForExecute(), leaf, proof.Index())
	require.NoError(t, err)
	require.Equal(t, genRoot, onchainRoot)

	// Execute the Message
	decodedMsg, err := abihelpers.DecodeCCIPMessage(reqs[0].Raw)
	require.NoError(t, err)
	abihelpers.MakeCCIPMsgArgs().PackValues([]interface{}{*decodedMsg})
	tx, err = ccipContracts.offRamp.ExecuteTransaction(ccipContracts.destUser, proof.PathForExecute(), *decodedMsg, proof.Index())
	require.NoError(t, err)
	ccipContracts.destChain.Commit()

	// We should now have the Message in the offchain receiver
	receivedMsg, err := ccipContracts.messageReceiver.SMessage(nil)
	require.NoError(t, err)
	assert.Equal(t, "hello xchain world", string(receivedMsg.Payload.Data))

	// Now let's send an EOA to EOA request
	// We can just use the sourceUser and destUser
	startBalanceSource, err := ccipContracts.sourceLinkToken.BalanceOf(nil, ccipContracts.sourceUser.From)
	require.NoError(t, err)
	startBalanceDest, err := ccipContracts.destLinkToken.BalanceOf(nil, ccipContracts.destUser.From)
	require.NoError(t, err)
	t.Log(startBalanceSource, startBalanceDest)

	ccipContracts.sourceUser.GasLimit = 500000
	// Approve the sender contract to take the tokens
	_, err = ccipContracts.sourceLinkToken.Approve(ccipContracts.sourceUser, ccipContracts.eoaTokenSender.Address(), big.NewInt(100))
	ccipContracts.sourceChain.Commit()
	// Send the tokens. Should invoke the onramp.
	// Only the destUser can execute.
	tx, err = ccipContracts.eoaTokenSender.SendTokens(ccipContracts.sourceUser, ccipContracts.destUser.From, big.NewInt(100), ccipContracts.destUser.From)
	require.NoError(t, err)
	ccipContracts.sourceChain.Commit()

	// DON should eventually send another report
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		report, err = ccipContracts.offRamp.GetLastReport(nil)
		require.NoError(t, err)
		ccipContracts.destChain.Commit()
		return report.MinSequenceNumber.String() == "2" && report.MaxSequenceNumber.String() == "2"
	}, 10*time.Second, 1*time.Second).Should(gomega.BeTrue())

	eoaReq, err := ccipReqORM.Requests(sourceChainID, destChainID, report.MinSequenceNumber, report.MaxSequenceNumber, "", nil, nil)
	require.NoError(t, err)
	root, proof = ccip.GenerateMerkleProof(32, [][]byte{eoaReq[0].Raw}, 0)
	// Root should match the report root
	require.True(t, bytes.Equal(root[:], report.MerkleRoot[:]))

	// Execute the Message
	decodedMsg, err = abihelpers.DecodeCCIPMessage(eoaReq[0].Raw)
	require.NoError(t, err)
	abihelpers.MakeCCIPMsgArgs().PackValues([]interface{}{*decodedMsg})
	tx, err = ccipContracts.offRamp.ExecuteTransaction(ccipContracts.destUser, proof.PathForExecute(), *decodedMsg, proof.Index())
	require.NoError(t, err)
	ccipContracts.destChain.Commit()

	// The destination user's balance should increase
	endBalanceSource, err := ccipContracts.sourceLinkToken.BalanceOf(nil, ccipContracts.sourceUser.From)
	require.NoError(t, err)
	endBalanceDest, err := ccipContracts.destLinkToken.BalanceOf(nil, ccipContracts.destUser.From)
	require.NoError(t, err)
	t.Log("Start balances", startBalanceSource, startBalanceDest)
	t.Log("End balances", endBalanceSource, endBalanceDest)
	assert.Equal(t, "100", big.NewInt(0).Sub(startBalanceSource, endBalanceSource).String())
	assert.Equal(t, "100", big.NewInt(0).Sub(endBalanceDest, startBalanceDest).String())

	// Now let's send a request flagged for oracle execution
	_, err = ccipContracts.sourceLinkToken.Approve(ccipContracts.sourceUser, ccipContracts.sourcePool.Address(), big.NewInt(100))
	require.NoError(t, err)
	ccipContracts.sourceChain.Commit()
	require.NoError(t, err)
	msg = single_token_onramp.CCIPMessagePayload{
		Receiver: ccipContracts.messageReceiver.Address(),
		Data:     []byte("hey DON, execute for me"),
		Tokens:   []common.Address{ccipContracts.sourceLinkToken.Address()},
		Amounts:  []*big.Int{big.NewInt(100)},
		Executor: ccipContracts.executor.Address(),
		Options:  []byte{},
	}
	_, err = ccipContracts.onRamp.RequestCrossChainSend(ccipContracts.sourceUser, msg)
	require.NoError(t, err)
	ccipContracts.sourceChain.Commit()

	// Should first be relayed, seq number 3
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		report, err = ccipContracts.offRamp.GetLastReport(nil)
		require.NoError(t, err)
		ccipContracts.destChain.Commit()
		return report.MinSequenceNumber.String() == "3" && report.MaxSequenceNumber.String() == "3"
	}, 10*time.Second, 1*time.Second).Should(gomega.BeTrue())

	// Should see the 3rd message be executed
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		it, err := ccipContracts.offRamp.FilterCrossChainMessageExecuted(nil, nil)
		require.NoError(t, err)
		ecount := 0
		for it.Next() {
			t.Log("executed", it.Event.SequenceNumber)
			ecount++
		}
		ccipContracts.destChain.Commit()
		return ecount == 3
	}, 20*time.Second, 1*time.Second).Should(gomega.BeTrue())
	// In total, we should see 3 relay reports containing seq 1,2,3
	// and 3 execution_confirmed messages
	reqs, err = ccipReqORM.Requests(sourceChainID, destChainID, big.NewInt(1), big.NewInt(3), ccip.RequestStatusExecutionConfirmed, nil, nil)
	require.NoError(t, err)
	require.Len(t, reqs, 3)
	_, err = ccipReqORM.RelayReport(big.NewInt(1))
	require.NoError(t, err)
	_, err = ccipReqORM.RelayReport(big.NewInt(2))
	require.NoError(t, err)
	_, err = ccipReqORM.RelayReport(big.NewInt(3))
	require.NoError(t, err)
}

func setupOnchainConfig(t *testing.T, ccipContracts CCIPContracts, oracles []confighelper2.OracleIdentityExtra, reportingPluginConfig []byte) {
	// Note We do NOT set the payees, payment is done in the OCR2Base implementation
	// Set the offramp config.
	signers, transmitters, threshold, onchainConfig, offchainConfigVersion, offchainConfig, err := confighelper2.ContractSetConfigArgsForTests(
		2*time.Second,        // deltaProgress
		1*time.Second,        // deltaResend
		1*time.Second,        // deltaRound
		500*time.Millisecond, // deltaGrace
		2*time.Second,        // deltaStage
		3,
		[]int{1, 1, 1, 1},
		oracles,
		reportingPluginConfig,
		50*time.Millisecond,
		50*time.Millisecond,
		50*time.Millisecond,
		50*time.Millisecond,
		50*time.Millisecond,
		1, // faults
		nil,
	)

	require.NoError(t, err)
	lggr := logger.TestLogger(t)
	lggr.Debugw("Setting Config on Oracle Contract",
		"signers", signers,
		"transmitters", transmitters,
		"threshold", threshold,
		"onchainConfig", onchainConfig,
		"encodedConfigVersion", offchainConfigVersion,
	)

	// Set the DON on the offramp
	_, err = ccipContracts.offRamp.SetConfig(
		ccipContracts.destUser,
		addressparser.OnchainPublicKeyToAddress(signers),
		addressparser.AccountToAddress(transmitters),
		threshold,
		onchainConfig,
		offchainConfigVersion,
		offchainConfig,
	)
	require.NoError(t, err)
	ccipContracts.destChain.Commit()

	// Same DON on the message executor
	_, err = ccipContracts.executor.SetConfig(
		ccipContracts.destUser,
		addressparser.OnchainPublicKeyToAddress(signers),
		addressparser.AccountToAddress(transmitters),
		threshold,
		onchainConfig,
		offchainConfigVersion,
		offchainConfig,
	)
	require.NoError(t, err)
	ccipContracts.destChain.Commit()
}
