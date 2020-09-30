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

type MsgCreateLiquidityPoolLegacy struct {
	PoolCreator         sdk.AccAddress // account address of the origin of this message
	PoolTypeIndex       uint32         // index of the liquidity pool type of this new liquidity pool
	ReserveTokenDenoms  []string       // list of reserve token denoms for this new liquidity pool, store alphabetical
	DepositTokensAmount sdk.Coins      // deposit token for initial pool deposit into this new liquidity pool
}

type MsgDepositToLiquidityPoolLegacy struct {
	Depositor           sdk.AccAddress // account address of the origin of this message
	PoolID              uint64         // id of the liquidity pool where this message is belong to
	DepositTokensAmount sdk.Coins      // deposit token of this pool deposit message
}

type MsgWithdrawFromLiquidityPoolLegacy struct {
	Withdrawer      sdk.AccAddress // account address of the origin of this message
	PoolID          uint64         // id of the liquidity pool where this message is belong to
	PoolTokenAmount sdk.Coins      // pool token sent for reserve token withdraw
}

type MsgSwapLegacy struct {
	SwapRequester sdk.AccAddress // account address of the origin of this message
	PoolID        uint64         // id of the liquidity pool where this message is belong to
	PoolTypeIndex uint32         // index of the liquidity pool type where this message is belong to
	SwapType      uint32         // swap type of this swap message, default 1: InstantSwap, requesting instant swap
	OfferToken    sdk.Coin       // offer token of this swap message
	DemandToken   sdk.Coin       // denom of demand token of this swap message
	OrderPrice    sdk.Dec        // order price of this swap message
}

// ------------------------------------------------------------------------
// MsgCreateLiquidityPool
// ------------------------------------------------------------------------

// NewMsgSwap creates a new MsgSwap object.
func NewMsgCreateLiquidityPool(
	poolCreator sdk.AccAddress,
	poolTypeIndex uint32,
	reserveTokenDenoms []string,
	depositTokensAmount sdk.Coins,
) *MsgCreateLiquidityPool {
	return &MsgCreateLiquidityPool{
		PoolCreator:         poolCreator,
		PoolTypeIndex:       poolTypeIndex,
		ReserveTokenDenoms:  reserveTokenDenoms,
		DepositTokensAmount: depositTokensAmount,
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
	depositTokensAmount sdk.Coins,
) *MsgDepositToLiquidityPool {
	return &MsgDepositToLiquidityPool{
		Depositor:           depositor,
		PoolID:              poolID,
		DepositTokensAmount: depositTokensAmount,
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
	poolTokenAmount sdk.Coins,
) *MsgWithdrawFromLiquidityPool {
	return &MsgWithdrawFromLiquidityPool{
		Withdrawer:      withdrawer,
		PoolID:          poolID,
		PoolTokenAmount: poolTokenAmount,
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
	offerToken sdk.Coin,
	demandToken sdk.Coin,
	orderPrice sdk.Dec,
) *MsgSwap {
	return &MsgSwap{
		SwapRequester: swapRequester,
		PoolID:        poolID,
		PoolTypeIndex: poolTypeIndex,
		SwapType:      swapType,
		OfferToken:    offerToken,
		DemandToken:   demandToken,
		OrderPrice:    orderPrice,
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
