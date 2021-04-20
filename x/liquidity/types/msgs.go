package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Messages Type of Liquidity module
var (
	_ sdk.Msg = &MsgCreatePool{}
	_ sdk.Msg = &MsgDepositWithinBatch{}
	_ sdk.Msg = &MsgWithdrawWithinBatch{}
	_ sdk.Msg = &MsgSwapWithinBatch{}
)

// Messages Type of Liquidity module
const (
	TypeMsgCreatePool          = "create_pool"
	TypeMsgDepositWithinBatch  = "deposit_within_batch"
	TypeMsgWithdrawWithinBatch = "withdraw_within_batch"
	TypeMsgSwapWithinBatch     = "swap_within_batch"
)

// ------------------------------------------------------------------------
// MsgCreatePool
// ------------------------------------------------------------------------

// NewMsgSwapWithinBatch creates a new MsgSwapWithinBatch object.
func NewMsgCreatePool(
	poolCreator sdk.AccAddress,
	poolTypeId uint32,
	depositCoins sdk.Coins,
) *MsgCreatePool {
	return &MsgCreatePool{
		PoolCreatorAddress: poolCreator.String(),
		PoolTypeId:         poolTypeId,
		DepositCoins:       depositCoins,
	}
}

// Route implements Msg.
func (msg MsgCreatePool) Route() string { return RouterKey }

// Type implements Msg.
func (msg MsgCreatePool) Type() string { return TypeMsgCreatePool }

// ValidateBasic implements Msg.
func (msg MsgCreatePool) ValidateBasic() error {
	if 1 > msg.PoolTypeId {
		return ErrBadPoolTypeId
	}
	if msg.PoolCreatorAddress == "" {
		return ErrEmptyPoolCreatorAddr
	}
	if err := msg.DepositCoins.Validate(); err != nil {
		return err
	}
	if uint32(msg.DepositCoins.Len()) > MaxReserveCoinNum ||
		MinReserveCoinNum > uint32(msg.DepositCoins.Len()) {
		return ErrNumOfReserveCoin
	}
	return nil
}

// GetSignBytes implements Msg.
func (msg MsgCreatePool) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners implements Msg.
func (msg MsgCreatePool) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.PoolCreatorAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgCreatePool) GetPoolCreator() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.PoolCreatorAddress)
	if err != nil {
		panic(err)
	}
	return addr
}

// ------------------------------------------------------------------------
// MsgDepositWithinBatch
// ------------------------------------------------------------------------

// NewMsgSwapWithinBatch creates a new MsgSwapWithinBatch object.
func NewMsgDepositWithinBatch(
	depositor sdk.AccAddress,
	poolId uint64,
	depositCoins sdk.Coins,
) *MsgDepositWithinBatch {
	return &MsgDepositWithinBatch{
		DepositorAddress: depositor.String(),
		PoolId:           poolId,
		DepositCoins:     depositCoins,
	}
}

// Route implements Msg.
func (msg MsgDepositWithinBatch) Route() string { return RouterKey }

// Type implements Msg.
func (msg MsgDepositWithinBatch) Type() string { return TypeMsgDepositWithinBatch }

// ValidateBasic implements Msg.
func (msg MsgDepositWithinBatch) ValidateBasic() error {
	if msg.DepositorAddress == "" {
		return ErrEmptyDepositorAddr
	}
	if err := msg.DepositCoins.Validate(); err != nil {
		return err
	}
	if !msg.DepositCoins.IsAllPositive() {
		return ErrBadDepositCoinsAmount
	}
	if uint32(msg.DepositCoins.Len()) > MaxReserveCoinNum ||
		MinReserveCoinNum > uint32(msg.DepositCoins.Len()) {
		return ErrNumOfReserveCoin
	}
	return nil
}

// GetSignBytes implements Msg.
func (msg MsgDepositWithinBatch) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners implements Msg.
func (msg MsgDepositWithinBatch) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.DepositorAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgDepositWithinBatch) GetDepositor() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.DepositorAddress)
	if err != nil {
		panic(err)
	}
	return addr
}

