go 1.15

module github.com/tendermint/liquidity

require (
	github.com/cosmos/cosmos-sdk v0.34.4-0.20201006170641-3e6089dc0ea6
	github.com/gogo/protobuf v1.3.1
	github.com/golang/protobuf v1.4.2
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/grpc-gateway v1.15.0
	github.com/spf13/cobra v1.0.0
	github.com/tendermint/tendermint v0.34.0-rc4
	google.golang.org/genproto v0.0.0-20200825200019-8632dd797987
	google.golang.org/grpc v1.32.0
	gopkg.in/yaml.v2 v2.3.0
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.4
