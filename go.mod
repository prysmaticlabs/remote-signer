module github.com/prysmaticlabs/remote-signer

go 1.14

require (
	github.com/gogo/protobuf v1.3.1
	github.com/pkg/errors v0.9.1
	github.com/prysmaticlabs/prysm v1.0.0-alpha.26
	github.com/sirupsen/logrus v1.6.0
	google.golang.org/grpc v1.30.0
)

replace github.com/ethereum/go-ethereum => github.com/prysmaticlabs/bazel-go-ethereum v0.0.0-20200626171358-a933315235ec
