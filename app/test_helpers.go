package app

// DONTCOVER

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/tendermint/liquidity/x/liquidity"
	"github.com/tendermint/liquidity/x/liquidity/keeper"
	"github.com/tendermint/liquidity/x/liquidity/types"
)

// DefaultConsensusParams defines the default Tendermint consensus params used in
// LiquidityApp testing.
var DefaultConsensusParams = &abci.ConsensusParams{
	Block: &abci.BlockParams{
		MaxBytes: 200000,
		MaxGas:   2000000,
	},
	Evidence: &tmproto.EvidenceParams{
		MaxAgeNumBlocks: 302400,
		MaxAgeDuration:  504 * time.Hour, // 3 weeks is the max duration
	},
	Validator: &tmproto.ValidatorParams{
		PubKeyTypes: []string{
			tmtypes.ABCIPubKeyTypeEd25519,
		},
	},
}

// Setup initializes a new LiquidityApp. A Nop logger is set in LiquidityApp.
func Setup(isCheckTx bool) *LiquidityApp {
	db := dbm.NewMemDB()
	app := NewLiquidityApp(log.NewNopLogger(), db, nil, true, map[int64]bool{}, DefaultNodeHome, 5, MakeEncodingConfig(), EmptyAppOptions{})
	if !isCheckTx {
		// init chain must be called to stop deliverState from being nil
		genesisState := NewDefaultGenesisState()
		stateBytes, err := json.MarshalIndent(genesisState, "", " ")
		if err != nil {
			panic(err)
		}

		// Initialize the chain
		app.InitChain(
			abci.RequestInitChain{
				Validators:      []abci.ValidatorUpdate{},
				ConsensusParams: DefaultConsensusParams,
				AppStateBytes:   stateBytes,
			},
		)
	}

	return app
}

type GenerateAccountStrategy func(int) []sdk.AccAddress

// createRandomAccounts is a strategy used by addTestAddrs() in order to generated addresses in random order.
func createRandomAccounts(accNum int) []sdk.AccAddress {
	testAddrs := make([]sdk.AccAddress, accNum)
	for i := 0; i < accNum; i++ {
		pk := ed25519.GenPrivKey().PubKey()
		testAddrs[i] = sdk.AccAddress(pk.Address())
	}

	return testAddrs
}

// createIncrementalAccounts is a strategy used by addTestAddrs() in order to generated addresses in ascending order.
func createIncrementalAccounts(accNum int) []sdk.AccAddress {
	var addresses []sdk.AccAddress
	var buffer bytes.Buffer

	// start at 100 so we can make up to 999 test addresses with valid test addresses
	for i := 100; i < (accNum + 100); i++ {
		numString := strconv.Itoa(i)
		buffer.WriteString("A58856F0FD53BF058B4909A21AEC019107BA6") //base address string

		buffer.WriteString(numString) //adding on final two digits to make addresses unique
		res, _ := sdk.AccAddressFromHex(buffer.String())
		bech := res.String()
		addr, _ := TestAddr(buffer.String(), bech)

		addresses = append(addresses, addr)
		buffer.Reset()
	}

	return addresses
}

// setTotalSupply provides the total supply based on accAmt * totalAccounts.
func setTotalSupply(app *LiquidityApp, ctx sdk.Context, accAmt sdk.Int, totalAccounts int) {
	totalSupply := sdk.NewCoins(sdk.NewCoin(app.StakingKeeper.BondDenom(ctx), accAmt.MulRaw(int64(totalAccounts))))
	prevSupply := app.BankKeeper.GetSupply(ctx)
	app.BankKeeper.SetSupply(ctx, banktypes.NewSupply(prevSupply.GetTotal().Add(totalSupply...)))
}

func addTotalSupply(app *LiquidityApp, ctx sdk.Context, coins sdk.Coins) {
	prevSupply := app.BankKeeper.GetSupply(ctx)
	app.BankKeeper.SetSupply(ctx, banktypes.NewSupply(prevSupply.GetTotal().Add(coins...)))
}

// AddRandomTestAddr creates new account with random address.
func AddRandomTestAddr(app *LiquidityApp, ctx sdk.Context, initCoins sdk.Coins) sdk.AccAddress {
	addr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	addTotalSupply(app, ctx, initCoins)
	SaveAccount(app, ctx, addr, initCoins)
	return addr
}

