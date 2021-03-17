package types_test

import (
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"

	"github.com/tendermint/liquidity/app"
	"github.com/tendermint/liquidity/x/liquidity/types"
)

func TestLiquidityPoolBatch(t *testing.T) {
	simapp, ctx := app.CreateTestInput()
	params := simapp.LiquidityKeeper.GetParams(ctx)
	pool := types.Pool{}
	require.Equal(t, types.ErrPoolNotExists, pool.Validate())
	pool.Id = 1
	require.Equal(t, types.ErrPoolTypeNotExists, pool.Validate())
	pool.TypeId = 1
	require.Equal(t, types.ErrNumOfReserveCoinDenoms, pool.Validate())
	pool.ReserveCoinDenoms = []string{DenomY, DenomX, DenomX}
	require.Equal(t, types.ErrNumOfReserveCoinDenoms, pool.Validate())
	pool.ReserveCoinDenoms = []string{DenomY, DenomX}
	require.Equal(t, types.ErrBadOrderingReserveCoinDenoms, pool.Validate())
	pool.ReserveCoinDenoms = []string{DenomX, DenomY}
	require.Equal(t, types.ErrEmptyReserveAccountAddress, pool.Validate())
	pool.ReserveAccountAddress = "badaddress"
	require.Equal(t, types.ErrBadReserveAccountAddress, pool.Validate())
	pool.ReserveAccountAddress = types.GetPoolReserveAcc(pool.Name()).String()
	add2, err := sdk.AccAddressFromBech32(pool.ReserveAccountAddress)
	require.Equal(t, add2, pool.GetReserveAccount())
	require.Equal(t, types.ErrEmptyPoolCoinDenom, pool.Validate())
	pool.PoolCoinDenom = "badPoolCoinDenom"
	require.Equal(t, types.ErrBadPoolCoinDenom, pool.Validate())
	pool.PoolCoinDenom = pool.Name()

	require.NoError(t, pool.Validate())

	require.Equal(t, pool.Name(), types.PoolName(pool.ReserveCoinDenoms, pool.TypeId))
	require.Equal(t, pool.Id, pool.GetPoolId())
	require.Equal(t, pool.PoolCoinDenom, pool.GetPoolCoinDenom())

	cdc := simapp.AppCodec()
	poolByte := types.MustMarshalPool(cdc, pool)
	require.Equal(t, pool, types.MustUnmarshalPool(cdc, poolByte))
	poolByte = types.MustMarshalPool(cdc, pool)
	poolMarshaled, err := types.UnmarshalPool(cdc, poolByte)
	require.NoError(t, err)
	require.Equal(t, pool, poolMarshaled)

	addr, err := sdk.AccAddressFromBech32(pool.ReserveAccountAddress)
	require.NoError(t, err)
	require.True(t, pool.GetReserveAccount().Equals(addr))

	require.Equal(t, strings.TrimSpace(pool.String()+"\n"+pool.String()), types.Pools{pool, pool}.String())

	simapp.LiquidityKeeper.SetPool(ctx, pool)
	batch := types.NewPoolBatch(pool.Id, 1)
	simapp.LiquidityKeeper.SetPoolBatch(ctx, batch)
	simapp.LiquidityKeeper.SetPoolBatchIndex(ctx, batch.PoolId, batch.Index)

	batchByte := types.MustMarshalPoolBatch(cdc, batch)
	require.Equal(t, batch, types.MustUnmarshalPoolBatch(cdc, batchByte))
	batchMarshaled, err := types.UnmarshalPoolBatch(cdc, batchByte)
	require.NoError(t, err)
	require.Equal(t, batch, batchMarshaled)

	batchDepositMsg := types.DepositMsgState{}
	batchWithdrawMsg := types.WithdrawMsgState{}
	batchSwapMsg := types.SwapMsgState{ExchangedOfferCoin: sdk.NewCoin("test", sdk.NewInt(1000)),
		RemainingOfferCoin: sdk.NewCoin("test", sdk.NewInt(1000)), ReservedOfferCoinFee: types.GetOfferCoinFee(sdk.NewCoin("test", sdk.NewInt(2000)), params.SwapFeeRate)}

	byte := types.MustMarshalDepositMsgState(cdc, batchDepositMsg)
	require.Equal(t, batchDepositMsg, types.MustUnmarshalDepositMsgState(cdc, byte))
	marshaled, err := types.UnmarshalDepositMsgState(cdc, byte)
	require.NoError(t, err)
	require.Equal(t, batchDepositMsg, marshaled)

	byte = types.MustMarshalWithdrawMsgState(cdc, batchWithdrawMsg)
	require.Equal(t, batchWithdrawMsg, types.MustUnmarshalWithdrawMsgState(cdc, byte))
	withdrawMsgMarshaled, err := types.UnmarshalWithdrawMsgState(cdc, byte)
	require.NoError(t, err)
	require.Equal(t, batchWithdrawMsg, withdrawMsgMarshaled)

	byte = types.MustMarshalSwapMsgState(cdc, batchSwapMsg)
	require.Equal(t, batchSwapMsg, types.MustUnmarshalSwapMsgState(cdc, byte))
	SwapMsgMarshaled, err := types.UnmarshalSwapMsgState(cdc, byte)
	require.NoError(t, err)
	require.Equal(t, batchSwapMsg, SwapMsgMarshaled)
}
