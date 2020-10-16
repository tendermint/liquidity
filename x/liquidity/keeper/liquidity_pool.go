package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	gogotypes "github.com/gogo/protobuf/types"
	"github.com/tendermint/liquidity/x/liquidity/types"
)

// return a specific liquidityPool
func (k Keeper) GetLiquidityPool(ctx sdk.Context, poolID uint64) (liquidityPool types.LiquidityPool, found bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetLiquidityPoolKey(poolID)

	value := store.Get(key)
	if value == nil {
		return liquidityPool, false
	}

	liquidityPool = types.MustUnmarshalLiquidityPool(k.cdc, value)

	return liquidityPool, true
}

func (k Keeper) SetLiquidityPool(ctx sdk.Context, liquidityPool types.LiquidityPool) {
	store := ctx.KVStore(k.storeKey)
	b := types.MustMarshalLiquidityPool(k.cdc, liquidityPool)
	store.Set(types.GetLiquidityPoolKey(liquidityPool.PoolID), b)
}

// IterateAllLiquidityPools iterate through all of the liquidityPools
func (k Keeper) IterateAllLiquidityPools(ctx sdk.Context, cb func(liquidityPool types.LiquidityPool) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.LiquidityPoolKeyPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		liquidityPool := types.MustUnmarshalLiquidityPool(k.cdc, iterator.Value())
		if cb(liquidityPool) {
			break
		}
	}
}

// GetAllLiquidityPools returns all liquidityPools used during genesis dump
func (k Keeper) GetAllLiquidityPools(ctx sdk.Context) (liquidityPools []types.LiquidityPool) {
	k.IterateAllLiquidityPools(ctx, func(liquidityPool types.LiquidityPool) bool {
		liquidityPools = append(liquidityPools, liquidityPool)
		return false
	})

	return liquidityPools
}

// GetNextLiquidityID returns and increments the global Liquidity Pool ID counter.
// If the global account number is not set, it initializes it with value 0.
func (k Keeper) GetNextLiquidityPoolID(ctx sdk.Context) uint64 {
	var poolID uint64
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GlobalLiquidityPoolIDKey)
	if bz == nil {
		// initialize the LiquidityPoolID
		poolID = 0
	} else {
		val := gogotypes.UInt64Value{}

		err := k.cdc.UnmarshalBinaryBare(bz, &val)
		if err != nil {
			panic(err)
		}

		poolID = val.GetValue()
	}

	bz = k.cdc.MustMarshalBinaryBare(&gogotypes.UInt64Value{Value: poolID + 1})
	store.Set(types.GlobalLiquidityPoolIDKey, bz)

	return poolID
}

func (k Keeper) GetLiquidityPoolByReserveAccIndex(ctx sdk.Context, reserveAcc sdk.AccAddress) (liquidityPool types.LiquidityPool, found bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetLiquidityPoolByReserveAccIndexKey(reserveAcc)

	value := store.Get(key)
	if value == nil {
		return liquidityPool, false
	}

	liquidityPool = types.MustUnmarshalLiquidityPool(k.cdc, value)

	return liquidityPool, true
}

func (k Keeper) SetLiquidityPoolByReserveAccIndex(ctx sdk.Context, liquidityPool types.LiquidityPool) {
	store := ctx.KVStore(k.storeKey)
	b := types.MustMarshalLiquidityPool(k.cdc, liquidityPool)
	store.Set(types.GetLiquidityPoolByReserveAccIndexKey(liquidityPool.ReserveAccount), b)
}

