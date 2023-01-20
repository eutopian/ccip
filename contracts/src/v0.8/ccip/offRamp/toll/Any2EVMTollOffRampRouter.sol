// SPDX-License-Identifier: MIT
pragma solidity 0.8.15;

import {TypeAndVersionInterface} from "../../../interfaces/TypeAndVersionInterface.sol";
import {IAny2EVMOffRampRouter} from "../../interfaces/offRamp/IAny2EVMOffRampRouter.sol";
import {IAny2EVMMessageReceiver} from "../../interfaces/applications/IAny2EVMMessageReceiver.sol";

import {Common} from "../../models/Common.sol";
import {OwnerIsCreator} from "../../access/OwnerIsCreator.sol";

contract Any2EVMTollOffRampRouter is IAny2EVMOffRampRouter, OwnerIsCreator, TypeAndVersionInterface {
  // solhint-disable-next-line chainlink-solidity/all-caps-constant-storage-variables
  string public constant override typeAndVersion = "Any2EVMTollOffRampRouter 1.0.0";

  uint256 private constant GAS_FOR_CALL_EXACT_CHECK = 5_000;

  // Mapping from offRamp to allowed status
  mapping(address => OffRampDetails) internal s_offRamps;
  // List of all offRamps that have  OffRampDetails
  address[] internal s_offRampsList;

  constructor(address[] memory offRamps) {
    s_offRampsList = offRamps;
    for (uint256 i = 0; i < offRamps.length; ++i) {
      s_offRamps[offRamps[i]] = OffRampDetails({listIndex: uint96(i), allowed: true});
    }
  }

  /// @inheritdoc IAny2EVMOffRampRouter
  function routeMessage(
    Common.Any2EVMMessage calldata message,
    bool manualExecution,
    uint256 gasLimit,
    address receiver
  ) external override onlyOffRamp returns (bool success) {
    bytes memory callData = abi.encodeWithSelector(IAny2EVMMessageReceiver.ccipReceive.selector, message);
    if (!manualExecution) {
      success = _callWithExactGas(gasLimit, receiver, 0, callData);
    } else {
      // solhint-disable-next-line avoid-low-level-calls
      (success, ) = receiver.call(callData);
    }
  }

  /**
   * @dev calls target address with exactly gasAmount gas and data as calldata
   * @param gasAmount gas limit for this call
   * @param target target address
   * @param value call ether value
   * @param data calldata
   */
  function _callWithExactGas(
    uint256 gasAmount,
    address target,
    uint256 value,
    bytes memory data
  ) internal returns (bool success) {
    // solhint-disable-next-line no-inline-assembly
    assembly {
      let g := gas()
      // Compute g -= GAS_FOR_CALL_EXACT_CHECK and check for underflow
      // The gas actually passed to the callee is _min(gasAmount, 63//64*gas available).
      // We want to ensure that we revert if gasAmount >  63//64*gas available
      // as we do not want to provide them with less, however that check itself costs
      // gas.  GAS_FOR_CALL_EXACT_CHECK ensures we have at least enough gas to be able
      // to revert if gasAmount >  63//64*gas available.
      if lt(g, GAS_FOR_CALL_EXACT_CHECK) {
        revert(0, 0)
      }
      g := sub(g, GAS_FOR_CALL_EXACT_CHECK)
      // if g - g//64 <= gasAmount, revert
      // (we subtract g//64 because of EIP-150)
      if iszero(gt(sub(g, div(g, 64)), gasAmount)) {
        revert(0, 0)
      }
      // solidity calls check that a contract actually exists at the destination, so we do the same
      if iszero(extcodesize(target)) {
        revert(0, 0)
      }
      // call and return whether we succeeded. ignore return data
      // call(gas,addr,value,argsOffset,argsLength,retOffset,retLength)
      success := call(gasAmount, target, value, add(data, 0x20), mload(data), 0, 0)
    }
    return (success);
  }

  /// @inheritdoc IAny2EVMOffRampRouter
  function addOffRamp(address offRamp) external onlyOwner {
    if (address(offRamp) == address(0)) revert InvalidAddress();
    OffRampDetails memory details = s_offRamps[offRamp];
    // Check if the offramp is already allowed
    if (details.allowed) revert AlreadyConfigured(offRamp);

    // Set the s_offRamps with the new offRamp
    details.allowed = true;
    details.listIndex = uint96(s_offRampsList.length);
    s_offRamps[offRamp] = details;

    // Add to the s_offRampsList
    s_offRampsList.push(offRamp);

    emit OffRampAdded(offRamp);
  }

  /// @inheritdoc IAny2EVMOffRampRouter
  function removeOffRamp(address offRamp) external onlyOwner {
    // Check that there are any feeds to remove
    uint256 listLength = s_offRampsList.length;
    if (listLength == 0) revert NoOffRampsConfigured();

    OffRampDetails memory oldDetails = s_offRamps[offRamp];
    // Check if it exists
    if (!oldDetails.allowed) revert OffRampNotAllowed(offRamp);

    // Swap the last item in the s_offRampsList with the item being removed,
    // update the index of the item moved from the end of the list to its new place,
    // then pop from the end of the list to remove.
    address lastItem = s_offRampsList[listLength - 1];
    // Perform swap
    s_offRampsList[listLength - 1] = s_offRampsList[oldDetails.listIndex];
    s_offRampsList[oldDetails.listIndex] = lastItem;
    // Update listIndex on moved item
    s_offRamps[lastItem].listIndex = oldDetails.listIndex;
    // Pop from list and delete from mapping
    s_offRampsList.pop();
    delete s_offRamps[offRamp];

    emit OffRampRemoved(address(offRamp));
  }

  /// @inheritdoc IAny2EVMOffRampRouter
  function getOffRamps() external view returns (address[] memory offRamps) {
    offRamps = s_offRampsList;
  }

  /// @inheritdoc IAny2EVMOffRampRouter
  function isOffRamp(address offRamp) external view returns (bool allowed) {
    return s_offRamps[offRamp].allowed;
  }

  // @notice only lets allowed offRamps execute
  modifier onlyOffRamp() {
    address offRamp = address(msg.sender);
    if (!s_offRamps[offRamp].allowed) revert MustCallFromOffRamp(msg.sender);
    _;
  }
}
