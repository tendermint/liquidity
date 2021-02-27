package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"

	"github.com/tendermint/liquidity/x/liquidity/types"
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
	coins := sdk.NewCoins(sdk.NewCoin(DenomX, sdk.NewInt(1000)), sdk.NewCoin(DenomY, sdk.NewInt(1000)))
	msg := types.NewMsgCreateLiquidityPool(addr, DefaultPoolTypeIndex, coins)
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
	msg = types.NewMsgCreateLiquidityPool(sdk.AccAddress{}, DefaultPoolTypeIndex, coins)
	err = msg.ValidateBasic()
	require.Error(t, err)
	coinsFail := sdk.NewCoins(sdk.NewCoin(DenomY, sdk.NewInt(1000)))
	msg = types.NewMsgCreateLiquidityPool(addr, DefaultPoolTypeIndex, coinsFail)
	err = msg.ValidateBasic()
	require.Error(t, err)
	coinsFail = sdk.NewCoins(sdk.NewCoin(DenomX, sdk.NewInt(1000)), sdk.NewCoin(DenomY, sdk.NewInt(1000)), sdk.NewCoin("Denomfail", sdk.NewInt(1000)))
	msg = types.NewMsgCreateLiquidityPool(addr, DefaultPoolTypeIndex, coinsFail)
	err = msg.ValidateBasic()
	require.Error(t, err)
	coinsFail = sdk.NewCoins(sdk.NewCoin(DenomX, sdk.NewInt(0)), sdk.NewCoin(DenomY, sdk.NewInt(1000)))
	msg = types.NewMsgCreateLiquidityPool(addr, DefaultPoolTypeIndex, coinsFail)
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
	msg := types.NewMsgSwap(addr, DefaultPoolId, DefaultSwapType, coin, DenomY, orderPrice, types.DefaultSwapFeeRate)
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

func TestMsgPanics(t *testing.T) {
	emptyMsgCreatePool := types.MsgCreateLiquidityPool{}
	emptyMsgDeposit := types.MsgDepositToLiquidityPool{}
	emptyMsgWithdraw := types.MsgWithdrawFromLiquidityPool{}
	emptyMsgSwap := types.MsgSwap{}
	for _, msg := range []sdk.Msg{&emptyMsgCreatePool, &emptyMsgDeposit, &emptyMsgWithdraw, &emptyMsgSwap} {
		require.PanicsWithError(t, "empty address string is not allowed", func() { msg.GetSigners() })
	}
	for _, tc := range []func() sdk.AccAddress{
		emptyMsgCreatePool.GetPoolCreator,
		emptyMsgDeposit.GetDepositor,
		emptyMsgWithdraw.GetWithdrawer,
		emptyMsgSwap.GetSwapRequester,
	} {
		require.PanicsWithError(t, "empty address string is not allowed", func() { tc() })
	}
}

