package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
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

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino = codec.NewLegacyAmino()

	// ModuleCdc references the global x/staking module codec. Note, the codec should
	// ONLY be used in certain instances of tests and for JSON encoding as Amino is
	// still used for that purpose.
	//
	// The actual codec used for serialization should be provided to x/staking and
	// defined at the application level.
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()
}
