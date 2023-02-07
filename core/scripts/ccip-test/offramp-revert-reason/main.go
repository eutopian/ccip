package main

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/afn_contract"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/commit_store"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/evm_2_evm_offramp"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/evm_2_evm_onramp"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/fee_manager"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/lock_release_token_pool"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/router"
	"github.com/smartcontractkit/chainlink/core/scripts/ccip-test/secrets"
)

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}

// You can either add an error string (like "0x4e487b710000000000000000000000000000000000000000000000000000000000000032")
// or you can specify an ethURL, txHash and requester.
func main() {
	errorCodeString := ""

	if errorCodeString == "" {
		// Need a node URL
		// NOTE: this node needs to run in archive mode
		ethUrl := secrets.GetRPC(5)
		txHash := "0x97be8559164442595aba46b5f849c23257905b78e72ee43d9b998b28eee78b84"
		requester := "0xe88ff73814fb891bb0e149f5578796fa41f20242"

		ec, ethErr := ethclient.Dial(ethUrl)
		panicErr(ethErr)
		errorString, _ := getErrorForTx(ec, txHash, requester)
		// Some nodes prepend "Reverted " and we also remove the 0x
		trimmed := strings.TrimPrefix(errorString, "Reverted ")[2:]

		contractABIs := getAllABIs()

		DecodeErrorStringFromABI(trimmed, contractABIs)
	} else {
		errorCodeString = strings.TrimPrefix(errorCodeString, "0x")
		DecodeErrorStringFromABI(errorCodeString, getAllABIs())
	}
}

func DecodeErrorStringFromABI(errorString string, contractABIs []string) {
	data, err := hex.DecodeString(errorString)
	panicErr(err)

	for _, contractABI := range contractABIs {
		parsedAbi, err2 := abi.JSON(strings.NewReader(contractABI))
		panicErr(err2)

		for k, abiError := range parsedAbi.Errors {
			if bytes.Equal(data[:4], abiError.ID.Bytes()[:4]) {
				// Found a matching error
				v, err3 := abiError.Unpack(data)
				panicErr(err3)
				fmt.Printf("Error is \"%v\" args %v\n", k, v)
				return
			}
		}
	}

	if len(errorString) > 8 && errorString[:8] == "4e487b71" {
		fmt.Println("Assertion failure")
		indicator := errorString[len(errorString)-2:]
		switch indicator {
		case "01":
			fmt.Printf("If you call assert with an argument that evaluates to false.")
		case "11":
			fmt.Printf("If an arithmetic operation results in underflow or overflow outside of an unchecked { ... } block.")
		case "12":
			fmt.Printf("If you divide or modulo by zero (e.g. 5 / 0 or 23 modulo 0).")
		case "21":
			fmt.Printf("If you convert a value that is too big or negative into an enum type.")
		case "31":
			fmt.Printf("If you call .pop() on an empty array.")
		case "32":
			fmt.Printf("If you access an array, bytesN or an array slice at an out-of-bounds or negative index (i.e. x[i] where i >= x.length or i < 0).")
		case "41":
			fmt.Printf("If you allocate too much memory or create an array that is too large.")
		case "51":
			fmt.Printf("If you call a zero-initialized variable of internal function type.")
		default:
			fmt.Printf("This is a revert produced by an assertion failure. Exact code not found \"%s\"", indicator)
		}
		return
	}

	fmt.Printf("Cannot match error with contract ABI. Error code \"%v\"\n", "trimmed")
}

func getAllABIs() []string {
	return []string{
		afn_contract.AFNContractABI,
		lock_release_token_pool.LockReleaseTokenPoolABI,
		commit_store.CommitStoreABI,
		fee_manager.FeeManagerABI,
		evm_2_evm_onramp.EVM2EVMOnRampABI,
		evm_2_evm_offramp.EVM2EVMOffRampABI,
		router.RouterABI,
	}
}

func getErrorForTx(client *ethclient.Client, txHash string, requester string) (string, common.Address) {
	tx, _, err := client.TransactionByHash(context.Background(), common.HexToHash(txHash))
	panicErr(err)
	re, err := client.TransactionReceipt(context.Background(), common.HexToHash(txHash))
	panicErr(err)

	call := ethereum.CallMsg{
		From:     common.HexToAddress(requester),
		To:       tx.To(),
		Data:     tx.Data(),
		Gas:      tx.Gas(),
		GasPrice: tx.GasPrice(),
	}
	_, err = client.CallContract(context.Background(), call, re.BlockNumber)
	if err == nil {
		panic("no error calling contract")
	}

	return parseError(err), *tx.To()
}

func parseError(txError error) string {
	b, err := json.Marshal(txError)
	panicErr(err)
	var callErr struct {
		Code    int
		Data    string `json:"data"`
		Message string `json:"message"`
	}
	err = json.Unmarshal(b, &callErr)
	panicErr(err)

	if callErr.Data == "" && strings.Contains(callErr.Message, "missing trie node") {
		panic("Use an archive node")
	}
	return callErr.Data
}
