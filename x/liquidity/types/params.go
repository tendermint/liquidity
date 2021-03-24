package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"gopkg.in/yaml.v2"
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
	KeyPoolTypes              = []byte("PoolTypes")
	KeyMinInitDepositAmount   = []byte("MinInitDepositAmount")
	KeyInitPoolCoinMintAmount = []byte("InitPoolCoinMintAmount")
	KeyMaxReserveCoinAmount   = []byte("MaxReserveCoinAmount")
	KeySwapFeeRate            = []byte("SwapFeeRate")
	KeyPoolCreationFee        = []byte("PoolCreationFee")
	KeyUnitBatchHeight        = []byte("UnitBatchHeight")
	KeyWithdrawFeeRate        = []byte("WithdrawFeeRate")
	KeyMaxOrderAmountRatio    = []byte("MaxOrderAmountRatio")
)

var (
	DefaultMinInitDepositAmount   = sdk.NewInt(1000000)
	DefaultInitPoolCoinMintAmount = sdk.NewInt(1000000)
	DefaultMaxReserveCoinAmount   = sdk.ZeroInt()
	DefaultSwapFeeRate            = sdk.NewDecWithPrec(3, 3) // "0.003000000000000000"
	DefaultWithdrawFeeRate        = sdk.NewDecWithPrec(3, 3) // "0.003000000000000000"
	DefaultMaxOrderAmountRatio    = sdk.NewDecWithPrec(1, 1) // "0.100000000000000000"
	DefaultPoolCreationFee        = sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(100000000)))
	DefaultPoolType               = PoolType{
		Id:                1,
		Name:              "DefaultPoolType",
		MinReserveCoinNum: MinReserveCoinNum,
		MaxReserveCoinNum: MaxReserveCoinNum,
	}
	DefaultPoolTypes = []PoolType{DefaultPoolType}

	MinOfferCoinAmount = sdk.NewInt(100) // TODO: move into parameters
)

var _ paramstypes.ParamSet = (*Params)(nil)

