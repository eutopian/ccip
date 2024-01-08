// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {IBridgeAdapter} from "./interfaces/IBridge.sol";
import {ILiquidityContainer} from "./interfaces/ILiquidityContainer.sol";
import {ILiquidityManager} from "./interfaces/ILiquidityManager.sol";

import {OCR3Base} from "../../ocr/OCR3Base.sol";

import {IERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/utils/SafeERC20.sol";

/// @notice Liquidity manager for a single token over multiple chains.
/// @dev This contract is designed to be used with the LockReleaseTokenPool contract but
/// isn't constraint to it. It can be used with any contract that implements the ILiquidityContainer
/// interface.
/// @dev The OCR3 DON should only be able to transfer funds to other pre-approved contracts
/// on other chains. Under no circumstances should it be able to transfer funds to arbitrary
/// addresses. The owner is therefore in full control of the funds in this contract, not the DON.
/// This is a security feature. The worst that can happen is that the DON can lock up funds in
/// bridges, but it can't steal them.
/// @dev References to local mean logic on the same chain as this contract is deployed on.
/// References to remote mean logic on other chains.
contract LiquidityManager is ILiquidityManager, OCR3Base {
  using SafeERC20 for IERC20;

  error ZeroAddress();
  error InvalidRemoteChain(uint64 chainSelector);
  error ZeroChainSelector();
  error InsufficientLiquidity(uint256 requested, uint256 available);

  event LiquidityTransferred(
    uint64 indexed fromChainSelector,
    uint64 indexed toChainSelector,
    address indexed to,
    uint256 amount
  );
  event LiquidityAdded(address indexed provider, uint256 indexed amount);
  event LiquidityRemoved(address indexed remover, uint256 indexed amount);

  struct CrossChainLiquidityManagerArgs {
    address remoteLiquidityManager;
    IBridgeAdapter localBridge;
    address remoteToken;
    uint64 remoteChainSelector;
    bool enabled;
  }

  struct CrossChainLiquidityManager {
    address remoteLiquidityManager;
    IBridgeAdapter localBridge;
    address remoteToken;
    bool enabled;
  }

  // solhint-disable-next-line chainlink-solidity/all-caps-constant-storage-variables
  string public constant override typeAndVersion = "LiquidityManager 1.3.0-dev";

  /// @notice The token that this pool manages liquidity for.
  IERC20 public immutable i_localToken;

  /// @notice The chain selector belonging to the chain this pool is deployed on.
  uint64 internal immutable i_localChainSelector;

  /// @notice Mapping of chain selector to liquidity container on other chains
  mapping(uint64 chainSelector => CrossChainLiquidityManager) private s_crossChainLiquidityManager;

  uint64[] private s_supportedDestChains;

  /// @notice The liquidity container on the local chain
  /// @dev In the case of CCIP, this would be the token pool.
  ILiquidityContainer private s_localLiquidityContainer;

  constructor(IERC20 token, uint64 localChainSelector, ILiquidityContainer localLiquidityContainer) OCR3Base() {
    if (localChainSelector == 0) {
      revert ZeroChainSelector();
    }

    if (address(token) == address(0)) {
      revert ZeroAddress();
    }
    i_localToken = token;
    i_localChainSelector = localChainSelector;
    s_localLiquidityContainer = localLiquidityContainer;
  }

  // ================================================================
  // │                    Liquidity management                      │
  // ================================================================

  /// @inheritdoc ILiquidityManager
  function getLiquidity() public view returns (uint256 currentLiquidity) {
    return i_localToken.balanceOf(address(s_localLiquidityContainer));
  }

  /// @notice Adds liquidity to the multi-chain system.
  /// @dev Anyone can call this function, but anyone other than the owner should regard
  /// adding liquidity as a donation to the system, as there is no way to get it out.
  /// This function is open to anyone to be able to quickly add funds to the system
  /// without having to go through potentially complicated multisig schemes to do it from
  /// the owner address.
  function addLiquidity(uint256 amount) external {
    i_localToken.safeTransferFrom(msg.sender, address(this), amount);

    // Make sure this is tether compatible, as they have strange approval requirements
    // Should be good since all approvals are always immediately used.
    i_localToken.approve(address(s_localLiquidityContainer), amount);
    s_localLiquidityContainer.provideLiquidity(amount);

    emit LiquidityAdded(msg.sender, amount);
  }

  /// @notice Removes liquidity from the system and sends it to the caller, so the owner.
  /// @dev Only the owner can call this function.
  function removeLiquidity(uint256 amount) external onlyOwner {
    uint256 currentBalance = i_localToken.balanceOf(address(s_localLiquidityContainer));
    if (currentBalance < amount) {
      revert InsufficientLiquidity(amount, currentBalance);
    }

    s_localLiquidityContainer.withdrawLiquidity(amount);
    i_localToken.safeTransfer(msg.sender, amount);

    emit LiquidityRemoved(msg.sender, amount);
  }

  /// @notice Transfers liquidity to another chain.
  /// @dev This function is a public version of the internal _rebalanceLiquidity function.
  /// to allow the owner to also initiate a rebalancing when needed.
  function rebalanceLiquidity(uint64 chainSelector, uint256 amount) external onlyOwner {
    _rebalanceLiquidity(chainSelector, amount);
  }

  /// @notice Transfers liquidity to another chain.
  /// @dev Called by both the owner and the DON.
  function _rebalanceLiquidity(uint64 chainSelector, uint256 amount) internal {
    uint256 currentBalance = getLiquidity();
    if (currentBalance < amount) {
      revert InsufficientLiquidity(amount, currentBalance);
    }

    CrossChainLiquidityManager memory remoteLiqManager = s_crossChainLiquidityManager[chainSelector];

    if (!remoteLiqManager.enabled) {
      revert InvalidRemoteChain(chainSelector);
    }

    uint256 nativeBridgeFee = remoteLiqManager.localBridge.getBridgeFeeInNative();

    // Could be optimized by withdrawing once and then sending to all destinations
    s_localLiquidityContainer.withdrawLiquidity(amount);
    i_localToken.approve(address(remoteLiqManager.localBridge), amount);

    remoteLiqManager.localBridge.sendERC20{value: nativeBridgeFee}(
      address(i_localToken),
      remoteLiqManager.remoteToken,
      remoteLiqManager.remoteLiquidityManager,
      amount
    );

    emit LiquidityTransferred(i_localChainSelector, chainSelector, remoteLiqManager.remoteLiquidityManager, amount);
  }

  function _receiveLiquidity(uint64 remoteChainSelector, uint256 amount, bytes memory bridgeData) internal {
    // TODO
  }

  function _report(bytes calldata report, uint64) internal override {
    ILiquidityManager.LiquidityInstructions memory instructions = abi.decode(
      report,
      (ILiquidityManager.LiquidityInstructions)
    );

    uint256 sendInstructions = instructions.sendLiquidityParams.length;
    for (uint256 i = 0; i < sendInstructions; ++i) {
      _rebalanceLiquidity(
        instructions.sendLiquidityParams[i].remoteChainSelector,
        instructions.sendLiquidityParams[i].amount
      );
    }

    uint256 receiveInstructions = instructions.receiveLiquidityParams.length;
    for (uint256 i = 0; i < receiveInstructions; ++i) {
      _receiveLiquidity(
        instructions.receiveLiquidityParams[i].remoteChainSelector,
        instructions.receiveLiquidityParams[i].amount,
        instructions.receiveLiquidityParams[i].bridgeData
      );
    }

    // todo emit?
  }

  // ================================================================
  // │                           Config                             │
  // ================================================================

  /// @notice Gets the cross chain liquidity manager
  function getCrossChainLiquidityManager(
    uint64 chainSelector
  ) external view returns (CrossChainLiquidityManager memory) {
    return s_crossChainLiquidityManager[chainSelector];
  }

  /// @notice Gets all cross chain liquidity managers
  /// @dev We don't care too much about gas since this function is intended for offchain usage.
  function getAllCrossChainLiquidityMangers() external view returns (CrossChainLiquidityManagerArgs[] memory) {
    CrossChainLiquidityManagerArgs[] memory managers;
    for (uint256 i = 0; i < s_supportedDestChains.length; ++i) {
      uint64 chainSelector = s_supportedDestChains[i];
      CrossChainLiquidityManager memory currentManager = s_crossChainLiquidityManager[chainSelector];
      managers[i] = CrossChainLiquidityManagerArgs({
        remoteLiquidityManager: currentManager.remoteLiquidityManager,
        localBridge: currentManager.localBridge,
        remoteToken: currentManager.remoteToken,
        remoteChainSelector: chainSelector,
        enabled: currentManager.enabled
      });
    }

    return managers;
  }

  /// @notice Sets a list of cross chain liquidity managers.
  /// @dev Will update the list of supported dest chains if the chain is new.
  function setCrossChainLiquidityManager(
    CrossChainLiquidityManagerArgs[] calldata crossChainLiquidityManagers
  ) external onlyOwner {
    for (uint256 i = 0; i < crossChainLiquidityManagers.length; ++i) {
      setCrossChainLiquidityManager(crossChainLiquidityManagers[i]);
    }
  }

  /// @notice Sets a single cross chain liquidity manager.
  /// @dev Will update the list of supported dest chains if the chain is new.
  function setCrossChainLiquidityManager(
    CrossChainLiquidityManagerArgs calldata crossChainLiqManager
  ) public onlyOwner {
    if (crossChainLiqManager.remoteChainSelector == 0) {
      revert ZeroChainSelector();
    }

    if (
      crossChainLiqManager.remoteLiquidityManager == address(0) ||
      address(crossChainLiqManager.localBridge) == address(0) ||
      crossChainLiqManager.remoteToken == address(0)
    ) {
      revert ZeroAddress();
    }

    // If the destination chain is new, add it to the list of supported chains
    if (s_crossChainLiquidityManager[crossChainLiqManager.remoteChainSelector].remoteToken == address(0)) {
      s_supportedDestChains.push(crossChainLiqManager.remoteChainSelector);
    }

    s_crossChainLiquidityManager[crossChainLiqManager.remoteChainSelector] = CrossChainLiquidityManager({
      remoteLiquidityManager: crossChainLiqManager.remoteLiquidityManager,
      localBridge: crossChainLiqManager.localBridge,
      remoteToken: crossChainLiqManager.remoteToken,
      enabled: crossChainLiqManager.enabled
    });
  }

  /// @notice Gets the local liquidity container.
  function getLocalLiquidityContainer() external view returns (address) {
    return address(s_localLiquidityContainer);
  }

  /// @notice Sets the local liquidity container.
  /// @dev Only the owner can call this function.
  function setLocalLiquidityContainer(ILiquidityContainer localLiquidityContainer) external onlyOwner {
    s_localLiquidityContainer = localLiquidityContainer;
  }
}
