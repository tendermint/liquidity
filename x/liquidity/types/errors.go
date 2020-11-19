package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// liquidity module sentinel errors
var (
	ErrPoolNotExists          = sdkerrors.Register(ModuleName, 1, "pool not exists")
	ErrPoolTypeNotExists      = sdkerrors.Register(ModuleName, 2, "pool type not exists")
	ErrEqualDenom             = sdkerrors.Register(ModuleName, 3, "reserve coin denomination are equal")
	ErrInvalidDenom           = sdkerrors.Register(ModuleName, 4, "invalid denom")
	ErrNumOfReserveCoin       = sdkerrors.Register(ModuleName, 5, "invalid number of reserve coin")
	ErrNumOfPoolCoin          = sdkerrors.Register(ModuleName, 6, "invalid number of pool coin")
	ErrInsufficientPool       = sdkerrors.Register(ModuleName, 7, "insufficient pool")
	ErrInsufficientBalance    = sdkerrors.Register(ModuleName, 8, "insufficient coin balance to escrow")
	ErrLessThanMinInitDeposit = sdkerrors.Register(ModuleName, 9, "deposit coin less than MinInitDepositToPool")
	ErrNotImplementedYet      = sdkerrors.Register(ModuleName, 10, "not implemented yet")
	ErrPoolAlreadyExists      = sdkerrors.Register(ModuleName, 11, "the pool already exists")
	ErrPoolBatchNotExists     = sdkerrors.Register(ModuleName, 12, "pool batch not exists")
	ErrOrderBookInvalidity    = sdkerrors.Register(ModuleName, 13, "orderbook is not validity")
	ErrBatchNotExecuted       = sdkerrors.Register(ModuleName, 14, "the liquidity pool batch is not executed")
)
