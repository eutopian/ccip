// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import {IEVM2AnyGEOnRamp} from "../../interfaces/onRamp/IEVM2AnyGEOnRamp.sol";
import {IBaseOnRampRouter} from "../../interfaces/onRamp/IBaseOnRampRouter.sol";
import {IGERouter} from "../../interfaces/router/IGERouter.sol";

import {PoolCollector} from "../../pools/PoolCollector.sol";
import {MockOffRamp} from "../mocks/MockOffRamp.sol";
import "../onRamp/ge/EVM2EVMGEOnRampSetup.t.sol";

/// @notice #constructor
contract GERouter_constructor is EVM2EVMGEOnRampSetup {
  // Success

  function testSuccess() public {
    // typeAndVersion
    assertEq("GERouter 1.0.0", s_sourceRouter.typeAndVersion());

    // owner
    assertEq(OWNER, s_sourceRouter.owner());
  }
}

/// @notice #ccipSend
contract GERouter_ccipSend is EVM2EVMGEOnRampSetup {
  event Burned(address indexed sender, uint256 amount);

  // Success

  function testCCIPSendOneTokenSuccess_gas() public {
    vm.pauseGasMetering();
    address sourceToken1Address = s_sourceTokens[1];
    IERC20 sourceToken1 = IERC20(sourceToken1Address);
    GEConsumer.EVM2AnyGEMessage memory message = _generateEmptyMessage();

    sourceToken1.approve(address(s_sourceRouter), 2**64);

    message.tokensAndAmounts = new Common.EVMTokenAndAmount[](1);
    message.tokensAndAmounts[0].amount = 2**64;
    message.tokensAndAmounts[0].token = sourceToken1Address;
    message.feeToken = s_sourceTokens[0];

    uint256 expectedFee = s_sourceRouter.getFee(DEST_CHAIN_ID, message);

    uint256 balanceBefore = sourceToken1.balanceOf(OWNER);

    // Assert that the tokens are burned
    vm.expectEmit(false, false, false, true);
    emit Burned(address(s_onRamp), message.tokensAndAmounts[0].amount);

    GE.EVM2EVMGEMessage memory msgEvent = _messageToEvent(message, 1, 1, expectedFee);

    vm.expectEmit(false, false, false, true);
    emit CCIPSendRequested(msgEvent);

    vm.resumeGasMetering();
    bytes32 messageId = s_sourceRouter.ccipSend(DEST_CHAIN_ID, message);
    vm.pauseGasMetering();

    assertEq(msgEvent.messageId, messageId);
    // Assert the user balance is lowered by the tokensAndAmounts sent and the fee amount
    uint256 expectedBalance = balanceBefore - (message.tokensAndAmounts[0].amount);
    assertEq(expectedBalance, sourceToken1.balanceOf(OWNER));
    vm.resumeGasMetering();
  }

  function testNonLinkFeeTokenSuccess() public {
    GE.FeeUpdate[] memory feeUpdates = new GE.FeeUpdate[](1);
    feeUpdates[0] = GE.FeeUpdate({sourceFeeToken: s_sourceTokens[1], destChainId: DEST_CHAIN_ID, linkPerUnitGas: 1000});
    s_IFeeManager.updateFees(feeUpdates);

    GEConsumer.EVM2AnyGEMessage memory message = _generateEmptyMessage();
    message.feeToken = s_sourceTokens[1];

    IERC20(s_sourceTokens[1]).approve(address(s_sourceRouter), 2**64);

    s_sourceRouter.ccipSend(DEST_CHAIN_ID, message);
  }

  // Reverts

  function testUnsupportedDestinationChainReverts() public {
    GEConsumer.EVM2AnyGEMessage memory message = _generateEmptyMessage();
    uint64 wrongChain = DEST_CHAIN_ID + 1;

    vm.expectRevert(abi.encodeWithSelector(IBaseOnRampRouter.UnsupportedDestinationChain.selector, wrongChain));

    s_sourceRouter.ccipSend(wrongChain, message);
  }

  function testUnsupportedFeeTokenReverts() public {
    GEConsumer.EVM2AnyGEMessage memory message = _generateEmptyMessage();
    address wrongFeeToken = address(1);
    message.feeToken = wrongFeeToken;

    vm.expectRevert(
      abi.encodeWithSelector(IFeeManager.TokenOrChainNotSupported.selector, wrongFeeToken, DEST_CHAIN_ID)
    );

    s_sourceRouter.ccipSend(DEST_CHAIN_ID, message);
  }

  function testFeeTokenAmountTooLowReverts() public {
    GEConsumer.EVM2AnyGEMessage memory message = _generateEmptyMessage();
    IERC20(s_sourceTokens[0]).approve(address(s_sourceRouter), 0);

    vm.expectRevert("ERC20: transfer amount exceeds allowance");

    s_sourceRouter.ccipSend(DEST_CHAIN_ID, message);
  }
}

