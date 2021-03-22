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
	// originally pool coin denom has prefix with / splitter, but remove prefix for pass validation of ibc-tranfer
	return fmt.Sprintf("%s%X", PoolCoinDenomPrefix, sha256.Sum256([]byte(poolName)))
}

// Find A string is exists in the given list
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// Safe Sub function for Coin with subtracting amount
func CoinSafeSubAmount(coinA sdk.Coin, coinBAmt sdk.Int) sdk.Coin {
	var resCoin sdk.Coin
	if coinA.Amount.Equal(coinBAmt) {
		resCoin = sdk.NewCoin(coinA.Denom, sdk.NewInt(0))
	} else {
		resCoin = coinA.Sub(sdk.NewCoin(coinA.Denom, coinBAmt))
	}
	return resCoin
}

// Check the decimals equal approximately
func CheckDecApproxEqual(a, b, threshold sdk.Dec) bool {
	if a.IsZero() && b.IsZero() {
		return true
	} else if a.IsZero() || b.IsZero() {
		return false
	} else if a.Quo(b).Sub(sdk.OneDec()).Abs().LTE(threshold) {
		return true
	} else {
		return false
	}
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
