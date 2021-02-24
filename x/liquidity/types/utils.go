package types

import (
	"crypto/sha256"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
	"strings"
)

// Get denom pair alphabetical ordered
func AlphabeticalDenomPair(denom1, denom2 string) (resDenom1, resDenom2 string) {
	if denom1 > denom2 {
		return denom2, denom1
	} else {
		return denom1, denom2
	}
}

// GetPoolReserveAcc returns the poor account for the provided poolKey (reserve denoms + poolType)
func GetPoolReserveAcc(poolKey string) sdk.AccAddress {
	return sdk.AccAddress(crypto.AddressHash([]byte(poolKey)))
}

// Generation absolute denomination of the Pool Coin. This rule will be changed on next milestone
func GetPoolCoinDenom(poolKey string) string {
	return fmt.Sprintf("%s/%X", PoolCoinDenomPrefix, sha256.Sum256([]byte(poolKey)))
}

// check is the denom poolcoin or not, need to additional checking the reserve account is existed
func IsPoolCoinDenom(denom string) bool {
	if err := sdk.ValidateDenom(denom); err != nil {
		return false
	}

	denomSplit := strings.SplitN(denom, "/", 2)
	switch {
	case strings.TrimSpace(denom) == "",
		len(denomSplit) == 1 && denomSplit[0] == PoolCoinDenomPrefix,
		len(denomSplit) == 2 && (denomSplit[0] != PoolCoinDenomPrefix || strings.TrimSpace(denomSplit[1]) == ""):
		return false

	case denomSplit[0] == denom && strings.TrimSpace(denom) != "":
		return false
	}
	return true
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
func CoinSafeSubAmount(coinA sdk.Coin, coinBamt sdk.Int) sdk.Coin {
	var resCoin sdk.Coin
	if coinA.Amount.Equal(coinBamt) {
		resCoin = sdk.NewCoin(coinA.Denom, sdk.NewInt(0))
	} else {
		resCoin = coinA.Sub(sdk.NewCoin(coinA.Denom, coinBamt))
	}
	return resCoin
}

//func CoinSafeSub(coinA, coinB sdk.Coin) sdk.Coin {
//	var resCoin sdk.Coin
//	if coinA.Denom != coinB.Denom {
//		return resCoin
//	}
//	if coinA.Equal(coinB) {
//		resCoin = sdk.NewCoin(coinA.Denom, sdk.ZeroInt())
//	} else {
//		coinA = coinA.Sub(sdk.NewCoin(coinA.Denom, coinB.Amount))
//	}
//	return resCoin
//}

// Check the decimals equal approximately
func CheckDecApproxEqual(a , b, threshold sdk.Dec) bool {
	if a.IsZero() && b.IsZero() {
		return true
	} else if a.IsZero() || b.IsZero() {
		return false
	} else if a.Quo(b).Sub(sdk.OneDec()).Abs().LTE(threshold){
		return true
	} else {
		fmt.Println(a, b)
		fmt.Println(a.Quo(b).Sub(sdk.OneDec()).Abs())
		return false
	}
}
