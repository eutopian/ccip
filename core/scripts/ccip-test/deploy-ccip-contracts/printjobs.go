package main

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	ccip_shared "github.com/smartcontractkit/chainlink/core/scripts/ccip-test/ccip-shared"
)

func PrintJobSpecs(onramp, offramp, executor common.Address) {
	jobs := fmt.Sprintf(bootstrapTemplate, offramp)
	for i, oracle := range ccip_shared.Oracles {
		jobs += "\n" + fmt.Sprintf(relayTemplate, i, offramp, onramp,
			ccip_shared.Kovan.ChainId.Int64(), ccip_shared.Rinkeby.ChainId.Int64(),
			oracle.OracleIdentity.TransmitAccount, ccip_shared.BootstrapPeerID)
		jobs += fmt.Sprintf(executionTemplate, onramp, offramp, executor,
			ccip_shared.Kovan.ChainId.Int64(), ccip_shared.Rinkeby.ChainId.Int64(),
			oracle.OracleIdentity.TransmitAccount, ccip_shared.BootstrapPeerID)
	}
	fmt.Println(jobs)
}

const bootstrapTemplate = `
// Bootstrap Node
# CCIPBootstrapSpec
type                                = "ccip-bootstrap"
name                                = "ccip-bootstrap"
schemaVersion                       = 1
contractAddress                     = "%s"
evmChainID                          = 4
isBootstrapPeer                     = true
contractConfigConfirmations         = 1
contractConfigTrackerPollInterval   = "60s"
`

const relayTemplate = `
// Node %d
# CCIPRelaySpec
type                = "ccip-relay"
name                = "ccip-relay"
schemaVersion       = 1
offRampAddress      = "%s"
onRampAddress       = "%s"
sourceEvmChainID    = "%d"
destEvmChainID      = "%d"
keyBundleID         = "<KEY-BUNDLE-ID>"
transmitterAddress  = "%s"
p2pBootstrapPeers   = ["%s@<BOOTSTRAP-HOST>:<PORT>"]
`

const executionTemplate = `
# CCIPExecutionSpec
type                = "ccip-execution"
name                = "ccip-execution"
schemaVersion       = 1
onRampAddress       = "%s"
offRampAddress      = "%s"
executorAddress     = "%s"
sourceEvmChainID    = "%d"
destEvmChainID      = "%d"
keyBundleID         = "<KEY-BUNDLE-ID>"
transmitterAddress  = "%s"
p2pBootstrapPeers   = ["%s@<BOOTSTRAP-HOST>:<PORT>"]
`
