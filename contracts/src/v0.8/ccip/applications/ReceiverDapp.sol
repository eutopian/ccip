// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import "../../vendor/SafeERC20.sol";
import "../../interfaces/TypeAndVersionInterface.sol";
import "../offRamp/interfaces/Any2EVMTollOffRampRouterInterface.sol";
import "./interfaces/CrossChainMessageReceiverInterface.sol";

/**
 * @notice Application contract for receiving messages from the OffRamp on behalf of an EOA
 */
contract ReceiverDapp is CrossChainMessageReceiverInterface, TypeAndVersionInterface {
  using SafeERC20 for IERC20;

  string public constant override typeAndVersion = "ReceiverDapp 1.0.0";

  Any2EVMTollOffRampRouterInterface public immutable ROUTER;

  address s_manager;

  error InvalidDeliverer(address deliverer);

  constructor(Any2EVMTollOffRampRouterInterface router) {
    ROUTER = router;
    s_manager = msg.sender;
  }

  function getSubscriptionManager() external view returns (address) {
    return s_manager;
  }

  /**
   * @notice Called by the OffRamp, this function receives a message and forwards
   * the tokens sent with it to the designated EOA
   * @param message CCIP Message
   */
  function ccipReceive(CCIP.Any2EVMTollMessage calldata message) external override onlyRouter {
    handleMessage(message.data, message.tokens, message.amounts);
  }

  /**
   * @notice Called by the OffRamp, this function receives a message and forwards
   * the tokens sent with it to the designated EOA
   * @param message CCIP Message
   */
  function ccipReceive(CCIP.Any2EVMSubscriptionMessage calldata message) external override onlyRouter {
    handleMessage(message.data, message.tokens, message.amounts);
  }

  function handleMessage(
    bytes memory data,
    IERC20[] memory tokens,
    uint256[] memory amounts
  ) internal {
    (
      ,
      /* address originalSender */
      address destinationAddress
    ) = abi.decode(data, (address, address));
    for (uint256 i = 0; i < tokens.length; ++i) {
      uint256 amount = amounts[i];
      if (destinationAddress != address(0) && amount != 0) {
        tokens[i].transfer(destinationAddress, amount);
      }
    }
  }

  /**
   * @dev only calls from the set router are accepted.
   */
  modifier onlyRouter() {
    if (msg.sender != address(ROUTER)) revert InvalidDeliverer(msg.sender);
    _;
  }
}
