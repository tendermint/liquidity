package types_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/liquidity/x/liquidity/types"
	"github.com/tendermint/tendermint/crypto"
	"testing"
)

const (
	DefaultPoolTypeIndex = uint32(1)
	DefaultPoolId = uint64(1)
	DefaultSwapType = uint32(1)
	DenomX = "denomX"
	DenomY = "denomY"
	DenomPoolCoin = "denomPoolCoin"

)

func TestMsgCreateLiquidityPool(t *testing.T) {
	addr := sdk.AccAddress(crypto.AddressHash([]byte("testAccount")))
	denoms := []string{DenomX, DenomY}
	coins := sdk.NewCoins(sdk.NewCoin(DenomX, sdk.NewInt(1000)), sdk.NewCoin(DenomY, sdk.NewInt(1000)))
	msg := types.NewMsgCreateLiquidityPool(addr, DefaultPoolTypeIndex, denoms, coins)
	require.IsType(t, &types.MsgCreateLiquidityPool{}, msg)
	require.Equal(t, types.RouterKey, msg.Route())
	require.Equal(t, types.TypeMsgCreateLiquidityPool, msg.Type())

	err := msg.ValidateBasic()
	require.NoError(t, err)
	signers := msg.GetSigners()
	require.Len(t, signers, 1)
	require.Equal(t, msg.GetPoolCreator(), signers[0])
	require.Equal(t, sdk.MustSortJSON(types.ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
}

func TestMsgDepositToLiquidityPool(t *testing.T) {
	addr := sdk.AccAddress(crypto.AddressHash([]byte("testAccount")))
	coins := sdk.NewCoins(sdk.NewCoin(DenomX, sdk.NewInt(1000)), sdk.NewCoin(DenomY, sdk.NewInt(1000)))
	msg := types.NewMsgDepositToLiquidityPool(addr, DefaultPoolId, coins)
	require.IsType(t, &types.MsgDepositToLiquidityPool{}, msg)
	require.Equal(t, types.RouterKey, msg.Route())
	require.Equal(t, types.TypeMsgDepositToLiquidityPool, msg.Type())

	err := msg.ValidateBasic()
	require.NoError(t, err)
	signers := msg.GetSigners()
	require.Len(t, signers, 1)
	require.Equal(t, msg.GetDepositor(), signers[0])
	require.Equal(t, sdk.MustSortJSON(types.ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
}
func TestMsgWithdrawFromLiquidityPool(t *testing.T) {
	addr := sdk.AccAddress(crypto.AddressHash([]byte("testAccount")))
	coin := sdk.NewCoin(DenomPoolCoin, sdk.NewInt(1000))
	msg := types.NewMsgWithdrawFromLiquidityPool(addr, DefaultPoolId, coin)
	require.IsType(t, &types.MsgWithdrawFromLiquidityPool{}, msg)
	require.Equal(t, types.RouterKey, msg.Route())
	require.Equal(t, types.TypeMsgWithdrawFromLiquidityPool, msg.Type())

	err := msg.ValidateBasic()
	require.NoError(t, err)
	signers := msg.GetSigners()
	require.Len(t, signers, 1)
	require.Equal(t, msg.GetWithdrawer(), signers[0])
	require.Equal(t, sdk.MustSortJSON(types.ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
}
func TestMsgSwap(t *testing.T) {
	addr := sdk.AccAddress(crypto.AddressHash([]byte("testAccount")))
	coin := sdk.NewCoin(DenomX, sdk.NewInt(1000))
	orderPrice, err := sdk.NewDecFromStr("0.1")
	require.NoError(t, err)
	msg := types.NewMsgSwap(addr, DefaultPoolId, DefaultPoolTypeIndex, DefaultSwapType, coin, DenomY, orderPrice)
	require.IsType(t, &types.MsgSwap{}, msg)
	require.Equal(t, types.RouterKey, msg.Route())
	require.Equal(t, types.TypeMsgSwap, msg.Type())

	err = msg.ValidateBasic()
	require.NoError(t, err)
	signers := msg.GetSigners()
	require.Len(t, signers, 1)
	require.Equal(t, msg.GetSwapRequester(), signers[0])
	require.Equal(t, sdk.MustSortJSON(types.ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
}