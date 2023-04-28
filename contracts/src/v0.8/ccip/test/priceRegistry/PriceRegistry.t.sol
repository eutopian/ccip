// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {IPriceRegistry} from "../../interfaces/IPriceRegistry.sol";

import {Internal} from "../../libraries/Internal.sol";
import {TokenSetup} from "../TokenSetup.t.sol";
import {RouterSetup} from "../router/RouterSetup.t.sol";
import {PriceRegistry} from "../../PriceRegistry.sol";

import {IERC20} from "../../../vendor/IERC20.sol";

contract PriceRegistrySetup is TokenSetup, RouterSetup {
  PriceRegistry internal s_priceRegistry;
  // Cheat to store the price updates in storage since struct arrays aren't supported.
  bytes internal s_encodedInitialPriceUpdates;
  address internal s_weth;

  function setUp() public virtual override(TokenSetup, RouterSetup) {
    TokenSetup.setUp();
    RouterSetup.setUp();

    s_weth = s_sourceRouter.getWrappedNative();

    address[] memory pricedSourceTokens = getPricedSourceToken();
    address[] memory pricedDestTokens = getPricedDestTokens();
    uint192[] memory sourceTokenPrices = getSourceTokenPrices();
    uint192[] memory destTokenPrices = getDestTokenPrices();

    uint256 sourceTokenCount = pricedSourceTokens.length;
    uint256 destTokenCount = pricedDestTokens.length;
    address[] memory pricedTokens = new address[](sourceTokenCount + destTokenCount);
    uint192[] memory tokenPrices = new uint192[](sourceTokenCount + destTokenCount);
    for (uint256 i = 0; i < sourceTokenCount; ++i) {
      pricedTokens[i] = pricedSourceTokens[i];
      tokenPrices[i] = sourceTokenPrices[i];
    }
    for (uint256 i = 0; i < destTokenCount; ++i) {
      pricedTokens[i + sourceTokenCount] = pricedDestTokens[i];
      tokenPrices[i + sourceTokenCount] = destTokenPrices[i];
    }

    Internal.PriceUpdates memory priceUpdates = getPriceUpdatesStruct(pricedTokens, tokenPrices);
    priceUpdates.destChainId = DEST_CHAIN_ID;
    priceUpdates.usdPerUnitGas = 1e6;

    s_encodedInitialPriceUpdates = abi.encode(priceUpdates);
    address[] memory priceUpdaters = new address[](0);
    address[] memory feeTokens = new address[](2);
    feeTokens[0] = s_sourceTokens[0];
    feeTokens[1] = s_weth;
    s_priceRegistry = new PriceRegistry(priceUpdates, priceUpdaters, feeTokens, uint32(TWELVE_HOURS));
  }

  function getPricedSourceToken() internal view returns (address[] memory) {
    address[] memory tokens = new address[](3);
    tokens[0] = s_sourceTokens[0];
    tokens[1] = s_sourceTokens[1];
    tokens[2] = s_sourceRouter.getWrappedNative();
    return tokens;
  }

  function getPricedDestTokens() internal view returns (address[] memory) {
    address[] memory tokens = new address[](3);
    tokens[0] = s_destTokens[0];
    tokens[1] = s_destTokens[1];
    tokens[2] = s_destRouter.getWrappedNative();
    return tokens;
  }

  function getSourceTokenPrices() internal pure returns (uint192[] memory) {
    uint192[] memory prices = new uint192[](3);
    prices[0] = 5e18;
    prices[1] = 2000e18;
    prices[2] = 2000e18;
    return prices;
  }

  function getDestTokenPrices() internal pure returns (uint192[] memory) {
    uint192[] memory prices = new uint192[](3);
    prices[0] = 5e18;
    prices[1] = 2000e18;
    prices[2] = 2000e18;
    return prices;
  }
}

