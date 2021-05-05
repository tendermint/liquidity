package types

import sdk "github.com/cosmos/cosmos-sdk/types"

type AirdropPair struct {
	TargetAddress string
	TargetAcc sdk.AccAddress
	DistributionCoins sdk.Coins
}

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

	// variables of example case for softfork airdrop with multiCoins
	Airdrop2SoftForkTargetHeight = int64(30000)
	Airdrop2ProviderAddr = "cosmos1f8s3n4lmlrancdrnnaky0j464prdr58d835yx2"
	Airdrop2Pairs = []AirdropPair {
		{"cosmos1w7xdwdllma6y2xhxwl3peurymx0tr95mk8urfp", nil, sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(100_000_000)), sdk.NewCoin("utest", sdk.NewInt(50_000_000)))},
		{"cosmos1uu9twaqca5f28ltdzqjlnklys4wcv97ke4038j", nil, sdk.NewCoins(sdk.NewCoin("uatom", sdk.NewInt(50_000_000)), sdk.NewCoin("utest", sdk.NewInt(100_000_000)))},
	}
)
