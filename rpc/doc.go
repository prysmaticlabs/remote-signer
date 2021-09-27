/*
Package rpc implements the Ethereum 2.0 Prysm remote signer spec.

It implements gRPC server for the following protobuf service definition:
https://github.com/prysmaticlabs/prysm/blob/develop/proto/prysm/v1alpha1/validator-client/keymanager.proto

Note: this implementation is meant to be a reference, and not meant to be used as is in production.
*/
package rpc
