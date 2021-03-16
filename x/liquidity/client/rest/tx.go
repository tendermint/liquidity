package rest

// DONTCOVER
// client is excluded from test coverage in the poc phase milestone 1 and will be included in milestone 2 with completeness

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	"github.com/tendermint/liquidity/x/liquidity/types"
	"net/http"
)

// TODO: Plans to increase completeness on Milestone 2
// using grpc
func registerTxRoutes(clientCtx client.Context, r *mux.Router) {
	//// create liquidityPool
	//r.HandleFunc(fmt.Sprintf("/liquidity/pool"), newLiquidityPoolHandlerFn(clientCtx)).Methods("POST")
	//// deposit to liquidityPool
	//r.HandleFunc(fmt.Sprintf("/liquidity/pool/{%s}/deposit", RestPoolId), newDepositLiquidityPoolHandlerFn(clientCtx)).Methods("POST")
	//// withdraw from liquidityPool
	//r.HandleFunc(fmt.Sprintf("/liquidity/pool/{%s}/withdraw", RestPoolId), newWithdrawLiquidityPoolHandlerFn(clientCtx)).Methods("POST")
}

func newLiquidityPoolHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.MsgCreatePoolRequest
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}

		//baseReq := req.BaseReq.Sanitize()
		//if !baseReq.ValidateBasic(w) {
		//	return
		//}
		//
		//poolCreator, e := sdk.AccAddressFromBech32(req.PoolCreator)
		//if e != nil {
		//	rest.WriteErrorResponse(w, http.StatusBadRequest, e.Error())
		//	return
		//}
		//
		//depositCoin, ok := sdk.NewIntFromString(req.DepositCoins)
		//if !ok || !depositCoin.IsPositive() {
		//	rest.WriteErrorResponse(w, http.StatusBadRequest, "coin amount: "+req.DepositCoins)
		//	return
		//}
		//
		//msg := types.NewMsgCreatePool()
		//if err := msg.ValidateBasic(); err != nil {
		//	rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		//	return
		//}
		//
		//tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}

// HTTP request handler to add liquidity.
func newDepositLiquidityPoolHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//vars := mux.Vars(r)
		//poolID := vars[RestPoolId]
		//
		//var req DepositLiquidityPoolReq
		//if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
		//	return
		//}
		//
		//baseReq := req.BaseReq.Sanitize()
		//if !baseReq.ValidateBasic(w) {
		//	return
		//}
		//
		//msg := types.NewMsgDepositWithinBatch()
		//if err := msg.ValidateBasic(); err != nil {
		//	rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		//	return
		//}
		//
		//tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}

// HTTP request handler to remove liquidity.
func newWithdrawLiquidityPoolHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//vars := mux.Vars(r)
		//poolID := vars[RestPoolId]
		//
		//var req WithdrawLiquidityPoolReq
		//if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
		//	return
		//}
		//
		//baseReq := req.BaseReq.Sanitize()
		//if !baseReq.ValidateBasic(w) {
		//	return
		//}
		//
		//withdrawer, err := sdk.AccAddressFromBech32(req.Withdrawer)
		//if err != nil {
		//	return
		//}
		//poolId, err := strconv.ParseUint(req.Id, 10, 64)
		//sdk.NewCoin
		//msg := types.NewMsgWithdrawWithinBatch(withdrawer, poolId, req.PoolCoin)
		//if err := msg.ValidateBasic(); err != nil {
		//	rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		//	return
		//}
		//
		//tx.WriteGeneratedTxResponse(clientCtx, w, req.BaseReq, msg)
	}
}

//// WithdrawLiquidityPoolReq defines the properties of a Deposit from liquidity Pool request's body
//type CreatePoolReq struct {
//	BaseReq           rest.BaseReq `json:"base_req" yaml:"base_req"`
//	PoolCreator       string       `json:"pool_creator" yaml:"pool_creator"`
//	Id     string       `json:"pool_type_id" yaml:"pool_type_id"`
//	ReserveCoinDenoms string       `json:"reserve_coin_denoms" yaml:"reserve_coin_denoms"`
//	DepositCoins      string       `json:"deposit_coins" yaml:"deposit_coins"`
//}
//
//// WithdrawLiquidityPoolReq defines the properties of a Deposit from liquidity Pool request's body
//type WithdrawLiquidityPoolReq struct {
//	BaseReq    rest.BaseReq `json:"base_req" yaml:"base_req"`
//	Withdrawer string       `json:"withdrawer" yaml:"withdrawer"`
//	Id     string       `json:"pool_id" yaml:"pool_id"`
//	PoolCoin   string       `json:"pool_coin_amount" yaml:"pool_coin"`
//}
//
//// DepositLiquidityPoolReq defines the properties of a Deposit liquidity request's body
//type DepositLiquidityPoolReq struct {
//	BaseReq      rest.BaseReq `json:"base_req" yaml:"base_req"`
//	Depositor    string       `json:"depositor" yaml:"depositor"`
//	Id       string       `json:"pool_id" yaml:"pool_id"`
//	DepositCoins string       `json:"deposit_coins_amount" yaml:"deposit_coins"`
//}
//
//// DepositLiquidityPoolReq defines the properties of a Deposit liquidity request's body
//type SwapReq struct {
//	BaseReq         rest.BaseReq `json:"base_req" yaml:"base_req"`
//	SwapRequester   string       `json:"swap_requester" yaml:"swap_requester"`
//	Id          string       `json:"pool_id" yaml:"pool_id"`
//	Id   string       `json:"pool_type_id" yaml:"pool_type_id"`
//	SwapType        string       `json:"swap_type" yaml:"swap_type"`
//	OfferCoin       string       `json:"offer_coin" yaml:"offer_coin"`
//	DemandCoinDenom string       `json:"demand_coin" yaml:"demand_coin"`
//	OrderPrice      string       `json:"order_price" yaml:"order_price"`
//}
