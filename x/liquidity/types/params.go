package types

import (
	"fmt"

	"gopkg.in/yaml.v2"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Const value of liquidity module
const (
	CancelOrderLifeSpan int64 = 0

	// min number of reserveCoins for PoolType only 2 is allowed on this spec
	MinReserveCoinNum uint32 = 2

	// max number of reserveCoins for PoolType only 2 is allowed on this spec
	MaxReserveCoinNum uint32 = 2

	// Number of blocks in one batch
	DefaultUnitBatchHeight uint32 = 1

	// index of target pool type, only 1 is allowed on this version.
	DefaultPoolTypeId uint32 = 1

	// swap type index of available swap request, only 1 (InstantSwap) is allowed on this version.
	DefaultSwapTypeId uint32 = 1
)

// Parameter store keys
var (
	KeyPoolTypes                = []byte("PoolTypes")
	KeyMinInitDepositAmount     = []byte("MinInitDepositAmount")
	KeyInitPoolCoinMintAmount   = []byte("InitPoolCoinMintAmount")
	KeyMaxReserveCoinAmount   = []byte("MaxReserveCoinAmount")
	KeySwapFeeRate              = []byte("SwapFeeRate")
	KeyPoolCreationFee = []byte("PoolCreationFee")
	KeyUnitBatchHeight            = []byte("UnitBatchHeight")
	KeyWithdrawFeeRate          = []byte("WithdrawFeeRate")
	KeyMaxOrderAmountRatio      = []byte("MaxOrderAmountRatio")

	DefaultMinInitDepositAmount     = sdk.NewInt(1000000)
	DefaultInitPoolCoinMintAmount   = sdk.NewInt(1000000)
	DefaultMaxReserveCoinAmount   = sdk.ZeroInt()
	DefaultSwapFeeRate              = sdk.NewDecWithPrec(3, 3) // "0.003000000000000000"
	DefaultWithdrawFeeRate          = sdk.NewDecWithPrec(3, 3) // "0.003000000000000000"
	DefaultMaxOrderAmountRatio      = sdk.NewDecWithPrec(1, 1) // "0.100000000000000000"
	DefaultPoolCreationFee = sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100000000)))
	MinOfferCoinAmount              = sdk.NewInt(100)

	HalfRatio = sdk.MustNewDecFromStr("0.5")

	DecimalErrThreshold3  = sdk.NewDecWithPrec(1, 3)
	DecimalErrThreshold10 = sdk.NewDecWithPrec(1, 10)

	DefaultPoolType = PoolType{
		Id:                1,
		Name:              "DefaultPoolType",
		MinReserveCoinNum: MinReserveCoinNum,
		MaxReserveCoinNum: MaxReserveCoinNum,
	}
)

// NewParams liquidity paramtypes constructor
func NewParams(poolTypes []PoolType, minInitDeposit, initPoolCoinMint, reserveCoinLimit sdk.Int, creationFee sdk.Coins,
	swapFeeRate, withdrawFeeRate, maxOrderAmtRatio sdk.Dec, unitBatchHeight uint32) Params {
	return Params{
		PoolTypes:                poolTypes,
		MinInitDepositAmount:     minInitDeposit,
		InitPoolCoinMintAmount:   initPoolCoinMint,
		MaxReserveCoinAmount:   reserveCoinLimit,
		PoolCreationFee: creationFee,
		SwapFeeRate:              swapFeeRate,
		WithdrawFeeRate:          withdrawFeeRate,
		MaxOrderAmountRatio:      maxOrderAmtRatio,
		UnitBatchHeight:            unitBatchHeight,
	}
}

// ParamTypeTable returns the TypeTable for liquidity module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// KeyValuePairs implements paramtypes.KeyValuePairs
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {

	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyPoolTypes, &p.PoolTypes, validatePoolTypes),
		paramtypes.NewParamSetPair(KeyMinInitDepositAmount, &p.MinInitDepositAmount, validateMinInitDepositAmount),
		paramtypes.NewParamSetPair(KeyInitPoolCoinMintAmount, &p.InitPoolCoinMintAmount, validateInitPoolCoinMintAmount),
		paramtypes.NewParamSetPair(KeyMaxReserveCoinAmount, &p.MaxReserveCoinAmount, validateMaxReserveCoinAmount),
		paramtypes.NewParamSetPair(KeyPoolCreationFee, &p.PoolCreationFee, validatePoolCreationFee),
		paramtypes.NewParamSetPair(KeySwapFeeRate, &p.SwapFeeRate, validateSwapFeeRate),
		paramtypes.NewParamSetPair(KeyWithdrawFeeRate, &p.WithdrawFeeRate, validateWithdrawFeeRate),
		paramtypes.NewParamSetPair(KeyMaxOrderAmountRatio, &p.MaxOrderAmountRatio, validateMaxOrderAmountRatio),
		paramtypes.NewParamSetPair(KeyUnitBatchHeight, &p.UnitBatchHeight, validateUnitBatchHeight),
	}
}

