// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import "../utils/CCIP.sol";
import "../interfaces/PoolInterface.sol";
import "../interfaces/AFNInterface.sol";

interface OnRampInterface {
  error TokenMismatch();
  error MessageTooLarge(uint256 maxSize, uint256 actualSize);
  error UnsupportedNumberOfTokens();
  error UnsupportedToken(IERC20 token);
  error UnsupportedFeeToken(IERC20 token);
  error SenderNotAllowed(address sender);
  error UnsupportedDestinationChain(uint256 destinationChainId);
  error MustBeCalledByRouter();

  event CrossChainSendRequested(CCIP.Message message);
  event AllowlistEnabledSet(bool enabled);
  event AllowlistSet(address[] allowlist);
  event NewTokenBucketConstructed(uint256 rate, uint256 capacity, bool full);
  event OnRampConfigSet(OnRampConfig config);
  event RouterSet(address router);
  event FeeCharged(address from, address to, uint256 fee);
  event FeesWithdrawn(IERC20 feeToken, address recipient, uint256 amount);

  struct OnRampConfig {
    address router;
    // Fee for sending message taken in this contract
    uint64 relayingFeeJuels;
    // maximum payload data size
    uint64 maxDataSize;
    // Maximum number of distinct ERC20 tokens that can be sent in a message
    uint64 maxTokensLength;
  }

  /**
   * @notice Request a message to be sent to the destination chain
   * @param payload The message payload
   * @param originalSender Original sender of the message if sent by a Router
   * @return The sequence number of the message
   */
  function requestCrossChainSend(CCIP.MessagePayload calldata payload, address originalSender)
    external
    returns (uint256);
}
