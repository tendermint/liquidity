package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"

	"github.com/tendermint/liquidity/x/liquidity/types"
)

const (
	DefaultPoolTypeId = uint32(1)
	DefaultPoolId     = uint64(1)
	DefaultSwapTypeId = uint32(1)
	DenomX            = "denomX"
	DenomY            = "denomY"
	DenomPoolCoin     = "denomPoolCoin"
)

func TestMsgCreatePool(t *testing.T) {
	addr := sdk.AccAddress(crypto.AddressHash([]byte("testAccount")))
	coins := sdk.NewCoins(sdk.NewCoin(DenomX, sdk.NewInt(1000)), sdk.NewCoin(DenomY, sdk.NewInt(1000)))
	msg := types.NewMsgCreatePool(addr, DefaultPoolTypeId, coins)
	require.IsType(t, &types.MsgCreatePool{}, msg)
	require.Equal(t, types.RouterKey, msg.Route())
	require.Equal(t, types.TypeMsgCreatePool, msg.Type())

	err := msg.ValidateBasic()
	require.NoError(t, err)
	signers := msg.GetSigners()
	require.Len(t, signers, 1)
	require.Equal(t, msg.GetPoolCreator(), signers[0])
	require.Equal(t, sdk.MustSortJSON(types.ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())

	// Fail cases
	msg = types.NewMsgCreatePool(sdk.AccAddress{}, DefaultPoolTypeId, coins)
	err = msg.ValidateBasic()
	require.Error(t, err)
	coinsFail := sdk.NewCoins(sdk.NewCoin(DenomY, sdk.NewInt(1000)))
	msg = types.NewMsgCreatePool(addr, DefaultPoolTypeId, coinsFail)
	err = msg.ValidateBasic()
	require.Error(t, err)
	coinsFail = sdk.NewCoins(sdk.NewCoin(DenomX, sdk.NewInt(1000)), sdk.NewCoin(DenomY, sdk.NewInt(1000)), sdk.NewCoin("Denomfail", sdk.NewInt(1000)))
	msg = types.NewMsgCreatePool(addr, DefaultPoolTypeId, coinsFail)
	err = msg.ValidateBasic()
	require.Error(t, err)
	coinsFail = sdk.NewCoins(sdk.NewCoin(DenomX, sdk.NewInt(0)), sdk.NewCoin(DenomY, sdk.NewInt(1000)))
	msg = types.NewMsgCreatePool(addr, DefaultPoolTypeId, coinsFail)
	err = msg.ValidateBasic()
	require.Error(t, err)
}

func TestMsgDepositWithinBatch(t *testing.T) {
	addr := sdk.AccAddress(crypto.AddressHash([]byte("testAccount")))
	coins := sdk.NewCoins(sdk.NewCoin(DenomX, sdk.NewInt(1000)), sdk.NewCoin(DenomY, sdk.NewInt(1000)))
	msg := types.NewMsgDepositWithinBatch(addr, DefaultPoolId, coins)
	require.IsType(t, &types.MsgDepositWithinBatch{}, msg)
	require.Equal(t, types.RouterKey, msg.Route())
	require.Equal(t, types.TypeMsgDepositWithinBatch, msg.Type())

	err := msg.ValidateBasic()
	require.NoError(t, err)
	signers := msg.GetSigners()
	require.Len(t, signers, 1)
	require.Equal(t, msg.GetDepositor(), signers[0])
	require.Equal(t, sdk.MustSortJSON(types.ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())

	// Fail case
	coinsFail := sdk.NewCoins(sdk.NewCoin(DenomX, sdk.NewInt(1000)), sdk.NewCoin(DenomY, sdk.NewInt(1000)), sdk.NewCoin("Denomfail", sdk.NewInt(1000)))
	msg = types.NewMsgDepositWithinBatch(addr, DefaultPoolId, coinsFail)
	err = msg.ValidateBasic()
	require.Error(t, err)
	coinsFail = sdk.NewCoins(sdk.NewCoin(DenomX, sdk.NewInt(0)), sdk.NewCoin(DenomY, sdk.NewInt(1000)))
	msg = types.NewMsgDepositWithinBatch(addr, DefaultPoolId, coinsFail)
	err = msg.ValidateBasic()
	require.Error(t, err)
	msg = types.NewMsgDepositWithinBatch(sdk.AccAddress{}, DefaultPoolId, coins)
	err = msg.ValidateBasic()
	require.Error(t, err)
}
func TestMsgWithdrawWithinBatch(t *testing.T) {
	addr := sdk.AccAddress(crypto.AddressHash([]byte("testAccount")))
	coin := sdk.NewCoin(DenomPoolCoin, sdk.NewInt(1000))
	msg := types.NewMsgWithdrawWithinBatch(addr, DefaultPoolId, coin)
	require.IsType(t, &types.MsgWithdrawWithinBatch{}, msg)
	require.Equal(t, types.RouterKey, msg.Route())
	require.Equal(t, types.TypeMsgWithdrawWithinBatch, msg.Type())

	err := msg.ValidateBasic()
	require.NoError(t, err)
	signers := msg.GetSigners()
	require.Len(t, signers, 1)
	require.Equal(t, msg.GetWithdrawer(), signers[0])
	require.Equal(t, sdk.MustSortJSON(types.ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())

	// Fail case
	coinFail := sdk.NewCoin("testPoolCoin", sdk.NewInt(0))
	msg = types.NewMsgWithdrawWithinBatch(addr, DefaultPoolId, coinFail)
	err = msg.ValidateBasic()
	require.Error(t, err)
	msg = types.NewMsgWithdrawWithinBatch(sdk.AccAddress{}, DefaultPoolId, coin)
	err = msg.ValidateBasic()
	require.Error(t, err)
}

func TestMsgSwapWithinBatch(t *testing.T) {
	addr := sdk.AccAddress(crypto.AddressHash([]byte("testAccount")))
	coin := sdk.NewCoin(DenomX, sdk.NewInt(1000))
	orderPrice, err := sdk.NewDecFromStr("0.1")
	require.NoError(t, err)
	msg := types.NewMsgSwapWithinBatch(addr, DefaultPoolId, DefaultSwapTypeId, coin, DenomY, orderPrice, types.DefaultSwapFeeRate)
	require.IsType(t, &types.MsgSwapWithinBatch{}, msg)
	require.Equal(t, types.RouterKey, msg.Route())
	require.Equal(t, types.TypeMsgSwapWithinBatch, msg.Type())

	err = msg.ValidateBasic()
	require.NoError(t, err)
	signers := msg.GetSigners()
	require.Len(t, signers, 1)
	require.Equal(t, msg.GetSwapRequester(), signers[0])
	require.Equal(t, sdk.MustSortJSON(types.ModuleCdc.MustMarshalJSON(msg)), msg.GetSignBytes())
}

func TestMsgPanics(t *testing.T) {
	emptyMsgCreatePool := types.MsgCreatePool{}
	emptyMsgDeposit := types.MsgDepositWithinBatch{}
	emptyMsgWithdraw := types.MsgWithdrawWithinBatch{}
	emptyMsgSwap := types.MsgSwapWithinBatch{}
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
	validPoolTypeId := DefaultPoolTypeId
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

	t.Run("MsgCreatePool", func(t *testing.T) {
		for _, tc := range []struct {
			msg    types.MsgCreatePool
			errMsg string
		}{
			{
				types.MsgCreatePool{},
				types.ErrBadPoolTypeId.Error(),
			},
			{
				types.MsgCreatePool{PoolTypeId: validPoolTypeId},
				types.ErrEmptyPoolCreatorAddr.Error(),
			},
			{
				types.MsgCreatePool{PoolCreatorAddress: validAddr, PoolTypeId: validPoolTypeId},
				types.ErrNumOfReserveCoin.Error(),
			},
			{
				types.MsgCreatePool{
					PoolCreatorAddress: validAddr,
					PoolTypeId:         validPoolTypeId,
					DepositCoins:       coinsWithInvalidDenom,
				},
				invalidDenomErrMsg,
			},
			{
				types.MsgCreatePool{
					PoolCreatorAddress: validAddr,
					PoolTypeId:         validPoolTypeId,
					DepositCoins:       coinsWithNegative,
				},
				negativeCoinErrMsg,
			},
			{
				types.MsgCreatePool{
					PoolCreatorAddress: validAddr,
					PoolTypeId:         validPoolTypeId,
					DepositCoins:       coinsWithZero,
				},
				zeroCoinErrMsg,
			},
			{
				types.MsgCreatePool{
					PoolCreatorAddress: validAddr,
					PoolTypeId:         validPoolTypeId,
					DepositCoins:       sdk.NewCoins(sdk.NewCoin(DenomX, sdk.NewInt(int64(types.MinReserveCoinNum)-1))),
				},
				types.ErrNumOfReserveCoin.Error(),
			},
			{
				types.MsgCreatePool{
					PoolCreatorAddress: validAddr,
					PoolTypeId:         validPoolTypeId,
					DepositCoins:       sdk.NewCoins(sdk.NewCoin(DenomX, sdk.NewInt(int64(types.MaxReserveCoinNum)+1))),
				},
				types.ErrNumOfReserveCoin.Error(),
			},
		} {
			err := tc.msg.ValidateBasic()
			require.EqualError(t, err, tc.errMsg)
		}
	})
	t.Run("MsgDepositWithinBatch", func(t *testing.T) {
		for _, tc := range []struct {
			msg    types.MsgDepositWithinBatch
			errMsg string
		}{
			{
				types.MsgDepositWithinBatch{},
				types.ErrEmptyDepositorAddr.Error(),
			},
			{
				types.MsgDepositWithinBatch{DepositorAddress: validAddr},
				types.ErrBadDepositCoinsAmount.Error(),
			},
			{
				types.MsgDepositWithinBatch{DepositorAddress: validAddr, DepositCoins: coinsWithInvalidDenom},
				invalidDenomErrMsg,
			},
			{
				types.MsgDepositWithinBatch{DepositorAddress: validAddr, DepositCoins: coinsWithNegative},
				negativeCoinErrMsg,
			},
			{
				types.MsgDepositWithinBatch{DepositorAddress: validAddr, DepositCoins: coinsWithZero},
				zeroCoinErrMsg,
			},
			{
				types.MsgDepositWithinBatch{
					DepositorAddress: validAddr,
					DepositCoins:     sdk.NewCoins(sdk.NewCoin(DenomX, sdk.NewInt(int64(types.MinReserveCoinNum)-1))),
				},
				types.ErrNumOfReserveCoin.Error(),
			},
			{
				types.MsgDepositWithinBatch{
					DepositorAddress: validAddr,
					DepositCoins:     sdk.NewCoins(sdk.NewCoin(DenomX, sdk.NewInt(int64(types.MaxReserveCoinNum)+1))),
				},
				types.ErrNumOfReserveCoin.Error(),
			},
		} {
			err := tc.msg.ValidateBasic()
			require.EqualError(t, err, tc.errMsg)
		}
	})
	t.Run("MsgWithdrawWithinBatch", func(t *testing.T) {
		for _, tc := range []struct {
			msg    types.MsgWithdrawWithinBatch
			errMsg string
		}{
			{
				types.MsgWithdrawWithinBatch{},
				types.ErrEmptyWithdrawerAddr.Error(),
			},
			{
				types.MsgWithdrawWithinBatch{WithdrawerAddress: validAddr, PoolCoin: invalidDenomCoin},
				invalidDenomErrMsg,
			},
			{
				types.MsgWithdrawWithinBatch{WithdrawerAddress: validAddr, PoolCoin: negativeCoin},
				negativeAmountErrMsg,
			},
			{
				types.MsgWithdrawWithinBatch{WithdrawerAddress: validAddr, PoolCoin: zeroCoin},
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
			msg    types.MsgSwapWithinBatch
			errMsg string
		}{
			{
				types.MsgSwapWithinBatch{},
				types.ErrEmptySwapRequesterAddr.Error(),
			},
			{
				types.MsgSwapWithinBatch{SwapRequesterAddress: validAddr, OfferCoin: invalidDenomCoin, OrderPrice: orderPrice},
				invalidDenomErrMsg,
			},
			{
				types.MsgSwapWithinBatch{SwapRequesterAddress: validAddr, OfferCoin: zeroCoin},
				types.ErrBadOfferCoinAmount.Error(),
			},
			{
				types.MsgSwapWithinBatch{SwapRequesterAddress: validAddr, OfferCoin: negativeCoin},
				negativeAmountErrMsg,
			},
			{
				types.MsgSwapWithinBatch{SwapRequesterAddress: validAddr, OfferCoin: offerCoin, OrderPrice: sdk.ZeroDec()},
				types.ErrBadOrderPrice.Error(),
			},
			{
				types.MsgSwapWithinBatch{SwapRequesterAddress: validAddr, OfferCoin: sdk.NewCoin(DenomX, sdk.OneInt()), OrderPrice: orderPrice},
				types.ErrLessThanMinOfferAmount.Error(),
			},
		} {
			err := tc.msg.ValidateBasic()
			require.EqualError(t, err, tc.errMsg)
		}
	})
}
