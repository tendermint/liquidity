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

### State Machine Breaking

* [\#436](https://github.com/tendermint/liquidity/pull/436) Validation `MsgSwapWithinBatch` and `OfferCoinFee` ceiling
  * When calculating `OfferCoinFee`, the decimal points are rounded up.
  * Fix reserveOfferCoinFee residual Issue due to decimal error
  
* [\#433](https://github.com/tendermint/liquidity/pull/433) (sdk) Bump SDK version to [v0.43.0](https://github.com/cosmos/cosmos-sdk/releases/tag/v0.43.0).  
  
## [v1.2.9](https://github.com/tendermint/liquidity/releases/tag/v1.2.9) - 2021-06-26
 * Liquidity module version 1 for Gravity-DEX
 * (sdk) Bump SDK version to [v0.42.9](https://github.com/cosmos/cosmos-sdk/releases/tag/v0.42.9). 
