package rest

// DONTCOVER
// client is excluded from test coverage in the poc phase milestone 1 and will be included in milestone 2 with completeness

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"

	"github.com/tendermint/liquidity/x/liquidity/types"
)

// using grpc server
func registerQueryRoutes(cliCtx client.Context, r *mux.Router) {
	// Get the liquidity pool
	r.HandleFunc(
		fmt.Sprintf("liquidity/legacy/pools/{%s}", RestPoolId),
		queryLiquidityPoolHandlerFn(cliCtx),
	).Methods("GET")

	// Get all liquidity pools
	r.HandleFunc(
		"liquidity/legacy/pools",
		queryLiquidityPoolsHandlerFn(cliCtx)).Methods("GET")
}

// HTTP request handler to query liquidity information.
func queryLiquidityPoolHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		strPoolId := vars[RestPoolId]

		poolId, ok := rest.ParseUint64OrReturnBadRequest(w, strPoolId)
		if !ok {
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		params := types.NewQueryLiquidityPoolParams(poolId)

		bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryLiquidityPool)
		res, height, err := cliCtx.QueryWithData(route, bz)
		if rest.CheckInternalServerError(w, err) {
			return
		}
		//if err != nil {
		//	rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		//	return
		//}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

// HTTP request handler to query list of validators
func queryLiquidityPoolsHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, page, limit, err := rest.ParseHTTPArgsWithLimit(r, 0)
		fmt.Println(page, limit, err)
		if rest.CheckBadRequestError(w, err) {
			fmt.Println("CheckBadRequestError", w, err)
			return
		}

		clientCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		fmt.Println("clientCtx", clientCtx, ok)
		if !ok {
			return
		}

		params := types.NewQueryLiquidityPoolsParams(page, limit)

		fmt.Println("params", params)
		bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
		fmt.Println("bz, err", bz, err)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryLiquidityPools)

		res, height, err := clientCtx.QueryWithData(route, bz)

		fmt.Println("route, res, height, err", route, res, height, err)
		if rest.CheckInternalServerError(w, err) {
			return
		}

		clientCtx = clientCtx.WithHeight(height)
		rest.PostProcessResponse(w, clientCtx, res)
	}
}

// HTTP request handler to query list of validators
func queryLiquidityBatchSwapsHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, page, limit, err := rest.ParseHTTPArgsWithLimit(r, 0)
		fmt.Println(page, limit, err)
		if rest.CheckBadRequestError(w, err) {
			fmt.Println("CheckBadRequestError", w, err)
			return
		}

		clientCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, clientCtx, r)
		fmt.Println("clientCtx", clientCtx, ok)
		if !ok {
			return
		}

		params := types.NewQueryLiquidityPoolsParams(page, limit)

		fmt.Println("params", params)
		bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
		fmt.Println("bz, err", bz, err)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryLiquidityPools)

		fmt.Println("route, res, height, err", route)
		res, height, err := clientCtx.QueryWithData(route, bz)

		fmt.Println("route, res, height, err", route, res, height, err)
		if rest.CheckInternalServerError(w, err) {
			return
		}

		clientCtx = clientCtx.WithHeight(height)
		rest.PostProcessResponse(w, clientCtx, res)
	}
}
