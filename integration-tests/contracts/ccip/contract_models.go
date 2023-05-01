package ccip

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	"github.com/smartcontractkit/chainlink-testing-framework/blockchain"

	networks "github.com/smartcontractkit/chainlink/integration-tests"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/commit_store"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/evm_2_evm_offramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/evm_2_evm_onramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/lock_release_token_pool"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/maybe_revert_message_receiver"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/mock_afn_contract"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/price_registry"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/simple_message_receiver"
)

var HundredCoins = new(big.Int).Mul(big.NewInt(1e18), big.NewInt(100))

type RateLimiterConfig struct {
	Rate     *big.Int
	Capacity *big.Int
}

type AFNConfig struct {
	AFNWeightsByParticipants map[string]*big.Int // mapping : AFN participant address => weight
	ThresholdForBlessing     *big.Int
	ThresholdForBadSignal    *big.Int
}

type LinkToken struct {
	client     blockchain.EVMClient
	instance   *link_token_interface.LinkToken
	EthAddress common.Address
}

func (token *LinkToken) Address() string {
	return token.EthAddress.Hex()
}

func (token *LinkToken) BalanceOf(ctx context.Context, addr string) (*big.Int, error) {
	opts := &bind.CallOpts{
		From:    common.HexToAddress(token.client.GetDefaultWallet().Address()),
		Context: ctx,
	}
	balance, err := token.instance.BalanceOf(opts, common.HexToAddress(addr))
	if err != nil {
		return nil, err
	}
	return balance, nil
}

