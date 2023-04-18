// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import "../helpers/MerkleHelper.sol";
import "../commitStore/CommitStore.t.sol";
import "../onRamp/EVM2EVMOnRampSetup.t.sol";
import "../offRamp/EVM2EVMOffRampSetup.t.sol";

contract E2E is EVM2EVMOnRampSetup, CommitStoreSetup, EVM2EVMOffRampSetup {
  using Internal for Internal.EVM2EVMMessage;

  function setUp() public virtual override(EVM2EVMOnRampSetup, CommitStoreSetup, EVM2EVMOffRampSetup) {
    EVM2EVMOnRampSetup.setUp();
    CommitStoreSetup.setUp();
    EVM2EVMOffRampSetup.setUp();

    deployOffRamp(s_commitStore, s_destRouter);
  }

  function testE2E_3MessagesSuccess_gas() public {
    vm.pauseGasMetering();
    IERC20 token0 = IERC20(s_sourceTokens[0]);
    IERC20 token1 = IERC20(s_sourceTokens[1]);
    uint256 balance0Pre = token0.balanceOf(OWNER);
    uint256 balance1Pre = token1.balanceOf(OWNER);

    Internal.EVM2EVMMessage[] memory messages = new Internal.EVM2EVMMessage[](3);
    messages[0] = sendRequest(1);
    messages[1] = sendRequest(2);
    messages[2] = sendRequest(3);

    uint256 expectedFee = s_sourceRouter.getFee(DEST_CHAIN_ID, _generateTokenMessage());
    // Asserts that the tokens have been sent and the fee has been paid.
    assertEq(balance0Pre - messages.length * (i_tokenAmount0 + expectedFee), token0.balanceOf(OWNER));
    assertEq(balance1Pre - messages.length * i_tokenAmount1, token1.balanceOf(OWNER));

    bytes32 metaDataHash = s_offRamp.metadataHash();

    bytes32[] memory hashedMessages = new bytes32[](3);
    hashedMessages[0] = messages[0]._hash(metaDataHash);
    hashedMessages[1] = messages[1]._hash(metaDataHash);
    hashedMessages[2] = messages[2]._hash(metaDataHash);

    bytes32[] memory merkleRoots = new bytes32[](1);
    merkleRoots[0] = MerkleHelper.getMerkleRoot(hashedMessages);

    address[] memory onRamps = new address[](1);
    onRamps[0] = ON_RAMP_ADDRESS;

    bytes memory commitReport = abi.encode(
      CommitStore.CommitReport({
        priceUpdates: getEmptyPriceUpdates(),
        interval: CommitStore.Interval(messages[0].sequenceNumber, messages[2].sequenceNumber),
        merkleRoot: merkleRoots[0]
      })
    );

    vm.resumeGasMetering();
    s_commitStore.report(commitReport);
    vm.pauseGasMetering();

    bytes32[] memory proofs = new bytes32[](0);
    uint256 timestamp = s_commitStore.verify(merkleRoots, proofs, 2**2 - 1);
    assertEq(BLOCK_TIME, timestamp);

    // We change the block time so when execute would e.g. use the current
    // block time instead of the committed block time the value would be
    // incorrect in the checks below.
    vm.warp(BLOCK_TIME + 2000);

    vm.expectEmit();
    emit ExecutionStateChanged(
      messages[0].sequenceNumber,
      messages[0].messageId,
      Internal.MessageExecutionState.SUCCESS
    );

    vm.expectEmit();
    emit ExecutionStateChanged(
      messages[1].sequenceNumber,
      messages[1].messageId,
      Internal.MessageExecutionState.SUCCESS
    );

    vm.expectEmit();
    emit ExecutionStateChanged(
      messages[2].sequenceNumber,
      messages[2].messageId,
      Internal.MessageExecutionState.SUCCESS
    );

    Internal.ExecutionReport memory execReport = _generateReportFromMessages(messages);
    vm.resumeGasMetering();
    s_offRamp.execute(execReport, false);
  }

  function sendRequest(uint64 expectedSeqNum) public returns (Internal.EVM2EVMMessage memory) {
    Client.EVM2AnyMessage memory message = _generateTokenMessage();
    uint256 expectedFee = s_sourceRouter.getFee(DEST_CHAIN_ID, message);

    IERC20(s_sourceTokens[0]).approve(address(s_sourceRouter), i_tokenAmount0 + expectedFee);
    IERC20(s_sourceTokens[1]).approve(address(s_sourceRouter), i_tokenAmount1);

    message.receiver = abi.encode(address(s_receiver));
    Internal.EVM2EVMMessage memory geEvent = _messageToEvent(
      message,
      expectedSeqNum,
      expectedSeqNum,
      expectedFee,
      OWNER
    );

    vm.expectEmit();
    emit CCIPSendRequested(geEvent);

    vm.resumeGasMetering();
    s_sourceRouter.ccipSend(DEST_CHAIN_ID, message);
    vm.pauseGasMetering();

    return geEvent;
  }
}
