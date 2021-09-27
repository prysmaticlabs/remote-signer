module github.com/prysmaticlabs/remote-signer

go 1.14

require (
	github.com/gogo/protobuf v1.3.2
	github.com/golang/protobuf v1.5.2
	github.com/kr/text v0.2.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/prysmaticlabs/prysm v1.4.2-0.20210927203955-328e3e6caf83
	github.com/sirupsen/logrus v1.7.0
	google.golang.org/grpc v1.37.0
)

replace github.com/ethereum/go-ethereum => github.com/prysmaticlabs/bazel-go-ethereum v0.0.0-20210707101027-e8523651bf6f

// See https://github.com/prysmaticlabs/grpc-gateway/issues/2
replace github.com/grpc-ecosystem/grpc-gateway/v2 => github.com/prysmaticlabs/grpc-gateway/v2 v2.3.1-0.20210702154020-550e1cd83ec1
