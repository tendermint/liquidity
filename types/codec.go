package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RegisterLegacyAminoCodec registers concrete types on the codec.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreateLiquidityPool{}, "liquidity/MsgCreateLiquidityPool", nil)
	cdc.RegisterConcrete(&MsgDepositToLiquidityPool{}, "liquidity/MsgDepositToLiquidityPool", nil)
	cdc.RegisterConcrete(&MsgWithdrawFromLiquidityPool{}, "liquidity/MsgWithdrawFromLiquidityPool", nil)
	cdc.RegisterConcrete(&MsgSwap{}, "liquidity/MsgSwap", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateLiquidityPool{},
		&MsgDepositToLiquidityPool{},
		&MsgWithdrawFromLiquidityPool{},
		&MsgSwap{},
	)
}

var (
	amino = codec.NewLegacyAmino()

	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()
}
