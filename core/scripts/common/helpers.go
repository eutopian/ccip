package common

import (
	"flag"
	"fmt"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
)

// PanicErr panics if an error is detected
func PanicErr(err error) {
	if err != nil {
		panic(err)
	}
}

// ParseArgs parses arguments and ensures the required args are set.
func ParseArgs(flagSet *flag.FlagSet, args []string, requiredArgs ...string) {
	PanicErr(flagSet.Parse(args))
	seen := map[string]bool{}
	argValues := map[string]string{}
	flagSet.Visit(func(f *flag.Flag) {
		seen[f.Name] = true
		argValues[f.Name] = f.Value.String()
	})
	for _, req := range requiredArgs {
		if !seen[req] {
			panic(fmt.Errorf("missing required -%s argument/flag", req))
		}
	}
}

// ExplorerLink creates a block explorer link for the given transaction hash. If the chain ID is
// unrecognized, the hash is returned as-is.
func ExplorerLink(chainID int64, txHash common.Hash) string {
	var fmtURL string
	switch chainID {
	case 1: // ETH mainnet
		fmtURL = "https://etherscan.io/tx/%s"
	case 4: // Rinkeby
		fmtURL = "https://rinkeby.etherscan.io/tx/%s"
	case 5: // Goerli
		fmtURL = "https://goerli.etherscan.io/tx/%s"
	case 42: // Kovan
		fmtURL = "https://kovan.etherscan.io/tx/%s"

	case 56: // BSC mainnet
		fmtURL = "https://bscscan.com/tx/%s"
	case 97: // BSC testnet
		fmtURL = "https://testnet.bscscan.com/tx/%s"

	case 137: // Polygon mainnet
		fmtURL = "https://polygonscan.com/tx/%s"
	case 80001: // Polygon Mumbai testnet
		fmtURL = "https://mumbai.polygonscan.com/tx/%s"

	case 250: // Fantom mainnet
		fmtURL = "https://ftmscan.com/tx/%s"
	case 4002: // Fantom testnet
		fmtURL = "https://testnet.ftmscan.com/tx/%s"

	case 43114: // Avalanche mainnet
		fmtURL = "https://snowtrace.io/tx/%s"
	case 43113: // Avalanche testnet
		fmtURL = "https://testnet.snowtrace.io/tx/%s"

	default: // Unknown chain, return TX as-is
		fmtURL = "%s"
	}

	return fmt.Sprintf(fmtURL, txHash.String())
}

// ChainName returns the name of the EVM network based on its chainID
func ChainName(chainID int64) string {
	switch chainID {
	case 1:
		return "Ethereum"
	case 4:
		return "Rinkeby"
	case 5:
		return "Goerli"
	case 42:
		return "Kovan"
	case 56:
		return "BSC"
	case 97:
		return "BSC Testnet"
	case 137:
		return "Polygon"
	case 4002:
		return "Fantom testnet"
	case 80001:
		return "Polygon Mumbai"
	default: // Unknown chain, return chainID as string
		return strconv.FormatInt(chainID, 10)
	}
}
