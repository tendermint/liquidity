package app

// DONTCOVER

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/tendermint/liquidity/x/liquidity"
	"math"
	"math/rand"
	"strconv"
	"testing"
	"time"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/liquidity/x/liquidity/keeper"
	"github.com/tendermint/liquidity/x/liquidity/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
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

// SetupWithGenesisValSet initializes a new LiquidityApp with a validator set and genesis accounts
// that also act as delegators. For simplicity, each validator is bonded with a delegation
// of one consensus engine unit (10^6) in the default token of the simapp from first genesis
// account. A Nop logger is set in LiquidityApp.
func SetupWithGenesisValSet(t *testing.T, valSet *tmtypes.ValidatorSet, genAccs []authtypes.GenesisAccount, balances ...banktypes.Balance) *LiquidityApp {
	db := dbm.NewMemDB()
	app := NewLiquidityApp(log.NewNopLogger(), db, nil, true, map[int64]bool{}, DefaultNodeHome, 5, MakeEncodingConfig(), EmptyAppOptions{})

	genesisState := NewDefaultGenesisState()

	// set genesis accounts
	authGenesis := authtypes.NewGenesisState(authtypes.DefaultParams(), genAccs)
	genesisState[authtypes.ModuleName] = app.AppCodec().MustMarshalJSON(authGenesis)

	validators := make([]stakingtypes.Validator, 0, len(valSet.Validators))
	delegations := make([]stakingtypes.Delegation, 0, len(valSet.Validators))

	bondAmt := sdk.NewInt(1000000)

	for _, val := range valSet.Validators {
		// Currently validator requires tmcrypto.ed25519 keys, which don't support
		// our Marshaling interfaces, so we need to pack them into our version of ed25519.
		// There is ongoing work to add secp256k1 keys (https://github.com/cosmos/cosmos-sdk/pull/7604).
		pk, err := cryptocodec.FromTmPubKeyInterface(val.PubKey)
		require.NoError(t, err)
		pkAny, err := codectypes.NewAnyWithValue(pk)
		require.NoError(t, err)
		validator := stakingtypes.Validator{
			OperatorAddress:   sdk.ValAddress(val.Address).String(),
			ConsensusPubkey:   pkAny,
			Jailed:            false,
			Status:            stakingtypes.Bonded,
			Tokens:            bondAmt,
			DelegatorShares:   sdk.OneDec(),
			Description:       stakingtypes.Description{},
			UnbondingHeight:   int64(0),
			UnbondingTime:     time.Unix(0, 0).UTC(),
			Commission:        stakingtypes.NewCommission(sdk.ZeroDec(), sdk.ZeroDec(), sdk.ZeroDec()),
			MinSelfDelegation: sdk.ZeroInt(),
		}
		validators = append(validators, validator)
		delegations = append(delegations, stakingtypes.NewDelegation(genAccs[0].GetAddress(), val.Address.Bytes(), sdk.OneDec()))

	}

	// set validators and delegations
	stakingGenesis := stakingtypes.NewGenesisState(stakingtypes.DefaultParams(), validators, delegations)
	genesisState[stakingtypes.ModuleName] = app.AppCodec().MustMarshalJSON(stakingGenesis)

	totalSupply := sdk.NewCoins()
	for _, b := range balances {
		// add genesis acc tokens and delegated tokens to total supply
		totalSupply = totalSupply.Add(b.Coins.Add(sdk.NewCoin(sdk.DefaultBondDenom, bondAmt))...)
	}

	// update total supply
	bankGenesis := banktypes.NewGenesisState(banktypes.DefaultGenesisState().Params, balances, totalSupply, []banktypes.Metadata{})
	genesisState[banktypes.ModuleName] = app.AppCodec().MustMarshalJSON(bankGenesis)

	stateBytes, err := json.MarshalIndent(genesisState, "", " ")
	require.NoError(t, err)

	// init chain will set the validator set and initialize the genesis accounts
	app.InitChain(
		abci.RequestInitChain{
			Validators:      []abci.ValidatorUpdate{},
			ConsensusParams: DefaultConsensusParams,
			AppStateBytes:   stateBytes,
		},
	)

	// commit genesis changes
	app.Commit()
	app.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{
		Height:             app.LastBlockHeight() + 1,
		AppHash:            app.LastCommitID().Hash,
		ValidatorsHash:     valSet.Hash(),
		NextValidatorsHash: valSet.Hash(),
	}})

	return app
}