// ------------------------------------------------------------------------
// MsgWithdrawWithinBatch
// ------------------------------------------------------------------------

// NewMsgWithdraw creates a new MsgWithdraw object.
func NewMsgWithdrawWithinBatch(
	withdrawer sdk.AccAddress,
	poolId uint64,
	poolCoin sdk.Coin,
) *MsgWithdrawWithinBatch {
	return &MsgWithdrawWithinBatch{
		WithdrawerAddress: withdrawer.String(),
		PoolId:            poolId,
		PoolCoin:          poolCoin,
	}
}

// Route implements Msg.
func (msg MsgWithdrawWithinBatch) Route() string { return RouterKey }

// Type implements Msg.
func (msg MsgWithdrawWithinBatch) Type() string { return TypeMsgWithdrawWithinBatch }

// ValidateBasic implements Msg.
func (msg MsgWithdrawWithinBatch) ValidateBasic() error {
	if msg.WithdrawerAddress == "" {
		return ErrEmptyWithdrawerAddr
	}
	if err := msg.PoolCoin.Validate(); err != nil {
		return err
	}
	if !msg.PoolCoin.IsPositive() {
		return ErrBadPoolCoinAmount
	}
	return nil
}

// GetSignBytes implements Msg.
func (msg MsgWithdrawWithinBatch) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners implements Msg.
func (msg MsgWithdrawWithinBatch) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.WithdrawerAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgWithdrawWithinBatch) GetWithdrawer() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.WithdrawerAddress)
	if err != nil {
		panic(err)
	}
	return addr
}

// ------------------------------------------------------------------------
// MsgSwapWithinBatch
// ------------------------------------------------------------------------

// NewMsgSwapWithinBatch creates a new MsgSwapWithinBatch object.
func NewMsgSwapWithinBatch(
	swapRequester sdk.AccAddress,
	poolId uint64,
	swapTypeId uint32,
	offerCoin sdk.Coin,
	demandCoinDenom string,
	orderPrice sdk.Dec,
	swapFeeRate sdk.Dec,
) *MsgSwapWithinBatch {
	return &MsgSwapWithinBatch{
		SwapRequesterAddress: swapRequester.String(),
		PoolId:               poolId,
		SwapTypeId:           swapTypeId,
		OfferCoin:            offerCoin,
		OfferCoinFee:         GetOfferCoinFee(offerCoin, swapFeeRate),
		DemandCoinDenom:      demandCoinDenom,
		OrderPrice:           orderPrice,
	}
}

func GetOfferCoinFee(offerCoin sdk.Coin, swapFeeRate sdk.Dec) sdk.Coin {
	// apply half-ratio swap fee rate
	// see https://github.com/tendermint/liquidity/issues/41 for details
	return sdk.NewCoin(offerCoin.Denom, offerCoin.Amount.ToDec().Mul(swapFeeRate.QuoInt64(2)).TruncateInt()) // offerCoin.Amount * (swapFeeRate/2)
}

// Route implements Msg.
func (msg MsgSwapWithinBatch) Route() string { return RouterKey }

// Type implements Msg.
func (msg MsgSwapWithinBatch) Type() string { return TypeMsgSwapWithinBatch }

// ValidateBasic implements Msg.
func (msg MsgSwapWithinBatch) ValidateBasic() error {
	if msg.SwapRequesterAddress == "" {
		return ErrEmptySwapRequesterAddr
	}
	if err := msg.OfferCoin.Validate(); err != nil {
		return err
	}
	if !msg.OfferCoin.IsPositive() {
		return ErrBadOfferCoinAmount
	}
	if !msg.OrderPrice.IsPositive() {
		return ErrBadOrderPrice
	}
	if !msg.OfferCoin.Amount.GTE(MinOfferCoinAmount) {
		return ErrLessThanMinOfferAmount
	}
	return nil
}

// GetSignBytes implements Msg.
func (msg MsgSwapWithinBatch) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners implements Msg.
func (msg MsgSwapWithinBatch) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.SwapRequesterAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgSwapWithinBatch) GetSwapRequester() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.SwapRequesterAddress)
	if err != nil {
		panic(err)
	}
	return addr
}
