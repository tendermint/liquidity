package rest

// DONTCOVER

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/rest"
	"github.com/gorilla/mux"
)

// Rest variable names
// nolint
const (
	RestPoolId = "pool-id"
)

// TODO: Plans to increase completeness on Milestone 2

// RegisterHandlers registers asset-related REST handlers to a router
func RegisterHandlers(cliCtx client.Context, r *mux.Router) {
	r = rest.WithHTTPDeprecationHeaders(r)
	registerQueryRoutes(cliCtx, r)
	//registerTxRoutes(cliCtx, r)
}