// ParamKeyTable returns the parameter key table.
func ParamKeyTable() paramstypes.KeyTable {
	return paramstypes.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams returns the default liquidity module parameters.
func DefaultParams() Params {
	return Params{
		PoolTypes:              DefaultPoolTypes,
		MinInitDepositAmount:   DefaultMinInitDepositAmount,
		InitPoolCoinMintAmount: DefaultInitPoolCoinMintAmount,
		MaxReserveCoinAmount:   DefaultMaxReserveCoinAmount,
		PoolCreationFee:        DefaultPoolCreationFee,
		SwapFeeRate:            DefaultSwapFeeRate,
		WithdrawFeeRate:        DefaultWithdrawFeeRate,
		MaxOrderAmountRatio:    DefaultMaxOrderAmountRatio,
		UnitBatchHeight:        DefaultUnitBatchHeight,
	}
}

// ParamSetPairs implements paramstypes.ParamSet.
func (p *Params) ParamSetPairs() paramstypes.ParamSetPairs {
	return paramstypes.ParamSetPairs{
		paramstypes.NewParamSetPair(KeyPoolTypes, &p.PoolTypes, validatePoolTypes),
		paramstypes.NewParamSetPair(KeyMinInitDepositAmount, &p.MinInitDepositAmount, validateMinInitDepositAmount),
		paramstypes.NewParamSetPair(KeyInitPoolCoinMintAmount, &p.InitPoolCoinMintAmount, validateInitPoolCoinMintAmount),
		paramstypes.NewParamSetPair(KeyMaxReserveCoinAmount, &p.MaxReserveCoinAmount, validateMaxReserveCoinAmount),
		paramstypes.NewParamSetPair(KeyPoolCreationFee, &p.PoolCreationFee, validatePoolCreationFee),
		paramstypes.NewParamSetPair(KeySwapFeeRate, &p.SwapFeeRate, validateSwapFeeRate),
		paramstypes.NewParamSetPair(KeyWithdrawFeeRate, &p.WithdrawFeeRate, validateWithdrawFeeRate),
		paramstypes.NewParamSetPair(KeyMaxOrderAmountRatio, &p.MaxOrderAmountRatio, validateMaxOrderAmountRatio),
		paramstypes.NewParamSetPair(KeyUnitBatchHeight, &p.UnitBatchHeight, validateUnitBatchHeight),
	}
}

// String returns a human readable string representation of the parameters.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// Validate validates parameters.
func (p Params) Validate() error {
	for _, v := range []struct {
		value     interface{}
		validator func(interface{}) error
	}{
		{p.PoolTypes, validatePoolTypes},
		{p.MinInitDepositAmount, validateMinInitDepositAmount},
		{p.InitPoolCoinMintAmount, validateInitPoolCoinMintAmount},
		{p.MaxReserveCoinAmount, validateMaxReserveCoinAmount},
		{p.PoolCreationFee, validatePoolCreationFee},
		{p.SwapFeeRate, validateSwapFeeRate},
		{p.WithdrawFeeRate, validateWithdrawFeeRate},
		{p.MaxOrderAmountRatio, validateMaxOrderAmountRatio},
		{p.UnitBatchHeight, validateUnitBatchHeight},
	} {
		if err := v.validator(v.value); err != nil {
			return err
		}
	}
	return nil
}

func validatePoolTypes(i interface{}) error {
	v, ok := i.([]PoolType)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if len(v) == 0 {
		return fmt.Errorf("pool types must be not empty")
	}
	for i, p := range v {
		if int(p.Id) != i+1 {
			return fmt.Errorf("pool type ids must be sorted")
		}
	}
	if len(v) > 1 || !v[0].Equal(DefaultPoolType) {
		return fmt.Errorf("only default pool type is allowed in this version of liquidity module")
	}

	return nil
}

func validateMinInitDepositAmount(i interface{}) error {
	v, ok := i.(sdk.Int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("minimum initial deposit amount must be not nil")
	}
	if !v.IsPositive() {
		return fmt.Errorf("minimum initial deposit amount must be positive: %s", v)
	}

	return nil
}

func validateInitPoolCoinMintAmount(i interface{}) error {
	v, ok := i.(sdk.Int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("initial pool coin mint amount must be not nil")
	}
	if !v.IsPositive() {
		return fmt.Errorf("initial pool coin mint amount must be positive: %s", v)
	}
	if v.LT(DefaultInitPoolCoinMintAmount) {
		return fmt.Errorf("initial pool coin mint amount must be greater or equal than 1000000: %s", v)
	}

	return nil
}

func validateMaxReserveCoinAmount(i interface{}) error {
	v, ok := i.(sdk.Int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("max reserve coin amount must be not nil")
	}
	if v.IsNegative() {
		return fmt.Errorf("max reserve coin amount must be not negative: %s", v)
	}

	return nil
}

func validateSwapFeeRate(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("swap fee rate must be not nil")
	}
	if v.IsNegative() {
		return fmt.Errorf("swap fee rate must be not negative: %s", v)
	}
	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("swap fee rate too large: %s", v)
	}

	return nil
}

func validateWithdrawFeeRate(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("withdraw fee rate must be not nil")
	}
	if v.IsNegative() {
		return fmt.Errorf("withdraw fee rate must be not negative: %s", v)
	}
	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("withdraw fee rate too large: %s", v)
	}

	return nil
}

func validateMaxOrderAmountRatio(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("max order amount ratio must be not nil")
	}
	if v.IsNegative() {
		return fmt.Errorf("max order amount ratio must be not negative: %s", v)
	}
	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("max order amount ratio too large: %s", v)
	}

	return nil
}

func validatePoolCreationFee(i interface{}) error {
	v, ok := i.(sdk.Coins)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if err := v.Validate(); err != nil {
		return err
	}
	if v.Empty() {
		return fmt.Errorf("pool creation fee must be not empty")
	}

	return nil
}

func validateUnitBatchHeight(i interface{}) error {
	v, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("unit batch height must be positive: %d", v)
	}

	return nil
}
