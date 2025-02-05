// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {IRebalancer} from "../interfaces/IRebalancer.sol";
import {IBridgeAdapter} from "../interfaces/IBridge.sol";

import {LockReleaseTokenPool} from "../../ccip/pools/LockReleaseTokenPool.sol";
import {Rebalancer} from "../Rebalancer.sol";
import {MockL1BridgeAdapter} from "./mocks/MockBridgeAdapter.sol";
import {RebalancerBaseTest} from "./RebalancerBaseTest.t.sol";
import {RebalancerHelper} from "./helpers/RebalancerHelper.sol";

contract RebalancerSetup is RebalancerBaseTest {
  event LiquidityTransferred(
    uint64 indexed ocrSeqNum,
    uint64 indexed fromChainSelector,
    uint64 indexed toChainSelector,
    address to,
    uint256 amount,
    bytes bridgeSpecificPayload,
    bytes bridgeReturnData
  );
  event FinalizationFailed(
    uint64 indexed ocrSeqNum,
    uint64 indexed remoteChainSelector,
    bytes bridgeSpecificData,
    bytes reason
  );
  event LiquidityAddedToContainer(address indexed provider, uint256 indexed amount);
  event LiquidityRemovedFromContainer(address indexed remover, uint256 indexed amount);
  // Liquidity container event
  event LiquidityAdded(address indexed provider, uint256 indexed amount);
  event LiquidityRemoved(address indexed remover, uint256 indexed amount);

  error NonceAlreadyUsed(uint256 nonce);

  RebalancerHelper internal s_rebalancer;
  LockReleaseTokenPool internal s_lockReleaseTokenPool;
  MockL1BridgeAdapter internal s_bridgeAdapter;

  function setUp() public override {
    RebalancerBaseTest.setUp();

    s_bridgeAdapter = new MockL1BridgeAdapter(s_l1Token);
    s_lockReleaseTokenPool = new LockReleaseTokenPool(s_l1Token, new address[](0), address(1), true, address(123));
    s_rebalancer = new RebalancerHelper(s_l1Token, i_localChainSelector, s_lockReleaseTokenPool);

    s_lockReleaseTokenPool.setRebalancer(address(s_rebalancer));
  }
}

contract Rebalancer_addLiquidity is RebalancerSetup {
  function test_addLiquiditySuccess() external {
    address caller = STRANGER;
    vm.startPrank(caller);

    uint256 amount = 12345679;
    deal(address(s_l1Token), caller, amount);

    s_l1Token.approve(address(s_rebalancer), amount);

    vm.expectEmit();
    emit LiquidityAddedToContainer(caller, amount);

    s_rebalancer.addLiquidity(amount);

    assertEq(s_l1Token.balanceOf(address(s_lockReleaseTokenPool)), amount);
  }
}

contract Rebalancer_removeLiquidity is RebalancerSetup {
  function test_removeLiquiditySuccess() external {
    uint256 amount = 12345679;
    deal(address(s_l1Token), address(s_lockReleaseTokenPool), amount);

    vm.expectEmit();
    emit LiquidityRemovedFromContainer(OWNER, amount);

    s_rebalancer.removeLiquidity(amount);

    assertEq(s_l1Token.balanceOf(address(s_rebalancer)), 0);
  }

  function test_InsufficientLiquidityReverts() external {
    uint256 balance = 923;
    uint256 requested = balance + 1;

    deal(address(s_l1Token), address(s_lockReleaseTokenPool), balance);

    vm.expectRevert(abi.encodeWithSelector(Rebalancer.InsufficientLiquidity.selector, requested, balance));

    s_rebalancer.removeLiquidity(requested);
  }

  function test_OnlyOwnerReverts() external {
    vm.stopPrank();

    vm.expectRevert("Only callable by owner");

    s_rebalancer.removeLiquidity(123);
  }
}

contract Rebalancer__report is RebalancerSetup {
  function test_EmptyReportReverts() external {
    IRebalancer.LiquidityInstructions memory instructions = IRebalancer.LiquidityInstructions({
      sendLiquidityParams: new IRebalancer.SendLiquidityParams[](0),
      receiveLiquidityParams: new IRebalancer.ReceiveLiquidityParams[](0)
    });

    vm.expectRevert(Rebalancer.EmptyReport.selector);

    s_rebalancer.report(abi.encode(instructions), 123);
  }
}

