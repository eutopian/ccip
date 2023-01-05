// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import {IEVM2EVMTollOnRamp} from "../interfaces/onRamp/IEVM2EVMTollOnRamp.sol";
import {IBaseOnRamp} from "../interfaces/onRamp/IBaseOnRamp.sol";
import {IPool} from "../interfaces/pools/IPool.sol";

import {OwnerIsCreator} from "../access/OwnerIsCreator.sol";
import {Common} from "../models/Common.sol";

import {SafeERC20} from "../../vendor/SafeERC20.sol";
import {IERC20} from "../../vendor/IERC20.sol";

contract PoolCollector is OwnerIsCreator {
  using SafeERC20 for IERC20;

  event FeeCharged(address from, address to, uint256 fee);
  event FeesWithdrawn(IERC20 feeToken, address recipient, uint256 amount);

  error FeeTokenAmountTooLow();

  /**
   * @notice Collect the fee
   * @param onRamp OnRamp to get the fee and pools from
   * @param feeToken the feeToken to be collected
   * @param feeTokenAmount the amount of feeToken that is available
   */
  function _chargeFee(
    IEVM2EVMTollOnRamp onRamp,
    IERC20 feeToken,
    uint256 feeTokenAmount
  ) internal returns (uint256 fee) {
    // Ensure fee token is valid.
    IPool feeTokenPool = onRamp.getPoolBySourceToken(feeToken);
    if (address(feeTokenPool) == address(0)) revert IBaseOnRamp.UnsupportedToken(feeToken);
    fee = onRamp.getRequiredFee(feeToken);
    address sender = msg.sender;
    if (fee > 0) {
      if (fee > feeTokenAmount) revert FeeTokenAmountTooLow();
      feeTokenAmount -= fee;
      feeToken.safeTransferFrom(sender, address(this), fee);
    }
    if (feeTokenAmount > 0) {
      // Send the fee token to the pool
      feeToken.safeTransferFrom(sender, address(feeTokenPool), feeTokenAmount);
    }
    emit FeeCharged(sender, address(this), fee);
  }

  /**
   * @notice Collect tokens and send them to the pools
   * @param onRamp OnRamp to get the fee and pools from
   * @param tokensAndAmounts the tokensAndAmounts to be collected
   */
  function _collectTokens(IBaseOnRamp onRamp, Common.EVMTokenAndAmount[] memory tokensAndAmounts) internal {
    // Send the tokens to the pools
    for (uint256 i = 0; i < tokensAndAmounts.length; ++i) {
      IERC20 token = IERC20(tokensAndAmounts[i].token);
      IPool pool = onRamp.getPoolBySourceToken(token);
      if (address(pool) == address(0)) revert IBaseOnRamp.UnsupportedToken(token);
      token.safeTransferFrom(msg.sender, address(pool), tokensAndAmounts[i].amount);
    }
  }

  /**
   * @notice Withdraw the fee tokens accumulated in this contract
   * @dev only callable by owner
   */
  function withdrawAccumulatedFees(
    IERC20 feeToken,
    address recipient,
    uint256 amount
  ) external onlyOwner {
    feeToken.safeTransfer(recipient, amount);
    emit FeesWithdrawn(feeToken, recipient, amount);
  }
}
