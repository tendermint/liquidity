package keeper

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/tendermint/liquidity/types"
)

// Keeper of the liquidity store
type Keeper struct {
	cdc           codec.Marshaler
	storeKey      sdk.StoreKey
	bankKeeper    types.BankKeeper
	accountKeeper types.AccountKeeper
	paramSpace    paramstypes.Subspace
}

// NewKeeper returns a liquidity keeper. It handles:
// - creating new ModuleAccounts for each pool ReserveAccount
// - sending to and from ModuleAccounts
// - minting, burning PoolCoins
func NewKeeper(cdc codec.Marshaler, key sdk.StoreKey, paramSpace paramstypes.Subspace, bankKeeper types.BankKeeper, accountKeeper types.AccountKeeper) Keeper {
	// ensure liquidity module account is set
	if addr := accountKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	// set KeyTable if it has not already been set
	if !paramSpace.HasKeyTable() {
		paramSpace = paramSpace.WithKeyTable(types.ParamKeyTable())
	}

	return Keeper{
		storeKey:      key,
		bankKeeper:    bankKeeper,
		accountKeeper: accountKeeper,
		cdc:           cdc,
		paramSpace:    paramSpace,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", types.ModuleName)
}

func (k Keeper) Swap(ctx sdk.Context, msg *types.MsgSwap) error {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSwap,
		),
	)
	return nil
}

func (k Keeper) CreateLiquidityPool(ctx sdk.Context, msg *types.MsgCreateLiquidityPool) error {
	params := k.GetParams(ctx)
	var poolType types.LiquidityPoolType

	// check poolType exist, get poolType from param
	if int(msg.PoolTypeIndex) > len(params.LiquidityPoolTypes)-1 {
		poolType = params.LiquidityPoolTypes[msg.PoolTypeIndex]
		if poolType.PoolTypeIndex != msg.PoolTypeIndex {
			return types.ErrPoolTypeNotExists
		}
	} else {
		return types.ErrPoolTypeNotExists
	}

	if poolType.MinReserveCoinNum != 2 && poolType.MaxReserveCoinNum != 2 {
		return types.ErrNotImplementedYet
	}

	if len(msg.ReserveCoinDenoms) > int(poolType.MaxReserveCoinNum) && int(poolType.MinReserveCoinNum) > len(msg.ReserveCoinDenoms) {
		return types.ErrNumOfReserveCoin
	}

	poolKey := types.GetPoolKey(msg.ReserveCoinDenoms, msg.PoolTypeIndex)
	reserveAcc := types.GetPoolReserveAcc(poolKey)

	if _, found := k.GetLiquidityPoolByReserveAccIndex(ctx, reserveAcc); found {
		return types.ErrPoolAlreadyExists
	}

	accPoolCreator := k.accountKeeper.GetAccount(ctx, msg.PoolCreator)
	poolCreatorBalances := k.bankKeeper.GetAllBalances(ctx, accPoolCreator.GetAddress())
	if !poolCreatorBalances.IsAllGTE(msg.DepositCoinsAmount) {
		return types.ErrInsufficientBalance
	}

	for _, coin := range msg.DepositCoinsAmount {
		if coin.Amount.LT(params.MinInitDepositToPool) {
			return types.ErrLessThanMinInitDeposit
		}
	}

	denom1, denom2 := types.AlphabeticalDenomPair(msg.ReserveCoinDenoms[0], msg.ReserveCoinDenoms[1])
	reserveCoinDenoms := []string{denom1, denom2}

	PoolCoinDenom := types.GetPoolCoinDenom(reserveAcc)

	liquidityPool := types.LiquidityPool{
		PoolTypeIndex:     msg.PoolTypeIndex,
		ReserveCoinDenoms: reserveCoinDenoms,
		ReserveAccount:    reserveAcc,
		PoolCoinDenom:     PoolCoinDenom,
	}

	mintPoolCoin := sdk.NewCoins(sdk.NewCoin(liquidityPool.PoolCoinDenom, params.InitPoolCoinMintAmount))
	if err := k.bankKeeper.SendCoins(ctx, msg.PoolCreator, liquidityPool.ReserveAccount, msg.DepositCoinsAmount); err != nil {
		return err
	}
	if err := k.bankKeeper.MintCoins(ctx, types.ModuleName, mintPoolCoin); err != nil {
		return err
	}
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, msg.PoolCreator, mintPoolCoin); err != nil {
		return err
	}

	k.SetLiquidityPoolAtomic(ctx, liquidityPool)

	// TODO: atomic transfer using like InputOutputCoins
	//var MultiSendInput []bankTypes.Input
	//var MultiSendOutput []bankTypes.Output
	//MultiSendInput = append(MultiSendInput, bankTypes.NewInput(msg.PoolCreator, msg.DepositCoinsAmount))

	// TODO: refactoring, LiquidityPoolCreationFee, check event on handler
	return nil
}

func (k Keeper) DepositLiquidityPoolToBatch(ctx sdk.Context, msg *types.MsgDepositToLiquidityPool) error {
	return types.ErrNotImplementedYet
}

func (k Keeper) WithdrawLiquidityPoolToBatch(ctx sdk.Context, msg *types.MsgWithdrawFromLiquidityPool) error {
	return types.ErrNotImplementedYet
}

func (k Keeper) SwapToBatch(ctx sdk.Context, msg *types.MsgSwap) error {
	return types.ErrNotImplementedYet
}

// GetParams gets the parameters for the liquidity module.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the parameters for the liquidity module.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}
