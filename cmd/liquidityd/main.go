package main

import (
	"os"

	_ "github.com/tendermint/liquidity/client/docs/statik"
	"github.com/tendermint/liquidity/cmd/liquidityd/cmd"
)

func main() {
	rootCmd, _ := cmd.NewRootCmd()
	if err := cmd.Execute(rootCmd); err != nil {
		os.Exit(1)
	}
}
