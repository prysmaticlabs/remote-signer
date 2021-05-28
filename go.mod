module github.com/prysmaticlabs/remote-signer

go 1.14

require (
	github.com/gogo/protobuf v1.3.1
	github.com/pkg/errors v0.9.1
	github.com/prysmaticlabs/prysm v1.0.0
	github.com/sirupsen/logrus v1.7.0
	github.com/tyler-smith/go-bip39 v1.0.2
	google.golang.org/grpc v1.33.1
)

replace github.com/ethereum/go-ethereum => github.com/prysmaticlabs/bazel-go-ethereum v0.0.0-20201113091623-013fd65b3791