// SetupWithGenesisAccounts initializes a new LiquidityApp with the provided genesis
// accounts and possible balances.
func SetupWithGenesisAccounts(genAccs []authtypes.GenesisAccount, balances ...banktypes.Balance) *LiquidityApp {
	db := dbm.NewMemDB()
	app := NewLiquidityApp(log.NewNopLogger(), db, nil, true, map[int64]bool{}, DefaultNodeHome, 0, MakeEncodingConfig(), EmptyAppOptions{})

	// initialize the chain with the passed in genesis accounts
	genesisState := NewDefaultGenesisState()

	authGenesis := authtypes.NewGenesisState(authtypes.DefaultParams(), genAccs)
	genesisState[authtypes.ModuleName] = app.AppCodec().MustMarshalJSON(authGenesis)

	totalSupply := sdk.NewCoins()
	for _, b := range balances {
		totalSupply = totalSupply.Add(b.Coins...)
	}

	bankGenesis := banktypes.NewGenesisState(banktypes.DefaultGenesisState().Params, balances, totalSupply, []banktypes.Metadata{})
	genesisState[banktypes.ModuleName] = app.AppCodec().MustMarshalJSON(bankGenesis)

	stateBytes, err := json.MarshalIndent(genesisState, "", " ")
	if err != nil {
		panic(err)
	}

	app.InitChain(
		abci.RequestInitChain{
			Validators:      []abci.ValidatorUpdate{},
			ConsensusParams: DefaultConsensusParams,
			AppStateBytes:   stateBytes,
		},
	)

	app.Commit()
	app.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{Height: app.LastBlockHeight() + 1}})

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

