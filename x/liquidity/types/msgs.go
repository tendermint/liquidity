package types

import (
	"fmt"
	rosettatypes "github.com/coinbase/rosetta-sdk-go/types"
	"github.com/cosmos/cosmos-sdk/server/rosetta"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
	op := func(account *rosettatypes.AccountIdentifier, amount string, index int) *rosettatypes.Operation {
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
			"pool_id":            msg.PoolId,
			"swap_type":          msg.SwapType,
			"demand_coind_denom": msg.DemandCoinDenom,
			"order_price":        msg.OrderPrice,
			"offer_coin_fee":     msg.OfferCoinFee,
		},
	}
	operations = append(operations,
		op(swapAcc, "-"+coin.Amount.String(), 0),
		op(poolAcc, coin.Amount.String(), 1),
	)
	return operations
}

func (msg *MsgSwap) FromOperations(ops []*rosettatypes.Operation) (sdk.Msg, error) {
	var (
		swapRequester   sdk.AccAddress
		poolId          uint64
		swapType        uint32
		offerCoin       sdk.Coin
		demandCoinDenom string
		orderPrice      sdk.Dec
		offerCoinFee    sdk.Coin
		err             error
	)
	for _, op := range ops {
		if strings.HasPrefix(op.Amount.Value, "-") {
			if op.Account == nil {
				return nil, fmt.Errorf("account identifier must be specified")
			}
			swapRequester, err = sdk.AccAddressFromBech32(op.Account.Address)
			if err != nil {
				return nil, err
			}
		} else {
			if op.Account == nil {
				return nil, fmt.Errorf("account identifier must be specified")
			}
			amount, err := strconv.ParseInt(op.Amount.Value, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid amount: %w", err)
			}
			offerCoin = sdk.NewCoin(op.Amount.Currency.Symbol, sdk.NewInt(amount))
			poolId = op.Account.Metadata["pool_id"].(uint64)
			swapType = op.Account.Metadata["swap_type"].(uint32)
			demandCoinDenom = op.Account.Metadata["demand_coind_denom"].(string)
			orderPrice = op.Account.Metadata["order_price"].(sdk.Dec)
			offerCoinFee = op.Account.Metadata["offer_coin_fee"].(sdk.Coin)
		}
	}
	offerCoin = offerCoin.Sub(offerCoinFee)
	return NewMsgSwapWithOfferCoinFee(swapRequester, poolId, swapType, offerCoin, demandCoinDenom, orderPrice, offerCoinFee), nil
}

func (msg *MsgCreateLiquidityPool) ToOperations(withStatus bool, hasError bool) []*rosettatypes.Operation {
	var operations []*rosettatypes.Operation
	poolCreator := msg.PoolCreatorAddress
	coins := msg.DepositCoins
	op := func(account *rosettatypes.AccountIdentifier, amount string, index int) *rosettatypes.Operation {
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
					Symbol: coins.GetDenomByIndex(index % 2),
				},
			},
		}
	}
	creatorAcc := &rosettatypes.AccountIdentifier{
		Address: poolCreator,
	}
	poolAcc := &rosettatypes.AccountIdentifier{
		Address: "liquidity_pool",
		Metadata: map[string]interface{}{
			"pool_type_index": msg.PoolTypeIndex,
			//"reserve_coin_denoms": GetCoinDenoms(coins),
		},
	}
	index := 0
	for _, coin := range coins {
		operations = append(operations,
			op(creatorAcc, "-"+coin.Amount.String(), index),
		)
		index += 1
	}
	for _, coin := range coins {
		operations = append(operations,
			op(poolAcc, coin.Amount.String(), index),
		)
		index += 1
	}
	return operations
}

func (msg *MsgCreateLiquidityPool) FromOperations(ops []*rosettatypes.Operation) (sdk.Msg, error) {
	var (
		poolCreator   sdk.AccAddress
		poolTypeIndex uint32
		depositCoins  sdk.Coins
		err           error
	)
	for _, op := range ops {
		if strings.HasPrefix(op.Amount.Value, "-") {
			if op.Account == nil {
				return nil, fmt.Errorf("account identifier must be specified")
			}
			poolCreator, err = sdk.AccAddressFromBech32(op.Account.Address)
			if err != nil {
				return nil, err
			}
			continue
		} else {
			if op.Account == nil {
				return nil, fmt.Errorf("account identifier must be specified")
			}
			poolTypeIndex = op.Account.Metadata["pool_type_index"].(uint32)
			//reserveCoinDenoms = op.Account.Metadata["reserve_coin_denoms"].([]string)
		}
		amount, err := strconv.ParseInt(op.Amount.Value, 10, 64)
		depositCoins = depositCoins.Add(sdk.NewCoin(op.Amount.Currency.Symbol, sdk.NewInt(amount)))
		if err != nil {
			return nil, fmt.Errorf("invalid amount: %w", err)
		}
	}
	return NewMsgCreateLiquidityPool(poolCreator, poolTypeIndex, depositCoins), nil
}

