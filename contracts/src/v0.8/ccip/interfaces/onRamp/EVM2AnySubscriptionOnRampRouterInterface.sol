// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {CCIP, IERC20} from "../../models/Models.sol";
import {BaseOnRampRouterInterface} from "./BaseOnRampRouterInterface.sol";
import {EVM2EVMSubscriptionOnRampInterface} from "./EVM2EVMSubscriptionOnRampInterface.sol";

interface EVM2AnySubscriptionOnRampRouterInterface is BaseOnRampRouterInterface {
  error OnRampAlreadySet(uint256 chainId, EVM2EVMSubscriptionOnRampInterface onRamp);
  error FundingTooLow(address sender);
  error OnlyCallableByFeeAdmin();

  event OnRampSet(uint256 indexed chainId, EVM2EVMSubscriptionOnRampInterface indexed onRamp);
  event FeeSet(uint96);
  event SubscriptionFunded(address indexed sender, uint256 amount);
  event SubscriptionUnfunded(address indexed sender, uint256 amount);

  struct RouterConfig {
    // The fee amount to be charged for each ccipSend
    uint96 fee;
    // The token to be charged for each ccipSend
    IERC20 feeToken;
    // The address that can change the fee properties of the router
    address feeAdmin;
  }

  /**
   * @notice Request a message to be sent to the destination chain
   * @param destinationChainId The destination chain ID
   * @param message The message payload
   * @return The sequence number assigned to message
   */
  function ccipSend(uint256 destinationChainId, CCIP.EVM2AnySubscriptionMessage calldata message)
    external
    returns (uint64);

  /**
   * @notice Set chainId => onRamp mapping
   * @dev only callable by owner
   * @param chainId destination chain ID
   * @param onRamp OnRamp to use for that destination chain
   */
  function setOnRamp(uint256 chainId, EVM2EVMSubscriptionOnRampInterface onRamp) external;

  /**
   * @notice Gets the current OnRamp for the specified chain ID
   * @param chainId Chain ID to get ramp details for
   * @return onRamp
   */
  function getOnRamp(uint256 chainId) external view returns (EVM2EVMSubscriptionOnRampInterface);

  /**
   * @notice Changes the fee amount
   * @param newFee The new fee amount
   * @dev only callable by the s_feeAdmin
   */
  function setFee(uint96 newFee) external;

  /**
   * @notice Gets the fee amount
   */
  function getFee() external returns (uint96);

  /**
   * @notice Adds funding to a subscription
   * @param amount The amount to add to the subscription
   */
  function fundSubscription(uint256 amount) external;

  /**
   * @notice Removes funding from a subscription
   * @param amount The amount to withdrawal from the subscription
   */
  function unfundSubscription(uint256 amount) external;

  /**
   * @notice Gets the balance for a given address
   * @param sender The address for which to get the balance
   */
  function getBalance(address sender) external returns (uint256 balance);
}
