# Eth2 Remote Signer

![Go](https://github.com/prysmaticlabs/remote-signer/workflows/Go/badge.svg?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/prysmaticlabs/remote-signer)](https://goreportcard.com/report/github.com/prysmaticlabs/remote-signer)
[![GoDoc](https://godoc.org/github.com/prysmaticlabs/remote-signer?status.svg)](https://godoc.org/github.com/prysmaticlabs/remote-signer)


Remote signing gRPC server reference implementation for the Go Ethereum 2.0 client [prysmaticlabs/prysm](https://github.com/prysmaticlabs/prysm)

## Overview

This is a simple, remote signing reference implementation to be used with the [Prysm](https://github.com/prysmaticlabs/prysm) project. It is **not** meant to be used in production deployments, but instead as an example of how to create a minimal remote-signer for eth2 validator keys in Go.

- Exposes a gRPC server implementation of the [RemoteSigner](https://github.com/prysmaticlabs/prysm/blob/master/proto/validator/accounts/v2/keymanager.proto) service defined in Prysm
- Exposes a gRPC gateway for JSON-HTTP requests to server 
- Allows for pluggable implementations of different `KeyVaults` to fetch eth2 validator private keys, making it easy to integrate secure enclaves such as [Hashicorp Vault](https://learn.hashicorp.com/vault)

## Installation

To install this package, you need to install [Go](https://golang.org/doc/install). The simplest way to then fetch the project is to run:

```
$ go get -u github.com/prysmaticlabs/remote-signer
```

## Usage

## Extending the Remote Signer

## Contributing

## License

[Apache License Version 2.0](https://github.com/prysmaticlabs/remote-signer/blob/master/LICENSE)
