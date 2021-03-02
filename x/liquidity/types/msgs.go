package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	rosettatypes "github.com/coinbase/rosetta-sdk-go/types"
	"github.com/cosmos/cosmos-sdk/server/rosetta"
	"github.com/gogo/protobuf/proto"
	"strconv"
	"strings"
)

// Messages Type of Liquidity module
var (
	_ sdk.Msg = &MsgCreateLiquidityPool{}
	_ sdk.Msg = &MsgDepositToLiquidityPool{}
	_ sdk.Msg = &MsgWithdrawFromLiquidityPool{}
	_ sdk.Msg = &MsgSwap{}
)

// Messages Type of Liquidity module
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
	depositCoins sdk.Coins,
) *MsgCreateLiquidityPool {
	return &MsgCreateLiquidityPool{
		PoolCreatorAddress: poolCreator.String(),
		PoolTypeIndex:      poolTypeIndex,
		DepositCoins:       depositCoins,
	}
}

// Route implements Msg.
func (msg MsgCreateLiquidityPool) Route() string { return RouterKey }

// Type implements Msg.
func (msg MsgCreateLiquidityPool) Type() string { return TypeMsgCreateLiquidityPool }

// ValidateBasic implements Msg.
func (msg MsgCreateLiquidityPool) ValidateBasic() error {
	if 1 > msg.PoolTypeIndex {
		return ErrBadPoolTypeIndex
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
	swapType uint32,
	offerCoin sdk.Coin,
	demandCoinDenom string,
	orderPrice sdk.Dec,
	swapFeeRate sdk.Dec,
) *MsgSwap {
	return &MsgSwap{
		SwapRequesterAddress: swapRequester.String(),
		PoolId:               poolId,
		SwapType:             swapType,
		OfferCoin:            offerCoin,
		OfferCoinFee:         GetOfferCoinFee(offerCoin, swapFeeRate),
		DemandCoinDenom:      demandCoinDenom,
		OrderPrice:           orderPrice,
	}
}

// NewMsgSwapWithOfferCoinFee creates a new MsgSwap object with explicit OfferCoinFee.
func NewMsgSwapWithOfferCoinFee(
	swapRequester sdk.AccAddress,
	poolId uint64,
	swapType uint32,
	offerCoin sdk.Coin,
	demandCoinDenom string,
	orderPrice sdk.Dec,
	offerCoinFee sdk.Coin,
) *MsgSwap {
	return &MsgSwap{
		SwapRequesterAddress: swapRequester.String(),
		PoolId:               poolId,
		SwapType:             swapType,
		OfferCoin:            offerCoin,
		OfferCoinFee:         offerCoinFee,
		DemandCoinDenom:      demandCoinDenom,
		OrderPrice:           orderPrice,
	}
}

//func (msg MsgSwap) GetOfferCoinFee() sdk.Coin {
//	return GetOfferCoinFee(msg.OfferCoin)
//}

func GetOfferCoinFee(offerCoin sdk.Coin, swapFeeRate sdk.Dec) sdk.Coin {
	return sdk.NewCoin(offerCoin.Denom, offerCoin.Amount.ToDec().Mul(swapFeeRate.Mul(HalfRatio)).TruncateInt())
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
		return ErrBadOrderPrice
	}
	if !msg.OfferCoin.Amount.GTE(MinOfferCoinAmount) {
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

// TODO: WIP implementation Rosetta Interface for liquidity module msgs
// Rosetta Msg interface.
func (msg *MsgSwap) ToOperations(withStatus bool, hasError bool) []*rosettatypes.Operation {
	var operations []*rosettatypes.Operation
	swapRequester := msg.SwapRequesterAddress
	coin := msg.OfferCoin.Add(msg.OfferCoinFee)
	swapOp := func(account *rosettatypes.AccountIdentifier, amount string, index int) *rosettatypes.Operation {
		var status string
		if withStatus {
			status = rosetta.StatusSuccess
			if hasError {
				status = rosetta.StatusReverted
			}
		}
		return &rosettatypes.Operation{
			OperationIdentifier: &rosettatypes.OperationIdentifier{
				Index: int64(index),
			},
			Type:    proto.MessageName(msg),
			Status:  status,
			Account: account,
			Amount: &rosettatypes.Amount{
				Value: amount,
				Currency: &rosettatypes.Currency{
					Symbol: coin.Denom,
				},
			},
		}
	}
	swapAcc := &rosettatypes.AccountIdentifier{
		Address: swapRequester,
	}
	poolAcc := &rosettatypes.AccountIdentifier{
		Address: "liquidity_pool",
		Metadata: map[string]interface{}{
			"pool_id": msg.PoolId,
			"swap_type": msg.SwapType,
			"demand_coind_denom": msg.DemandCoinDenom,
			"order_price": msg.OfferCoin,
			"offer_coin_fee": msg.OfferCoinFee,
		},
	}
	operations = append(operations,
		swapOp(swapAcc, "-"+coin.Amount.String(), 0),
		swapOp(poolAcc, coin.Amount.String(), 1),
	)
	return operations
}

func (msg *MsgSwap) FromOperations(ops []*rosettatypes.Operation) (sdk.Msg, error) {
	var (
		swapRequester sdk.AccAddress
		poolId uint64
		swapType uint32
		offerCoin sdk.Coin
		demandCoinDenom string
		orderPrice sdk.Dec
		offerCoinFee sdk.Coin
		err     error
	)
	var m map[string]interface{}
	for _, op := range ops {
		if strings.HasPrefix(op.Amount.Value, "-") {
			if op.Account == nil {
				return nil, fmt.Errorf("account identifier must be specified")
			}
			swapRequester, err = sdk.AccAddressFromBech32(op.Account.Address)
			if err != nil {
				return nil, err
			}
			continue
		} else {
			if op.Account == nil {
				return nil, fmt.Errorf("account identifier must be specified")
			}
			m = op.Metadata
		}

		amount, err := strconv.ParseInt(op.Amount.Value, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid amount: %w", err)
		}

		poolId = m["pool_id"].(uint64)
		swapType = m["swap_type"].(uint32)
		demandCoinDenom = m["demand_coind_denom"].(string)
		orderPrice = m["order_price"].(sdk.Dec)
		offerCoinFee = m["offer_coin_fee"].(sdk.Coin)
		offerCoin = sdk.NewCoin(op.Amount.Currency.Symbol, sdk.NewInt(amount)).Sub(offerCoinFee)
	}

	return NewMsgSwapWithOfferCoinFee(swapRequester, poolId, swapType, offerCoin, demandCoinDenom, orderPrice, offerCoinFee), nil
}
