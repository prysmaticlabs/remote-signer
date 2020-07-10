package rpc

import (
	"context"

	ptypes "github.com/gogo/protobuf/types"
	validatorpb "github.com/prysmaticlabs/prysm/proto/validator/accounts/v2"
	"github.com/prysmaticlabs/prysm/shared/bls"
	"github.com/prysmaticlabs/remote-signer/keyvault"
)

type RemoteSigner struct {
	keyVault keyvault.Store
}

func NewRemoteSigner(ctx context.Context, keyVault keyvault.Store) *RemoteSigner {
	return &RemoteSigner{
		keyVault: keyVault,
	}
}

func (r *RemoteSigner) Sign(ctx context.Context, req *validatorpb.SignRequest) (*validatorpb.SignResponse, error) {
	pubKey, err := bls.PublicKeyFromBytes(req.PublicKey)
	if err != nil {
		return nil, err
	}
	secretKey, err := r.keyVault.GetSecretKey(ctx, pubKey)
	if err != nil {
		return nil, err
	}
	sig := secretKey.Sign(req.SigningRoot)
	return &validatorpb.SignResponse{
		Signature: sig.Marshal(),
		Status:    validatorpb.SignResponse_SUCCEEDED,
	}, nil
}

func (r *RemoteSigner) ListValidatingPublicKeys(
	ctx context.Context, _ *ptypes.Empty,
) (*validatorpb.ListPublicKeysResponse, error) {
	pubKeys, err := r.keyVault.GetPublicKeys(ctx)
	if err != nil {
		return nil, err
	}
	rawKeys := make([][]byte, len(pubKeys))
	for i, k := range pubKeys {
		rawKeys[i] = k.Marshal()
	}
	return &validatorpb.ListPublicKeysResponse{
		ValidatingPublicKeys: rawKeys,
	}, nil
}
