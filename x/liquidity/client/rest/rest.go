package rest

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/gorilla/mux"
)

// Rest variable names
// nolint
const (
	RestPoolID = "pool-id"
)

// RegisterHandlers registers asset-related REST handlers to a router
func RegisterHandlers(cliCtx client.Context, r *mux.Router) {
	registerQueryRoutes(cliCtx, r)
	registerTxRoutes(cliCtx, r)
}
