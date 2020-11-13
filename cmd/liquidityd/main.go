package main

import (
	"os"

	"github.com/tendermint/liquidity/cmd/liquidityd/cmd"
)

func main() {
	rootCmd, _ := cmd.NewRootCmd()
	if err := cmd.Execute(rootCmd); err != nil {
		os.Exit(1)
	}
}
