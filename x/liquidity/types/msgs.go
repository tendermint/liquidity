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
		PoolCreator:       poolCreator,
		PoolTypeIndex:     poolTypeIndex,
		ReserveCoinDenoms: reserveCoinDenoms,
		DepositCoins:      depositCoins,
	}
}

// Route implements Msg.
func (msg MsgCreateLiquidityPool) Route() string { return RouterKey }

// Type implements Msg.
func (msg MsgCreateLiquidityPool) Type() string { return TypeMsgCreateLiquidityPool }

// ValidateBasic implements Msg.
func (msg MsgCreateLiquidityPool) ValidateBasic() error {
	return nil
}

// GetSignBytes implements Msg.
func (msg MsgCreateLiquidityPool) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
	return nil
}

// GetSigners implements Msg.
func (msg MsgCreateLiquidityPool) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.PoolCreator}
}

// ------------------------------------------------------------------------
// MsgDepositToLiquidityPool
// ------------------------------------------------------------------------

// NewMsgSwap creates a new MsgSwap object.
func NewMsgDepositToLiquidityPool(
	depositor sdk.AccAddress,
	poolID uint64,
	depositCoins sdk.Coins,
) *MsgDepositToLiquidityPool {
	return &MsgDepositToLiquidityPool{
		Depositor:    depositor,
		PoolID:       poolID,
		DepositCoins: depositCoins,
	}
}

// Route implements Msg.
func (msg MsgDepositToLiquidityPool) Route() string { return RouterKey }

// Type implements Msg.
func (msg MsgDepositToLiquidityPool) Type() string { return TypeMsgDepositToLiquidityPool }

// ValidateBasic implements Msg.
func (msg MsgDepositToLiquidityPool) ValidateBasic() error {
	return nil
}

// GetSignBytes implements Msg.
func (msg MsgDepositToLiquidityPool) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
	return nil
}

// GetSigners implements Msg.
func (msg MsgDepositToLiquidityPool) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Depositor}
}

// ------------------------------------------------------------------------
// MsgWithdrawFromLiquidityPool
// ------------------------------------------------------------------------

// NewMsgSwap creates a new MsgSwap object.
func NewMsgWithdrawFromLiquidityPool(
	withdrawer sdk.AccAddress,
	poolID uint64,
	poolCoin sdk.Coins,
) *MsgWithdrawFromLiquidityPool {
	return &MsgWithdrawFromLiquidityPool{
		Withdrawer: withdrawer,
		PoolID:     poolID,
		PoolCoin:   poolCoin,
	}
}

// Route implements Msg.
func (msg MsgWithdrawFromLiquidityPool) Route() string { return RouterKey }

// Type implements Msg.
func (msg MsgWithdrawFromLiquidityPool) Type() string { return TypeMsgWithdrawFromLiquidityPool }

// ValidateBasic implements Msg.
func (msg MsgWithdrawFromLiquidityPool) ValidateBasic() error {
	return nil
}

// GetSignBytes implements Msg.
func (msg MsgWithdrawFromLiquidityPool) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
	return nil
}

// GetSigners implements Msg.
func (msg MsgWithdrawFromLiquidityPool) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Withdrawer}
}

// ------------------------------------------------------------------------
// MsgSwap
// ------------------------------------------------------------------------

// NewMsgSwap creates a new MsgSwap object.
func NewMsgSwap(
	swapRequester sdk.AccAddress,
	poolID uint64,
	poolTypeIndex uint32,
	swapType uint32,
	offerCoin sdk.Coin,
	demandCoinDenom string,
	orderPrice sdk.Dec,
) *MsgSwap {
	return &MsgSwap{
		SwapRequester:   swapRequester,
		PoolID:          poolID,
		PoolTypeIndex:   poolTypeIndex,
		SwapType:        swapType,
		OfferCoin:       offerCoin,
		DemandCoinDenom: demandCoinDenom,
		OrderPrice:      orderPrice,
	}
}

// Route implements Msg.
func (msg MsgSwap) Route() string { return RouterKey }

// Type implements Msg.
func (msg MsgSwap) Type() string { return TypeMsgSwap }

// ValidateBasic implements Msg.
func (msg MsgSwap) ValidateBasic() error {
	return nil
}

// GetSignBytes implements Msg.
func (msg MsgSwap) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
	return nil
}

// GetSigners implements Msg.
func (msg MsgSwap) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.SwapRequester}
}
