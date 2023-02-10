// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import "../../interfaces/ICommitStore.sol";

contract MockCommitStore is ICommitStore {
  /// @inheritdoc ICommitStore
  function getExpectedNextSequenceNumber() external pure returns (uint64) {
    return 1;
  }

  /// @inheritdoc ICommitStore
  function setMinSeqNr(uint64 minSeqNr) external {}

  /// @inheritdoc ICommitStore
  function getConfig() external view override returns (ICommitStore.CommitStoreConfig memory) {
    return ICommitStore.CommitStoreConfig({chainId: 1, sourceChainId: 2, onRamp: address(1)});
  }

  /// @inheritdoc ICommitStore
  function verify(
    bytes32[] calldata,
    bytes32[] calldata,
    uint256
  ) external pure returns (uint256 timestamp) {
    return 1;
  }

  /// @inheritdoc ICommitStore
  function merkleRoot(
    bytes32[] memory leaves,
    bytes32[] memory,
    uint256
  ) public pure returns (bytes32) {
    return leaves[0];
  }

  /// @inheritdoc ICommitStore
  function getMerkleRoot(bytes32) external pure returns (uint256) {
    return 1;
  }

  /// @inheritdoc ICommitStore
  function isBlessed(bytes32) external pure returns (bool) {
    return true;
  }
}