contract PriceRegistry_constructor is PriceRegistrySetup {
  function testSetupSuccess() public virtual {
    assertEq(s_priceRegistry.getTokenPrice(s_sourceTokens[0]).value, 5e18);
    assertEq(s_priceRegistry.getTokenPrice(s_sourceTokens[1]).value, 2000e18);
    assertEq(s_priceRegistry.getDestinationChainGasPrice(DEST_CHAIN_ID).value, 1e6);
    assertEq(s_priceRegistry.getFeeTokenBaseUnitsPerUnitGas(s_sourceTokens[0], DEST_CHAIN_ID), 2e5);
    assertEq(s_priceRegistry.getFeeTokenBaseUnitsPerUnitGas(s_weth, DEST_CHAIN_ID), 500);
  }

  function testInvalidStalenessThresholdReverts() public {
    Internal.PriceUpdates memory priceUpdates = Internal.PriceUpdates({
      tokenPriceUpdates: new Internal.TokenPriceUpdate[](0),
      destChainId: DEST_CHAIN_ID,
      usdPerUnitGas: 1e6
    });

    vm.expectRevert(PriceRegistry.InvalidStalenessThreshold.selector);
    s_priceRegistry = new PriceRegistry(priceUpdates, new address[](0), new address[](0), 0);
  }
}

contract PriceRegistry_getTokenPrices is PriceRegistrySetup {
  function testGetTokenPricesSuccess() public {
    Internal.PriceUpdates memory priceUpdates = abi.decode(s_encodedInitialPriceUpdates, (Internal.PriceUpdates));

    address[] memory tokens = new address[](3);
    tokens[0] = s_sourceTokens[0];
    tokens[1] = s_sourceTokens[1];
    tokens[2] = s_weth;

    IPriceRegistry.TimestampedUint192Value[] memory tokenPrices = s_priceRegistry.getTokenPrices(tokens);

    assertEq(tokenPrices.length, 3);
    assertEq(tokenPrices[0].value, priceUpdates.tokenPriceUpdates[0].usdPerToken);
    assertEq(tokenPrices[1].value, priceUpdates.tokenPriceUpdates[1].usdPerToken);
    assertEq(tokenPrices[2].value, priceUpdates.tokenPriceUpdates[2].usdPerToken);
  }
}

contract PriceRegistry_applyPriceUpdatersUpdates is PriceRegistrySetup {
  event PriceUpdaterSet(address indexed priceUpdater);
  event PriceUpdaterRemoved(address indexed priceUpdater);

  function testApplyPriceUpdaterUpdatesSuccess() public {
    address[] memory priceUpdaters = new address[](1);
    priceUpdaters[0] = STRANGER;

    vm.expectEmit();
    emit PriceUpdaterSet(STRANGER);

    s_priceRegistry.applyPriceUpdatersUpdates(priceUpdaters, new address[](0));
    assertEq(s_priceRegistry.getPriceUpdaters().length, 1);
    assertEq(s_priceRegistry.getPriceUpdaters()[0], STRANGER);

    // add same priceUpdater is no-op
    s_priceRegistry.applyPriceUpdatersUpdates(priceUpdaters, new address[](0));
    assertEq(s_priceRegistry.getPriceUpdaters().length, 1);
    assertEq(s_priceRegistry.getPriceUpdaters()[0], STRANGER);

    vm.expectEmit();
    emit PriceUpdaterRemoved(STRANGER);

    s_priceRegistry.applyPriceUpdatersUpdates(new address[](0), priceUpdaters);
    assertEq(s_priceRegistry.getPriceUpdaters().length, 0);

    // removing already removed priceUpdater is no-op
    s_priceRegistry.applyPriceUpdatersUpdates(new address[](0), priceUpdaters);
    assertEq(s_priceRegistry.getPriceUpdaters().length, 0);
  }

  function testOnlyCallableByOwnerReverts() public {
    address[] memory priceUpdaters = new address[](1);
    priceUpdaters[0] = STRANGER;
    changePrank(STRANGER);
    vm.expectRevert("Only callable by owner");
    s_priceRegistry.applyPriceUpdatersUpdates(priceUpdaters, new address[](0));
  }
}