// DefaultParams returns the default liquidity module parameters
func DefaultParams() Params {
	var defaultPoolTypes []PoolType
	defaultPoolTypes = append(defaultPoolTypes, DefaultPoolType)

	return NewParams(
		defaultPoolTypes,
		DefaultMinInitDepositAmount,
		DefaultInitPoolCoinMintAmount,
		DefaultMaxReserveCoinAmount,
		DefaultPoolCreationFee,
		DefaultSwapFeeRate,
		DefaultWithdrawFeeRate,
		DefaultMaxOrderAmountRatio,
		DefaultUnitBatchHeight)
}

// String returns a human readable string representation of the parameters.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// Validate returns err if Params is invalid
func (p Params) Validate() error {
	if err := validatePoolTypes(p.PoolTypes); err != nil {
		return err
	}

	if err := validateMinInitDepositAmount(p.MinInitDepositAmount); err != nil {
		return err
	}

	if err := validateInitPoolCoinMintAmount(p.InitPoolCoinMintAmount); err != nil {
		return err
	}

	if err := validateMaxReserveCoinAmount(p.MaxReserveCoinAmount); err != nil {
		return err
	}

	if err := validatePoolCreationFee(p.PoolCreationFee); err != nil {
		return err
	}

	if err := validateSwapFeeRate(p.SwapFeeRate); err != nil {
		return err
	}

	if err := validateWithdrawFeeRate(p.WithdrawFeeRate); err != nil {
		return err
	}

	if err := validateMaxOrderAmountRatio(p.MaxOrderAmountRatio); err != nil {
		return err
	}

	if err := validateUnitBatchHeight(p.UnitBatchHeight); err != nil {
		return err
	}
	return nil
}

// check validity of the list of liquidity pool type
func validatePoolTypes(i interface{}) error {
	v, ok := i.([]PoolType)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v == nil {
		return fmt.Errorf("empty parameter: PoolTypes")
	}
	for i, p := range v {
		if i+1 != int(p.Id) {
			return fmt.Errorf("PoolTypes index must be sorted")
		}
	}
	if len(v) > 1 {
		return fmt.Errorf("only default pool type allowed on this version")
	}
	if len(v) < 1 {
		return fmt.Errorf("need to default pool type")
	}
	if !v[0].Equal(DefaultPoolType) {
		return fmt.Errorf("only default pool type allowed")
	}
	return nil
}

// Validate that the minimum deposit.
func validateMinInitDepositAmount(i interface{}) error {
	v, ok := i.(sdk.Int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if !v.IsPositive() {
		return fmt.Errorf("MinInitDepositAmount must be positive: %s", v)
	}

	return nil
}

// Validate that the minimum deposit for initiating pool.
func validateInitPoolCoinMintAmount(i interface{}) error {
	v, ok := i.(sdk.Int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if !v.IsPositive() {
		return fmt.Errorf("InitPoolCoinMintAmount must be positive: %s", v)
	}
	if v.LT(DefaultInitPoolCoinMintAmount) {
		return fmt.Errorf("InitPoolCoinMintAmount should over default value: %s", v)
	}
	return nil
}

// Validate that the Limit the size of each liquidity pool.
func validateMaxReserveCoinAmount(i interface{}) error {
	v, ok := i.(sdk.Int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v.IsNegative() {
		return fmt.Errorf("InitPoolCoinMintAmount must be positive or zero: %s", v)
	}
	return nil
}

// Check if the swap fee rate is between 0 and 1
func validateSwapFeeRate(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v.IsNegative() {
		return fmt.Errorf("SwapFeeRate cannot be negative: %s", v)
	}
	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("SwapFeeRate too large: %s", v)
	}
	return nil
}

// Check if the withdraw fee rate is between 0 and 1
func validateWithdrawFeeRate(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v.IsNegative() {
		return fmt.Errorf("WithdrawFeeRate cannot be negative: %s", v)
	}
	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("WithdrawFeeRate too large: %s", v)
	}
	return nil
}

// Check if the max order amount ratio is between 0 and 1
func validateMaxOrderAmountRatio(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if v.IsNegative() {
		return fmt.Errorf("WithdrawFeeRate cannot be negative: %s", v)
	}
	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("WithdrawFeeRate too large: %s", v)
	}
	return nil
}

// Check if the pool creation fee is valid
func validatePoolCreationFee(i interface{}) error {
	coins, ok := i.(sdk.Coins)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if err := coins.Validate(); err != nil {
		return err
	}
	if coins.Empty() {
		return fmt.Errorf("PoolCreationFee cannot be Empty: %s", coins)
	}
	return nil
}

// Check if the liquidity Msg fee is valid
func validateUnitBatchHeight(i interface{}) error {
	int, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if int == 0 {
		return fmt.Errorf("UnitBatchHeight cannot be zero")
	}
	return nil
}