func (msg *MsgDepositToLiquidityPool) ToOperations(withStatus bool, hasError bool) []*rosettatypes.Operation {
	var operations []*rosettatypes.Operation
	depositor := msg.DepositorAddress
	poolId := msg.PoolId
	coins := msg.DepositCoins
	op := func(account *rosettatypes.AccountIdentifier, amount string, index int) *rosettatypes.Operation {
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
					Symbol: coins.GetDenomByIndex(index % 2),
				},
			},
		}
	}
	creatorAcc := &rosettatypes.AccountIdentifier{
		Address: depositor,
	}
	poolAcc := &rosettatypes.AccountIdentifier{
		Address: "liquidity_pool",
		Metadata: map[string]interface{}{
			"pool_id": poolId,
		},
	}
	index := 0
	for _, coin := range coins {
		operations = append(operations,
			op(creatorAcc, "-"+coin.Amount.String(), index),
		)
		index += 1
	}
	for _, coin := range coins {
		operations = append(operations,
			op(poolAcc, coin.Amount.String(), index),
		)
		index += 1
	}
	return operations
}

func (msg *MsgDepositToLiquidityPool) FromOperations(ops []*rosettatypes.Operation) (sdk.Msg, error) {
	var (
		depositor    sdk.AccAddress
		poolId       uint64
		depositCoins sdk.Coins
		err          error
	)
	for _, op := range ops {
		if strings.HasPrefix(op.Amount.Value, "-") {
			if op.Account == nil {
				return nil, fmt.Errorf("account identifier must be specified")
			}
			depositor, err = sdk.AccAddressFromBech32(op.Account.Address)
			if err != nil {
				return nil, err
			}
			continue
		} else {
			if op.Account == nil {
				return nil, fmt.Errorf("account identifier must be specified")
			}
			poolId = op.Account.Metadata["pool_id"].(uint64)
		}
		amount, err := strconv.ParseInt(op.Amount.Value, 10, 64)
		depositCoins = depositCoins.Add(sdk.NewCoin(op.Amount.Currency.Symbol, sdk.NewInt(amount)))
		if err != nil {
			return nil, fmt.Errorf("invalid amount: %w", err)
		}
	}
	return NewMsgDepositToLiquidityPool(depositor, poolId, depositCoins), nil
}

func (msg *MsgWithdrawFromLiquidityPool) ToOperations(withStatus bool, hasError bool) []*rosettatypes.Operation {
	var operations []*rosettatypes.Operation
	withdrawer := msg.WithdrawerAddress
	poolId := msg.PoolId
	coin := msg.PoolCoin
	op := func(account *rosettatypes.AccountIdentifier, amount string, index int) *rosettatypes.Operation {
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
	withdrawerAcc := &rosettatypes.AccountIdentifier{
		Address: withdrawer,
	}
	poolAcc := &rosettatypes.AccountIdentifier{
		Address: "liquidity_pool",
		Metadata: map[string]interface{}{
			"pool_id": poolId,
		},
	}
	operations = append(operations,
		op(withdrawerAcc, "-"+coin.Amount.String(), 0),
		op(poolAcc, coin.Amount.String(), 1),
	)
	return operations
}

func (msg *MsgWithdrawFromLiquidityPool) FromOperations(ops []*rosettatypes.Operation) (sdk.Msg, error) {
	var (
		withdrawer sdk.AccAddress
		poolId     uint64
		poolCoin   sdk.Coin
		err        error
	)
	for _, op := range ops {
		if strings.HasPrefix(op.Amount.Value, "-") {
			if op.Account == nil {
				return nil, fmt.Errorf("account identifier must be specified")
			}
			withdrawer, err = sdk.AccAddressFromBech32(op.Account.Address)
			if err != nil {
				return nil, err
			}
			continue
		} else {
			if op.Account == nil {
				return nil, fmt.Errorf("account identifier must be specified")
			}
			poolId = op.Account.Metadata["pool_id"].(uint64)
		}
		amount, err := strconv.ParseInt(op.Amount.Value, 10, 64)
		poolCoin = sdk.NewCoin(op.Amount.Currency.Symbol, sdk.NewInt(amount))
		if err != nil {
			return nil, fmt.Errorf("invalid amount: %w", err)
		}
	}
	return NewMsgWithdrawFromLiquidityPool(withdrawer, poolId, poolCoin), nil
}
