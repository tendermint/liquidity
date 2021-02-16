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
	poolKey := types.GetPoolKey([]string{"denomX", "denomY"}, 1)
	require.Equal(t, "denomX/denomY/1", poolKey)
	reserveAcc := types.GetPoolReserveAcc(poolKey)
	require.NotNil(t, reserveAcc)
	require.Equal(t, "cosmos16ddqestwukv0jzcyfn3fdfq9h2wrs83cr4rfm3", reserveAcc.String())
	require.Equal(t, "pool/D35A0CC16EE598F90B044CE296A405BA9C381E38837599D96F2F70C2F02A23A4", types.GetPoolCoinDenom(poolKey))
}

func TestGetPoolReserveAcc2(t *testing.T) {
	poolKey := types.GetPoolKey([]string{"stake", "token"}, 1)
	require.Equal(t, "stake/token/1", poolKey)
	reserveAcc := types.GetPoolReserveAcc(poolKey)
	require.NotNil(t, reserveAcc)
	require.Equal(t, "cosmos1unfxz7l7q0s3gmmthgwe3yljk0thhg57ym3p6u", reserveAcc.String())
	require.Equal(t, "pool/E4D2617BFE03E1146F6BBA1D9893F2B3D77BA29E7ED532BB721A39FF1ECC1B07", types.GetPoolCoinDenom(poolKey))
}

func TestGetPoolReserveAcc3(t *testing.T) {
	poolKey := types.GetPoolKey([]string{"acoin", "bcoin"}, 1)
	require.Equal(t, "acoin/bcoin/1", poolKey)
	reserveAcc := types.GetPoolReserveAcc(poolKey)
	require.NotNil(t, reserveAcc)
	require.Equal(t, "cosmos19cwhfmgmdwv2tntlr5l30cwv6njjgsyd2528kv", reserveAcc.String())
	require.Equal(t, "pool/2E1D74ED1B6B98A5CD7F1D3F17E1CCD4E524408D5860FBD5A87CBC07C1BB9967", types.GetPoolCoinDenom(poolKey))
}

func TestIsPoolCoinDenom(t *testing.T) {
	poolKey := types.GetPoolKey([]string{"denomX", "denomY"}, 1)
	require.Equal(t, "denomX/denomY/1", poolKey)
	poolCoinDenom := types.GetPoolCoinDenom(poolKey)
	require.True(t, types.IsPoolCoinDenom(poolCoinDenom))
	require.True(t, types.IsPoolCoinDenom("pool/D35A0CC16EE598F90B044CE296A405BA9C381E38837599D96F2F70C2F02A23A4"))
	require.False(t, types.IsPoolCoinDenom("D35A0CC16EE598F90B044CE296A405BA9C381E38837599D96F2F70C2F02A23A4"))
	require.False(t, types.IsPoolCoinDenom("ibc/D35A0CC16EE598F90B044CE296A405BA9C381E38837599D96F2F70C2F02A23A4"))
	require.False(t, types.IsPoolCoinDenom("denomX/denomY/1"))
}
