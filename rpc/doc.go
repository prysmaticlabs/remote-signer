/*
Package rpc implements the Ethereum 2.0 Prysm remote signer spec.

It implements gRPC server for the following protobuf service definition:
https://github.com/prysmaticlabs/prysm/blob/master/proto/validator/accounts/v2/keymanager.proto

Note: this implementation is meant to be a reference, and not likely something used extensively in production.
*/
package rpc
