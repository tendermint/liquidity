package types

import sdk "github.com/cosmos/cosmos-sdk/types"

var (
	// variables of example case for airdrop softfork
	Airdrop1SoftForkTargetHeight = int64(20000)
	Airdrop1ProviderAddr = "cosmos1f8s3n4lmlrancdrnnaky0j464prdr58d835yx2"
	Airdrop1DistributionCoin = sdk.NewCoin("uatom", sdk.NewInt(100_000_000))
	Airdrop1TargetAddrs = []string {
		"cosmos1w7xdwdllma6y2xhxwl3peurymx0tr95mk8urfp",
		"cosmos1uu9twaqca5f28ltdzqjlnklys4wcv97ke4038j",
		//"cosmos1...",
	}
)
