package types

import (
	"fmt"
	"gopkg.in/yaml.v2"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

const DefaultBatchSize uint32 = 1

// Parameter store keys
var (
	KeyLiquidityPoolTypes      = []byte("LiquidityPoolTypes")
	KeyMinInitDepositToPool    = []byte("MinInitDepositToPool")
	KeyInitPoolTokenMintAmount = []byte("InitPoolTokenMintAmount")
	KeySwapFeeRate             = []byte("SwapFeeRate")
	KeyLiquidityPoolFeeRate    = []byte("LiquidityPoolFeeRate")

	LiquidityPoolTypeConstantProduct = LiquidityPoolType{
		PoolTypeIndex:         0,
		NumOfReserveTokens:    2,
		SwapPriceFunctionName: ConstantProductFunctionName,
		Description:           "Default Constant Product Liquidity Pool",
	}
)

type ParamsLegacy struct {
	LiquidityPoolTypes      []LiquidityPoolType
	MinInitDepositToPool    sdk.Int
	InitPoolTokenMintAmount sdk.Int
	SwapFeeRate             sdk.Dec
	LiquidityPoolFeeRate    sdk.Dec
}

// NewParams liquidity paramtypes constructor
func NewParams(liquidityPoolTypes []LiquidityPoolType, minInitDeposit, initPoolTokenMint sdk.Int, swapFeeRate, poolFeeRate sdk.Dec) Params {
	return Params{
		LiquidityPoolTypes:      liquidityPoolTypes,
		MinInitDepositToPool:    minInitDeposit,
		InitPoolTokenMintAmount: initPoolTokenMint,
		SwapFeeRate:             swapFeeRate,
		LiquidityPoolFeeRate:    poolFeeRate,
	}
}

// ParamTypeTable returns the TypeTable for liquidity module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// KeyValuePairs implements paramtypes.KeyValuePairs
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {

	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyLiquidityPoolTypes, &p.LiquidityPoolTypes, validateLiquidityPoolTypes),
		paramtypes.NewParamSetPair(KeyMinInitDepositToPool, &p.MinInitDepositToPool, validateMinInitDepositToPool),
		paramtypes.NewParamSetPair(KeyInitPoolTokenMintAmount, &p.InitPoolTokenMintAmount, validateInitPoolTokenMintAmount),
		paramtypes.NewParamSetPair(KeySwapFeeRate, &p.SwapFeeRate, validateSwapFeeRate),
		paramtypes.NewParamSetPair(KeyLiquidityPoolFeeRate, &p.LiquidityPoolFeeRate, validateLiquidityPoolFeeRate),
	}
}

// DefaultParams returns the default liquidity module parameters
func DefaultParams() Params {
	var defaultLiquidityPoolTypes []LiquidityPoolType
	defaultLiquidityPoolTypes = append(defaultLiquidityPoolTypes, LiquidityPoolTypeConstantProduct)

	return Params{
		LiquidityPoolTypes:      defaultLiquidityPoolTypes,
		MinInitDepositToPool:    sdk.NewInt(1000000),
		InitPoolTokenMintAmount: sdk.NewInt(1000000),
		SwapFeeRate:             sdk.NewDecWithPrec(1, 3), // "0.001000000000000000"
		LiquidityPoolFeeRate:    sdk.NewDecWithPrec(1, 3), // "0.001000000000000000"

	}
}

// String returns a human readable string representation of the parameters.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// Validate returns err if Params is invalid
func (p Params) Validate() error {
	// TODO: add detail validate logic
	return nil
}

func validateLiquidityPoolTypes(i interface{}) error {
	v, ok := i.([]LiquidityPoolType)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	for i, p := range v {
		if i != int(p.PoolTypeIndex) {
			return fmt.Errorf("LiquidityPoolTypes index must be sorted")
		}
	}
	return nil
}
func validateMinInitDepositToPool(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("MinInitDepositToPool must be positive: %d", v)
	}

	return nil
}

func validateInitPoolTokenMintAmount(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("InitPoolTokenMintAmount must be positive: %d", v)
	}

	return nil
}

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

func validateLiquidityPoolFeeRate(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("LiquidityPoolFeeRate cannot be negative: %s", v)
	}

	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("LiquidityPoolFeeRate too large: %s", v)
	}

	return nil
}
