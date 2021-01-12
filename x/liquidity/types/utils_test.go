package types_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/liquidity/x/liquidity/types"
	"testing"
)

func TestAlphabeticalDenomPair(t *testing.T) {
	denomA := "uCoinA"
	denomB := "uCoinB"
	afterDenomA, afterDenomB := types.AlphabeticalDenomPair(denomA, denomB)
	require.Equal(t, denomA, afterDenomA)
	require.Equal(t, denomB, afterDenomB)

	afterDenomA, afterDenomB = types.AlphabeticalDenomPair(denomB, denomA)
	require.Equal(t, denomA, afterDenomA)
	require.Equal(t, denomB, afterDenomB)
}

func TestStringInSlice(t *testing.T) {
	denomA := "uCoinA"
	denomB := "uCoinB"
	denomC := "uCoinC"
	denoms := []string{denomA, denomB}
	require.True(t, types.StringInSlice(denomA, denoms))
	require.True(t, types.StringInSlice(denomB, denoms))
	require.False(t, types.StringInSlice(denomC, denoms))
}

func TestCoinSafeSubAmount(t *testing.T) {
	denom := "uCoinA"
	a := sdk.NewCoin(denom, sdk.NewInt(100))
	b := sdk.NewCoin(denom, sdk.NewInt(100))
	res := types.CoinSafeSubAmount(a, b.Amount)
	require.Equal(t, sdk.NewCoin(denom, sdk.NewInt(0)), res)

	a = sdk.NewCoin(denom, sdk.NewInt(100))
	b = sdk.NewCoin(denom, sdk.NewInt(50))
	res = types.CoinSafeSubAmount(a, b.Amount)
	require.Equal(t, sdk.NewCoin(denom, sdk.NewInt(50)), res)

	require.Panics(t, func() {
		res = types.CoinSafeSubAmount(b, a.Amount)
	})
}

func TestGetPoolReserveAcc(t *testing.T) {
	reserveAcc := types.GetPoolReserveAcc("denomX-denomY-1")
	require.NotNil(t, reserveAcc)
	require.Equal(t, "cosmos18gvpvm3lrzx6rs6yq5c6klnye2t5qumm3v3re8", reserveAcc.String())
	require.Equal(t, "cosmos18gvpvm3lrzx6rs6yq5c6klnye2t5qumm3v3re8", types.GetPoolCoinDenom(reserveAcc))
}

// WIP for check equal approximately OfferCoinFee
//func TestEqualApprox(t *testing.T) {
//	decA, err := sdk.NewDecFromStr("10000000000000000.0000000001")
//	require.NoError(t, err)
//	decB, err := sdk.NewDecFromStr("10000000000000000.00000000021")
//	require.NoError(t, err)
//
//	decC, err := sdk.NewDecFromStr("10000000000000000.00000000001")
//	require.NoError(t, err)
//	decD, err := sdk.NewDecFromStr("10000000000000000.00000000002")
//	require.NoError(t, err)
//
//	decE, err := sdk.NewDecFromStr("10000000000000000.000000001")
//	require.NoError(t, err)
//	decF, err := sdk.NewDecFromStr("10000000000000000.000000002")
//	require.NoError(t, err)
//
//	// TODO: fix
//	require.False(t, decA.Equal(decB))
//	require.False(t, types.EqualApprox(decA, decB))
//
//	require.False(t, decC.Equal(decD))
//	require.True(t, types.EqualApprox(decC, decD))
//
//	require.False(t, decE.Equal(decF))
//	require.True(t, types.EqualApprox(decE, decF))
//
//	fmt.Println(decA)
//	fmt.Println(decB)
//	fmt.Println(sdk.NewDecWithPrec(1, 10).MulInt64(10000000000))
//	fmt.Println(sdk.NewDecWithPrec(1, 10).MulInt64(10000000000))
//	fmt.Println(sdk.NewDecWithPrec(1, 11).MulInt64(10000000000))
//}