// AddTestAddrs constructs and returns accNum amount of accounts with an
// initial balance of accAmt in random order
func AddTestAddrs(app *LiquidityApp, ctx sdk.Context, accNum int, initCoins sdk.Coins) []sdk.AccAddress {
	testAddrs := createIncrementalAccounts(accNum)
	for _, addr := range testAddrs {
		addTotalSupply(app, ctx, initCoins)
		SaveAccount(app, ctx, addr, initCoins)
	}
	return testAddrs
}

// AddTestAddrs constructs and returns accNum amount of accounts with an
// initial balance of accAmt in random order
func AddTestAddrsIncremental(app *LiquidityApp, ctx sdk.Context, accNum int, accAmt sdk.Int) []sdk.AccAddress {
	return addTestAddrs(app, ctx, accNum, accAmt, createIncrementalAccounts)
}

func addTestAddrs(app *LiquidityApp, ctx sdk.Context, accNum int, accAmt sdk.Int, strategy GenerateAccountStrategy) []sdk.AccAddress {
	testAddrs := strategy(accNum)

	initCoins := sdk.NewCoins(sdk.NewCoin(app.StakingKeeper.BondDenom(ctx), accAmt))
	setTotalSupply(app, ctx, accAmt, accNum)

	// fill all the addresses with some coins, set the loose pool tokens simultaneously
	for _, addr := range testAddrs {
		SaveAccount(app, ctx, addr, initCoins)
	}

	return testAddrs
}

// SaveAccount saves the provided account into the simapp with balance based on initCoins.
func SaveAccount(app *LiquidityApp, ctx sdk.Context, addr sdk.AccAddress, initCoins sdk.Coins) {
	acc := app.AccountKeeper.NewAccountWithAddress(ctx, addr)
	app.AccountKeeper.SetAccount(ctx, acc)
	err := app.BankKeeper.AddCoins(ctx, addr, initCoins)
	if err != nil {
		panic(err)
	}
}

func SaveAccountWithFee(app *LiquidityApp, ctx sdk.Context, addr sdk.AccAddress, initCoins sdk.Coins, offerCoin sdk.Coin) {
	SaveAccount(app, ctx, addr, initCoins)
	//acc := app.AccountKeeper.GetAccount(ctx, addr)
	params := app.LiquidityKeeper.GetParams(ctx)
	offerCoinFee := types.GetOfferCoinFee(offerCoin, params.SwapFeeRate)
	err := app.BankKeeper.AddCoins(ctx, addr, sdk.NewCoins(offerCoinFee))
	if err != nil {
		panic(err)
	}
}

func TestAddr(addr string, bech string) (sdk.AccAddress, error) {
	res, err := sdk.AccAddressFromHex(addr)
	if err != nil {
		return nil, err
	}
	bechexpected := res.String()
	if bech != bechexpected {
		return nil, fmt.Errorf("bech encoding doesn't match reference")
	}

	bechres, err := sdk.AccAddressFromBech32(bech)
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(bechres, res) {
		return nil, err
	}

	return res, nil
}

// CreateTestInput returns a simapp with custom LiquidityKeeper to avoid
// messing with the hooks.
func CreateTestInput() (*LiquidityApp, sdk.Context) {
	cdc := codec.NewLegacyAmino()
	types.RegisterLegacyAminoCodec(cdc)

	app := Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	appCodec := app.AppCodec()

	app.LiquidityKeeper = keeper.NewKeeper(
		appCodec,
		app.GetKey(types.StoreKey),
		app.GetSubspace(types.ModuleName),
		app.BankKeeper,
		app.AccountKeeper,
		app.DistrKeeper,
	)

	return app, ctx
}

func GetRandPoolAmt(r *rand.Rand, minInitDepositAmt sdk.Int) (X, Y sdk.Int) {
	X = GetRandRange(r, int(minInitDepositAmt.Int64()), 100000000000000).MulRaw(int64(math.Pow10(r.Intn(10))))
	Y = GetRandRange(r, int(minInitDepositAmt.Int64()), 100000000000000).MulRaw(int64(math.Pow10(r.Intn(10))))
	//fmt.Println(X, Y, X.ToDec().Quo(Y.ToDec()))
	return
}

