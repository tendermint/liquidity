package rest

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"

	"github.com/tendermint/liquidity/x/liquidity/types"
)

func registerQueryRoutes(cliCtx client.Context, r *mux.Router) {
	// query liquidity
	r.HandleFunc(fmt.Sprintf("/liquidity/pool/{%s}", RestPoolID), queryLiquidityHandlerFn(cliCtx)).Methods("GET")
}

// HTTP request handler to query liquidity information.
func queryLiquidityHandlerFn(cliCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		strPoolID := vars[RestPoolID]

		poolID, ok := rest.ParseUint64OrReturnBadRequest(w, strPoolID)
		if !ok {
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		params := types.QueryLiquidityPoolParams{
			PoolID: poolID,
		}

		bz, err := cliCtx.LegacyAmino.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryLiquidityPool)
		res, height, err := cliCtx.QueryWithData(route, bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}
