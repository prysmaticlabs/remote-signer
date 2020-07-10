package rpc

import (
	"context"

	ptypes "github.com/gogo/protobuf/types"
	validatorpb "github.com/prysmaticlabs/prysm/proto/validator/accounts/v2"
)

type RemoteSigner struct{}

func (r *RemoteSigner) Sign(ctx context.Context, req *validatorpb.SignRequest) (*validatorpb.SignResponse, error) {
	return nil, nil
}

func (r *RemoteSigner) ListValidatingPublicKeys(
	ctx context.Context, _ *ptypes.Empty,
) (*validatorpb.ListPublicKeysResponse, error) {
	return nil, nil
}
