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
		return denom1, denom1
	}
}

// GetPoolReserveAcc returns the poor account for the provided poolKey (reserve denoms + poolType)
func GetPoolReserveAcc(poolKey string) sdk.AccAddress {
	return sdk.AccAddress(crypto.AddressHash([]byte(poolKey)))
}

// tmp denom rule, TBD
func GetPoolCoinDenom(reserveAcc sdk.AccAddress) string {
	return reserveAcc.String()
}
