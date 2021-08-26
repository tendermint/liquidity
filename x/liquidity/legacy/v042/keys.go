// Package v042 is copy-pasted from:
// https://github.com/tendermint/liquidity/blob/v1.2.9/x/liquidity/types/keys.go
package v042

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName is the name of the liquidity module
	ModuleName = "liquidity"
)

var (
	PoolByReserveAccIndexKeyPrefix = []byte{0x12}
)

// - PoolByReserveAccIndex: `0x12 | ReserveAcc -> Id`
// GetPoolByReserveAccIndexKey returns kv indexing key of the pool indexed by reserve account
func GetPoolByReserveAccIndexKey(reserveAcc sdk.AccAddress) []byte {
	return append(PoolByReserveAccIndexKeyPrefix, reserveAcc.Bytes()...)
}