contract PriceRegistry_applyFeeTokensUpdates is PriceRegistrySetup {
  event FeeTokenAdded(address indexed feeToken);
  event FeeTokenRemoved(address indexed feeToken);

  function testApplyFeeTokensUpdatesSuccess() public {
    address[] memory feeTokens = new address[](1);
    feeTokens[0] = s_sourceTokens[1];

    vm.expectEmit();
    emit FeeTokenAdded(feeTokens[0]);

    s_priceRegistry.applyFeeTokensUpdates(feeTokens, new address[](0));
    assertEq(s_priceRegistry.getFeeTokens().length, 3);
    assertEq(s_priceRegistry.getFeeTokens()[2], feeTokens[0]);

    // add same feeToken is no-op
    s_priceRegistry.applyFeeTokensUpdates(feeTokens, new address[](0));
    assertEq(s_priceRegistry.getFeeTokens().length, 3);
    assertEq(s_priceRegistry.getFeeTokens()[2], feeTokens[0]);

    vm.expectEmit();
    emit FeeTokenRemoved(feeTokens[0]);

    s_priceRegistry.applyFeeTokensUpdates(new address[](0), feeTokens);
    assertEq(s_priceRegistry.getFeeTokens().length, 2);

    // removing already removed feeToken is no-op
    s_priceRegistry.applyFeeTokensUpdates(new address[](0), feeTokens);
    assertEq(s_priceRegistry.getFeeTokens().length, 2);
  }

  function testOnlyCallableByOwnerReverts() public {
    address[] memory feeTokens = new address[](1);
    feeTokens[0] = STRANGER;
    changePrank(STRANGER);
    vm.expectRevert("Only callable by owner");
    s_priceRegistry.applyFeeTokensUpdates(feeTokens, new address[](0));
  }
}

contract PriceRegistry_updatePrices is PriceRegistrySetup {
  // Cheat to store the price updates in storage since struct arrays aren't supported.
  bytes internal s_encodedNewPriceUpdates;

  function setUp() public virtual override {
    PriceRegistrySetup.setUp();
    Internal.TokenPriceUpdate[] memory tokenPriceUpdates = new Internal.TokenPriceUpdate[](2);
    tokenPriceUpdates[0] = Internal.TokenPriceUpdate({sourceToken: s_sourceTokens[0], usdPerToken: 4e18});
    tokenPriceUpdates[1] = Internal.TokenPriceUpdate({sourceToken: s_sourceTokens[1], usdPerToken: 1800e18});
    Internal.PriceUpdates memory priceUpdates = Internal.PriceUpdates({
      tokenPriceUpdates: tokenPriceUpdates,
      destChainId: DEST_CHAIN_ID,
      usdPerUnitGas: 2e6
    });
    s_encodedNewPriceUpdates = abi.encode(priceUpdates);
  }

  function testUpdatePricesSuccess() public {
    Internal.PriceUpdates memory priceUpdates = abi.decode(s_encodedNewPriceUpdates, (Internal.PriceUpdates));
    s_priceRegistry.updatePrices(priceUpdates);

    assertEq(s_priceRegistry.getTokenPrice(s_sourceTokens[0]).value, priceUpdates.tokenPriceUpdates[0].usdPerToken);
    assertEq(s_priceRegistry.getTokenPrice(s_sourceTokens[1]).value, priceUpdates.tokenPriceUpdates[1].usdPerToken);
    assertEq(s_priceRegistry.getDestinationChainGasPrice(DEST_CHAIN_ID).value, priceUpdates.usdPerUnitGas);
  }

  // Reverts

  function testOnlyCallableByUpdaterOrOwnerReverts() public {
    Internal.PriceUpdates memory priceUpdates = abi.decode(s_encodedNewPriceUpdates, (Internal.PriceUpdates));
    changePrank(STRANGER);
    vm.expectRevert(abi.encodeWithSelector(PriceRegistry.OnlyCallableByUpdaterOrOwner.selector));
    s_priceRegistry.updatePrices(priceUpdates);
  }
}