func GetRandFloats(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func GetRandRange(r *rand.Rand, min, max int) sdk.Int {
	return sdk.NewInt(int64(r.Intn(max-min) + min))
}

func GetRandomSizeOrders(denomX, denomY string, X, Y sdk.Int, r *rand.Rand, sizeXtoY, sizeYtoX int32) (XtoY, YtoX []*types.MsgSwapWithinBatch) {
	randomSizeXtoY := int(r.Int31n(sizeXtoY))
	randomSizeYtoX := int(r.Int31n(sizeYtoX))
	return GetRandomOrders(denomX, denomY, X, Y, r, randomSizeXtoY, randomSizeYtoX)
}

func GetRandomOrders(denomX, denomY string, X, Y sdk.Int, r *rand.Rand, sizeXtoY, sizeYtoX int) (XtoY, YtoX []*types.MsgSwapWithinBatch) {
	currentPrice := X.ToDec().Quo(Y.ToDec())

	for len(XtoY) < sizeXtoY {
		orderPrice := currentPrice.Mul(sdk.NewDecFromIntWithPrec(GetRandRange(r, 991, 1009), 3))
		//offerAmt := X.ToDec().Mul(sdk.NewDecFromIntWithPrec(GetRandRange(r, 1, 100), 4))
		orderAmt := sdk.ZeroDec()
		if r.Intn(2) == 1 {
			orderAmt = X.ToDec().Mul(sdk.NewDecFromIntWithPrec(GetRandRange(r, 1, 100), 4))
		} else {
			orderAmt = sdk.NewDecFromIntWithPrec(GetRandRange(r, 1000, 10000), 0)
		}
		if orderAmt.Quo(orderPrice).TruncateInt().IsZero() {
			continue
		}
		orderCoin := sdk.NewCoin(denomX, orderAmt.Ceil().TruncateInt())

		XtoY = append(XtoY, &types.MsgSwapWithinBatch{
			OfferCoin:       orderCoin,
			DemandCoinDenom: denomY,
			OrderPrice:      orderPrice,
		})
	}

	for len(YtoX) < sizeYtoX {
		orderPrice := currentPrice.Mul(sdk.NewDecFromIntWithPrec(GetRandRange(r, 991, 1009), 3))
		//offerAmt := Y.ToDec().Mul(sdk.NewDecFromIntWithPrec(GetRandRange(r, 1, 100), 4))
		orderAmt := sdk.ZeroDec()
		if r.Intn(2) == 1 {
			orderAmt = Y.ToDec().Mul(sdk.NewDecFromIntWithPrec(GetRandRange(r, 1, 100), 4))
		} else {
			orderAmt = sdk.NewDecFromIntWithPrec(GetRandRange(r, 1000, 10000), 0)
		}
		if orderAmt.Mul(orderPrice).TruncateInt().IsZero() {
			continue
		}
		orderCoin := sdk.NewCoin(denomY, orderAmt.Ceil().TruncateInt())

		YtoX = append(YtoX, &types.MsgSwapWithinBatch{
			OfferCoin:       orderCoin,
			DemandCoinDenom: denomX,
			OrderPrice:      orderPrice,
		})
	}
	return
}

func TestCreatePool(t *testing.T, simapp *LiquidityApp, ctx sdk.Context, X, Y sdk.Int, denomX, denomY string, addr sdk.AccAddress) uint64 {
	deposit := sdk.NewCoins(sdk.NewCoin(denomX, X), sdk.NewCoin(denomY, Y))
	params := simapp.LiquidityKeeper.GetParams(ctx)
	// set accounts for creator, depositor, withdrawer, balance for deposit
	SaveAccount(simapp, ctx, addr, deposit.Add(params.PoolCreationFee...)) // pool creator
	depositX := simapp.BankKeeper.GetBalance(ctx, addr, denomX)
	depositY := simapp.BankKeeper.GetBalance(ctx, addr, denomY)
	depositBalance := sdk.NewCoins(depositX, depositY)
	require.Equal(t, deposit, depositBalance)

	// create Liquidity pool
	poolTypeId := types.DefaultPoolTypeId
	poolId := simapp.LiquidityKeeper.GetNextPoolId(ctx)
	msg := types.NewMsgCreatePool(addr, poolTypeId, depositBalance)
	_, err := simapp.LiquidityKeeper.CreatePool(ctx, msg)
	require.NoError(t, err)

	// verify created liquidity pool
	pool, found := simapp.LiquidityKeeper.GetPool(ctx, poolId)
	require.True(t, found)
	require.Equal(t, poolId, pool.Id)
	require.Equal(t, denomX, pool.ReserveCoinDenoms[0])
	require.Equal(t, denomY, pool.ReserveCoinDenoms[1])

	// verify minted pool coin
	poolCoin := simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, pool)
	creatorBalance := simapp.BankKeeper.GetBalance(ctx, addr, pool.PoolCoinDenom)
	require.Equal(t, poolCoin, creatorBalance.Amount)
	return poolId
}

