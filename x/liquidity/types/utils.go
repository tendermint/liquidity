package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
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

// TODO: tmp denom rule, It will fixed on milestone 2
func GetPoolCoinDenom(reserveAcc sdk.AccAddress) string {
	return reserveAcc.String()
}

// check is poolcoin or not when poolcoin denom rule fixed
//func IsPoolCoin(coin sdk.Coin) bool {
//}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
