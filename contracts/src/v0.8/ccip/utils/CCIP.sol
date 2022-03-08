// SPDX-License-Identifier: MIT
pragma solidity 0.8.12;

import "../../vendor/IERC20.sol";

contract CCIP {
  /// @notice High level message
  struct Message {
    uint256 sequenceNumber;
    uint256 sourceChainId;
    address sender;
    MessagePayload payload;
  }

  /// @notice Payload within the message
  struct MessagePayload {
    IERC20[] tokens;
    uint256[] amounts;
    uint256 destinationChainId;
    address receiver;
    address executor;
    bytes data;
    bytes options;
  }

  /// @notice Report that is relayed by the observing DON at the relay phase
  struct RelayReport {
    bytes32 merkleRoot;
    uint256 minSequenceNumber;
    uint256 maxSequenceNumber;
  }

  struct MerkleProof {
    bytes32[] path;
    uint256 index;
  }
}
