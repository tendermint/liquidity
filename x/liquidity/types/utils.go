package types

import (
	"crypto/sha256"
	"fmt"
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
)

// AlphabeticalDenomPair returns denom pairs that are alphabetically sorted.
func AlphabeticalDenomPair(denom1, denom2 string) (resDenom1, resDenom2 string) {
	if denom1 > denom2 {
		return denom2, denom1
	} else {
		return denom1, denom2
	}
}

// SortDenoms sorts denoms in alphabetical order.
func SortDenoms(denoms []string) []string {
	sort.Strings(denoms)
	return denoms
}

// GetPoolReserveAcc returns the address of the pool's reserve account.
func GetPoolReserveAcc(poolName string) sdk.AccAddress {
	return sdk.AccAddress(crypto.AddressHash([]byte(poolName)))
}

// GetPoolCoinDenom returns the denomination of the pool coin.
func GetPoolCoinDenom(poolName string) string {
	// Originally pool coin denom has prefix with / splitter, but removed prefix for pass validation of ibc-transfer
	return fmt.Sprintf("%s%X", PoolCoinDenomPrefix, sha256.Sum256([]byte(poolName)))
}

// GetCoinsTotalAmount returns total amount of all coins in sdk.Coins.
func GetCoinsTotalAmount(coins sdk.Coins) sdk.Int {
	totalAmount := sdk.ZeroInt()
	for _, coin := range coins {
		totalAmount = totalAmount.Add(coin.Amount)
	}
	return totalAmount
}

// ValidateReserveCoinLimit checks if total amounts of depositCoins exceed maxReserveCoinAmount.
func ValidateReserveCoinLimit(maxReserveCoinAmount sdk.Int, depositCoins sdk.Coins) error {
	totalAmount := GetCoinsTotalAmount(depositCoins)
	if maxReserveCoinAmount.IsZero() {
		return nil
	} else if totalAmount.GT(maxReserveCoinAmount) {
		return ErrExceededReserveCoinLimit
	} else {
		return nil
	}
}

func GetOfferCoinFee(offerCoin sdk.Coin, swapFeeRate sdk.Dec) sdk.Coin {
	// apply half-ratio swap fee rate
	// see https://github.com/tendermint/liquidity/issues/41 for details
	return sdk.NewCoin(offerCoin.Denom, offerCoin.Amount.ToDec().Mul(swapFeeRate.QuoInt64(2)).TruncateInt()) // offerCoin.Amount * (swapFeeRate/2)
}
