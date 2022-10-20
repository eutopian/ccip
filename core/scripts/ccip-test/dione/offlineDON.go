package dione

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"math/big"
	"strings"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	confighelper2 "github.com/smartcontractkit/libocr/offchainreporting2/confighelper"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2/types"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/scripts/ccip-test/rhea"
	"github.com/smartcontractkit/chainlink/core/scripts/common"
	helpers "github.com/smartcontractkit/chainlink/core/scripts/common"
)

type OfflineDON struct {
	Config NodesConfig
	env    Environment
	lggr   logger.Logger
}

func NewOfflineDON(env Environment, lggr logger.Logger) OfflineDON {
	config := MustReadNodeConfig(env)

	return OfflineDON{
		Config: config,
		env:    env,
		lggr:   lggr,
	}
}

func (don *OfflineDON) GenerateOracleIdentities(chain string) []confighelper2.OracleIdentityExtra {
	var oracles []confighelper2.OracleIdentityExtra

	for _, node := range don.Config.Nodes {
		evmKeys := GetOCRkeysForChainType(node.OCRKeys, "evm")

		oracles = append(oracles,
			confighelper2.OracleIdentityExtra{
				OracleIdentity: confighelper2.OracleIdentity{
					TransmitAccount:   ocr2types.Account(node.EthKeys[chain]),
					OnchainPublicKey:  gethcommon.HexToAddress(strings.TrimPrefix(evmKeys.Attributes.OnChainPublicKey, "ocr2on_evm_")).Bytes(),
					OffchainPublicKey: common.ToOffchainPublicKey("0x" + strings.TrimPrefix(evmKeys.Attributes.OffChainPublicKey, "ocr2off_evm_")),
					PeerID:            node.PeerID,
				},
				ConfigEncryptionPublicKey: common.StringTo32Bytes("0x" + strings.TrimPrefix(evmKeys.Attributes.ConfigPublicKey, "ocr2cfg_evm_")),
			})
	}
	return oracles
}

func (don *OfflineDON) FundNodeKeys(chainConfig rhea.EvmChainConfig, ownerPrivKey string, amount *big.Int) {
	nonce, err := chainConfig.Client.PendingNonceAt(context.Background(), chainConfig.Owner.From)
	helpers.PanicErr(err)
	var gasTipCap *big.Int
	if chainConfig.GasSettings.EIP1559 {
		gasTipCap, err = chainConfig.Client.SuggestGasTipCap(context.Background())
		helpers.PanicErr(err)
	}
	gasPrice, err := chainConfig.Client.SuggestGasPrice(context.Background())
	helpers.PanicErr(err)

	ownerKey, err := crypto.HexToECDSA(ownerPrivKey)
	helpers.PanicErr(err)

	for i, node := range don.Config.Nodes {
		to := gethcommon.HexToAddress(node.EthKeys[chainConfig.ChainId.String()])
		if to == gethcommon.HexToAddress("0x") {
			don.lggr.Warnf("Node %2d has no sending key configured. Skipping funding")
			continue
		}
		if chainConfig.GasSettings.EIP1559 {
			sendEthEIP1559(to, chainConfig, nonce+uint64(i), gasTipCap, ownerKey, amount)
		} else {
			sendEth(to, chainConfig, nonce+uint64(i), gasPrice, ownerKey, amount)
		}
		don.lggr.Infof("Sent %s wei to %s", amount.String(), to.Hex())
	}
}

func (don *OfflineDON) WriteToFile() error {
	path := getFileLocation(don.env, NODES_FOLDER)
	file, err := json.MarshalIndent(don.Config, "", "  ")
	if err != nil {
		return err
	}
	return WriteJSON(path, file)
}

func (don *OfflineDON) PrintConfig() {
	file, err := json.MarshalIndent(don.Config, "", "  ")
	common.PanicErr(err)

	don.lggr.Infof(string(file))
}

func sendEth(to gethcommon.Address, chainConfig rhea.EvmChainConfig, nonce uint64, gasPrice *big.Int, ownerKey *ecdsa.PrivateKey, amount *big.Int) {
	tx := types.NewTx(
		&types.LegacyTx{
			Nonce:    nonce,
			GasPrice: gasPrice,
			Gas:      21_000,
			To:       &to,
			Value:    amount,
			Data:     []byte{},
		},
	)

	signedTx, err := types.SignTx(tx, types.NewLondonSigner(chainConfig.ChainId), ownerKey)
	helpers.PanicErr(err)
	err = chainConfig.Client.SendTransaction(context.Background(), signedTx)
	helpers.PanicErr(err)
}

func sendEthEIP1559(to gethcommon.Address, chainConfig rhea.EvmChainConfig, nonce uint64, gasTipCap *big.Int, ownerKey *ecdsa.PrivateKey, amount *big.Int) {
	tx := types.NewTx(
		&types.DynamicFeeTx{
			ChainID:    chainConfig.ChainId,
			Nonce:      nonce,
			GasTipCap:  gasTipCap,
			GasFeeCap:  big.NewInt(2e9),
			Gas:        uint64(21_000),
			To:         &to,
			Value:      amount,
			Data:       []byte{},
			AccessList: types.AccessList{},
		},
	)

	signedTx, err := types.SignTx(tx, types.NewLondonSigner(chainConfig.ChainId), ownerKey)
	helpers.PanicErr(err)
	err = chainConfig.Client.SendTransaction(context.Background(), signedTx)
	helpers.PanicErr(err)
}