contract Rebalancer_rebalanceLiquidity is RebalancerSetup {
  function test_rebalanceLiquiditySuccess() external {
    uint256 amount = 12345679;
    deal(address(s_l1Token), address(s_lockReleaseTokenPool), amount);

    Rebalancer.CrossChainRebalancerArgs[] memory args = new Rebalancer.CrossChainRebalancerArgs[](1);
    args[0] = IRebalancer.CrossChainRebalancerArgs({
      remoteRebalancer: address(s_rebalancer),
      localBridge: s_bridgeAdapter,
      remoteToken: address(s_l2Token),
      remoteChainSelector: i_remoteChainSelector,
      enabled: true
    });
    s_rebalancer.setCrossChainRebalancers(args);

    vm.expectEmit();
    emit Transfer(address(s_lockReleaseTokenPool), address(s_rebalancer), amount);

    vm.expectEmit();
    emit Approval(address(s_rebalancer), address(s_bridgeAdapter), amount);

    vm.expectEmit();
    emit Transfer(address(s_rebalancer), address(s_bridgeAdapter), amount);

    vm.expectEmit();
    bytes memory encodedNonce = abi.encode(uint256(1));
    emit LiquidityTransferred(
      type(uint64).max,
      i_localChainSelector,
      i_remoteChainSelector,
      address(s_rebalancer),
      amount,
      bytes(""),
      encodedNonce
    );

    s_rebalancer.rebalanceLiquidity(i_remoteChainSelector, amount, 0, bytes(""));

    assertEq(s_l1Token.balanceOf(address(s_rebalancer)), 0);
    assertEq(s_l1Token.balanceOf(address(s_bridgeAdapter)), amount);
    assertEq(s_l1Token.allowance(address(s_rebalancer), address(s_bridgeAdapter)), 0);
  }

  /// @notice this test sets up a circular system where the liquidity container of
  /// the local Liquidity manager is the bridge adapter of the remote liquidity manager
  /// and the other way around for the remote liquidity manager. This allows us to
  /// rebalance funds between the two liquidity managers on the same chain.
  function test_rebalanceBetweenPoolsSuccess() external {
    uint256 amount = 12345670;

    s_rebalancer = new RebalancerHelper(s_l1Token, i_localChainSelector, s_bridgeAdapter);

    MockL1BridgeAdapter mockRemoteBridgeAdapter = new MockL1BridgeAdapter(s_l1Token);
    Rebalancer mockRemoteRebalancer = new Rebalancer(s_l1Token, i_remoteChainSelector, mockRemoteBridgeAdapter);

    Rebalancer.CrossChainRebalancerArgs[] memory args = new Rebalancer.CrossChainRebalancerArgs[](1);
    args[0] = IRebalancer.CrossChainRebalancerArgs({
      remoteRebalancer: address(mockRemoteRebalancer),
      localBridge: mockRemoteBridgeAdapter,
      remoteToken: address(s_l1Token),
      remoteChainSelector: i_remoteChainSelector,
      enabled: true
    });

    s_rebalancer.setCrossChainRebalancers(args);

    args[0] = IRebalancer.CrossChainRebalancerArgs({
      remoteRebalancer: address(s_rebalancer),
      localBridge: s_bridgeAdapter,
      remoteToken: address(s_l1Token),
      remoteChainSelector: i_localChainSelector,
      enabled: true
    });

    mockRemoteRebalancer.setCrossChainRebalancers(args);

    deal(address(s_l1Token), address(s_bridgeAdapter), amount);

    s_rebalancer.rebalanceLiquidity(i_remoteChainSelector, amount, 0, bytes(""));

    assertEq(s_l1Token.balanceOf(address(s_bridgeAdapter)), 0);
    assertEq(s_l1Token.balanceOf(address(mockRemoteBridgeAdapter)), amount);
    assertEq(s_l1Token.allowance(address(s_rebalancer), address(s_bridgeAdapter)), 0);

    // attach a bridge fee and see the relevant adapter's ether balance change.
    // the bridge fee is sent along with the sendERC20 call.
    uint256 bridgeFee = 123;
    vm.deal(address(mockRemoteRebalancer), bridgeFee);
    mockRemoteRebalancer.rebalanceLiquidity(i_localChainSelector, amount, bridgeFee, bytes(""));

    assertEq(s_l1Token.balanceOf(address(s_bridgeAdapter)), amount);
    assertEq(s_l1Token.balanceOf(address(mockRemoteBridgeAdapter)), 0);
    assertEq(address(s_bridgeAdapter).balance, bridgeFee);

    // Assert partial rebalancing works correctly
    s_rebalancer.rebalanceLiquidity(i_remoteChainSelector, amount / 2, 0, bytes(""));

    assertEq(s_l1Token.balanceOf(address(s_bridgeAdapter)), amount / 2);
    assertEq(s_l1Token.balanceOf(address(mockRemoteBridgeAdapter)), amount / 2);
  }

  function test_rebalanceBetweenPoolsSuccess_AlreadyFinalized() external {
    uint256 amount = 12345670;
    // set up a rebalancer on another chain, an "L2".
    // note we use the L1 bridge adapter because it has the reverting logic
    // when finalization is already done.
    MockL1BridgeAdapter remoteBridgeAdapter = new MockL1BridgeAdapter(s_l2Token);
    LockReleaseTokenPool remotePool = new LockReleaseTokenPool(
      s_l2Token,
      new address[](0),
      address(1),
      true,
      address(123)
    );
    Rebalancer remoteRebalancer = new Rebalancer(s_l2Token, i_remoteChainSelector, remotePool);

    // set rebalancer role on the pool.
    remotePool.setRebalancer(address(remoteRebalancer));

    // set up the cross chain rebalancer on "L1".
    Rebalancer.CrossChainRebalancerArgs[] memory args = new Rebalancer.CrossChainRebalancerArgs[](1);
    args[0] = IRebalancer.CrossChainRebalancerArgs({
      remoteRebalancer: address(remoteRebalancer),
      localBridge: s_bridgeAdapter,
      remoteToken: address(s_l2Token),
      remoteChainSelector: i_remoteChainSelector,
      enabled: true
    });

    s_rebalancer.setCrossChainRebalancers(args);

    // set up the cross chain rebalancer on "L2".
    args[0] = IRebalancer.CrossChainRebalancerArgs({
      remoteRebalancer: address(s_rebalancer),
      localBridge: remoteBridgeAdapter,
      remoteToken: address(s_l1Token),
      remoteChainSelector: i_localChainSelector,
      enabled: true
    });

    remoteRebalancer.setCrossChainRebalancers(args);

    // deal some L1 tokens to the L1 bridge adapter so that it can send them to the rebalancer
    // when the withdrawal gets finalized.
    deal(address(s_l1Token), address(s_bridgeAdapter), amount);
    // deal some L2 tokens to the remote token pool so that we can withdraw it when we rebalance.
    deal(address(s_l2Token), address(remotePool), amount);

    uint256 nonce = 1;
    uint64 maxSeqNum = type(uint64).max;
    bytes memory bridgeSendReturnData = abi.encode(nonce);
    bytes memory bridgeSpecificPayload = bytes("");
    vm.expectEmit();
    emit LiquidityRemoved(address(remoteRebalancer), amount);
    vm.expectEmit();
    emit LiquidityTransferred(
      maxSeqNum,
      i_remoteChainSelector,
      i_localChainSelector,
      address(s_rebalancer),
      amount,
      bridgeSpecificPayload,
      bridgeSendReturnData
    );
    remoteRebalancer.rebalanceLiquidity(i_localChainSelector, amount, 0, bridgeSpecificPayload);

    // available liquidity has been moved to the remote bridge adapter from the token pool.
    assertEq(s_l2Token.balanceOf(address(remoteBridgeAdapter)), amount, "remoteBridgeAdapter balance");
    assertEq(s_l2Token.balanceOf(address(remotePool)), 0, "remotePool balance");

    // finalize manually on the L1 bridge adapter.
    // this should transfer the funds to the rebalancer.
    bytes memory finalizationData = abi.encode(amount, nonce);
    s_bridgeAdapter.finalizeWithdrawERC20(address(0), address(s_rebalancer), finalizationData);

    // available balance on the L1 bridge adapter has been moved to the rebalancer.
    assertEq(s_l1Token.balanceOf(address(s_rebalancer)), amount, "rebalancer balance 1");
    assertEq(s_l1Token.balanceOf(address(s_bridgeAdapter)), 0, "bridgeAdapter balance");

    // try to finalize on L1 again
    bytes memory revertData = abi.encodeWithSelector(NonceAlreadyUsed.selector, nonce);
    vm.expectEmit();
    emit FinalizationFailed(maxSeqNum, i_remoteChainSelector, finalizationData, revertData);
    vm.expectEmit();
    emit LiquidityAdded(address(s_rebalancer), amount);
    vm.expectEmit();
    emit LiquidityTransferred(
      maxSeqNum,
      i_remoteChainSelector,
      i_localChainSelector,
      address(s_rebalancer),
      amount,
      finalizationData,
      bytes("")
    );
    s_rebalancer.receiveLiquidity(i_remoteChainSelector, amount, finalizationData);

    // available balance on the rebalancer has been injected into the token pool.
    assertEq(s_l1Token.balanceOf(address(s_rebalancer)), 0, "rebalancer balance 2");
    assertEq(s_l1Token.balanceOf(address(s_lockReleaseTokenPool)), amount, "lockReleaseTokenPool balance");
  }

  // Reverts

  function test_InsufficientLiquidityReverts() external {
    uint256 amount = 1245;

    vm.expectRevert(abi.encodeWithSelector(Rebalancer.InsufficientLiquidity.selector, amount, 0));

    s_rebalancer.rebalanceLiquidity(0, amount, 0, bytes(""));
  }

  function test_InvalidRemoteChainReverts() external {
    uint256 amount = 12345679;
    deal(address(s_l1Token), address(s_lockReleaseTokenPool), amount);

    vm.expectRevert(abi.encodeWithSelector(Rebalancer.InvalidRemoteChain.selector, i_remoteChainSelector));

    s_rebalancer.rebalanceLiquidity(i_remoteChainSelector, amount, 0, bytes(""));
  }
}