func (k Keeper) SetLiquidityPoolAtomic(ctx sdk.Context, liquidityPool types.LiquidityPool) {
	liquidityPool.PoolID = k.GetNextLiquidityPoolID(ctx)
	k.SetLiquidityPool(ctx, liquidityPool)
	k.SetLiquidityPoolByReserveAccIndex(ctx, liquidityPool)
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
	if !poolCreatorBalances.IsAllGTE(msg.DepositCoins) {
		return types.ErrInsufficientBalance
	}

	for _, coin := range msg.DepositCoins {
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
	if err := k.bankKeeper.SendCoins(ctx, msg.PoolCreator, liquidityPool.ReserveAccount, msg.DepositCoins); err != nil {
		return err
	}
	// TODO: fix module Name as poolKey or moduleName
	if err := k.bankKeeper.MintCoins(ctx, poolKey, mintPoolCoin); err != nil {
		return err
	}
	if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, poolKey, msg.PoolCreator, mintPoolCoin); err != nil {
		return err
	}

	k.SetLiquidityPoolAtomic(ctx, liquidityPool)

	// TODO: atomic transfer using like InputOutputCoins
	//var MultiSendInput []bankTypes.Input
	//var MultiSendOutput []bankTypes.Output
	//MultiSendInput = append(MultiSendInput, bankTypes.NewInput(msg.PoolCreator, msg.DepositCoins))

	// TODO: params.LiquidityPoolCreationFee logic
	// TODO: refactoring, LiquidityPoolCreationFee, check event on handler
	return nil
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func (k Keeper) GetReserveCoins(ctx sdk.Context, lp types.LiquidityPool) (reserveCoins sdk.Coins) {
	for _, denom := range lp.ReserveCoinDenoms {
		reserveCoins = reserveCoins.Add(k.bankKeeper.GetBalance(ctx, lp.ReserveAccount, denom))
	}
	return
}

func (k Keeper) GetPoolCoinTotalSupply(ctx sdk.Context, lp types.LiquidityPool) sdk.Int {
	supply := k.bankKeeper.GetSupply(ctx)
	total := supply.GetTotal()
	return total.AmountOf(lp.PoolCoinDenom)
	//TODO: add totalSupply on liquidityPool
}

func (k Keeper) DepositLiquidityPool(ctx sdk.Context, msg *types.MsgDepositToLiquidityPool) error {
	depositCoins := (sdk.Coins)(msg.DepositCoins)
	depositCoins = depositCoins.Sort()
	if err := depositCoins.Validate(); err != nil {
		return err
	}

	pool, found := k.GetLiquidityPool(ctx, msg.PoolID)
	if !found {
		return types.ErrPoolNotExists
	}
	for _, coin := range depositCoins {
		if !stringInSlice(coin.Denom, pool.ReserveCoinDenoms) {
			return types.ErrInvalidDenom
		}
	}

	if depositCoins.Len() != 2 || len(pool.ReserveCoinDenoms) != 2 {
		return types.ErrNotImplementedYet
	}

	if depositCoins.Len() != len(pool.ReserveCoinDenoms) {
		return types.ErrNumOfReserveCoin
	}

	reserveCoins := k.GetReserveCoins(ctx, pool)
	reserveCoins.Sort()
	// TODO: coin sort alphabetical
	coinA := depositCoins[0]
	coinB := depositCoins[1]
	//lastReserveRatio := coinA.Amount.Quo(coinB.Amount)
	lastReserveRatio := reserveCoins[0].Amount.Quo(reserveCoins[1].Amount)
	depositableCoinA := coinA.Amount.Mul(lastReserveRatio)
	depositableCoinB := coinB.Amount.Mul(lastReserveRatio)
	var inputs []banktypes.Input
	var outputs []banktypes.Output

	// TODO: msg.Depositor to escrowModAcc
	batchEscrowAcc := msg.Depositor

	if coinB.Amount.GT(depositableCoinA) {
		depositableCoinB = depositableCoinA
		refundAmtB := coinA.Amount.Sub(depositableCoinB)
		inputs = append(inputs, banktypes.NewInput(batchEscrowAcc, sdk.NewCoins(coinA)))
		inputs = append(inputs, banktypes.NewInput(batchEscrowAcc, sdk.NewCoins(sdk.NewCoin(coinB.Denom, depositableCoinB))))
		outputs = append(outputs, banktypes.NewOutput(pool.ReserveAccount, sdk.NewCoins(coinA)))
		outputs = append(outputs, banktypes.NewOutput(pool.ReserveAccount, sdk.NewCoins(sdk.NewCoin(coinB.Denom, depositableCoinB))))
		// refund
		inputs = append(inputs, banktypes.NewInput(batchEscrowAcc, sdk.NewCoins(sdk.NewCoin(coinB.Denom, refundAmtB))))
		outputs = append(outputs, banktypes.NewOutput(msg.Depositor, sdk.NewCoins(sdk.NewCoin(coinB.Denom, refundAmtB))))
	} else if coinB.Amount.LT(depositableCoinA) {
		depositableCoinA = coinB.Amount.Quo(lastReserveRatio)
		refundAmtA := coinA.Amount.Sub(depositableCoinA)
		inputs = append(inputs, banktypes.NewInput(batchEscrowAcc, sdk.NewCoins(sdk.NewCoin(coinA.Denom, depositableCoinA))))
		inputs = append(inputs, banktypes.NewInput(batchEscrowAcc, sdk.NewCoins(coinB)))
		outputs = append(outputs, banktypes.NewOutput(pool.ReserveAccount, sdk.NewCoins(sdk.NewCoin(coinA.Denom, depositableCoinA))))
		outputs = append(outputs, banktypes.NewOutput(pool.ReserveAccount, sdk.NewCoins(coinB)))
		// refund
		inputs = append(inputs, banktypes.NewInput(batchEscrowAcc, sdk.NewCoins(sdk.NewCoin(coinA.Denom, refundAmtA))))
		outputs = append(outputs, banktypes.NewOutput(msg.Depositor, sdk.NewCoins(sdk.NewCoin(coinA.Denom, refundAmtA))))
	} else {
		inputs = append(inputs, banktypes.NewInput(batchEscrowAcc, sdk.NewCoins(coinA)))
		inputs = append(inputs, banktypes.NewInput(batchEscrowAcc, sdk.NewCoins(coinB)))
		outputs = append(outputs, banktypes.NewOutput(pool.ReserveAccount, sdk.NewCoins(coinA)))
		outputs = append(outputs, banktypes.NewOutput(pool.ReserveAccount, sdk.NewCoins(coinB)))
	}
	// TODO: InputOutputCoins impossible for moduleAcc, SendCoinsFromModuleToAccount need string for moduleName, fix to SendCoinsFromModuleToAccount
	//k.bankKeeper.SendCoinsFromModuleToAccount()

	// calculate pool token mint amount
	poolCoinAmt := k.GetPoolCoinTotalSupply(ctx, pool).Mul(depositableCoinA).Quo(coinA.Amount)
	poolCoin := sdk.NewCoins(sdk.NewCoin(pool.PoolCoinDenom, poolCoinAmt))
	// mint pool token to Depositor
	if err := k.bankKeeper.MintCoins(ctx, pool.GetPoolKey(), poolCoin); err != nil {
		return err
	}

	// TODO: fix to SendCoinsFromModuleToAccount
	//k.bankKeeper.SendCoinsFromModuleToAccount()
	inputs = append(inputs, banktypes.NewInput(pool.ReserveAccount, poolCoin))
	outputs = append(outputs, banktypes.NewOutput(msg.Depositor, poolCoin))

	// execute multi-send
	if err := k.bankKeeper.InputOutputCoins(ctx, inputs, outputs); err != nil {
		return err
	}
	return nil
}

func (k Keeper) WithdrawLiquidityPool(ctx sdk.Context, msg *types.MsgWithdrawFromLiquidityPool) error {
	poolCoins := (sdk.Coins)(msg.PoolCoin)
	poolCoins = poolCoins.Sort()
	if err := poolCoins.Validate(); err != nil {
		return err
	}
	if poolCoins.Len() != 1 {
		return types.ErrNumOfReserveCoin
	}
	poolCoin := poolCoins[0]

	pool, found := k.GetLiquidityPool(ctx, msg.PoolID)
	if !found {
		return types.ErrPoolNotExists
	}

	totalSupply := k.GetPoolCoinTotalSupply(ctx, pool)
	reserveCoins := k.GetReserveCoins(ctx, pool)
	reserveCoins.Sort()

	var inputs []banktypes.Input
	var outputs []banktypes.Output

	//params := k.GetParams(ctx)
	for _, reserveCoin := range reserveCoins {
		withdrawAmt := reserveCoin.Amount.Mul(poolCoin.Amount).Quo(totalSupply)
		// TODO: apply fee, (sdk.NewDec(1).Sub(params.LiquidityPoolFeeRate)
		// TODO: to using k.bankKeeper.SendCoinsFromModuleToAccount() with poolKey
		inputs = append(inputs, banktypes.NewInput(pool.ReserveAccount, sdk.NewCoins(sdk.NewCoin(reserveCoin.Denom, withdrawAmt))))
		outputs = append(outputs, banktypes.NewOutput(msg.Withdrawer, sdk.NewCoins(sdk.NewCoin(reserveCoin.Denom, withdrawAmt))))
	}

	// execute multi-send
	if err := k.bankKeeper.InputOutputCoins(ctx, inputs, outputs); err != nil {
		return err
	}
	// TODO: apply pool.GetPoolKey() as moduleName
	k.bankKeeper.BurnCoins(ctx, pool.GetPoolKey(), poolCoins)

	return nil
}

func (k Keeper) DepositLiquidityPoolToBatch(ctx sdk.Context, msg *types.MsgDepositToLiquidityPool) error {
	return types.ErrNotImplementedYet
}

func (k Keeper) WithdrawLiquidityPoolToBatch(ctx sdk.Context, msg *types.MsgWithdrawFromLiquidityPool) error {
	return types.ErrNotImplementedYet
}
