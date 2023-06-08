// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/// @notice Implements Token Bucket rate limiting.
/// @dev uint128 is safe for rate limiter state.
/// For USD value rate limiting, it can adequately store USD value in 18 decimals.
/// For ERC20 token amount rate limiting, all tokens that will be listed will have at most
/// a supply of uint128.max tokens, and it will therefore not overflow the bucket.
/// In exceptional scenarios where tokens consumed may be larger than uint128,
/// e.g. compromised issuer, an enabled RateLimiter will check and revert.
library RateLimiter {
  error BucketOverfilled();
  error ConsumingMoreThanMaxCapacity(uint256 capacity, uint256 requested);
  error RateLimitReached(uint256 waitInSeconds);
  error OnlyCallableByAdminOrOwner();

  event TokensConsumed(uint256 tokens);
  event ConfigChanged(Config config);

  struct TokenBucket {
    uint128 tokens; // ------┐ Current number of tokens that are in the bucket.
    uint32 lastUpdated; //   | Timestamp in seconds of the last token refill, good for 100+ years.
    bool isEnabled; // ------┘ Indication whether the rate limiting is enabled or not
    uint128 capacity; // ----┐ Maximum number of tokens that can be in the bucket.
    uint128 rate; // --------┘ Number of tokens per second that the bucket is refilled.
  }

  struct Config {
    bool isEnabled; // Indication whether the rate limiting should be enabled
    uint128 capacity; // ----┐ Specifies the capacity of the rate limiter
    uint128 rate; //  -------┘ Specifies the rate of the rate limiter
  }

  /// @notice _consume removes the given tokens from the pool, lowering the
  /// rate tokens allowed to be consumed for subsequent calls.
  /// @param requestTokens The total tokens to be consumed from the bucket.
  /// @dev Reverts when requestTokens exceeds bucket capacity or available tokens in the bucket
  /// @dev emits removal of requestTokens if requestTokens is > 0
  function _consume(TokenBucket storage s_bucket, uint256 requestTokens) internal {
    // If there is no value to remove or rate limiting is turned off, skip this step to reduce gas usage
    if (!s_bucket.isEnabled || requestTokens == 0) {
      return;
    }

    uint256 tokens = s_bucket.tokens;
    uint256 capacity = s_bucket.capacity;
    uint256 timeDiff = block.timestamp - s_bucket.lastUpdated;

    if (timeDiff != 0) {
      if (tokens > capacity) revert BucketOverfilled();

      // Refill tokens when arriving at a new block time
      tokens = _calculateRefill(capacity, tokens, timeDiff, s_bucket.rate);

      s_bucket.lastUpdated = uint32(block.timestamp);
    }

    if (capacity < requestTokens) revert ConsumingMoreThanMaxCapacity(capacity, requestTokens);
    if (tokens < requestTokens) {
      uint256 rate = s_bucket.rate;
      // Seconds wait required until the bucket is refilled enough to accept this value
      // Round up to the next higher second
      revert RateLimitReached(((requestTokens - tokens) + (rate - 1)) / rate);
    }
    tokens -= requestTokens;

    // Downcast is safe here, as tokens is not larger than capacity
    s_bucket.tokens = uint128(tokens);
    emit TokensConsumed(requestTokens);
  }

  /// @notice Gets the token bucket with its values for the block it was requested at.
  /// @return The token bucket.
  function _currentTokenBucketState(TokenBucket memory bucket) internal view returns (TokenBucket memory) {
    // We update the bucket to reflect the status at the exact time of the
    // call. This means we might need to refill a part of the bucket based
    // on the time that has passed since the last update.
    bucket.tokens = uint128(
      _calculateRefill(bucket.capacity, bucket.tokens, block.timestamp - bucket.lastUpdated, bucket.rate)
    );
    bucket.lastUpdated = uint32(block.timestamp);
    return bucket;
  }

  /// @notice Sets the rate limited config.
  /// @param s_bucket The token bucket
  /// @param config The new config
  function _setTokenBucketConfig(TokenBucket storage s_bucket, Config memory config) internal {
    // First update the bucket to make sure the proper rate is used for all the time
    // up until the config change.
    uint256 timeDiff = block.timestamp - s_bucket.lastUpdated;
    if (timeDiff != 0) {
      s_bucket.tokens = uint128(_calculateRefill(s_bucket.capacity, s_bucket.tokens, timeDiff, s_bucket.rate));

      s_bucket.lastUpdated = uint32(block.timestamp);
    }

    s_bucket.tokens = uint128(_min(config.capacity, s_bucket.tokens));
    s_bucket.isEnabled = config.isEnabled;
    s_bucket.capacity = config.capacity;
    s_bucket.rate = config.rate;

    emit ConfigChanged(config);
  }

  /// @notice Calculate refilled tokens
  /// @param capacity bucket capacity
  /// @param tokens current bucket tokens
  /// @param timeDiff block time difference since last refill
  /// @param rate bucket refill rate
  /// @return the value of tokens after refill
  function _calculateRefill(
    uint256 capacity,
    uint256 tokens,
    uint256 timeDiff,
    uint256 rate
  ) private pure returns (uint256) {
    return _min(capacity, tokens + timeDiff * rate);
  }

  /// @notice Return the smallest of two integers
  /// @param a first int
  /// @param b second int
  /// @return smallest
  function _min(uint256 a, uint256 b) internal pure returns (uint256) {
    return a < b ? a : b;
  }
}
