package rest

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/gorilla/mux"
)

// Rest variable names
// nolint
const (
	RestPoolId = "pool-id"
)

// TODO: after rebase latest stable sdk 0.40.0 for other endpoints
// RegisterHandlers registers asset-related REST handlers to a router
func RegisterHandlers(cliCtx client.Context, r *mux.Router) {
	//registerQueryRoutes(cliCtx, r)
	//registerTxRoutes(cliCtx, r)
}
