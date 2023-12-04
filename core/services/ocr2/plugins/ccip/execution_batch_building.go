package ccip

import (
	"context"

	"github.com/pkg/errors"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/abihelpers"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/commit_store"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/offramp"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/onramp"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/hashlib"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/merklemulti"
)

func getProofData(
	ctx context.Context,
	sourceReader onramp.OnRampReader,
	interval commit_store.CommitStoreInterval,
) (sendReqsInRoot []ccipdata.Event[internal.EVM2EVMMessage], leaves [][32]byte, tree *merklemulti.Tree[[32]byte], err error) {
	sendReqs, err := sourceReader.GetSendRequestsBetweenSeqNums(ctx, interval.Min, interval.Max)
	if err != nil {
		return nil, nil, nil, err
	}
	leaves = make([][32]byte, 0, len(sendReqs))
	for _, req := range sendReqs {
		leaves = append(leaves, req.Data.Hash)
	}
	tree, err = merklemulti.NewTree(hashlib.NewKeccakCtx(), leaves)
	if err != nil {
		return nil, nil, nil, err
	}
	return sendReqs, leaves, tree, nil
}

func buildExecutionReportForMessages(
	msgsInRoot []ccipdata.Event[internal.EVM2EVMMessage],
	leaves [][32]byte,
	tree *merklemulti.Tree[[32]byte],
	commitInterval commit_store.CommitStoreInterval,
	observedMessages []ObservedMessage,
) (offramp.ExecReport, error) {
	innerIdxs := make([]int, 0, len(observedMessages))
	var messages []internal.EVM2EVMMessage
	var offchainTokenData [][][]byte
	for _, observedMessage := range observedMessages {
		if observedMessage.SeqNr < commitInterval.Min || observedMessage.SeqNr > commitInterval.Max {
			// We only return messages from a single root (the root of the first message).
			continue
		}
		innerIdx := int(observedMessage.SeqNr - commitInterval.Min)
		messages = append(messages, msgsInRoot[innerIdx].Data)
		offchainTokenData = append(offchainTokenData, observedMessage.TokenData)
		innerIdxs = append(innerIdxs, innerIdx)
	}

	merkleProof, err := tree.Prove(innerIdxs)
	if err != nil {
		return offramp.ExecReport{}, err
	}

	// any capped proof will have length <= this one, so we reuse it to avoid proving inside loop, and update later if changed
	return offramp.ExecReport{
		Messages:          messages,
		Proofs:            merkleProof.Hashes,
		ProofFlagBits:     abihelpers.ProofFlagsToBits(merkleProof.SourceFlags),
		OffchainTokenData: offchainTokenData,
	}, nil
}

// Validates the given message observations do not exceed the committed sequence numbers
// in the commitStoreReader.
func validateSeqNumbers(serviceCtx context.Context, commitStore commit_store.CommitStoreReader, observedMessages []ObservedMessage) error {
	nextMin, err := commitStore.GetExpectedNextSequenceNumber(serviceCtx)
	if err != nil {
		return err
	}
	// observedMessages are always sorted by SeqNr and never empty, so it's safe to take last element
	maxSeqNumInBatch := observedMessages[len(observedMessages)-1].SeqNr

	if maxSeqNumInBatch >= nextMin {
		return errors.Errorf("Cannot execute uncommitted seq num. nextMin %v, seqNums %v", nextMin, observedMessages)
	}
	return nil
}

// Gets the commit report from the saved logs for a given sequence number.
func getCommitReportForSeqNum(ctx context.Context, commitStoreReader commit_store.CommitStoreReader, seqNum uint64) (commit_store.CommitStoreReport, error) {
	acceptedReports, err := commitStoreReader.GetCommitReportMatchingSeqNum(ctx, seqNum, 0)
	if err != nil {
		return commit_store.CommitStoreReport{}, err
	}

	if len(acceptedReports) == 0 {
		return commit_store.CommitStoreReport{}, errors.Errorf("seq number not committed")
	}

	return acceptedReports[0].Data, nil
}
