package types

import (
	"fmt"
	"gopkg.in/yaml.v2"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

const (
	UnitBatchSize uint32 = 1
	CancelOrderLifeSpan int64 = 0
	MinReserveCoinNum uint32 = 2
	MaxReserveCoinNum uint32 = 2
)

// Parameter store keys
var (
	KeyLiquidityPoolTypes       = []byte("LiquidityPoolTypes")
	KeyMinInitDepositToPool     = []byte("MinInitDepositToPool")
	KeyInitPoolCoinMintAmount   = []byte("InitPoolCoinMintAmount")
	KeySwapFeeRate              = []byte("SwapFeeRate")
	KeyLiquidityPoolFeeRate     = []byte("LiquidityPoolFeeRate")
	KeyLiquidityPoolCreationFee = []byte("LiquidityPoolCreationFee")
	KeyUnitBatchSize            = []byte("UnitBatchSize")

	DefaultLiquidityPoolType = LiquidityPoolType{
		PoolTypeIndex:     0,
		Name:              "DefaultPoolType",
		MinReserveCoinNum: MinReserveCoinNum,
		MaxReserveCoinNum: MaxReserveCoinNum,
	}
)

// NewParams liquidity paramtypes constructor
func NewParams(liquidityPoolTypes []LiquidityPoolType, minInitDeposit, initPoolCoinMint sdk.Int, swapFeeRate sdk.Dec, creationFee sdk.Coins) Params {
	return Params{
		LiquidityPoolTypes:       liquidityPoolTypes,
		MinInitDepositToPool:     minInitDeposit,
		InitPoolCoinMintAmount:   initPoolCoinMint,
		SwapFeeRate:              swapFeeRate,
		LiquidityPoolCreationFee: creationFee,
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
		paramtypes.NewParamSetPair(KeyInitPoolCoinMintAmount, &p.InitPoolCoinMintAmount, validateInitPoolCoinMintAmount),
		paramtypes.NewParamSetPair(KeySwapFeeRate, &p.SwapFeeRate, validateSwapFeeRate),
		paramtypes.NewParamSetPair(KeyLiquidityPoolCreationFee, &p.LiquidityPoolCreationFee, validateLiquidityPoolCreationFee),
	}
}

// DefaultParams returns the default liquidity module parameters
func DefaultParams() Params {
	var defaultLiquidityPoolTypes []LiquidityPoolType
	defaultLiquidityPoolTypes = append(defaultLiquidityPoolTypes, DefaultLiquidityPoolType)

	return NewParams(
		defaultLiquidityPoolTypes,
		sdk.NewInt(1000000),
		sdk.NewInt(1000000),
		sdk.NewDecWithPrec(3, 3),  // "0.001000000000000000"
		sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(100000000))))
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
	v, ok := i.(sdk.Int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == sdk.ZeroInt() {
		return fmt.Errorf("MinInitDepositToPool must be positive: %d", v)
	}

	return nil
}

func validateInitPoolCoinMintAmount(i interface{}) error {
	v, ok := i.(sdk.Int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == sdk.ZeroInt() {
		return fmt.Errorf("InitPoolCoinMintAmount must be positive: %d", v)
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

func validateLiquidityPoolCreationFee(i interface{}) error {
	coins, ok := i.(sdk.Coins)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	if coins.Empty() {
		return fmt.Errorf("LiquidityPoolCreationFee cannot be Empty: %s", coins)
	}
	return nil
}