func TestDepositPool(t *testing.T, simapp *LiquidityApp, ctx sdk.Context, X, Y sdk.Int, addrs []sdk.AccAddress, poolId uint64, withEndblock bool) {
	pool, found := simapp.LiquidityKeeper.GetPool(ctx, poolId)
	require.True(t, found)
	denomX, denomY := pool.ReserveCoinDenoms[0], pool.ReserveCoinDenoms[1]
	deposit := sdk.NewCoins(sdk.NewCoin(denomX, X), sdk.NewCoin(denomY, Y))

	moduleAccAddress := simapp.AccountKeeper.GetModuleAddress(types.ModuleName)
	moduleAccEscrowAmtX := simapp.BankKeeper.GetBalance(ctx, moduleAccAddress, denomX)
	moduleAccEscrowAmtY := simapp.BankKeeper.GetBalance(ctx, moduleAccAddress, denomY)
	iterNum := len(addrs)
	for i := 0; i < iterNum; i++ {
		SaveAccount(simapp, ctx, addrs[i], deposit) // pool creator

		depositMsg := types.NewMsgDepositWithinBatch(addrs[i], poolId, deposit)
		_, err := simapp.LiquidityKeeper.DepositLiquidityPoolToBatch(ctx, depositMsg)
		require.NoError(t, err)

		depositorBalanceX := simapp.BankKeeper.GetBalance(ctx, addrs[i], pool.ReserveCoinDenoms[0])
		depositorBalanceY := simapp.BankKeeper.GetBalance(ctx, addrs[i], pool.ReserveCoinDenoms[1])
		require.Equal(t, denomX, depositorBalanceX.Denom)
		require.Equal(t, denomY, depositorBalanceY.Denom)

		// check escrow balance of module account
		moduleAccEscrowAmtX = moduleAccEscrowAmtX.Add(deposit[0])
		moduleAccEscrowAmtY = moduleAccEscrowAmtY.Add(deposit[1])
		moduleAccEscrowAmtXAfter := simapp.BankKeeper.GetBalance(ctx, moduleAccAddress, denomX)
		moduleAccEscrowAmtYAfter := simapp.BankKeeper.GetBalance(ctx, moduleAccAddress, denomY)
		require.Equal(t, moduleAccEscrowAmtX, moduleAccEscrowAmtXAfter)
		require.Equal(t, moduleAccEscrowAmtY, moduleAccEscrowAmtYAfter)
	}
	batch, bool := simapp.LiquidityKeeper.GetPoolBatch(ctx, poolId)
	require.True(t, bool)
	msgs := simapp.LiquidityKeeper.GetAllPoolBatchDepositMsgs(ctx, batch)

	// endblock
	if withEndblock {
		liquidity.EndBlocker(ctx, simapp.LiquidityKeeper)
		msgs = simapp.LiquidityKeeper.GetAllPoolBatchDepositMsgs(ctx, batch)
		for i := 0; i < iterNum; i++ {
			// verify minted pool coin
			poolCoin := simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, pool)
			depositorPoolCoinBalance := simapp.BankKeeper.GetBalance(ctx, addrs[i], pool.PoolCoinDenom)
			require.NotEqual(t, sdk.ZeroInt(), depositorPoolCoinBalance)
			require.NotEqual(t, sdk.ZeroInt(), poolCoin)

			require.True(t, msgs[i].Executed)
			require.True(t, msgs[i].Succeeded)
			require.True(t, msgs[i].ToBeDeleted)

			// error balance after endblock
			depositorBalanceX := simapp.BankKeeper.GetBalance(ctx, addrs[i], pool.ReserveCoinDenoms[0])
			depositorBalanceY := simapp.BankKeeper.GetBalance(ctx, addrs[i], pool.ReserveCoinDenoms[1])
			require.Equal(t, denomX, depositorBalanceX.Denom)
			require.Equal(t, denomY, depositorBalanceY.Denom)
		}
	}
}

