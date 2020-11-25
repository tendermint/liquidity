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
	DefaultPoolId        = uint64(1)
	DefaultSwapType      = uint32(1)
	DenomX               = "denomX"
	DenomY               = "denomY"
	DenomPoolCoin        = "denomPoolCoin"
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

	// Fail cases
	msg = types.NewMsgCreateLiquidityPool(sdk.AccAddress{}, DefaultPoolTypeIndex, denoms, coins)
	err = msg.ValidateBasic()
	require.Error(t, err)
	denomsFail := []string{DenomX, DenomY, DenomY}
	msg = types.NewMsgCreateLiquidityPool(addr, DefaultPoolTypeIndex, denomsFail, coins)
	err = msg.ValidateBasic()
	require.Error(t, err)
	coinsFail := sdk.NewCoins(sdk.NewCoin(DenomY, sdk.NewInt(1000)))
	msg = types.NewMsgCreateLiquidityPool(addr, DefaultPoolTypeIndex, denoms, coinsFail)
	err = msg.ValidateBasic()
	require.Error(t, err)
	coinsFail = sdk.NewCoins(sdk.NewCoin(DenomX, sdk.NewInt(1000)), sdk.NewCoin(DenomY, sdk.NewInt(1000)), sdk.NewCoin("Denomfail", sdk.NewInt(1000)))
	msg = types.NewMsgCreateLiquidityPool(addr, DefaultPoolTypeIndex, denoms, coinsFail)
	err = msg.ValidateBasic()
	require.Error(t, err)
	coinsFail = sdk.NewCoins(sdk.NewCoin(DenomX, sdk.NewInt(0)), sdk.NewCoin(DenomY, sdk.NewInt(1000)))
	msg = types.NewMsgCreateLiquidityPool(addr, DefaultPoolTypeIndex, denoms, coinsFail)
	err = msg.ValidateBasic()
	require.Error(t, err)
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

	// Fail case
	coinsFail := sdk.NewCoins(sdk.NewCoin(DenomX, sdk.NewInt(1000)), sdk.NewCoin(DenomY, sdk.NewInt(1000)), sdk.NewCoin("Denomfail", sdk.NewInt(1000)))
	msg = types.NewMsgDepositToLiquidityPool(addr, DefaultPoolId, coinsFail)
	err = msg.ValidateBasic()
	require.Error(t, err)
	coinsFail = sdk.NewCoins(sdk.NewCoin(DenomX, sdk.NewInt(0)), sdk.NewCoin(DenomY, sdk.NewInt(1000)))
	msg = types.NewMsgDepositToLiquidityPool(addr, DefaultPoolId, coinsFail)
	err = msg.ValidateBasic()
	require.Error(t, err)
	msg = types.NewMsgDepositToLiquidityPool(sdk.AccAddress{}, DefaultPoolId, coins)
	err = msg.ValidateBasic()
	require.Error(t, err)
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

	// Fail case
	coinFail := sdk.NewCoin("testPoolCoin", sdk.NewInt(0))
	msg = types.NewMsgWithdrawFromLiquidityPool(addr, DefaultPoolId, coinFail)
	err = msg.ValidateBasic()
	require.Error(t, err)
	msg = types.NewMsgWithdrawFromLiquidityPool(sdk.AccAddress{}, DefaultPoolId, coin)
	err = msg.ValidateBasic()
	require.Error(t, err)
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

	// Fail case
	msg = types.NewMsgSwap(addr, DefaultPoolId, DefaultPoolTypeIndex, DefaultSwapType, coin, DenomY, sdk.ZeroDec())
	err = msg.ValidateBasic()
	require.Error(t, err)
	coinFail := sdk.NewCoin("testPoolCoin", sdk.NewInt(0))
	msg = types.NewMsgSwap(addr, DefaultPoolId, DefaultPoolTypeIndex, DefaultSwapType, coinFail, DenomY, orderPrice)
	err = msg.ValidateBasic()
	require.Error(t, err)
	msg = types.NewMsgSwap(sdk.AccAddress{}, DefaultPoolId, DefaultPoolTypeIndex, DefaultSwapType, coin, DenomY, orderPrice)
	err = msg.ValidateBasic()
	require.Error(t, err)
}