func (l *LinkToken) Approve(to string, amount *big.Int) error {
	opts, err := l.client.TransactionOpts(l.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	log.Info().
		Str("From", l.client.GetDefaultWallet().Address()).
		Str("To", to).
		Str("Token", l.Address()).
		Str("Amount", amount.String()).
		Uint64("Nonce", opts.Nonce.Uint64()).
		Str("Network Name", l.client.GetNetworkConfig().Name).
		Msg("Approving LINK Transfer")
	tx, err := l.instance.Approve(opts, common.HexToAddress(to), amount)
	if err != nil {
		return err
	}
	return l.client.ProcessTransaction(tx)
}

func (l *LinkToken) Transfer(to string, amount *big.Int) error {
	opts, err := l.client.TransactionOpts(l.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	log.Info().
		Str("From", l.client.GetDefaultWallet().Address()).
		Str("To", to).
		Str("Amount", amount.String()).
		Uint64("Nonce", opts.Nonce.Uint64()).
		Str("Network Name", l.client.GetNetworkConfig().Name).
		Msg("Transferring LINK")
	tx, err := l.instance.Transfer(opts, common.HexToAddress(to), amount)
	if err != nil {
		return err
	}
	return l.client.ProcessTransaction(tx)
}

// LockReleaseTokenPool represents a LockReleaseTokenPool address
type LockReleaseTokenPool struct {
	client     blockchain.EVMClient
	instance   *lock_release_token_pool.LockReleaseTokenPool
	EthAddress common.Address
}

func (pool *LockReleaseTokenPool) Address() string {
	return pool.EthAddress.Hex()
}

func (pool *LockReleaseTokenPool) RemoveLiquidity(amount *big.Int) error {
	opts, err := pool.client.TransactionOpts(pool.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	log.Info().
		Str("Token Pool", pool.Address()).
		Str("Amount", amount.String()).
		Msg("Initiating removing funds from pool")
	tx, err := pool.instance.RemoveLiquidity(opts, amount)
	if err != nil {
		return err
	}
	log.Info().
		Str("Token Pool", pool.Address()).
		Str("Amount", amount.String()).
		Str("Network Name", pool.client.GetNetworkConfig().Name).
		Msg("Liquidity removed")
	return pool.client.ProcessTransaction(tx)
}

func (pool *LockReleaseTokenPool) AddLiquidity(linkToken *LinkToken, amount *big.Int) error {
	log.Info().
		Str("Link Token", linkToken.Address()).
		Str("Token Pool", pool.Address()).
		Msg("Initiating transferring of token to token pool")
	err := linkToken.Approve(pool.Address(), amount)
	if err != nil {
		return err
	}
	err = pool.client.WaitForEvents()
	if err != nil {
		return err
	}
	opts, err := pool.client.TransactionOpts(pool.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	log.Info().
		Str("Token Pool", pool.Address()).
		Msg("Initiating adding Tokens in pool")
	tx, err := pool.instance.AddLiquidity(opts, amount)
	if err != nil {
		return err
	}
	log.Info().
		Str("Token Pool", pool.Address()).
		Str("Link Token", linkToken.Address()).
		Str("Network Name", pool.client.GetNetworkConfig().Name).
		Msg("Liquidity added")
	return pool.client.ProcessTransaction(tx)
}

func (pool *LockReleaseTokenPool) SetOnRamp(onRamp common.Address) error {
	opts, err := pool.client.TransactionOpts(pool.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	log.Info().
		Str("Token Pool", pool.Address()).
		Msg("Setting on ramp for onramp router")
	tx, err := pool.instance.ApplyRampUpdates(opts, []lock_release_token_pool.TokenPoolRampUpdate{{Ramp: onRamp, Allowed: true}}, []lock_release_token_pool.TokenPoolRampUpdate{})

	if err != nil {
		return err
	}
	log.Info().
		Str("Token Pool", pool.Address()).
		Str("OnRamp", onRamp.Hex()).
		Str("Network Name", pool.client.GetNetworkConfig().Name).
		Msg("OnRamp is set")
	return pool.client.ProcessTransaction(tx)
}

func (pool *LockReleaseTokenPool) SetOffRamp(offRamp common.Address) error {
	opts, err := pool.client.TransactionOpts(pool.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	log.Info().
		Str("Token Pool", pool.Address()).
		Msg("Setting off ramp for Token Pool")
	tx, err := pool.instance.ApplyRampUpdates(opts, []lock_release_token_pool.TokenPoolRampUpdate{}, []lock_release_token_pool.TokenPoolRampUpdate{{Ramp: offRamp, Allowed: true}})
	if err != nil {
		return err
	}
	log.Info().
		Str("Token Pool", pool.Address()).
		Str("OffRamp", offRamp.Hex()).
		Str("Network Name", pool.client.GetNetworkConfig().Name).
		Msg("OffRamp is set")
	return pool.client.ProcessTransaction(tx)
}

type AFN struct {
	client     blockchain.EVMClient
	instance   *mock_afn_contract.MockAFNContract
	EthAddress common.Address
}

func (afn *AFN) Address() string {
	return afn.EthAddress.Hex()
}

type CommitStore struct {
	client     blockchain.EVMClient
	Instance   *commit_store.CommitStore
	EthAddress common.Address
}

func (bv *CommitStore) Address() string {
	return bv.EthAddress.Hex()
}

// SetOCR2Config sets the offchain reporting protocol configuration
func (b *CommitStore) SetOCR2Config(
	signers []common.Address,
	transmitters []common.Address,
	f uint8,
	onchainConfig []byte,
	offchainConfigVersion uint64,
	offchainConfig []byte,
) error {
	log.Info().Str("Contract Address", b.Address()).Msg("Configuring OCR config for CommitStore Contract")
	// Set Config
	opts, err := b.client.TransactionOpts(b.client.GetDefaultWallet())
	if err != nil {
		return err
	}

	log.Info().
		Interface("signerAddresses", signers).
		Interface("transmitterAddresses", transmitters).
		Str("Network Name", b.client.GetNetworkConfig().Name).
		Msg("Configuring CommitStore")
	tx, err := b.Instance.SetOCR2Config(
		opts,
		signers,
		transmitters,
		f,
		onchainConfig,
		offchainConfigVersion,
		offchainConfig,
	)

	if err != nil {
		return err
	}
	return b.client.ProcessTransaction(tx)
}

type MessageReceiver struct {
	client     blockchain.EVMClient
	instance   *simple_message_receiver.SimpleMessageReceiver
	EthAddress common.Address
}

type ReceiverDapp struct {
	client     blockchain.EVMClient
	instance   *maybe_revert_message_receiver.MaybeRevertMessageReceiver
	EthAddress common.Address
}

func (rDapp *ReceiverDapp) Address() string {
	return rDapp.EthAddress.Hex()
}

func (rDapp *ReceiverDapp) ToggleRevert(revert bool) error {
	opts, err := rDapp.client.TransactionOpts(rDapp.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := rDapp.instance.SetRevert(opts, revert)
	if err != nil {
		return err
	}
	log.Info().
		Bool("revert", revert).
		Str("tx", tx.Hash().String()).
		Str("ReceiverDapp", rDapp.Address()).
		Str("Network Name", rDapp.client.GetNetworkConfig().Name).
		Msg("ReceiverDapp revert set")
	return rDapp.client.ProcessTransaction(tx)
}

type PriceRegistry struct {
	client     blockchain.EVMClient
	instance   *price_registry.PriceRegistry
	EthAddress common.Address
}

func (c *PriceRegistry) Address() string {
	return c.EthAddress.Hex()
}

func (c *PriceRegistry) AddPriceUpdater(addr common.Address) error {
	opts, err := c.client.TransactionOpts(c.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := c.instance.ApplyPriceUpdatersUpdates(opts, []common.Address{addr}, []common.Address{})
	if err != nil {
		return err
	}
	log.Info().
		Str("updaters", addr.Hex()).
		Str("Network Name", c.client.GetNetworkConfig().Name).
		Msg("PriceRegistry updater added")
	return c.client.ProcessTransaction(tx)
}

func (c *PriceRegistry) AddFeeToken(addr common.Address) error {
	opts, err := c.client.TransactionOpts(c.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := c.instance.ApplyFeeTokensUpdates(opts, []common.Address{addr}, []common.Address{})
	if err != nil {
		return err
	}
	log.Info().
		Str("feeTokens", addr.Hex()).
		Str("Network Name", c.client.GetNetworkConfig().Name).
		Msg("PriceRegistry feeToken set")
	return c.client.ProcessTransaction(tx)
}

func (c *PriceRegistry) UpdatePrices(priceUpdates price_registry.InternalPriceUpdates) error {
	opts, err := c.client.TransactionOpts(c.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	tx, err := c.instance.UpdatePrices(opts, priceUpdates)
	if err != nil {
		return err
	}
	log.Info().
		Str("Network Name", c.client.GetNetworkConfig().Name).
		Interface("PriceUpdates", priceUpdates).
		Msg("Prices updated")
	return c.client.ProcessTransaction(tx)
}

type Router struct {
	client     blockchain.EVMClient
	Instance   *router.Router
	EthAddress common.Address
}

func (r *Router) Copy() *Router {
	ri := *r.Instance
	return &Router{
		client:     r.client,
		Instance:   &ri,
		EthAddress: r.EthAddress,
	}
}

func (r *Router) Address() string {
	return r.EthAddress.Hex()
}

func (r *Router) SetOnRamp(chainID uint64, onRamp common.Address) error {
	opts, err := r.client.TransactionOpts(r.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	log.Info().
		Str("Router", r.Address()).
		Msg("Setting on ramp for r")

	tx, err := r.Instance.ApplyRampUpdates(opts, []router.RouterOnRampUpdate{{DestChainId: chainID, OnRamp: onRamp}}, nil)
	if err != nil {
		return err
	}
	log.Info().
		Str("onRamp", onRamp.Hex()).
		Str("Network Name", r.client.GetNetworkConfig().Name).
		Msg("Router is configured")
	return r.client.ProcessTransaction(tx)
}

func (r *Router) CCIPSend(destChainId uint64, msg router.ClientEVM2AnyMessage, valueForNative *big.Int) (*types.Transaction, error) {
	opts, err := r.client.TransactionOpts(r.client.GetDefaultWallet())
	if err != nil {
		return nil, err
	}
	if valueForNative != nil {
		opts.Value = valueForNative
	}
	opts.GasLimit = 500000
	chain := r.client.GetNetworkConfig().Name
	if chain == networks.ArbitrumGoerli.Name {
		opts.GasLimit = 100000000
	}

	tx, err := r.Instance.CcipSend(opts, destChainId, msg)
	if err != nil {
		return nil, err
	}
	log.Info().
		Str("router", r.Address()).
		Str("txHash", tx.Hash().Hex()).
		Str("Network Name", r.client.GetNetworkConfig().Name).
		Msg("msg is sent")
	return tx, r.client.ProcessTransaction(tx)
}

func (r *Router) AddOffRamp(offRamp common.Address, sourceChainId uint64) (*types.Transaction, error) {
	opts, err := r.client.TransactionOpts(r.client.GetDefaultWallet())
	if err != nil {
		return nil, err
	}
	tx, err := r.Instance.ApplyRampUpdates(opts, nil, []router.RouterOffRampUpdate{
		{SourceChainId: sourceChainId, OffRamps: []common.Address{offRamp}}})
	if err != nil {
		return nil, err
	}
	log.Info().
		Str("offRamp", offRamp.Hex()).
		Str("Network Name", r.client.GetNetworkConfig().Name).
		Msg("offRamp is added to Router")
	return tx, r.client.ProcessTransaction(tx)
}

func (r *Router) SetWrappedNative(wNative common.Address) (*types.Transaction, error) {
	opts, err := r.client.TransactionOpts(r.client.GetDefaultWallet())
	if err != nil {
		return nil, err
	}
	tx, err := r.Instance.SetWrappedNative(opts, wNative)
	if err != nil {
		return nil, err
	}
	log.Info().
		Str("wrapped native", wNative.Hex()).
		Str("router", r.Address()).
		Str("Network Name", r.client.GetNetworkConfig().Name).
		Msg("wrapped native is added for Router")
	return tx, r.client.ProcessTransaction(tx)
}

func (r *Router) GetFee(destinationChainId uint64, message router.ClientEVM2AnyMessage) (*big.Int, error) {
	return r.Instance.GetFee(nil, destinationChainId, message)
}

type OnRamp struct {
	client     blockchain.EVMClient
	Instance   *evm_2_evm_onramp.EVM2EVMOnRamp
	EthAddress common.Address
}

func (onRamp *OnRamp) Address() string {
	return onRamp.EthAddress.Hex()
}

type OffRamp struct {
	client     blockchain.EVMClient
	Instance   *evm_2_evm_offramp.EVM2EVMOffRamp
	EthAddress common.Address
}

func (offRamp *OffRamp) Address() string {
	return offRamp.EthAddress.Hex()
}

// SetOCR2Config sets the offchain reporting protocol configuration
func (offRamp *OffRamp) SetOCR2Config(
	signers []common.Address,
	transmitters []common.Address,
	f uint8,
	onchainConfig []byte,
	offchainConfigVersion uint64,
	offchainConfig []byte,
) error {
	log.Info().Str("Contract Address", offRamp.Address()).Msg("Configuring OffRamp Contract")
	// Set Config
	opts, err := offRamp.client.TransactionOpts(offRamp.client.GetDefaultWallet())
	if err != nil {
		return err
	}
	log.Info().
		Interface("signerAddresses", signers).
		Interface("transmitterAddresses", transmitters).
		Str("Network Name", offRamp.client.GetNetworkConfig().Name).
		Msg("Configuring OffRamp")
	tx, err := offRamp.Instance.SetOCR2Config(
		opts,
		signers,
		transmitters,
		f,
		onchainConfig,
		offchainConfigVersion,
		offchainConfig,
	)

	if err != nil {
		return err
	}
	return offRamp.client.ProcessTransaction(tx)
}