func TestWithdrawPool(t *testing.T, simapp *LiquidityApp, ctx sdk.Context, poolCoinAmt sdk.Int, addrs []sdk.AccAddress, poolId uint64, withEndblock bool) {
	pool, found := simapp.LiquidityKeeper.GetPool(ctx, poolId)
	require.True(t, found)
	//denomX, denomY := pool.ReserveCoinDenoms[0], pool.ReserveCoinDenoms[1]
	moduleAccAddress := simapp.AccountKeeper.GetModuleAddress(types.ModuleName)
	moduleAccEscrowAmtPool := simapp.BankKeeper.GetBalance(ctx, moduleAccAddress, pool.PoolCoinDenom)

	iterNum := len(addrs)
	for i := 0; i < iterNum; i++ {
		balancePoolCoin := simapp.BankKeeper.GetBalance(ctx, addrs[i], pool.PoolCoinDenom)
		require.True(t, balancePoolCoin.Amount.GTE(poolCoinAmt))

		withdrawCoin := sdk.NewCoin(pool.PoolCoinDenom, poolCoinAmt)
		withdrawMsg := types.NewMsgWithdrawWithinBatch(addrs[i], poolId, withdrawCoin)
		_, err := simapp.LiquidityKeeper.WithdrawLiquidityPoolToBatch(ctx, withdrawMsg)
		require.NoError(t, err)

		moduleAccEscrowAmtPoolAfter := simapp.BankKeeper.GetBalance(ctx, moduleAccAddress, pool.PoolCoinDenom)
		moduleAccEscrowAmtPool.Amount = moduleAccEscrowAmtPool.Amount.Add(withdrawMsg.PoolCoin.Amount)
		require.Equal(t, moduleAccEscrowAmtPool, moduleAccEscrowAmtPoolAfter)

		balancePoolCoinAfter := simapp.BankKeeper.GetBalance(ctx, addrs[i], pool.PoolCoinDenom)
		if balancePoolCoin.Amount.Equal(withdrawCoin.Amount) {

		} else {
			require.Equal(t, balancePoolCoin.Sub(withdrawCoin).Amount, balancePoolCoinAfter.Amount)
		}

	}
	batch, bool := simapp.LiquidityKeeper.GetPoolBatch(ctx, poolId)

	if withEndblock {
		poolCoinBefore := simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, pool)

		// endblock
		liquidity.EndBlocker(ctx, simapp.LiquidityKeeper)

		batch, bool = simapp.LiquidityKeeper.GetPoolBatch(ctx, poolId)
		require.True(t, bool)

		// verify burned pool coin
		poolCoinAfter := simapp.LiquidityKeeper.GetPoolCoinTotalSupply(ctx, pool)
		fmt.Println(poolCoinAfter, poolCoinBefore)
		require.True(t, poolCoinAfter.LT(poolCoinBefore))

		for i := 0; i < iterNum; i++ {
			withdrawerBalanceX := simapp.BankKeeper.GetBalance(ctx, addrs[i], pool.ReserveCoinDenoms[0])
			withdrawerBalanceY := simapp.BankKeeper.GetBalance(ctx, addrs[i], pool.ReserveCoinDenoms[1])
			require.True(t, withdrawerBalanceX.IsPositive())
			require.True(t, withdrawerBalanceY.IsPositive())

			withdrawMsgs := simapp.LiquidityKeeper.GetAllPoolBatchWithdrawMsgStates(ctx, batch)
			require.True(t, withdrawMsgs[i].Executed)
			require.True(t, withdrawMsgs[i].Succeeded)
			require.True(t, withdrawMsgs[i].ToBeDeleted)
		}
	}
}

