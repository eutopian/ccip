package ccip

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/lib/pq"

	"github.com/ethereum/go-ethereum/common"
)

type Request struct {
	SeqNum        utils.Big
	SourceChainID string // TODO Note this will be some super set which includes evm_chain_id
	DestChainID   string // TODO Note this will be some super set which includes evm_chain_id
	Sender        common.Address
	Receiver      common.Address
	Data          []byte
	Tokens        pq.StringArray
	Amounts       pq.StringArray
	Executor      common.Address
	Options       []byte
	Raw           []byte // Full ABI-encoded event for merkle tree
	Status        RequestStatus
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (Request) TableName() string {
	return "ccip_requests"
}

func MakeOptions() abi.Arguments {
	mustType := func(ts string) abi.Type {
		ty, _ := abi.NewType(ts, "", nil)
		return ty
	}
	return []abi.Argument{
		{
			Type: mustType("bool"),
			Name: "oracleExecute",
		},
	}
}

func mustStringToBigInt(s string) *big.Int {
	i, ok := big.NewInt(0).SetString(s, 10)
	if !ok {
		panic(fmt.Sprintf("invalid big.Int string %v", s))
	}
	return i
}

func (r Request) ToMessage() Message {
	var tokens []common.Address
	for _, t := range r.Tokens {
		tokens = append(tokens, common.HexToAddress(t))
	}
	var amounts []*big.Int
	for _, a := range r.Amounts {
		amounts = append(amounts, mustStringToBigInt(a))
	}
	return Message{
		SequenceNumber:     r.SeqNum.ToInt(),
		SourceChainId:      mustStringToBigInt(r.SourceChainID),
		DestinationChainId: mustStringToBigInt(r.DestChainID),
		Sender:             r.Sender,
		Payload: struct {
			Receiver common.Address   `json:"receiver"`
			Data     []uint8          `json:"data"`
			Tokens   []common.Address `json:"tokens"`
			Amounts  []*big.Int       `json:"amounts"`
			Executor common.Address   `json:"executor"`
			Options  []uint8          `json:"options"`
		}{
			Receiver: r.Receiver,
			Data:     r.Data,
			Tokens:   tokens,
			Amounts:  amounts,
			Executor: r.Executor,
			Options:  r.Options,
		},
	}
}

type RequestStatus string

const (
	RequestStatusUnstarted    RequestStatus = "unstarted"
	RequestStatusRelayPending RequestStatus = "relay_pending"
	// We only mark relay confirmed after we've seen the report accepted log with sufficient
	// number of confirmations
	RequestStatusRelayConfirmed   RequestStatus = "relay_confirmed"
	RequestStatusExecutionPending RequestStatus = "execution_pending"
	// We only mark execution confirmed after we've seen the Message executed log with sufficient
	// number of confirmations
	RequestStatusExecutionConfirmed RequestStatus = "execution_confirmed"
)

type RelayReport struct {
	Root      []byte
	MinSeqNum utils.Big
	MaxSeqNum utils.Big
	CreatedAt time.Time
}
