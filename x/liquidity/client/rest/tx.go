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

// using grpc
func registerTxRoutes(clientCtx client.Context, r *mux.Router) {
}

func newLiquidityPoolHandlerFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.MsgCreatePoolRequest
		if !rest.ReadRESTReq(w, r, clientCtx.LegacyAmino, &req) {
			return
		}
	}
}