// AddTestAddrsFromPubKeys adds the addresses into the LiquidityApp providing only the public keys.
func AddTestAddrsFromPubKeys(app *LiquidityApp, ctx sdk.Context, pubKeys []cryptotypes.PubKey, accAmt sdk.Int) {
	initCoins := sdk.NewCoins(sdk.NewCoin(app.StakingKeeper.BondDenom(ctx), accAmt))

	setTotalSupply(app, ctx, accAmt, len(pubKeys))

	// fill all the addresses with some coins, set the loose pool tokens simultaneously
	for _, pubKey := range pubKeys {
		SaveAccount(app, ctx, sdk.AccAddress(pubKey.Address()), initCoins)
	}
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

// ConvertAddrsToValAddrs converts the provided addresses to ValAddress.
func ConvertAddrsToValAddrs(addrs []sdk.AccAddress) []sdk.ValAddress {
	valAddrs := make([]sdk.ValAddress, len(addrs))

	for i, addr := range addrs {
		valAddrs[i] = sdk.ValAddress(addr)
	}

	return valAddrs
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

// CheckBalance checks the balance of an account.
func CheckBalance(t *testing.T, app *LiquidityApp, addr sdk.AccAddress, balances sdk.Coins) {
	ctxCheck := app.BaseApp.NewContext(true, tmproto.Header{})
	require.True(t, balances.IsEqual(app.BankKeeper.GetAllBalances(ctxCheck, addr)))
}

func incrementAllSequenceNumbers(initSeqNums []uint64) {
	for i := 0; i < len(initSeqNums); i++ {
		initSeqNums[i]++
	}
}

// CreateTestPubKeys returns a total of numPubKeys public keys in ascending order.
func CreateTestPubKeys(numPubKeys int) []cryptotypes.PubKey {
	var publicKeys []cryptotypes.PubKey
	var buffer bytes.Buffer

	// start at 10 to avoid changing 1 to 01, 2 to 02, etc
	for i := 100; i < (numPubKeys + 100); i++ {
		numString := strconv.Itoa(i)
		buffer.WriteString("0B485CFC0EECC619440448436F8FC9DF40566F2369E72400281454CB552AF") // base pubkey string
		buffer.WriteString(numString)                                                       // adding on final two digits to make pubkeys unique
		publicKeys = append(publicKeys, NewPubKeyFromHex(buffer.String()))
		buffer.Reset()
	}

	return publicKeys
}

// NewPubKeyFromHex returns a PubKey from a hex string.
func NewPubKeyFromHex(pk string) (res cryptotypes.PubKey) {
	pkBytes, err := hex.DecodeString(pk)
	if err != nil {
		panic(err)
	}
	if len(pkBytes) != ed25519.PubKeySize {
		panic(errors.Wrap(errors.ErrInvalidPubKey, "invalid pubkey size"))
	}
	return &ed25519.PubKey{Key: pkBytes}
}

// CreateTestInput Returns a simapp with custom LiquidityKeeper
// to avoid messing with the hooks.
func CreateTestInput() (*LiquidityApp, sdk.Context) {
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

func GetRandPoolAmtLegacy(r *rand.Rand) (X, Y sdk.Int) {
	X = sdk.NewInt(int64(r.Float32() * 1000000000000))
	Y = sdk.NewInt(int64(r.Float32() * 1000000000000))
	return
}

func GetRandFloats(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func GetRandRange(r *rand.Rand, min, max int) sdk.Int {
	return sdk.NewInt(int64(r.Intn(max-min) + min))
}

func GetRandomSizeOrders(denomX, denomY string, X, Y sdk.Int, r *rand.Rand, sizeXtoY, sizeYtoX int32) (XtoY, YtoX []*types.MsgSwap) {
	randomSizeXtoY := int(r.Int31n(sizeXtoY))
	randomSizeYtoX := int(r.Int31n(sizeYtoX))
	return GetRandomOrders(denomX, denomY, X, Y, r, randomSizeXtoY, randomSizeYtoX)
}

func GetRandomOrders(denomX, denomY string, X, Y sdk.Int, r *rand.Rand, sizeXtoY, sizeYtoX int) (XtoY, YtoX []*types.MsgSwap) {
	currentPrice := X.ToDec().Quo(Y.ToDec())

	for i := 0; i < sizeXtoY; i++ {
		orderPrice := currentPrice.Mul(sdk.NewDecFromIntWithPrec(GetRandRange(r, 991, 1009), 3))
		//offerAmt := X.ToDec().Mul(sdk.NewDecFromIntWithPrec(GetRandRange(r, 1, 100), 4))
		orderAmt := sdk.ZeroDec()
		if r.Intn(2) == 1 {
			orderAmt = X.ToDec().Mul(sdk.NewDecFromIntWithPrec(GetRandRange(r, 1, 100), 4))
		} else {
			orderAmt = sdk.NewDecFromIntWithPrec(GetRandRange(r, 1000, 10000), 0)
		}
		orderCoin := sdk.NewCoin(denomX, orderAmt.RoundInt())

		XtoY = append(XtoY, &types.MsgSwap{
			OfferCoin:       orderCoin,
			DemandCoinDenom: denomY,
			OrderPrice:      orderPrice,
		})
	}

	for i := 0; i < sizeYtoX; i++ {
		orderPrice := currentPrice.Mul(sdk.NewDecFromIntWithPrec(GetRandRange(r, 991, 1009), 3))
		//offerAmt := Y.ToDec().Mul(sdk.NewDecFromIntWithPrec(GetRandRange(r, 1, 100), 4))
		orderAmt := sdk.ZeroDec()
		if r.Intn(2) == 1 {
			orderAmt = Y.ToDec().Mul(sdk.NewDecFromIntWithPrec(GetRandRange(r, 1, 100), 4))
		} else {
			orderAmt = sdk.NewDecFromIntWithPrec(GetRandRange(r, 1000, 10000), 0)
		}
		orderCoin := sdk.NewCoin(denomY, orderAmt.RoundInt())

		YtoX = append(YtoX, &types.MsgSwap{
			OfferCoin:       orderCoin,
			DemandCoinDenom: denomX,
			OrderPrice:      orderPrice,
		})
	}
	return
}

func GetRandomBatchSwapOrders(denomX, denomY string, X, Y sdk.Int, r *rand.Rand) (XtoY, YtoX []*types.SwapMsgState) {
	currentPrice := X.ToDec().Quo(Y.ToDec())

	XtoYnewSize := int(r.Int31n(20)) // 0~19
	YtoXnewSize := int(r.Int31n(20)) // 0~19

	for i := 0; i < XtoYnewSize; i++ {
		GetRandFloats(0.1, 0.9)
		orderPrice := currentPrice.Mul(sdk.NewDecFromIntWithPrec(GetRandRange(r, 991, 1009), 3))
		offerAmt := X.ToDec().Mul(sdk.NewDecFromIntWithPrec(GetRandRange(r, 1, 100), 4))
		orderCoin := sdk.NewCoin(denomX, offerAmt.RoundInt())

		XtoY = append(XtoY, &types.SwapMsgState{
			Msg: &types.MsgSwap{
				OfferCoin:       orderCoin,
				DemandCoinDenom: denomY,
				OrderPrice:      orderPrice,
				OfferCoinFee:    types.GetOfferCoinFee(orderCoin, types.DefaultSwapFeeRate),
			},
		})
	}

	for i := 0; i < YtoXnewSize; i++ {
		orderPrice := currentPrice.Mul(sdk.NewDecFromIntWithPrec(GetRandRange(r, 991, 1009), 3))
		offerAmt := Y.ToDec().Mul(sdk.NewDecFromIntWithPrec(GetRandRange(r, 1, 100), 4))
		orderCoin := sdk.NewCoin(denomY, offerAmt.RoundInt())

		YtoX = append(YtoX, &types.SwapMsgState{
			Msg: &types.MsgSwap{
				OfferCoin:       orderCoin,
				DemandCoinDenom: denomX,
				OrderPrice:      orderPrice,
				OfferCoinFee:    types.GetOfferCoinFee(orderCoin, types.DefaultSwapFeeRate),
			},
		})
	}
	return
}

func TestCreatePool(t *testing.T, simapp *LiquidityApp, ctx sdk.Context, X, Y sdk.Int, denomX, denomY string, addr sdk.AccAddress) uint64 {
	deposit := sdk.NewCoins(sdk.NewCoin(denomX, X), sdk.NewCoin(denomY, Y))
	params := simapp.LiquidityKeeper.GetParams(ctx)
	// set accounts for creator, depositor, withdrawer, balance for deposit
	SaveAccount(simapp, ctx, addr, deposit.Add(params.LiquidityPoolCreationFee...)) // pool creator
	depositX := simapp.BankKeeper.GetBalance(ctx, addr, denomX)
	depositY := simapp.BankKeeper.GetBalance(ctx, addr, denomY)
	depositBalance := sdk.NewCoins(depositX, depositY)
	require.Equal(t, deposit, depositBalance)

	// create Liquidity pool
	poolTypeIndex := types.DefaultPoolTypeIndex
	poolId := simapp.LiquidityKeeper.GetNextPoolId(ctx)
	msg := types.NewMsgCreateLiquidityPool(addr, poolTypeIndex, depositBalance)
	_, err := simapp.LiquidityKeeper.CreatePool(ctx, msg)
	require.NoError(t, err)

	// verify created liquidity pool
	pool, found := simapp.LiquidityKeeper.GetPool(ctx, poolId)
	require.True(t, found)
	require.Equal(t, poolId, pool.PoolId)
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

		depositMsg := types.NewMsgDepositToLiquidityPool(addrs[i], poolId, deposit)
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
		fmt.Println(balancePoolCoin, poolCoinAmt)
		require.True(t, balancePoolCoin.Amount.GTE(poolCoinAmt))

		withdrawCoin := sdk.NewCoin(pool.PoolCoinDenom, poolCoinAmt)
		withdrawMsg := types.NewMsgWithdrawFromLiquidityPool(addrs[i], poolId, withdrawCoin)
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

func TestSwapPool(t *testing.T, simapp *LiquidityApp, ctx sdk.Context, offerCoinList []sdk.Coin, orderPrices []sdk.Dec,
	addrs []sdk.AccAddress, poolId uint64, withEndblock bool) ([]*types.SwapMsgState, types.PoolBatch) {
	if len(offerCoinList) != len(orderPrices) || len(orderPrices) != len(addrs) {
		require.True(t, false)
	}

	pool, found := simapp.LiquidityKeeper.GetPool(ctx, poolId)
	require.True(t, found)

	moduleAccAddress := simapp.AccountKeeper.GetModuleAddress(types.ModuleName)

	var batchPoolSwapMsgList []*types.SwapMsgState

	params := simapp.LiquidityKeeper.GetParams(ctx)

	iterNum := len(addrs)
	for i := 0; i < iterNum; i++ {
		moduleAccEscrowAmtPool := simapp.BankKeeper.GetBalance(ctx, moduleAccAddress, offerCoinList[i].Denom)
		currentBalance := simapp.BankKeeper.GetBalance(ctx, addrs[i], offerCoinList[i].Denom)
		if currentBalance.IsLT(offerCoinList[i]) {
			SaveAccountWithFee(simapp, ctx, addrs[i], sdk.NewCoins(offerCoinList[i]), offerCoinList[i])
			//SaveAccount(simapp, ctx, addrs[i], sdk.NewCoins(offerCoinList[i]))
		}
		var demandCoinDenom string
		if pool.ReserveCoinDenoms[0] == offerCoinList[i].Denom {
			demandCoinDenom = pool.ReserveCoinDenoms[1]
		} else if pool.ReserveCoinDenoms[1] == offerCoinList[i].Denom {
			demandCoinDenom = pool.ReserveCoinDenoms[0]
		} else {
			require.True(t, false)
		}

		swapMsg := types.NewMsgSwap(addrs[i], poolId, types.DefaultSwapType, offerCoinList[i], demandCoinDenom, orderPrices[i], params.SwapFeeRate)
		batchPoolSwapMsg, err := simapp.LiquidityKeeper.SwapLiquidityPoolToBatch(ctx, swapMsg, 0)
		require.NoError(t, err)

		batchPoolSwapMsgList = append(batchPoolSwapMsgList, batchPoolSwapMsg)
		moduleAccEscrowAmtPoolAfter := simapp.BankKeeper.GetBalance(ctx, moduleAccAddress, offerCoinList[i].Denom)
		moduleAccEscrowAmtPool.Amount = moduleAccEscrowAmtPool.Amount.Add(offerCoinList[i].Amount).Add(types.GetOfferCoinFee(offerCoinList[i], params.SwapFeeRate).Amount)
		require.Equal(t, moduleAccEscrowAmtPool, moduleAccEscrowAmtPoolAfter)

	}
	batch, bool := simapp.LiquidityKeeper.GetPoolBatch(ctx, poolId)

	if withEndblock {
		// endblock
		liquidity.EndBlocker(ctx, simapp.LiquidityKeeper)

		batch, bool = simapp.LiquidityKeeper.GetPoolBatch(ctx, poolId)
		require.True(t, bool)
	}
	return batchPoolSwapMsgList, batch
}

func GetSwapMsg(t *testing.T, simapp *LiquidityApp, ctx sdk.Context, offerCoinList []sdk.Coin, orderPrices []sdk.Dec,
	addrs []sdk.AccAddress, poolId uint64) []*types.MsgSwap {
	if len(offerCoinList) != len(orderPrices) || len(orderPrices) != len(addrs) {
		require.True(t, false)
	}

	var msgList []*types.MsgSwap
	pool, found := simapp.LiquidityKeeper.GetPool(ctx, poolId)
	require.True(t, found)

	params := simapp.LiquidityKeeper.GetParams(ctx)

	iterNum := len(addrs)
	for i := 0; i < iterNum; i++ {
		currentBalance := simapp.BankKeeper.GetBalance(ctx, addrs[i], offerCoinList[i].Denom)
		if currentBalance.IsLT(offerCoinList[i]) {
			SaveAccountWithFee(simapp, ctx, addrs[i], sdk.NewCoins(offerCoinList[i]), offerCoinList[i])
			//SaveAccount(simapp, ctx, addrs[i], sdk.NewCoins(offerCoinList[i]))
		}
		var demandCoinDenom string
		if pool.ReserveCoinDenoms[0] == offerCoinList[i].Denom {
			demandCoinDenom = pool.ReserveCoinDenoms[1]
		} else if pool.ReserveCoinDenoms[1] == offerCoinList[i].Denom {
			demandCoinDenom = pool.ReserveCoinDenoms[0]
		} else {
			require.True(t, false)
		}

		msgList = append(msgList, types.NewMsgSwap(addrs[i], poolId, types.DefaultSwapType, offerCoinList[i], demandCoinDenom, orderPrices[i], params.SwapFeeRate))
	}
	return msgList
}

// EmptyAppOptions is a stub implementing AppOptions
type EmptyAppOptions struct{}

// Get implements AppOptions
func (ao EmptyAppOptions) Get(o string) interface{} {
	return nil
}