func TestMsgValidateBasic(t *testing.T) {
	validPoolTypeIndex := DefaultPoolTypeIndex
	validAddr := sdk.AccAddress(crypto.AddressHash([]byte("testAccount"))).String()
	validCoin := sdk.NewCoin(DenomY, sdk.NewInt(10000))

	invalidDenomCoin := sdk.Coin{Denom: "-", Amount: sdk.NewInt(10000)}
	negativeCoin := sdk.Coin{Denom: DenomX, Amount: sdk.NewInt(-1)}
	zeroCoin := sdk.Coin{Denom: DenomX, Amount: sdk.ZeroInt()}

	coinsWithInvalidDenom := sdk.Coins{invalidDenomCoin, validCoin}
	coinsWithNegative := sdk.Coins{negativeCoin, validCoin}
	coinsWithZero := sdk.Coins{zeroCoin, validCoin}

	invalidDenomErrMsg := "invalid denom: -"
	negativeCoinErrMsg := "coin -1denomX amount is not positive"
	negativeAmountErrMsg := "negative coin amount: -1"
	zeroCoinErrMsg := "coin 0denomX amount is not positive"

	t.Run("MsgCreateLiquidityPool", func(t *testing.T) {
		for _, tc := range []struct {
			msg    types.MsgCreateLiquidityPool
			errMsg string
		}{
			{
				types.MsgCreateLiquidityPool{},
				types.ErrBadPoolTypeIndex.Error(),
			},
			{
				types.MsgCreateLiquidityPool{PoolTypeIndex: validPoolTypeIndex},
				types.ErrEmptyPoolCreatorAddr.Error(),
			},
			{
				types.MsgCreateLiquidityPool{PoolCreatorAddress: validAddr, PoolTypeIndex: validPoolTypeIndex},
				types.ErrNumOfReserveCoin.Error(),
			},
			{
				types.MsgCreateLiquidityPool{
					PoolCreatorAddress: validAddr,
					PoolTypeIndex:      validPoolTypeIndex,
					DepositCoins:       coinsWithInvalidDenom,
				},
				invalidDenomErrMsg,
			},
			{
				types.MsgCreateLiquidityPool{
					PoolCreatorAddress: validAddr,
					PoolTypeIndex:      validPoolTypeIndex,
					DepositCoins:       coinsWithNegative,
				},
				negativeCoinErrMsg,
			},
			{
				types.MsgCreateLiquidityPool{
					PoolCreatorAddress: validAddr,
					PoolTypeIndex:      validPoolTypeIndex,
					DepositCoins:       coinsWithZero,
				},
				zeroCoinErrMsg,
			},
			{
				types.MsgCreateLiquidityPool{
					PoolCreatorAddress: validAddr,
					PoolTypeIndex:      validPoolTypeIndex,
					DepositCoins:       sdk.NewCoins(sdk.NewCoin(DenomX, sdk.NewInt(int64(types.MinReserveCoinNum)-1))),
				},
				types.ErrNumOfReserveCoin.Error(),
			},
			{
				types.MsgCreateLiquidityPool{
					PoolCreatorAddress: validAddr,
					PoolTypeIndex:      validPoolTypeIndex,
					DepositCoins:       sdk.NewCoins(sdk.NewCoin(DenomX, sdk.NewInt(int64(types.MaxReserveCoinNum)+1))),
				},
				types.ErrNumOfReserveCoin.Error(),
			},
		} {
			err := tc.msg.ValidateBasic()
			require.EqualError(t, err, tc.errMsg)
		}
	})
	t.Run("MsgDepositToLiquidityPool", func(t *testing.T) {
		for _, tc := range []struct {
			msg    types.MsgDepositToLiquidityPool
			errMsg string
		}{
			{
				types.MsgDepositToLiquidityPool{},
				types.ErrEmptyDepositorAddr.Error(),
			},
			{
				types.MsgDepositToLiquidityPool{DepositorAddress: validAddr},
				types.ErrBadDepositCoinsAmount.Error(),
			},
			{
				types.MsgDepositToLiquidityPool{DepositorAddress: validAddr, DepositCoins: coinsWithInvalidDenom},
				invalidDenomErrMsg,
			},
			{
				types.MsgDepositToLiquidityPool{DepositorAddress: validAddr, DepositCoins: coinsWithNegative},
				negativeCoinErrMsg,
			},
			{
				types.MsgDepositToLiquidityPool{DepositorAddress: validAddr, DepositCoins: coinsWithZero},
				zeroCoinErrMsg,
			},
			{
				types.MsgDepositToLiquidityPool{
					DepositorAddress: validAddr,
					DepositCoins: sdk.NewCoins(sdk.NewCoin(DenomX, sdk.NewInt(int64(types.MinReserveCoinNum)-1))),
				},
				types.ErrNumOfReserveCoin.Error(),
			},
			{
				types.MsgDepositToLiquidityPool{
					DepositorAddress: validAddr,
					DepositCoins: sdk.NewCoins(sdk.NewCoin(DenomX, sdk.NewInt(int64(types.MaxReserveCoinNum)+1))),
				},
				types.ErrNumOfReserveCoin.Error(),
			},
		} {
			err := tc.msg.ValidateBasic()
			require.EqualError(t, err, tc.errMsg)
		}
	})
	t.Run("MsgWithdrawFromLiquidityPool", func(t *testing.T) {
		for _, tc := range []struct {
			msg types.MsgWithdrawFromLiquidityPool
			errMsg string
		} {
			{
				types.MsgWithdrawFromLiquidityPool{ },
				types.ErrEmptyWithdrawerAddr.Error(),
			},
			{
				types.MsgWithdrawFromLiquidityPool{WithdrawerAddress: validAddr, PoolCoin: invalidDenomCoin},
				invalidDenomErrMsg,
			},
			{
				types.MsgWithdrawFromLiquidityPool{WithdrawerAddress: validAddr, PoolCoin: negativeCoin},
				negativeAmountErrMsg,
			},
			{
				types.MsgWithdrawFromLiquidityPool{WithdrawerAddress: validAddr, PoolCoin: zeroCoin},
				types.ErrBadPoolCoinAmount.Error(),
			},
		} {
			err := tc.msg.ValidateBasic()
			require.EqualError(t, err, tc.errMsg)
		}
	})
	t.Run("MsgSwap", func(t *testing.T) {
		offerCoin := sdk.NewCoin(DenomX, sdk.NewInt(10000))
		orderPrice := sdk.MustNewDecFromStr("1.0")

		for _, tc := range []struct {
			msg    types.MsgSwap
			errMsg string
		}{
			{
				types.MsgSwap{},
				types.ErrEmptySwapRequesterAddr.Error(),
			},
			{
				types.MsgSwap{SwapRequesterAddress: validAddr, OfferCoin: invalidDenomCoin, OrderPrice: orderPrice},
				invalidDenomErrMsg,
			},
			{
				types.MsgSwap{SwapRequesterAddress: validAddr, OfferCoin: zeroCoin},
				types.ErrBadOfferCoinAmount.Error(),
			},
			{
				types.MsgSwap{SwapRequesterAddress: validAddr, OfferCoin: negativeCoin},
				negativeAmountErrMsg,
			},
			{
				types.MsgSwap{SwapRequesterAddress: validAddr, OfferCoin: offerCoin, OrderPrice: sdk.ZeroDec()},
				types.ErrBadOrderPrice.Error(),
			},
			{
				types.MsgSwap{SwapRequesterAddress: validAddr, OfferCoin: sdk.NewCoin(DenomX, sdk.OneInt()), OrderPrice: orderPrice},
				types.ErrLessThanMinOfferAmount.Error(),
			},
		} {
			err := tc.msg.ValidateBasic()
			require.EqualError(t, err, tc.errMsg)
		}
	})
}
