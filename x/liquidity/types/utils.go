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

func CoinSafeSubAmount(coinA sdk.Coin, coinBamt sdk.Int) sdk.Coin {
	//fmt.Println("CoinSafeSubAmount")
	//fmt.Println(coinA)
	//fmt.Println(coinBamt)
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
