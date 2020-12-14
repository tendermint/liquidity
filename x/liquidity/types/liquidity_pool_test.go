package types_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/liquidity/app"
	"github.com/tendermint/liquidity/x/liquidity/types"
	"strings"
	"testing"
)

func TestLiquidityPoolBatch(t *testing.T) {
	simapp, ctx := app.CreateTestInput()
	pool := types.LiquidityPool{}
	require.Equal(t, types.ErrPoolNotExists, pool.Validate())
	pool.PoolId = 1
	require.Equal(t, types.ErrPoolTypeNotExists, pool.Validate())
	pool.PoolTypeIndex = 1
	require.Equal(t, types.ErrNumOfReserveCoinDenoms, pool.Validate())
	pool.ReserveCoinDenoms = []string{DenomY, DenomX, DenomX}
	require.Equal(t, types.ErrNumOfReserveCoinDenoms, pool.Validate())
	pool.ReserveCoinDenoms = []string{DenomY, DenomX}
	require.Equal(t, types.ErrBadOrderingReserveCoinDenoms, pool.Validate())
	pool.ReserveCoinDenoms = []string{DenomX, DenomY}
	require.Equal(t, types.ErrEmptyReserveAccountAddress, pool.Validate())
	pool.ReserveAccountAddress = "badaddress"
	require.Equal(t, types.ErrBadReserveAccountAddress, pool.Validate())
	pool.ReserveAccountAddress = types.GetPoolReserveAcc(pool.GetPoolKey()).String()
	add2, err := sdk.AccAddressFromBech32(pool.ReserveAccountAddress)
	require.Equal(t, add2, pool.GetReserveAccount())
	require.Equal(t, types.ErrEmptyPoolCoinDenom, pool.Validate())
	pool.PoolCoinDenom = "badPoolCoinDenom"
	require.Equal(t, types.ErrBadPoolCoinDenom, pool.Validate())
	pool.PoolCoinDenom = pool.GetPoolKey()

	require.NoError(t, pool.Validate())

	require.Equal(t, pool.GetPoolKey(), types.GetPoolKey(pool.ReserveCoinDenoms, pool.PoolTypeIndex))
	require.Equal(t, pool.PoolId, pool.GetPoolId())
	require.Equal(t, pool.PoolCoinDenom, pool.GetPoolCoinDenom())

	cdc := simapp.AppCodec()
	poolByte := types.MustMarshalLiquidityPool(cdc, pool)
	require.Equal(t, pool, types.MustUnmarshalLiquidityPool(cdc, poolByte))
	poolByte = types.MustMarshalLiquidityPool(cdc, pool)
	poolMarshaled, err := types.UnmarshalLiquidityPool(cdc, poolByte)
	require.NoError(t, err)
	require.Equal(t, pool, poolMarshaled)

	addr, err := sdk.AccAddressFromBech32(pool.ReserveAccountAddress)
	require.NoError(t, err)
	require.True(t, pool.GetReserveAccount().Equals(addr))

	require.Equal(t, strings.TrimSpace(pool.String()+"\n"+pool.String()), types.LiquidityPools{pool, pool}.String())

	simapp.LiquidityKeeper.SetLiquidityPool(ctx, pool)
	batch := types.NewLiquidityPoolBatch(pool.PoolId, 1)
	simapp.LiquidityKeeper.SetLiquidityPoolBatch(ctx, batch)
	simapp.LiquidityKeeper.SetLiquidityPoolBatchIndex(ctx, batch.PoolId, batch.BatchIndex)

	batchByte := types.MustMarshalLiquidityPoolBatch(cdc, batch)
	require.Equal(t, batch, types.MustUnmarshalLiquidityPoolBatch(cdc, batchByte))
	batchMarshaled, err := types.UnmarshalLiquidityPoolBatch(cdc, batchByte)
	require.NoError(t, err)
	require.Equal(t, batch, batchMarshaled)

	batchDepositMsg := types.BatchPoolDepositMsg{}
	batchWithdrawMsg := types.BatchPoolWithdrawMsg{}
	batchSwapMsg := types.BatchPoolSwapMsg{ExchangedOfferCoin: sdk.NewCoin("test", sdk.NewInt(1000)),
		RemainingOfferCoin: sdk.NewCoin("test", sdk.NewInt(1000)), OfferCoinFeeReserve: types.GetOfferCoinFee(sdk.NewCoin("test", sdk.NewInt(2000)))}

	byte := types.MustMarshalBatchPoolDepositMsg(cdc, batchDepositMsg)
	require.Equal(t, batchDepositMsg, types.MustUnmarshalBatchPoolDepositMsg(cdc, byte))
	marshaled, err := types.UnmarshalBatchPoolDepositMsg(cdc, byte)
	require.NoError(t, err)
	require.Equal(t, batchDepositMsg, marshaled)

	byte = types.MustMarshalBatchPoolWithdrawMsg(cdc, batchWithdrawMsg)
	require.Equal(t, batchWithdrawMsg, types.MustUnmarshalBatchPoolWithdrawMsg(cdc, byte))
	withdrawMsgMarshaled, err := types.UnmarshalBatchPoolWithdrawMsg(cdc, byte)
	require.NoError(t, err)
	require.Equal(t, batchWithdrawMsg, withdrawMsgMarshaled)

	byte = types.MustMarshalBatchPoolSwapMsg(cdc, batchSwapMsg)
	require.Equal(t, batchSwapMsg, types.MustUnmarshalBatchPoolSwapMsg(cdc, byte))
	SwapMsgMarshaled, err := types.UnmarshalBatchPoolSwapMsg(cdc, byte)
	require.NoError(t, err)
	require.Equal(t, batchSwapMsg, SwapMsgMarshaled)
}
