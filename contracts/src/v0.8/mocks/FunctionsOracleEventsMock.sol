// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

contract FunctionsOracleEventsMock {
  event AuthorizedSendersActive(address account);
  event AuthorizedSendersChanged(address[] senders, address changedBy);
  event AuthorizedSendersDeactive(address account);
  event ConfigSet(
    uint32 previousConfigBlockNumber,
    bytes32 configDigest,
    uint64 configCount,
    address[] signers,
    address[] transmitters,
    uint8 f,
    bytes onchainConfig,
    uint64 offchainConfigVersion,
    bytes offchainConfig
  );
  event Initialized(uint8 version);
  event InvalidRequestID(bytes32 indexed requestId);
  event OracleRequest(
    bytes32 indexed requestId,
    address requestingContract,
    address requestInitiator,
    uint64 subscriptionId,
    address subscriptionOwner,
    bytes data
  );
  event OracleResponse(bytes32 indexed requestId);
  event OwnershipTransferRequested(address indexed from, address indexed to);
  event OwnershipTransferred(address indexed from, address indexed to);
  event Transmitted(bytes32 configDigest, uint32 epoch);
  event UserCallbackError(bytes32 indexed requestId, string reason);
  event UserCallbackRawError(bytes32 indexed requestId, bytes lowLevelData);

  function emitAuthorizedSendersActive(address account) public {
    emit AuthorizedSendersActive(account);
  }

  function emitAuthorizedSendersChanged(address[] memory senders, address changedBy) public {
    emit AuthorizedSendersChanged(senders, changedBy);
  }

  function emitAuthorizedSendersDeactive(address account) public {
    emit AuthorizedSendersDeactive(account);
  }

  function emitConfigSet(
    uint32 previousConfigBlockNumber,
    bytes32 configDigest,
    uint64 configCount,
    address[] memory signers,
    address[] memory transmitters,
    uint8 f,
    bytes memory onchainConfig,
    uint64 offchainConfigVersion,
    bytes memory offchainConfig
  ) public {
    emit ConfigSet(
      previousConfigBlockNumber,
      configDigest,
      configCount,
      signers,
      transmitters,
      f,
      onchainConfig,
      offchainConfigVersion,
      offchainConfig
    );
  }

  function emitInitialized(uint8 version) public {
    emit Initialized(version);
  }

  function emitInvalidRequestID(bytes32 requestId) public {
    emit InvalidRequestID(requestId);
  }

  function emitOracleRequest(
    bytes32 requestId,
    address requestingContract,
    address requestInitiator,
    uint64 subscriptionId,
    address subscriptionOwner,
    bytes memory data
  ) public {
    emit OracleRequest(requestId, requestingContract, requestInitiator, subscriptionId, subscriptionOwner, data);
  }

  function emitOracleResponse(bytes32 requestId) public {
    emit OracleResponse(requestId);
  }

  function emitOwnershipTransferRequested(address from, address to) public {
    emit OwnershipTransferRequested(from, to);
  }

  function emitOwnershipTransferred(address from, address to) public {
    emit OwnershipTransferred(from, to);
  }

  function emitTransmitted(bytes32 configDigest, uint32 epoch) public {
    emit Transmitted(configDigest, epoch);
  }

  function emitUserCallbackError(bytes32 requestId, string memory reason) public {
    emit UserCallbackError(requestId, reason);
  }

  function emitUserCallbackRawError(bytes32 requestId, bytes memory lowLevelData) public {
    emit UserCallbackRawError(requestId, lowLevelData);
  }
}
