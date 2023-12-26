// SPDX-License-Identifier: BUSL-1.1
pragma solidity 0.8.19;

import {IL1Bridge} from "./interfaces/IBridge.sol";
import {IWrappedNative} from "../../interfaces/IWrappedNative.sol";

import {IL1GatewayRouter} from "@arbitrum/token-bridge-contracts/contracts/tokenbridge/ethereum/gateway/IL1GatewayRouter.sol";
import {IOutbox} from "@arbitrum/nitro-contracts/src/bridge/IOutbox.sol";
import {IERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/IERC20.sol";
import {SafeERC20} from "../../../vendor/openzeppelin-solidity/v4.8.3/contracts/token/ERC20/utils/SafeERC20.sol";

/// @notice Arbitrum L1 Bridge adapter
/// @dev Auto unwraps and re-wraps wrapped eth in the bridge.
contract ArbitrumL1BridgeAdapter is IL1Bridge {
  using SafeERC20 for IERC20;

  error InsufficientEthValue(uint256 wanted, uint256 got);

  IL1GatewayRouter internal immutable i_l1GatewayRouter;
  address internal immutable i_l1ERC20Gateway;
  IOutbox internal immutable i_l1Outbox;

  // TODO not static?
  uint256 public constant MAX_GAS = 100_000;
  uint256 public constant GAS_PRICE_BID = 300_000_000;
  uint256 public constant MAX_SUBMISSION_COST = 8e14;

  // Nonce to use for L2 deposits to allow for better tracking offchain.
  uint64 private s_nonce = 0;

  constructor(IL1GatewayRouter l1GatewayRouter, IOutbox l1Outbox, address l1ERC20Gateway) {
    if (
      address(l1GatewayRouter) == address(0) || address(l1Outbox) == address(0) || address(l1ERC20Gateway) == address(0)
    ) {
      revert BridgeAddressCannotBeZero();
    }
    i_l1GatewayRouter = l1GatewayRouter;
    i_l1Outbox = l1Outbox;
    i_l1ERC20Gateway = l1ERC20Gateway;
  }

  function sendERC20(address l1Token, address, address recipient, uint256 amount) external payable {
    IERC20(l1Token).safeTransferFrom(msg.sender, address(this), amount);

    IERC20(l1Token).approve(i_l1ERC20Gateway, amount);

    uint256 wantedNativeFeeCoin = MAX_SUBMISSION_COST + MAX_GAS * GAS_PRICE_BID;
    if (msg.value < wantedNativeFeeCoin) {
      revert InsufficientEthValue(wantedNativeFeeCoin, msg.value);
    }

    i_l1GatewayRouter.outboundTransferCustomRefund{value: wantedNativeFeeCoin}(
      l1Token,
      recipient,
      recipient,
      amount,
      MAX_GAS,
      GAS_PRICE_BID,
      abi.encode(MAX_SUBMISSION_COST, bytes(""))
    );
  }

  /// @param proof Merkle proof of message inclusion in send root
  /// @param index Merkle path to message
  /// @param l2Block l2 block number at which sendTxToL1 call was made
  /// @param l1Block l1 block number at which sendTxToL1 call was made
  /// @param l2Timestamp l2 Timestamp at which sendTxToL1 call was made
  /// @param value wei in L1 message
  /// @param data abi-encoded L1 message data
  struct ArbitrumFinalizationPayload {
    bytes32[] proof;
    uint256 index;
    uint256 l2Block;
    uint256 l1Block;
    uint256 l2Timestamp;
    uint256 value;
    bytes data;
  }

  /// @param l2Sender sender if original message (i.e., caller of ArbSys.sendTxToL1)
  /// @param l1Receiver destination address for L1 contract call
  function finalizeWithdrawERC20FromL2(
    address l2Sender,
    address l1Receiver,
    bytes calldata arbitrumFinalizationPayload
  ) external {
    ArbitrumFinalizationPayload memory payload = abi.decode(arbitrumFinalizationPayload, (ArbitrumFinalizationPayload));
    i_l1Outbox.executeTransaction(
      payload.proof,
      payload.index,
      l2Sender,
      l1Receiver,
      payload.l2Block,
      payload.l1Block,
      payload.l2Timestamp,
      payload.value,
      payload.data
    );
  }

  function getL2Token(address l1Token) external view returns (address) {
    return i_l1GatewayRouter.calculateL2TokenAddress(l1Token);
  }
}
