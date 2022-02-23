<!--
Guiding Principles:

Changelogs are for humans, not machines.
There should be an entry for every single version.
The same types of changes should be grouped.
Versions and sections should be linkable.
The latest version comes first.
The release date of each version is displayed.
Mention whether you follow Semantic Versioning.

Usage:

Change log entries are to be added to the Unreleased section under the
appropriate stanza (see below). Each entry should ideally include a tag and
the Github issue reference in the following format:

* (<tag>) \#<issue-number> message

The issue numbers will later be link-ified during the release process so you do
not have to worry about including a link manually, but you can if you wish.

Types of changes (Stanzas):

"Features" for new features.
"Improvements" for changes in existing functionality.
"Deprecated" for soon-to-be removed features.
"Bug Fixes" for any bug fixes.
"Client Breaking" for breaking Protobuf, gRPC and REST routes used by end-users.
"CLI Breaking" for breaking CLI commands.
"API Breaking" for breaking exported APIs used by developers building on SDK.
"State Machine Breaking" for any changes that result in a different AppState given same genesisState and txList.
Ref: https://keepachangelog.com/en/1.0.0/
-->

# Changelog

## [Unreleased]

## [v1.4.6](https://github.com/tendermint/liquidity/releases) - 2022.02.23
* (sdk) Bump SDK version to [v0.44.6](https://github.com/cosmos/cosmos-sdk/releases/tag/v0.44.6)

## [v1.4.5](https://github.com/tendermint/liquidity/releases/tag/v1.4.5) - 2022.01.30
* Unusable release

## [v1.4.4](https://github.com/tendermint/liquidity/releases/tag/v1.4.4) - 2022.01.26
* Unusable release

## [v1.4.2](https://github.com/tendermint/liquidity/releases) - 2021.11.11

* [\#461](https://github.com/tendermint/liquidity/pull/461) (sdk) Bump SDK version to [v0.44.3](https://github.com/cosmos/cosmos-sdk/releases/tag/v0.44.3)

## [v1.4.1](https://github.com/tendermint/liquidity/releases/tag/v1.4.1) - 2021.10.25

* [\#455](https://github.com/tendermint/liquidity/pull/455) (sdk) Bump SDK version to [v0.44.2](https://github.com/cosmos/cosmos-sdk/releases/tag/v0.44.2)
* [\#446](https://github.com/tendermint/liquidity/pull/446) Fix: Pool Coin Decimal Truncation During Deposit
* [\#448](https://github.com/tendermint/liquidity/pull/448) Fix: add overflow checking and test codes for cover edge cases


## [v1.4.0](https://github.com/tendermint/liquidity/releases/tag/v1.4.0) - 2021.09.07

* [\#440](https://github.com/tendermint/liquidity/pull/440) (sdk) Bump SDK version to [v0.44.0](https://github.com/cosmos/cosmos-sdk/releases/tag/v0.44.0)

## [v1.3.0](https://github.com/tendermint/liquidity/releases/tag/v1.3.0) - 2021-08-31

### State Machine Breaking

* [\#433](https://github.com/tendermint/liquidity/pull/433) (sdk) Bump SDK version to [v0.43.0](https://github.com/cosmos/cosmos-sdk/releases/tag/v0.43.0).

* [\#436](https://github.com/tendermint/liquidity/pull/436) Validation `MsgSwapWithinBatch` and `OfferCoinFee` ceiling
  * When calculating `OfferCoinFee`, the decimal points are rounded up.
    - before (v1.2.x):  `MsgSwapWithinBatch.OfferCoinFee` should be `OfferCoin` * `params.SwapFeeRate/2` with Truncate or 0
    - after (v1.3.x):  `MsgSwapWithinBatch.OfferCoinFee` should be `OfferCoin` * `params.SwapFeeRate/2` with Ceil
  * Fix reserveOfferCoinFee residual Issue due to decimal error
  
* [\#438](https://github.com/tendermint/liquidity/pull/438) Fix PoolBatch index, beginHeight issues and genesis logic
  * Remove `PoolBatchIndex`
  * Fix `PoolBatch.Index` duplicated bug
  * Fix `PoolBatch.BeginHeight` consistency issue on genesis init logic  
  
## [v1.2.9](https://github.com/tendermint/liquidity/releases/tag/v1.2.9) - 2021-06-26
 * Liquidity module version 1 for Gravity-DEX
 * (sdk) Bump SDK version to [v0.42.9](https://github.com/cosmos/cosmos-sdk/releases/tag/v0.42.9). 
