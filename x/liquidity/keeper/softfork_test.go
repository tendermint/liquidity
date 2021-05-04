package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	lapp "github.com/tendermint/liquidity/app"
	"github.com/tendermint/liquidity/x/liquidity"
	"github.com/tendermint/liquidity/x/liquidity/types"
)

func TestSoftFork(t *testing.T) {
	simapp, ctx := createTestInput()
	simapp.LiquidityKeeper.SetParams(ctx, types.DefaultParams())
	params := simapp.LiquidityKeeper.GetParams(ctx)

	poolTypeId := types.DefaultPoolTypeId
	addrs := lapp.AddTestAddrs(simapp, ctx, 3, params.PoolCreationFee)

	denomA := "uatom"
	denomB := "utest"
	denomA, denomB = types.AlphabeticalDenomPair(denomA, denomB)

	deposit := sdk.NewCoins(sdk.NewCoin(denomA, sdk.NewInt(100*1000000)), sdk.NewCoin(denomB, sdk.NewInt(2000*1000000)))
	lapp.SaveAccount(simapp, ctx, addrs[0], deposit)

	depositA := simapp.BankKeeper.GetBalance(ctx, addrs[0], denomA)
	depositB := simapp.BankKeeper.GetBalance(ctx, addrs[0], denomB)
	depositBalance := sdk.NewCoins(depositA, depositB)

	require.Equal(t, deposit, depositBalance)

	msg := types.NewMsgCreatePool(addrs[0], poolTypeId, depositBalance)
	_, err := simapp.LiquidityKeeper.CreatePool(ctx, msg)
	require.NoError(t, err)

	invalidMsg := types.NewMsgCreatePool(addrs[0], 0, depositBalance)
	_, err = simapp.LiquidityKeeper.CreatePool(ctx, invalidMsg)
	require.Error(t, err, types.ErrBadPoolTypeId)

	pools := simapp.LiquidityKeeper.GetAllPools(ctx)
	require.Equal(t, 1, len(pools))
	require.Equal(t, uint64(1), pools[0].Id)
	require.Equal(t, uint64(1), simapp.LiquidityKeeper.GetNextPoolId(ctx)-1)
	require.Equal(t, denomA, pools[0].ReserveCoinDenoms[0])
	require.Equal(t, denomB, pools[0].ReserveCoinDenoms[1])

	poolCoin := simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, pools[0])
	creatorBalance := simapp.BankKeeper.GetBalance(ctx, addrs[0], pools[0].PoolCoinDenom)
	require.Equal(t, poolCoin, creatorBalance.Amount)

	_, err = simapp.LiquidityKeeper.CreatePool(ctx, msg)
	require.Error(t, err, types.ErrPoolAlreadyExists)

	providerAcc, err := sdk.AccAddressFromBech32(types.Airdrop1ProviderAddr)
	require.Nil(t, err)
	totalDistribution := sdk.NewCoin(types.Airdrop1DistributionCoin.Denom, types.Airdrop1DistributionCoin.Amount.MulRaw(int64(len(types.Airdrop1TargetAddrs))))
	lapp.SaveAccount(simapp, ctx, providerAcc, sdk.Coins{totalDistribution})

	// Airdrop1 Success Case
	liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)
	providerBalanceBeforeDistribution := simapp.BankKeeper.GetBalance(ctx, providerAcc, types.Airdrop1DistributionCoin.Denom)
	err = simapp.LiquidityKeeper.SoftForkAirdrop(ctx, types.Airdrop1ProviderAddr, types.Airdrop1TargetAddrs, types.Airdrop1DistributionCoin)
	require.NoError(t, err)
	liquidity.EndBlocker(ctx, simapp.LiquidityKeeper)
	providerBalanceAfterDistribution := simapp.BankKeeper.GetBalance(ctx, providerAcc, types.Airdrop1DistributionCoin.Denom)
	require.Equal(t, providerBalanceBeforeDistribution.Sub(providerBalanceAfterDistribution), totalDistribution)
	for _, addr := range types.Airdrop1TargetAddrs {
		acc, _ := sdk.AccAddressFromBech32(addr)
		balance := simapp.BankKeeper.GetBalance(ctx, acc, types.Airdrop1DistributionCoin.Denom)
		require.Equal(t, balance, types.Airdrop1DistributionCoin)
	}

	// Airdrop1 Fail Case, insufficient balances of provider account for softfork distribution
	bankState := simapp.BankKeeper.ExportGenesis(ctx)
	liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)
	providerBalanceBeforeDistribution = simapp.BankKeeper.GetBalance(ctx, providerAcc, types.Airdrop1DistributionCoin.Denom)
	err = simapp.LiquidityKeeper.SoftForkAirdrop(ctx, types.Airdrop1ProviderAddr, types.Airdrop1TargetAddrs, types.Airdrop1DistributionCoin)
	require.Error(t, err)
	liquidity.EndBlocker(ctx, simapp.LiquidityKeeper)
	// assert no changes
	bankStateAfterFail := simapp.BankKeeper.ExportGenesis(ctx)
	require.Equal(t, bankState, bankStateAfterFail)

	// Airdrop1 Fail Case, wrong address
	lapp.SaveAccount(simapp, ctx, providerAcc, sdk.Coins{totalDistribution})
	bankState = simapp.BankKeeper.ExportGenesis(ctx)
	liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)
	providerBalanceBeforeDistribution = simapp.BankKeeper.GetBalance(ctx, providerAcc, types.Airdrop1DistributionCoin.Denom)
	err = simapp.LiquidityKeeper.SoftForkAirdrop(ctx, "cosmos1...wrongAddr", types.Airdrop1TargetAddrs, types.Airdrop1DistributionCoin)
	require.Error(t, err)
	liquidity.EndBlocker(ctx, simapp.LiquidityKeeper)
	// assert no changes
	bankStateAfterFail = simapp.BankKeeper.ExportGenesis(ctx)
	require.Equal(t, bankState, bankStateAfterFail)

	// Airdrop1 Fail Case, wrong address, empty address string is not allowed
	lapp.SaveAccount(simapp, ctx, providerAcc, sdk.Coins{totalDistribution})
	bankState = simapp.BankKeeper.ExportGenesis(ctx)
	liquidity.BeginBlocker(ctx, simapp.LiquidityKeeper)
	providerBalanceBeforeDistribution = simapp.BankKeeper.GetBalance(ctx, providerAcc, types.Airdrop1DistributionCoin.Denom)
	err = simapp.LiquidityKeeper.SoftForkAirdrop(ctx, types.Airdrop1ProviderAddr, []string{""}, types.Airdrop1DistributionCoin)
	require.Error(t, err)
	liquidity.EndBlocker(ctx, simapp.LiquidityKeeper)
	// assert no changes
	bankStateAfterFail = simapp.BankKeeper.ExportGenesis(ctx)
	require.Equal(t, bankState, bankStateAfterFail)
}

