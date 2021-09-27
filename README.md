# Eth2 Remote Signer

![Go](https://github.com/prysmaticlabs/remote-signer/workflows/Go/badge.svg?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/prysmaticlabs/remote-signer)](https://goreportcard.com/report/github.com/prysmaticlabs/remote-signer)
[![GoDoc](https://godoc.org/github.com/prysmaticlabs/remote-signer?status.svg)](https://godoc.org/github.com/prysmaticlabs/remote-signer)


Remote signing gRPC server reference implementation for the [Prysm](https://github.com/prysmaticlabs/prysm) Ethereum 2.0 client created by [Prysmatic Labs](https://prysmaticlabs.com)

## Overview

This is a simple, remote signing reference implementation to be used with the [Prysm](https://github.com/prysmaticlabs/prysm) project. It is **not** meant to be used in production deployments, but instead as an example of how to create a minimal remote-signer for eth2 validator keys in Go.

- Exposes a gRPC server implementation of the [RemoteSigner](https://github.com/prysmaticlabs/prysm/blob/develop/proto/prysm/v1alpha1/validator-client/keymanager.proto) service defined in Prysm secured by TLS certificates
- Exposes a gRPC gateway for JSON-HTTP requests to server 
- Allows for pluggable implementations different ways to load eth2 validator private keys, making it easy to integrate secure enclaves such as [Hashicorp Vault](https://learn.hashicorp.com/vault)

## Installation

To install this package, you need to install [Go](https://golang.org/doc/install). The simplest way to then fetch the project is to run:

```
$ go get -u github.com/prysmaticlabs/remote-signer
```

## Usage

Available parameters:

- **--grpc-server-host**: (required) host for the gRPC server, default 127.0.0.1
- **--grpc-port**: (required) port for the gRPC server, default 4000
- **--tls-crt-path**: (required) /path/to/server.crt for secure TLS connections
- **--tls-key-path**: (required) /path/to/server.key for secure TLS connections
- **--keyvault**: (required) type of [keyvault](https://github.com/prysmaticlabs/remote-signer/blob/master/keyvault/vault.go) to retrieve secret keys from, either: deterministic (default and unsafe) | mnemonic | s3 (unimplemented) | hashicorp (unimplemented)
- **--num-deterministic-keys**: number of deterministic keys to generate if using a deterministic keyvault
- **--num-mnemonic-keys**: number of keys to generate if using a mnemonic keyvault
- **--start-index**: starting index to generate the keys if using a mnemonic keyvault
- **--mnemonic-file**: file where mnemonic is placed if using a mnemonic keyvault
- **--mnemonic-password**: password of the mnemonic file if using a mnemonic keyvault

For local testing, example TLS cert key files for `localhost` are provided: [example-server.crt](https://github.com/prysmaticlabs/remote-signer/blob/master/example-server.crt) and [example-server.key](https://github.com/prysmaticlabs/remote-signer/blob/master/example-server.key) and [ca.crt](https://github.com/prysmaticlabs/remote-signer/blob/master/ca.crt). It is recommended you create your own TLS certificates using a tool such as [openssl](https://www.openssl.org/) or obtain new ones from a trusted certificate authority. For a tutorial on how to generate these certs for our use case, please see [securing your gRPC connection](https://docs.prylabs.network/docs/prysm-usage/secure-grpc) in our documentation portal.

```bash
$ go build -o server
```

### `deterministic` keyvault

```bash
$ ./server --tls-crt-path=./example-server.crt --tls-key-path=./example-server.key
```

Will output:

```text
INFO[0000] Loaded TLS certificates                       crt-path=./example-server.crt key-path=./example-server.key prefix=rpc
INFO[0000] gRPC server listening on address              address="127.0.0.1:4000" prefix=rpc
```


### `mnemonic` keyvault

Note that mnemonic and private keys are logged, making this tool useful for tests and debugging but not meant for production.

```bash
./server --tls-crt-path=example-server.crt --tls-key-path=example-server.key --keyvault=mnemonic --num-mnemonic-keys=3 --start-index=0 --mnemonic-file=sample-mnemonic.txt
```

Will output:
```text
INFO[0000] Contents of file: voice gospel easy verb front diesel sense worth sword equip giggle jeans shoe defy kid degree van frost like blush chef silk spoil obtain
INFO[0000] Generating keys from mnemonic                 prefix=mnemonic-keyvault
INFO[0000] Key from mnemonic at index 0, b2e17b5de68f5e929425028437daf3df6a6e7c9332c0cdbb9eb99d1cc115a56afd73b8fd04e8a1d11e53eb75b54d4176  prefix=mnemonic-keyvault
INFO[0000] Private Key from mnemonic 0, 2898511c2a36206a5a8e5c67dfb3549c8020edb0f4d855aa2f65106e720dce73  prefix=mnemonic-keyvault
INFO[0000] Key from mnemonic at index 1, a7d1d71d5e45f328ad5744341fa7a8f773fcaf3881c9b417479015f6f18326b702f1e13ce385cf0dc5db5558955a0e6e  prefix=mnemonic-keyvault
INFO[0000] Private Key from mnemonic 1, 1ecba74f354aac4c8563c69499b48846edb827953252223138a317802d9dbc10  prefix=mnemonic-keyvault
INFO[0000] Key from mnemonic at index 2, a2ec0b1deff9e6766a80c5e91130499fd00b5db6b607d3cf3b37e51423a8f32097c7fc69dd63d0a8cf14e17d491b0cec  prefix=mnemonic-keyvault
INFO[0000] Private Key from mnemonic 2, 6894d7f91519372074b475c50c65ff52c6870611c90a5e72c32975a97537960f  prefix=mnemonic-keyvault
INFO[0000] Loaded TLS certificates                       crt-path=example-server.crt key-path=example-server.key prefix=rpc
INFO[0000] gRPC server listening on address              address="127.0.0.1:4000" prefix=rpc
```


## Extending the Remote Signer

This reference implementation supports the retrieval of eth2 validator secrets from your desired secure enclave such as [Hashicorp Vault](https://learn.hashicorp.com/vault). You can define a new implementation of a `Store` ([keyvault/vault.go](https://github.com/prysmaticlabs/remote-signer/blob/master/keyvault/vault.go)) and add a new handler to the `--keyvault` flag in [main.go](https://github.com/prysmaticlabs/remote-signer/blob/master/main.go#L70) and the remote signer server will automatically be able to use it.

```go
// Store defines a struct which has capabilities of retrieving
// BLS12-381 eth2 secret keys and public keys from a secure source.
type Store interface {
	GetSecretKey(context.Context, bls.PublicKey) (bls.SecretKey, error)
	GetPublicKeys(context.Context) ([]bls.PublicKey, error)
}
```

By default, this reference implementation uses an **unsafe**, **deterministic** keyvault implementation which is meant to be there for demonstrative purposes. It is **not meant for production deployments** and merely an example on how to create a remote signer server to interact with a [Prysm validator client](https://github.com/prysmaticlabs/prysm).

Launching the remote server with default parameters and a deterministic keyvault:

```text
WARN[0000] You are using a deterministic keyvault (only for reference purposes) DO NOT USE in production  prefix=main
INFO[0000] Generating 1 determistic keys...              prefix=deterministic-keyvault
INFO[0000] Initialized deterministic keyvault            numKeys=1 prefix=deterministic-keyvault
INFO[0000] Loaded TLS certificates                       crt-path=./example-server.crt key-path=./example-server.key prefix=rpc
INFO[0000] gRPC server listening on address              address="127.0.0.1:4000" prefix=rpc
```

## Contributing

Contributions are very much welcome! Please fork the repository and create a pull request clearly explaining your feature, add tests, and sign our contributor licensing agreement which will automatically show up as a comment in your pull request. 

## License

[Apache License Version 2.0](https://github.com/prysmaticlabs/remote-signer/blob/master/LICENSE)
