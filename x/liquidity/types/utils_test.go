package types_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"

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
	testCases := []struct {
		coins        sdk.Coins
		expectResult sdk.Int
	}{
		{
			coins:        sdk.NewCoins(sdk.NewCoin("uCoinA", sdk.NewInt(100)), sdk.NewCoin("uCoinB", sdk.NewInt(100))),
			expectResult: sdk.NewInt(200),
		},
		{
			coins:        sdk.NewCoins(sdk.NewCoin("uCoinA", sdk.NewInt(100)), sdk.NewCoin("uCoinB", sdk.NewInt(300))),
			expectResult: sdk.NewInt(400),
		},
		{
			coins:        sdk.NewCoins(sdk.NewCoin("uCoinA", sdk.NewInt(500))),
			expectResult: sdk.NewInt(500),
		},
	}

	for _, tc := range testCases {
		totalAmount := types.GetCoinsTotalAmount(tc.coins)
		require.Equal(t, tc.expectResult, totalAmount)
	}
}

func TestValidateReserveCoinLimit(t *testing.T) {
	testCases := []struct {
		name                 string
		maxReserveCoinAmount sdk.Int
		depositCoins         sdk.Coins
		expectErr            bool
	}{
		{
			name:                 "valid case",
			maxReserveCoinAmount: sdk.ZeroInt(), // 0 means unlimited amount
			depositCoins:         sdk.NewCoins(sdk.NewCoin("uCoinA", sdk.NewInt(100_000_000_000)), sdk.NewCoin("uCoinB", sdk.NewInt(100))),
			expectErr:            false,
		},
		{
			name:                 "valid case",
			maxReserveCoinAmount: sdk.NewInt(1_000_000_000_000),
			depositCoins:         sdk.NewCoins(sdk.NewCoin("uCoinA", sdk.NewInt(500_000_000_000)), sdk.NewCoin("uCoinB", sdk.NewInt(500_000_000_000))),
			expectErr:            false,
		},
		{
			name:                 "negative value of max reserve coin amount",
			maxReserveCoinAmount: sdk.NewInt(-100),
			depositCoins:         sdk.NewCoins(sdk.NewCoin("uCoinA", sdk.NewInt(100_000_000_000)), sdk.NewCoin("uCoinB", sdk.NewInt(100))),
			expectErr:            true,
		},
		{
			name:                 "cannot exceed reserve coin limit amount",
			maxReserveCoinAmount: sdk.NewInt(1_000_000_000_000),
			depositCoins:         sdk.NewCoins(sdk.NewCoin("uCoinA", sdk.NewInt(1_000_000_000_000)), sdk.NewCoin("uCoinB", sdk.NewInt(100))),
			expectErr:            true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectErr {
				err := types.ValidateReserveCoinLimit(tc.maxReserveCoinAmount, tc.depositCoins)
				require.Equal(t, types.ErrExceededReserveCoinLimit, err)
			} else {
				err := types.ValidateReserveCoinLimit(tc.maxReserveCoinAmount, tc.depositCoins)
				require.NoError(t, err)
			}
		})
	}
}
