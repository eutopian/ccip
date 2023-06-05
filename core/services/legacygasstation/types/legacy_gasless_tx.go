package types

import (
	"errors"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/core/utils"
)

type Status uint8

const (
	Submitted       Status = iota // meta-transaction is submitted to tx manager
	SourceFinalized               // cross-chain meta-transaction is finalized on the source chain
	Finalized                     // same-chain meta-transaction is finalized on the source chain
	Failure                       // same-chain or cross-chain meta-transaction failed
)

func (s Status) String() string {
	switch s {
	case Submitted:
		return "submitted"
	case SourceFinalized:
		return "sourceFinalized"
	case Finalized:
		return "finalized"
	case Failure:
		return "failure"
	}
	return "unknown"
}

func (s *Status) Scan(value interface{}) error {
	status, ok := value.(string)
	if !ok {
		return fmt.Errorf("unable to set %v of %T to Status enum", value, value)
	}
	switch status {
	case "submitted":
		*s = Submitted
	case "sourceFinalized":
		*s = SourceFinalized
	case "finalized":
		*s = Finalized
	case "failure":
		*s = Failure
	default:
		return fmt.Errorf("string to enum conversion not found for \"%s\"", status)
	}
	return nil
}

// LegacyGaslessTx encapsulates data related to a meta-transaction
type LegacyGaslessTx struct {
	ID                 string         `db:"legacy_gasless_tx_id"` // UUID
	Forwarder          common.Address `db:"forwarder_address"`    // forwarder contract
	From               common.Address `db:"from_address"`         // token sender
	Target             common.Address `db:"target_address"`       // token contract
	Receiver           common.Address `db:"receiver_address"`     // token receiver
	Nonce              *utils.Big     // forwarder nonce
	Amount             *utils.Big     // token amount to be transferred
	SourceChainID      uint64         // meta-transaction source chain ID. This is CCIP chain selector instead of EVM chain ID.
	DestinationChainID uint64         // meta-transaction destination chain ID. This is CCIP chain selector instead of EVM chain ID.
	ValidUntilTime     *utils.Big     // unix timestamp of meta-transaction expiry in seconds
	Signature          []byte         `db:"tx_signature"` // EIP712 signature
	Status             Status         `db:"tx_status"`    // status of meta-transaction
	FailureReason      *string        // failure reason of meta-transaction
	TokenName          string         // name of token used to generate EIP712 domain separator hash
	TokenVersion       string         // version of token used to generate EIP712 domain separator hash
	CCIPMessageID      *common.Hash   `db:"ccip_message_id"` // CCIP message ID
	EthTxID            string         `db:"eth_tx_id"`       // tx ID in transaction manager
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

func (gt *LegacyGaslessTx) Key() (*string, error) {
	if gt.Forwarder == utils.ZeroAddress {
		return nil, errors.New("empty forwarder address")
	}
	if gt.From == utils.ZeroAddress {
		return nil, errors.New("empty from address")
	}
	if gt.Nonce == nil {
		return nil, errors.New("nil nonce")
	}

	key := fmt.Sprintf("%s/%s/%s", gt.Forwarder, gt.From, gt.Nonce.String())
	return &key, nil
}