/// @notice #setOnRamp
contract GERouter_setOnRamp is EVM2EVMGEOnRampSetup {
  event OnRampSet(uint64 indexed chainId, IEVM2AnyGEOnRamp indexed onRamp);

  // Success

  // Asserts that setOnRamp changes the configured onramp. Also tests getOnRamp
  // and isChainSupported.
  function testSuccess() public {
    IEVM2AnyGEOnRamp onramp = IEVM2AnyGEOnRamp(address(1));
    uint64 chainId = 1337;
    IEVM2AnyGEOnRamp before = s_sourceRouter.getOnRamp(chainId);
    assertEq(address(0), address(before));
    assertFalse(s_sourceRouter.isChainSupported(chainId));

    vm.expectEmit(true, true, false, true);
    emit OnRampSet(chainId, onramp);

    s_sourceRouter.setOnRamp(chainId, onramp);
    IEVM2AnyGEOnRamp afterSet = s_sourceRouter.getOnRamp(chainId);
    assertEq(address(onramp), address(afterSet));
    assertTrue(s_sourceRouter.isChainSupported(chainId));
  }

  // Reverts

  // Asserts that setOnRamp reverts when the config was already set to
  // the same onRamp.
  function testAlreadySetReverts() public {
    vm.expectRevert(abi.encodeWithSelector(IGERouter.OnRampAlreadySet.selector, DEST_CHAIN_ID, s_onRamp));
    s_sourceRouter.setOnRamp(DEST_CHAIN_ID, s_onRamp);
  }

  // Asserts that setOnRamp can only be called by the owner.
  function testOnlyOwnerReverts() public {
    vm.stopPrank();
    vm.expectRevert("Only callable by owner");
    s_sourceRouter.setOnRamp(1337, IEVM2AnyGEOnRamp(address(1)));
  }
}

/// @notice #isChainSupported
contract GERouter_isChainSupported is EVM2EVMGEOnRampSetup {
  // Success
  function testSuccess() public {
    assertTrue(s_sourceRouter.isChainSupported(DEST_CHAIN_ID));
    assertFalse(s_sourceRouter.isChainSupported(DEST_CHAIN_ID + 1));
    assertFalse(s_sourceRouter.isChainSupported(0));
  }
}

/// @notice #getSupportedTokens
contract GERouter_getSupportedTokens is EVM2EVMGEOnRampSetup {
  // Success

  function testGetSupportedTokensSuccess() public {
    assertEq(s_sourceTokens, s_sourceRouter.getSupportedTokens(DEST_CHAIN_ID));
  }

  function testUnknownChainSuccess() public {
    address[] memory supportedTokens = s_sourceRouter.getSupportedTokens(DEST_CHAIN_ID + 10);
    assertEq(0, supportedTokens.length);
  }
}