contract PriceRegistry_convertTokenAmount is PriceRegistrySetup {
  function testConvertTokenAmountSuccess() public {
    Internal.PriceUpdates memory initialPriceUpdates = abi.decode(
      s_encodedInitialPriceUpdates,
      (Internal.PriceUpdates)
    );
    uint256 amount = 3e16;
    uint256 conversionRate = (uint256(initialPriceUpdates.tokenPriceUpdates[2].usdPerToken) * 1e18) /
      uint256(initialPriceUpdates.tokenPriceUpdates[0].usdPerToken);
    uint256 expected = (amount * conversionRate) / 1e18;
    assertEq(s_priceRegistry.convertTokenAmount(s_weth, amount, s_sourceTokens[0]), expected);
  }

  function test_fuzz_ConvertTokenAmountSuccess(
    uint256 feeTokenAmount,
    uint192 usdPerFeeToken,
    uint160 usdPerLinkToken,
    uint192 usdPerUnitGas
  ) public {
    vm.assume(usdPerFeeToken > 0);
    vm.assume(usdPerLinkToken > 0);
    // We bound the max fees to be at most uint96.max link.
    feeTokenAmount = bound(feeTokenAmount, 0, (uint256(type(uint96).max) * usdPerLinkToken) / usdPerFeeToken);

    address feeToken = address(1);
    address linkToken = address(2);
    address[] memory feeTokens = new address[](1);
    feeTokens[0] = feeToken;
    s_priceRegistry.applyFeeTokensUpdates(feeTokens, new address[](0));

    Internal.TokenPriceUpdate[] memory tokenPriceUpdates = new Internal.TokenPriceUpdate[](2);
    tokenPriceUpdates[0] = Internal.TokenPriceUpdate({sourceToken: feeToken, usdPerToken: usdPerFeeToken});
    tokenPriceUpdates[1] = Internal.TokenPriceUpdate({sourceToken: linkToken, usdPerToken: usdPerLinkToken});
    Internal.PriceUpdates memory priceUpdates = Internal.PriceUpdates({
      tokenPriceUpdates: tokenPriceUpdates,
      destChainId: DEST_CHAIN_ID,
      usdPerUnitGas: usdPerUnitGas
    });

    s_priceRegistry.updatePrices(priceUpdates);

    uint256 linkFee = s_priceRegistry.convertTokenAmount(feeToken, feeTokenAmount, linkToken);
    assertEq(linkFee, (feeTokenAmount * usdPerFeeToken) / usdPerLinkToken);
  }

  // Reverts

  function testStaleFeeTokenReverts() public {
    vm.warp(block.timestamp + TWELVE_HOURS + 1);

    Internal.TokenPriceUpdate[] memory tokenPriceUpdates = new Internal.TokenPriceUpdate[](1);
    tokenPriceUpdates[0] = Internal.TokenPriceUpdate({sourceToken: s_sourceTokens[0], usdPerToken: 4e18});
    Internal.PriceUpdates memory priceUpdates = Internal.PriceUpdates({
      tokenPriceUpdates: tokenPriceUpdates,
      destChainId: 0,
      usdPerUnitGas: 0
    });
    s_priceRegistry.updatePrices(priceUpdates);

    vm.expectRevert(
      abi.encodeWithSelector(
        PriceRegistry.StaleTokenPrice.selector,
        s_weth,
        uint128(TWELVE_HOURS),
        uint128(TWELVE_HOURS + 1)
      )
    );
    s_priceRegistry.convertTokenAmount(s_weth, 3e16, s_sourceTokens[0]);
  }

  function testLinkTokenNotSupportedReverts() public {
    vm.expectRevert(abi.encodeWithSelector(PriceRegistry.TokenNotSupported.selector, DUMMY_CONTRACT_ADDRESS));
    s_priceRegistry.convertTokenAmount(DUMMY_CONTRACT_ADDRESS, 3e16, s_sourceTokens[0]);

    vm.expectRevert(abi.encodeWithSelector(PriceRegistry.TokenNotSupported.selector, DUMMY_CONTRACT_ADDRESS));
    s_priceRegistry.convertTokenAmount(s_sourceTokens[0], 3e16, DUMMY_CONTRACT_ADDRESS);
  }

  function testStaleLinkTokenReverts() public {
    vm.warp(block.timestamp + TWELVE_HOURS + 1);

    Internal.TokenPriceUpdate[] memory tokenPriceUpdates = new Internal.TokenPriceUpdate[](1);
    tokenPriceUpdates[0] = Internal.TokenPriceUpdate({sourceToken: s_weth, usdPerToken: 18e17});
    Internal.PriceUpdates memory priceUpdates = Internal.PriceUpdates({
      tokenPriceUpdates: tokenPriceUpdates,
      destChainId: 0,
      usdPerUnitGas: 0
    });
    s_priceRegistry.updatePrices(priceUpdates);

    vm.expectRevert(
      abi.encodeWithSelector(
        PriceRegistry.StaleTokenPrice.selector,
        s_sourceTokens[0],
        uint128(TWELVE_HOURS),
        uint128(TWELVE_HOURS + 1)
      )
    );
    s_priceRegistry.convertTokenAmount(s_weth, 3e16, s_sourceTokens[0]);
  }
}

