package types

import (
	"crypto/sha256"
	"fmt"
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
)

// Get denom pair alphabetical ordered
// [NOTE] WILL BE DEPRECATED in v2
func AlphabeticalDenomPair(denom1, denom2 string) (resDenom1, resDenom2 string) {
	if denom1 > denom2 {
		return denom2, denom1
	} else {
		return denom1, denom2
	}
}

// SortDenoms sorts denoms in an alphabetical order
func SortDenoms(denoms []string) []string {
	sort.Strings(denoms)
	return denoms
}

// GetPoolReserveAcc returns the poor account for the provided poolName (reserve denoms + poolType)
func GetPoolReserveAcc(poolName string) sdk.AccAddress {
	return sdk.AccAddress(crypto.AddressHash([]byte(poolName)))
}

// Generation absolute denomination of the Pool Coin. This rule will be changed on next milestone
func GetPoolCoinDenom(poolName string) string {
	// originally pool coin denom has prefix with / splitter, but remove prefix for pass validation of ibc-transfer
	return fmt.Sprintf("%s%X", PoolCoinDenomPrefix, sha256.Sum256([]byte(poolName)))
}

// Get Total amount of the coins
func GetCoinsTotalAmount(coins sdk.Coins) sdk.Int {
	totalAmount := sdk.ZeroInt()
	for _, coin := range coins {
		totalAmount = totalAmount.Add(coin.Amount)
	}
	return totalAmount
}

// Check Validity of the depositCoins exceed maxReserveCoinAmount
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
