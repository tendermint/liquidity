package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ sdk.Msg = &MsgCreateLiquidityPool{}
	_ sdk.Msg = &MsgDepositToLiquidityPool{}
	_ sdk.Msg = &MsgWithdrawFromLiquidityPool{}
	_ sdk.Msg = &MsgSwap{}
)

const (
	TypeMsgCreateLiquidityPool       = "create_liquidity_pool"
	TypeMsgDepositToLiquidityPool    = "deposit_to_liquidity_pool"
	TypeMsgWithdrawFromLiquidityPool = "withdraw_from_liquidity_pool"
	TypeMsgSwap                      = "swap"
)

// ------------------------------------------------------------------------
// MsgCreateLiquidityPool
// ------------------------------------------------------------------------

// NewMsgSwap creates a new MsgSwap object.
func NewMsgCreateLiquidityPool(
	poolCreator sdk.AccAddress,
	poolTypeIndex uint32,
	reserveCoinDenoms []string,
	depositCoins sdk.Coins,
) *MsgCreateLiquidityPool {
	return &MsgCreateLiquidityPool{
		PoolCreatorAddress: poolCreator.String(),
		PoolTypeIndex:      poolTypeIndex,
		ReserveCoinDenoms:  reserveCoinDenoms,
		DepositCoins:       depositCoins,
	}
}

// Route implements Msg.
func (msg MsgCreateLiquidityPool) Route() string { return RouterKey }

// Type implements Msg.
func (msg MsgCreateLiquidityPool) Type() string { return TypeMsgCreateLiquidityPool }

// ValidateBasic implements Msg.
func (msg MsgCreateLiquidityPool) ValidateBasic() error {
	if msg.PoolCreatorAddress == "" {
		return ErrEmptyPoolCreatorAddr
	}
	if len(msg.ReserveCoinDenoms) != msg.DepositCoins.Len() {
		return ErrNumOfReserveCoin
	}
	if err := msg.DepositCoins.Validate(); err != nil {
		return err
	}
	if !msg.DepositCoins.IsAllPositive() {
		return ErrBadPoolCoinAmount
	}
	if uint32(msg.DepositCoins.Len()) > MaxReserveCoinNum ||
		MinReserveCoinNum > uint32(msg.DepositCoins.Len()) {
		return ErrNumOfReserveCoin
	}
	return nil
}

// GetSignBytes implements Msg.
func (msg MsgCreateLiquidityPool) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners implements Msg.
func (msg MsgCreateLiquidityPool) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.PoolCreatorAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgCreateLiquidityPool) GetPoolCreator() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.PoolCreatorAddress)
	if err != nil {
		panic(err)
	}
	return addr
}

// ------------------------------------------------------------------------
// MsgDepositToLiquidityPool
// ------------------------------------------------------------------------

// NewMsgSwap creates a new MsgSwap object.
func NewMsgDepositToLiquidityPool(
	depositor sdk.AccAddress,
	poolId uint64,
	depositCoins sdk.Coins,
) *MsgDepositToLiquidityPool {
	return &MsgDepositToLiquidityPool{
		DepositorAddress: depositor.String(),
		PoolId:           poolId,
		DepositCoins:     depositCoins,
	}
}

// Route implements Msg.
func (msg MsgDepositToLiquidityPool) Route() string { return RouterKey }

// Type implements Msg.
func (msg MsgDepositToLiquidityPool) Type() string { return TypeMsgDepositToLiquidityPool }

// ValidateBasic implements Msg.
func (msg MsgDepositToLiquidityPool) ValidateBasic() error {
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
func (msg MsgDepositToLiquidityPool) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners implements Msg.
func (msg MsgDepositToLiquidityPool) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.DepositorAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgDepositToLiquidityPool) GetDepositor() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.DepositorAddress)
	if err != nil {
		panic(err)
	}
	return addr
}

// ------------------------------------------------------------------------
// MsgWithdrawFromLiquidityPool
// ------------------------------------------------------------------------

// NewMsgWithdraw creates a new MsgWithdraw object.
func NewMsgWithdrawFromLiquidityPool(
	withdrawer sdk.AccAddress,
	poolId uint64,
	poolCoin sdk.Coin,
) *MsgWithdrawFromLiquidityPool {
	return &MsgWithdrawFromLiquidityPool{
		WithdrawerAddress: withdrawer.String(),
		PoolId:            poolId,
		PoolCoin:          poolCoin,
	}
}

// Route implements Msg.
func (msg MsgWithdrawFromLiquidityPool) Route() string { return RouterKey }

// Type implements Msg.
func (msg MsgWithdrawFromLiquidityPool) Type() string { return TypeMsgWithdrawFromLiquidityPool }

// ValidateBasic implements Msg.
func (msg MsgWithdrawFromLiquidityPool) ValidateBasic() error {
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
func (msg MsgWithdrawFromLiquidityPool) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners implements Msg.
func (msg MsgWithdrawFromLiquidityPool) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.WithdrawerAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgWithdrawFromLiquidityPool) GetWithdrawer() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.WithdrawerAddress)
	if err != nil {
		panic(err)
	}
	return addr
}

// ------------------------------------------------------------------------
// MsgSwap
// ------------------------------------------------------------------------

// NewMsgSwap creates a new MsgSwap object.
func NewMsgSwap(
	swapRequester sdk.AccAddress,
	poolId uint64,
	poolTypeIndex uint32,
	swapType uint32,
	offerCoin sdk.Coin,
	demandCoinDenom string,
	orderPrice sdk.Dec,
) *MsgSwap {
	return &MsgSwap{
		SwapRequesterAddress: swapRequester.String(),
		PoolId:               poolId,
		PoolTypeIndex:        poolTypeIndex,
		SwapType:             swapType,
		OfferCoin:            offerCoin,
		DemandCoinDenom:      demandCoinDenom,
		OrderPrice:           orderPrice,
	}
}

// Route implements Msg.
func (msg MsgSwap) Route() string { return RouterKey }

// Type implements Msg.
func (msg MsgSwap) Type() string { return TypeMsgSwap }

// ValidateBasic implements Msg.
func (msg MsgSwap) ValidateBasic() error {
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
		return ErrBadOderPrice
	}
	if !msg.OfferCoin.Amount.GTE(DefaultOfferCoinAmount) {
		return ErrLessThanMinOfferAmount
	}
	return nil
}

// GetSignBytes implements Msg.
func (msg MsgSwap) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}

// GetSigners implements Msg.
func (msg MsgSwap) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.SwapRequesterAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr}
}

func (msg MsgSwap) GetSwapRequester() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.SwapRequesterAddress)
	if err != nil {
		panic(err)
	}
	return addr
}
