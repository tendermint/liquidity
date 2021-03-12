package types_test

import (
	"testing"

	"github.com/tendermint/liquidity/x/liquidity/types"

	"github.com/stretchr/testify/suite"
)

type keysTestSuite struct {
	suite.Suite
}

func TestKeysTestSuite(t *testing.T) {
	suite.Run(t, new(keysTestSuite))
}

func (s *keysTestSuite) TestGetLiquidityPoolKey() {
	s.Require().Equal([]byte{0x11, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x9}, types.GetPoolKey(9))
	s.Require().Equal([]byte{0x11, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}, types.GetPoolKey(0))
}

func (s *keysTestSuite) TestGetLiquidityPoolByReserveAccIndexKey() {
	s.Require().Equal([]byte{18, 116, 101, 115, 116}, types.GetPoolByReserveAccIndexKey([]byte("test")))
	s.Require().Equal([]byte{18}, types.GetPoolByReserveAccIndexKey(nil))
}

func (s *keysTestSuite) TestGetLiquidityPoolBatchIndexKey() {
	s.Require().Equal([]byte{0x21, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xa}, types.GetPoolBatchIndexKey(10))
}

func (s *keysTestSuite) TestGetLiquidityPoolBatchKey() {
	s.Require().Equal([]byte{0x22, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xa}, types.GetPoolBatchKey(10))
	s.Require().Equal([]byte{0x22, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}, types.GetPoolBatchKey(0))
}

func (s *keysTestSuite) TestGetLiquidityPoolBatchDepositMsgsPrefix() {
	s.Require().Equal([]byte{0x31, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xa}, types.GetPoolBatchDepositMsgStatesPrefix(10))
	s.Require().Equal([]byte{0x31, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}, types.GetPoolBatchDepositMsgStatesPrefix(0))
}

func (s *keysTestSuite) TestGetLiquidityPoolBatchWithdrawMsgsPrefix() {
	s.Require().Equal([]byte{0x32, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xa}, types.GetPoolBatchWithdrawMsgsPrefix(10))
	s.Require().Equal([]byte{0x32, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}, types.GetPoolBatchWithdrawMsgsPrefix(0))
}

func (s *keysTestSuite) TestGetLiquidityPoolBatchSwapMsgsPrefix() {
	s.Require().Equal([]byte{0x33, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xa}, types.GetPoolBatchSwapMsgStatesPrefix(10))
	s.Require().Equal([]byte{0x33, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0}, types.GetPoolBatchSwapMsgStatesPrefix(0))
}

func (s *keysTestSuite) TestGetLiquidityPoolBatchDepositMsgIndex() {
	s.Require().Equal([]byte{0x31, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xa, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xa},
		types.GetPoolBatchDepositMsgStateIndexKey(10, 10))
	s.Require().Equal([]byte{0x31, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
		types.GetPoolBatchDepositMsgStateIndexKey(0, 0))
}

func (s *keysTestSuite) TestGetLiquidityPoolBatchWithdrawMsgIndex() {
	s.Require().Equal([]byte{0x32, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xa, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xa},
		types.GetPoolBatchWithdrawMsgStateIndexKey(10, 10))
	s.Require().Equal([]byte{0x32, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
		types.GetPoolBatchWithdrawMsgStateIndexKey(0, 0))
}

func (s *keysTestSuite) TestGetLiquidityPoolBatchSwapMsgIndex() {
	s.Require().Equal([]byte{0x33, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xa, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xa},
		types.GetPoolBatchSwapMsgStateIndexKey(10, 10))
	s.Require().Equal([]byte{0x33, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
		types.GetPoolBatchSwapMsgStateIndexKey(0, 0))
}