func TestSwapPool(t *testing.T, simapp *LiquidityApp, ctx sdk.Context, offerCoins []sdk.Coin, orderPrices []sdk.Dec,
	addrs []sdk.AccAddress, poolId uint64, withEndblock bool) ([]*types.SwapMsgState, types.PoolBatch) {
	if len(offerCoins) != len(orderPrices) || len(orderPrices) != len(addrs) {
		require.True(t, false)
	}

	pool, found := simapp.LiquidityKeeper.GetPool(ctx, poolId)
	require.True(t, found)

	moduleAccAddress := simapp.AccountKeeper.GetModuleAddress(types.ModuleName)

	var swapMsgStates []*types.SwapMsgState

	params := simapp.LiquidityKeeper.GetParams(ctx)

	iterNum := len(addrs)
	for i := 0; i < iterNum; i++ {
		moduleAccEscrowAmtPool := simapp.BankKeeper.GetBalance(ctx, moduleAccAddress, offerCoins[i].Denom)
		currentBalance := simapp.BankKeeper.GetBalance(ctx, addrs[i], offerCoins[i].Denom)
		if currentBalance.IsLT(offerCoins[i]) {
			SaveAccountWithFee(simapp, ctx, addrs[i], sdk.NewCoins(offerCoins[i]), offerCoins[i])
			//SaveAccount(simapp, ctx, addrs[i], sdk.NewCoins(offerCoins[i]))
		}
		var demandCoinDenom string
		if pool.ReserveCoinDenoms[0] == offerCoins[i].Denom {
			demandCoinDenom = pool.ReserveCoinDenoms[1]
		} else if pool.ReserveCoinDenoms[1] == offerCoins[i].Denom {
			demandCoinDenom = pool.ReserveCoinDenoms[0]
		} else {
			require.True(t, false)
		}

		swapMsg := types.NewMsgSwapWithinBatch(addrs[i], poolId, types.DefaultSwapTypeId, offerCoins[i], demandCoinDenom, orderPrices[i], params.SwapFeeRate)
		batchPoolSwapMsg, err := simapp.LiquidityKeeper.SwapLiquidityPoolToBatch(ctx, swapMsg, 0)
		require.NoError(t, err)

		swapMsgStates = append(swapMsgStates, batchPoolSwapMsg)
		moduleAccEscrowAmtPoolAfter := simapp.BankKeeper.GetBalance(ctx, moduleAccAddress, offerCoins[i].Denom)
		moduleAccEscrowAmtPool.Amount = moduleAccEscrowAmtPool.Amount.Add(offerCoins[i].Amount).Add(types.GetOfferCoinFee(offerCoins[i], params.SwapFeeRate).Amount)
		require.Equal(t, moduleAccEscrowAmtPool, moduleAccEscrowAmtPoolAfter)

	}
	batch, bool := simapp.LiquidityKeeper.GetPoolBatch(ctx, poolId)

	if withEndblock {
		// endblock
		liquidity.EndBlocker(ctx, simapp.LiquidityKeeper)

		batch, bool = simapp.LiquidityKeeper.GetPoolBatch(ctx, poolId)
		require.True(t, bool)
	}
	return swapMsgStates, batch
}

func GetSwapMsg(t *testing.T, simapp *LiquidityApp, ctx sdk.Context, offerCoins []sdk.Coin, orderPrices []sdk.Dec,
	addrs []sdk.AccAddress, poolId uint64) []*types.MsgSwapWithinBatch {
	if len(offerCoins) != len(orderPrices) || len(orderPrices) != len(addrs) {
		require.True(t, false)
	}

	var msgs []*types.MsgSwapWithinBatch
	pool, found := simapp.LiquidityKeeper.GetPool(ctx, poolId)
	require.True(t, found)

	params := simapp.LiquidityKeeper.GetParams(ctx)

	iterNum := len(addrs)
	for i := 0; i < iterNum; i++ {
		currentBalance := simapp.BankKeeper.GetBalance(ctx, addrs[i], offerCoins[i].Denom)
		if currentBalance.IsLT(offerCoins[i]) {
			SaveAccountWithFee(simapp, ctx, addrs[i], sdk.NewCoins(offerCoins[i]), offerCoins[i])
			//SaveAccount(simapp, ctx, addrs[i], sdk.NewCoins(offerCoins[i]))
		}
		var demandCoinDenom string
		if pool.ReserveCoinDenoms[0] == offerCoins[i].Denom {
			demandCoinDenom = pool.ReserveCoinDenoms[1]
		} else if pool.ReserveCoinDenoms[1] == offerCoins[i].Denom {
			demandCoinDenom = pool.ReserveCoinDenoms[0]
		} else {
			require.True(t, false)
		}

		msgs = append(msgs, types.NewMsgSwapWithinBatch(addrs[i], poolId, types.DefaultSwapTypeId, offerCoins[i], demandCoinDenom, orderPrices[i], params.SwapFeeRate))
	}
	return msgs
}

// EmptyAppOptions is a stub implementing AppOptions
type EmptyAppOptions struct{}

// Get implements AppOptions
func (ao EmptyAppOptions) Get(o string) interface{} {
	return nil
}
