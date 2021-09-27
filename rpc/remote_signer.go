package rpc

import (
	"context"

	emptypb "github.com/golang/protobuf/ptypes/empty"
	"github.com/prysmaticlabs/prysm/crypto/bls"
	validatorpb "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1/validator-client"
	"github.com/prysmaticlabs/remote-signer/keyvault"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const blsPublicKeyLength = 48 // 48 byte public keys.

// RemoteSigner capable of signing requests by using
// BLS secret keys retrieved from a keyvault.
type RemoteSigner struct {
	keyVault keyvault.Store
}

// NewRemoteSigner instantiates a new server instance using
// a keyvault for retrieving.
func NewRemoteSigner(ctx context.Context, keyVault keyvault.Store) *RemoteSigner {
	return &RemoteSigner{
		keyVault: keyVault,
	}
}

// Sign a remote request by retrieving the corresponding secret key for
// the public key in the request from a keyvault. If we have already signed
// the data in the request, we return a DENIED signing response.
func (r *RemoteSigner) Sign(ctx context.Context, req *validatorpb.SignRequest) (*validatorpb.SignResponse, error) {
	if req.PublicKey == nil {
		return &validatorpb.SignResponse{
			Status: validatorpb.SignResponse_FAILED,
		}, status.Error(codes.InvalidArgument, "Expected public key in request")
	}
	if len(req.PublicKey) != blsPublicKeyLength {
		return &validatorpb.SignResponse{
				Status: validatorpb.SignResponse_FAILED,
			}, status.Errorf(
				codes.InvalidArgument,
				"Wrong public key byte size: %d, expected %d",
				len(req.PublicKey),
				blsPublicKeyLength,
			)
	}
	pubKey, err := bls.PublicKeyFromBytes(req.PublicKey)
	if err != nil {
		return &validatorpb.SignResponse{
			Status: validatorpb.SignResponse_FAILED,
		}, status.Errorf(codes.InvalidArgument, "Could not parse public key: %v", err)
	}
	secretKey, err := r.keyVault.GetSecretKey(ctx, pubKey)
	if err != nil {
		return &validatorpb.SignResponse{
			Status: validatorpb.SignResponse_FAILED,
		}, status.Errorf(codes.Internal, "Could not fetch secret key from vault: %v", err)
	}
	// TODO: Handle naive slashing protection.
	sig := secretKey.Sign(req.SigningRoot)
	return &validatorpb.SignResponse{
		Signature: sig.Marshal(),
		Status:    validatorpb.SignResponse_SUCCEEDED,
	}, nil
}

// ListValidatingPublicKeys retrieves the BLS public keys
// available for signing in the remote signer.
func (r *RemoteSigner) ListValidatingPublicKeys(
	ctx context.Context, _ *emptypb.Empty,
) (*validatorpb.ListPublicKeysResponse, error) {
	pubKeys, err := r.keyVault.GetPublicKeys(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Could not retrieve public keys: %v", err)
	}
	rawKeys := make([][]byte, len(pubKeys))
	for i, k := range pubKeys {
		rawKeys[i] = k.Marshal()
	}
	return &validatorpb.ListPublicKeysResponse{
		ValidatingPublicKeys: rawKeys,
	}, nil
}
