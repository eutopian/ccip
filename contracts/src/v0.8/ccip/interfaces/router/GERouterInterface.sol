// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {GEConsumer} from "../../models/GEConsumer.sol";
import {BaseOnRampRouterInterface} from "../onRamp/BaseOnRampRouterInterface.sol";
import {EVM2EVMGEOnRampInterface} from "../onRamp/EVM2EVMGEOnRampInterface.sol";
import {Any2EVMOffRampRouterInterface} from "../offRamp/Any2EVMOffRampRouterInterface.sol";

interface GERouterInterface is BaseOnRampRouterInterface, Any2EVMOffRampRouterInterface {
  error OnRampAlreadySet(uint64 chainId, EVM2EVMGEOnRampInterface onRamp);

  event OnRampSet(uint64 indexed chainId, EVM2EVMGEOnRampInterface indexed onRamp);

  /**
   * @notice Request a message to be sent to the destination chain
   * @param destinationChainId The destination chain ID
   * @param message The message payload
   * @return The sequence number assigned to message
   */
  function ccipSend(uint64 destinationChainId, GEConsumer.EVM2AnyGEMessage calldata message) external returns (bytes32);

  function getFee(uint64 destinationChainId, GEConsumer.EVM2AnyGEMessage memory message)
    external
    view
    returns (uint256 fee);

  /**
   * @notice Set chainId => onRamp mapping
   * @dev only callable by owner
   * @param chainId destination chain ID
   * @param onRamp OnRamp to use for that destination chain
   */
  function setOnRamp(uint64 chainId, EVM2EVMGEOnRampInterface onRamp) external;

  /**
   * @notice Gets the current OnRamp for the specified chain ID
   * @param chainId Chain ID to get ramp details for
   * @return onRamp
   */
  function getOnRamp(uint64 chainId) external view returns (EVM2EVMGEOnRampInterface);
}
