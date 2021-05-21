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
		{
			denoms:         []string{"uCoinB", "uCoinA"},
			expectedDenoms: []string{"uCoinA", "uCoinB"},
		},
		{
			denoms:         []string{"uCoinC", "uCoinA", "uCoinB"},
			expectedDenoms: []string{"uCoinA", "uCoinB", "uCoinC"},
		},
		{
			denoms:         []string{"uCoinC", "uCoinA", "uCoinD", "uCoinB"},
			expectedDenoms: []string{"uCoinA", "uCoinB", "uCoinC", "uCoinD"},
		},
	}

	for _, tc := range tests {
		sortedDenoms := types.SortDenoms(tc.denoms)
		require.Equal(t, tc.expectedDenoms, sortedDenoms)
	}
}

func TestGetPoolInformation(t *testing.T) {
	testCases := []struct {
		reserveCoinDenoms     []string
		poolTypeId            uint32
		expectedPoolName      string
		expectedReserveAcc    string
		expectedPoolCoinDenom string
	}{
		{
			reserveCoinDenoms:     []string{"denomX", "denomY"},
			poolTypeId:            uint32(1),
			expectedPoolName:      "denomX/denomY/1",
			expectedReserveAcc:    "cosmos16ddqestwukv0jzcyfn3fdfq9h2wrs83cr4rfm3",
			expectedPoolCoinDenom: "poolD35A0CC16EE598F90B044CE296A405BA9C381E38837599D96F2F70C2F02A23A4",
		},
		{
			reserveCoinDenoms:     []string{"stake", "token"},
			poolTypeId:            uint32(1),
			expectedPoolName:      "stake/token/1",
			expectedReserveAcc:    "cosmos1unfxz7l7q0s3gmmthgwe3yljk0thhg57ym3p6u",
			expectedPoolCoinDenom: "poolE4D2617BFE03E1146F6BBA1D9893F2B3D77BA29E7ED532BB721A39FF1ECC1B07",
		},
		{
			reserveCoinDenoms:     []string{"uatom", "uusd"},
			poolTypeId:            uint32(2),
			expectedPoolName:      "uatom/uusd/2",
			expectedReserveAcc:    "cosmos1xqm0g09czvdp5c7jk0fmz85u7maz52m040eh8g",
			expectedPoolCoinDenom: "pool3036F43CB8131A1A63D2B3D3B11E9CF6FA2A2B6FEC17D5AD283C25C939614A8C",
		},
	}

	for _, tc := range testCases {
		poolName := types.PoolName(tc.reserveCoinDenoms, tc.poolTypeId)
		require.Equal(t, tc.expectedPoolName, poolName)

		reserveAcc := types.GetPoolReserveAcc(poolName)
		require.Equal(t, tc.expectedReserveAcc, reserveAcc.String())

		poolCoinDenom := types.GetPoolCoinDenom(poolName)
		require.Equal(t, tc.expectedPoolCoinDenom, poolCoinDenom)
	}
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
