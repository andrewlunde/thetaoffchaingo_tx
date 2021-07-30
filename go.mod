module github.com/thetatoken/theta/tx

go 1.14

require (
	github.com/onsi/gomega v1.14.0 // indirect
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.5.0
	github.com/thetatoken/theta v0.0.0
	github.com/thetatoken/theta/common v0.0.0
	github.com/thetatoken/theta/utils v0.0.0
	github.com/ybbus/jsonrpc v1.1.1
)

replace (
	// github.com/thetatoken/theta v0.0.0 => ../theta-protocol-ledger
	github.com/thetatoken/theta v0.0.0 => github.com/andrewlunde/theta-protocol-ledger v0.1.0
	// github.com/thetatoken/theta/common v0.0.0 => ../theta-protocol-ledger/common
	github.com/thetatoken/theta/common v0.0.0 => github.com/andrewlunde/thetaoffchaingo_common v0.1.1
	// github.com/thetatoken/theta/rpc/lib/rpc-codec/jsonrpc2 v0.0.0 => ../theta-protocol-ledger/rpc/lib/rpc-codec/jsonrpc2/
	github.com/thetatoken/theta/rpc/lib/rpc-codec/jsonrpc2 v0.0.0 => github.com/andrewlunde/thetaoffchaingo_rpc v0.1.1
	// github.com/thetatoken/theta/utils v0.0.0 => ../thetaoffchaingo_utils
	github.com/thetatoken/theta/utils v0.0.0 => github.com/andrewlunde/thetaoffchaingo_utils v0.1.0
)