contract PriceRegistry_getFeeTokenBaseUnitsPerUnitGas is PriceRegistrySetup {
  function testGetFeeSuccess() public {
    // 1 unit of gas costs 0.000001 USD -> 1e6
    // 1 LINK costs 5 USD -> 5e18
    // gasPrice / linkPrice = (1e6 * 1e18) / 5e18 = 2e5
    assertEq(s_priceRegistry.getFeeTokenBaseUnitsPerUnitGas(s_sourceTokens[0], DEST_CHAIN_ID), 2e5);
  }

  function testZeroGasPriceSuccess() public {
    uint64 zeroGasDestChainId = 345678;
    Internal.PriceUpdates memory priceUpdates = Internal.PriceUpdates({
      tokenPriceUpdates: new Internal.TokenPriceUpdate[](0),
      destChainId: zeroGasDestChainId,
      usdPerUnitGas: 0
    });
    s_priceRegistry.updatePrices(priceUpdates);

    assertEq(s_priceRegistry.getFeeTokenBaseUnitsPerUnitGas(s_sourceTokens[0], zeroGasDestChainId), 0);
  }

  function testUnsupportedTokenReverts() public {
    vm.expectRevert(abi.encodeWithSelector(PriceRegistry.NotAFeeToken.selector, DUMMY_CONTRACT_ADDRESS));
    s_priceRegistry.getFeeTokenBaseUnitsPerUnitGas(DUMMY_CONTRACT_ADDRESS, DEST_CHAIN_ID);
  }

  function testUnsupportedChainReverts() public {
    vm.expectRevert(abi.encodeWithSelector(PriceRegistry.ChainNotSupported.selector, DEST_CHAIN_ID + 1));
    s_priceRegistry.getFeeTokenBaseUnitsPerUnitGas(s_sourceTokens[0], DEST_CHAIN_ID + 1);
  }

  function testStaleGasPriceReverts() public {
    uint256 diff = TWELVE_HOURS + 1;
    vm.warp(block.timestamp + diff);
    vm.expectRevert(abi.encodeWithSelector(PriceRegistry.StaleGasPrice.selector, DEST_CHAIN_ID, TWELVE_HOURS, diff));
    s_priceRegistry.getFeeTokenBaseUnitsPerUnitGas(s_sourceTokens[0], DEST_CHAIN_ID);
  }

  function testStaleTokenPriceReverts() public {
    uint256 diff = TWELVE_HOURS + 1;
    vm.warp(block.timestamp + diff);

    Internal.PriceUpdates memory priceUpdates = Internal.PriceUpdates({
      tokenPriceUpdates: new Internal.TokenPriceUpdate[](0),
      destChainId: DEST_CHAIN_ID,
      usdPerUnitGas: 1e6
    });
    s_priceRegistry.updatePrices(priceUpdates);

    vm.expectRevert(
      abi.encodeWithSelector(PriceRegistry.StaleTokenPrice.selector, s_sourceTokens[0], TWELVE_HOURS, diff)
    );
    s_priceRegistry.getFeeTokenBaseUnitsPerUnitGas(s_sourceTokens[0], DEST_CHAIN_ID);
  }
}