/// @notice #addOffRamp
contract GERouter_addOffRamp is EVM2EVMGEOnRampSetup {
  address internal s_newOffRamp;

  event OffRampAdded(address indexed offRamp);

  function setUp() public virtual override {
    EVM2EVMGEOnRampSetup.setUp();

    s_newOffRamp = address(new MockOffRamp());
  }

  // Success

  function testSuccess() public {
    assertFalse(s_sourceRouter.isOffRamp(s_newOffRamp));
    uint256 lengthBefore = s_sourceRouter.getOffRamps().length;

    vm.expectEmit(true, false, false, true);
    emit OffRampAdded(s_newOffRamp);
    s_sourceRouter.addOffRamp(s_newOffRamp);

    assertTrue(s_sourceRouter.isOffRamp(s_newOffRamp));
    assertEq(lengthBefore + 1, s_sourceRouter.getOffRamps().length);
  }

  // Reverts

  function testOwnerReverts() public {
    changePrank(STRANGER);
    vm.expectRevert("Only callable by owner");
    s_sourceRouter.addOffRamp(s_newOffRamp);
  }

  function testAlreadyConfiguredReverts() public {
    address existingOffRamp = s_offRamps[0];
    vm.expectRevert(abi.encodeWithSelector(IAny2EVMOffRampRouter.AlreadyConfigured.selector, existingOffRamp));
    s_sourceRouter.addOffRamp(existingOffRamp);
  }

  function testZeroAddressReverts() public {
    vm.expectRevert(IAny2EVMOffRampRouter.InvalidAddress.selector);
    s_sourceRouter.addOffRamp(address(0));
  }
}

/// @notice #removeOffRamp
contract GERouter_removeOffRamp is EVM2EVMGEOnRampSetup {
  event OffRampRemoved(address indexed offRamp);

  // Success

  function testSuccess() public {
    uint256 lengthBefore = s_sourceRouter.getOffRamps().length;

    vm.expectEmit(true, false, false, true);
    emit OffRampRemoved(s_offRamps[0]);
    s_sourceRouter.removeOffRamp(s_offRamps[0]);

    assertFalse(s_sourceRouter.isOffRamp(s_offRamps[0]));
    assertEq(lengthBefore - 1, s_sourceRouter.getOffRamps().length);
  }

  // Reverts

  function testOwnerReverts() public {
    changePrank(STRANGER);
    vm.expectRevert("Only callable by owner");
    s_sourceRouter.removeOffRamp(s_offRamps[0]);
  }

  function testNoOffRampsReverts() public {
    s_sourceRouter.removeOffRamp(s_offRamps[0]);
    s_sourceRouter.removeOffRamp(s_offRamps[1]);

    assertEq(0, s_sourceRouter.getOffRamps().length);

    vm.expectRevert(IAny2EVMOffRampRouter.NoOffRampsConfigured.selector);
    s_sourceRouter.removeOffRamp(s_offRamps[0]);
  }

  function testOffRampNotAllowedReverts() public {
    address newRamp = address(1234678);
    vm.expectRevert(abi.encodeWithSelector(IAny2EVMOffRampRouter.OffRampNotAllowed.selector, newRamp));
    s_sourceRouter.removeOffRamp(newRamp);
  }
}

/// @notice #getOffRamps
contract GERouter_getOffRamps is EVM2EVMGEOnRampSetup {
  // Success
  function testGetOffRampsSuccess() public {
    address[] memory offRamps = s_sourceRouter.getOffRamps();
    assertEq(2, offRamps.length);
    assertEq(address(s_offRamps[0]), address(offRamps[0]));
    assertEq(address(s_offRamps[1]), address(offRamps[1]));
  }
}

/// @notice #isOffRamp
contract GERouter_isOffRamp is EVM2EVMGEOnRampSetup {
  // Success
  function testIsOffRampSuccess() public {
    assertTrue(s_sourceRouter.isOffRamp(s_offRamps[0]));
    assertTrue(s_sourceRouter.isOffRamp(s_offRamps[1]));
    assertFalse(s_sourceRouter.isOffRamp(address(1)));
  }
}
