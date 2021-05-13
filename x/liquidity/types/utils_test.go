package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/tendermint/liquidity/x/liquidity/types"
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

func TestSortDenoms(t *testing.T) {
	tests := []struct {
		denoms         []string
		expectedDenoms []string
	}{
		{[]string{"uCoinB", "uCoinA"}, []string{"uCoinA", "uCoinB"}},
		{[]string{"uCoinC", "uCoinA", "uCoinB"}, []string{"uCoinA", "uCoinB", "uCoinC"}},
		{[]string{"uCoinC", "uCoinA", "uCoinD", "uCoinB"}, []string{"uCoinA", "uCoinB", "uCoinC", "uCoinD"}},
	}

	for _, tc := range tests {
		sortedDenoms := types.SortDenoms(tc.denoms)
		require.Equal(t, tc.expectedDenoms, sortedDenoms)
	}
}

func TestGetPoolReserveAcc(t *testing.T) {
	poolName := types.PoolName([]string{"denomX", "denomY"}, 1)
	require.Equal(t, "denomX/denomY/1", poolName)
	reserveAcc := types.GetPoolReserveAcc(poolName)
	require.NotNil(t, reserveAcc)
	require.Equal(t, "cosmos16ddqestwukv0jzcyfn3fdfq9h2wrs83cr4rfm3", reserveAcc.String())
	require.Equal(t, "poolD35A0CC16EE598F90B044CE296A405BA9C381E38837599D96F2F70C2F02A23A4", types.GetPoolCoinDenom(poolName))
}

func TestGetPoolReserveAcc2(t *testing.T) {
	poolName := types.PoolName([]string{"stake", "token"}, 1)
	require.Equal(t, "stake/token/1", poolName)
	reserveAcc := types.GetPoolReserveAcc(poolName)
	require.NotNil(t, reserveAcc)
	require.Equal(t, "cosmos1unfxz7l7q0s3gmmthgwe3yljk0thhg57ym3p6u", reserveAcc.String())
	require.Equal(t, "poolE4D2617BFE03E1146F6BBA1D9893F2B3D77BA29E7ED532BB721A39FF1ECC1B07", types.GetPoolCoinDenom(poolName))
}

func TestGetPoolReserveAcc3(t *testing.T) {
	poolName := types.PoolName([]string{"uusd", "uatom"}, 1)
	require.Equal(t, "uatom/uusd/1", poolName)
	reserveAcc := types.GetPoolReserveAcc(poolName)
	require.NotNil(t, reserveAcc)
	require.Equal(t, "cosmos1jmhkafh94jpgakr735r70t32sxq9wzkayzs9we", reserveAcc.String())
	require.Equal(t, "pool96EF6EA6E5AC828ED87E8D07E7AE2A8180570ADD212117B2DA6F0B75D17A6295", types.GetPoolCoinDenom(poolName))
}

func TestGetCoinsTotalAmount(t *testing.T) {
	denomA := "uCoinA"
	denomB := "uCoinB"
	a := sdk.NewCoin(denomA, sdk.NewInt(100))
	b := sdk.NewCoin(denomB, sdk.NewInt(100))
	sum := types.GetCoinsTotalAmount(sdk.NewCoins(a, b))
	require.Equal(t, sdk.NewInt(200), sum)

	a = sdk.NewCoin(denomA, sdk.NewInt(100))
	b = sdk.NewCoin(denomB, sdk.NewInt(300))
	sum = types.GetCoinsTotalAmount(sdk.NewCoins(a, b))
	require.Equal(t, sdk.NewInt(400), sum)

	a = sdk.NewCoin(denomA, sdk.NewInt(500))
	sum = types.GetCoinsTotalAmount(sdk.NewCoins(a))
	require.Equal(t, sdk.NewInt(500), sum)
}

func TestValidateReserveCoinLimit(t *testing.T) {
	denomA := "uCoinA"
	denomB := "uCoinB"

	a := sdk.NewCoin(denomA, sdk.NewInt(1000000000000))
	b := sdk.NewCoin(denomB, sdk.NewInt(100))

	err := types.ValidateReserveCoinLimit(sdk.ZeroInt(), sdk.NewCoins(a, b))
	require.NoError(t, err)

	err = types.ValidateReserveCoinLimit(sdk.NewInt(1000000000000), sdk.NewCoins(a, b))
	require.Equal(t, types.ErrExceededReserveCoinLimit, err)

	a = sdk.NewCoin(denomA, sdk.NewInt(500000000000))
	b = sdk.NewCoin(denomB, sdk.NewInt(500000000000))
	err = types.ValidateReserveCoinLimit(sdk.NewInt(1000000000000), sdk.NewCoins(a, b))
	require.NoError(t, err)

	a = sdk.NewCoin(denomA, sdk.NewInt(500000000001))
	b = sdk.NewCoin(denomB, sdk.NewInt(500000000000))
	err = types.ValidateReserveCoinLimit(sdk.NewInt(1000000000000), sdk.NewCoins(a, b))
	require.Equal(t, types.ErrExceededReserveCoinLimit, err)
}