contract Rebalancer_setCrossChainRebalancer is RebalancerSetup {
  event CrossChainRebalancerSet(
    uint64 indexed remoteChainSelector,
    IBridgeAdapter localBridge,
    address remoteToken,
    address remoteRebalancer,
    bool enabled
  );

  function test_setCrossChainRebalancerSuccess() external {
    address newRebalancer = address(23892423);
    uint64 remoteChainSelector = 12301293;

    uint64[] memory supportedChains = s_rebalancer.getSupportedDestChains();
    assertEq(supportedChains.length, 0);

    Rebalancer.CrossChainRebalancerArgs[] memory args = new Rebalancer.CrossChainRebalancerArgs[](1);
    args[0] = IRebalancer.CrossChainRebalancerArgs({
      remoteRebalancer: newRebalancer,
      localBridge: s_bridgeAdapter,
      remoteToken: address(190490124908),
      remoteChainSelector: remoteChainSelector,
      enabled: true
    });

    vm.expectEmit();
    emit CrossChainRebalancerSet(
      remoteChainSelector,
      args[0].localBridge,
      args[0].remoteToken,
      newRebalancer,
      args[0].enabled
    );

    s_rebalancer.setCrossChainRebalancers(args);

    assertEq(s_rebalancer.getCrossChainRebalancer(remoteChainSelector).remoteRebalancer, newRebalancer);

    Rebalancer.CrossChainRebalancerArgs[] memory got = s_rebalancer.getAllCrossChainRebalancers();
    assertEq(got.length, 1);
    assertEq(got[0].remoteRebalancer, args[0].remoteRebalancer);
    assertEq(address(got[0].localBridge), address(args[0].localBridge));
    assertEq(got[0].remoteToken, args[0].remoteToken);
    assertEq(got[0].remoteChainSelector, args[0].remoteChainSelector);
    assertEq(got[0].enabled, args[0].enabled);

    supportedChains = s_rebalancer.getSupportedDestChains();
    assertEq(supportedChains.length, 1);
    assertEq(supportedChains[0], remoteChainSelector);

    address anotherRebalancer = address(123);
    args[0].remoteRebalancer = anotherRebalancer;

    vm.expectEmit();
    emit CrossChainRebalancerSet(
      remoteChainSelector,
      args[0].localBridge,
      args[0].remoteToken,
      anotherRebalancer,
      args[0].enabled
    );

    s_rebalancer.setCrossChainRebalancer(args[0]);

    assertEq(s_rebalancer.getCrossChainRebalancer(remoteChainSelector).remoteRebalancer, anotherRebalancer);

    supportedChains = s_rebalancer.getSupportedDestChains();
    assertEq(supportedChains.length, 1);
    assertEq(supportedChains[0], remoteChainSelector);
  }

  function test_ZeroChainSelectorReverts() external {
    Rebalancer.CrossChainRebalancerArgs memory arg = IRebalancer.CrossChainRebalancerArgs({
      remoteRebalancer: address(9),
      localBridge: s_bridgeAdapter,
      remoteToken: address(190490124908),
      remoteChainSelector: 0,
      enabled: true
    });

    vm.expectRevert(Rebalancer.ZeroChainSelector.selector);

    s_rebalancer.setCrossChainRebalancer(arg);
  }

  function test_ZeroAddressReverts() external {
    Rebalancer.CrossChainRebalancerArgs memory arg = IRebalancer.CrossChainRebalancerArgs({
      remoteRebalancer: address(0),
      localBridge: s_bridgeAdapter,
      remoteToken: address(190490124908),
      remoteChainSelector: 123,
      enabled: true
    });

    vm.expectRevert(Rebalancer.ZeroAddress.selector);

    s_rebalancer.setCrossChainRebalancer(arg);

    arg.remoteRebalancer = address(9);
    arg.localBridge = IBridgeAdapter(address(0));

    vm.expectRevert(Rebalancer.ZeroAddress.selector);

    s_rebalancer.setCrossChainRebalancer(arg);

    arg.localBridge = s_bridgeAdapter;
    arg.remoteToken = address(0);

    vm.expectRevert(Rebalancer.ZeroAddress.selector);

    s_rebalancer.setCrossChainRebalancer(arg);
  }

  function test_OnlyOwnerReverts() external {
    vm.stopPrank();

    vm.expectRevert("Only callable by owner");

    // Test the entrypoint that takes a list
    s_rebalancer.setCrossChainRebalancers(new Rebalancer.CrossChainRebalancerArgs[](0));

    vm.expectRevert("Only callable by owner");

    // Test the entrypoint that takes a single item
    s_rebalancer.setCrossChainRebalancer(
      IRebalancer.CrossChainRebalancerArgs({
        remoteRebalancer: address(9),
        localBridge: s_bridgeAdapter,
        remoteToken: address(190490124908),
        remoteChainSelector: 124,
        enabled: true
      })
    );
  }
}

contract Rebalancer_setLocalLiquidityContainer is RebalancerSetup {
  event LiquidityContainerSet(address indexed newLiquidityContainer);

  function test_setLocalLiquidityContainerSuccess() external {
    LockReleaseTokenPool newPool = new LockReleaseTokenPool(
      s_l1Token,
      new address[](0),
      address(1),
      true,
      address(123)
    );

    vm.expectEmit();
    emit LiquidityContainerSet(address(newPool));

    s_rebalancer.setLocalLiquidityContainer(newPool);

    assertEq(s_rebalancer.getLocalLiquidityContainer(), address(newPool));
  }

  function test_OnlyOwnerReverts() external {
    vm.stopPrank();

    vm.expectRevert("Only callable by owner");

    s_rebalancer.setLocalLiquidityContainer(LockReleaseTokenPool(address(1)));
  }
}
