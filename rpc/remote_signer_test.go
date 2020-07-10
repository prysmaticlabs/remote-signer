package rpc

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	ptypes "github.com/gogo/protobuf/types"
	validatorpb "github.com/prysmaticlabs/prysm/proto/validator/accounts/v2"
	"github.com/prysmaticlabs/prysm/shared/bls"
	"github.com/prysmaticlabs/remote-signer/keyvault"
)

type mockKeyVault struct {
	pubKeys []bls.PublicKey
	wantErr bool
}

func (m *mockKeyVault) GetSecretKey(context.Context, bls.PublicKey) (bls.SecretKey, error) {
	if m.wantErr {
		return nil, errors.New("failed")
	}
	return bls.RandKey(), nil
}

func (m *mockKeyVault) GetPublicKeys(context.Context) ([]bls.PublicKey, error) {
	if m.wantErr {
		return nil, errors.New("failed")
	}
	return m.pubKeys, nil
}

func TestRemoteSigner_Sign(t *testing.T) {
	ctx := context.Background()
	badPubKey := make([]byte, blsPublicKeyLength)
	copy(badPubKey, "hello-world")
	tests := []struct {
		name     string
		keyVault keyvault.Store
		req      *validatorpb.SignRequest
		want     validatorpb.SignResponse_Status
		wantErr  bool
	}{
		{
			name:     "Fails with nil public key",
			keyVault: &mockKeyVault{},
			req: &validatorpb.SignRequest{
				PublicKey: nil,
			},
			want:    validatorpb.SignResponse_FAILED,
			wantErr: true,
		},
		{
			name:     "Fails with public key of wrong size",
			keyVault: &mockKeyVault{},
			req: &validatorpb.SignRequest{
				PublicKey: make([]byte, blsPublicKeyLength-1),
			},
			want:    validatorpb.SignResponse_FAILED,
			wantErr: true,
		},
		{
			name:     "Fails with bad public key",
			keyVault: &mockKeyVault{},
			req: &validatorpb.SignRequest{
				PublicKey: badPubKey,
			},
			want:    validatorpb.SignResponse_FAILED,
			wantErr: true,
		},
		{
			name:     "Fails with if vault retrieval fails",
			keyVault: &mockKeyVault{},
			req: &validatorpb.SignRequest{
				PublicKey: badPubKey,
			},
			want:    validatorpb.SignResponse_FAILED,
			wantErr: true,
		},
		{
			name: "Fails with if vault retrieval fails",
			keyVault: &mockKeyVault{
				wantErr: true,
			},
			req: &validatorpb.SignRequest{
				PublicKey: bls.RandKey().Marshal(),
			},
			want:    validatorpb.SignResponse_FAILED,
			wantErr: true,
		},
		{
			name:     "Succeeds with proper request",
			keyVault: &mockKeyVault{},
			req: &validatorpb.SignRequest{
				PublicKey:   bls.RandKey().PublicKey().Marshal(),
				SigningRoot: make([]byte, 32),
			},
			want:    validatorpb.SignResponse_SUCCEEDED,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RemoteSigner{
				keyVault: tt.keyVault,
			}
			got, err := r.Sign(ctx, tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Sign() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got.Status != tt.want {
				t.Errorf("Incorrect response status got = %v, want %v", got.Status, tt.want)
			}
		})
	}
}

func TestRemoteSigner_ListValidatingPublicKeys(t *testing.T) {
	ctx := context.Background()
	r := &RemoteSigner{
		keyVault: &mockKeyVault{wantErr: true},
	}
	_, err := r.ListValidatingPublicKeys(ctx, &ptypes.Empty{})
	if err == nil {
		t.Fatal("Wanted error, received nil")
	}
	want := "Could not retrieve public keys"
	if !strings.Contains(err.Error(), want) {
		t.Errorf("Wanted %v, received %v", want, err)
	}
	keys := make([]bls.PublicKey, 10)
	for i := 0; i < len(keys); i++ {
		keys[i] = bls.RandKey().PublicKey()
	}
	r.keyVault = &mockKeyVault{
		pubKeys: keys,
		wantErr: false,
	}
	res, err := r.ListValidatingPublicKeys(ctx, &ptypes.Empty{})
	if err != nil {
		t.Fatal(err)
	}
	numReceivedKeys := len(res.ValidatingPublicKeys)
	if numReceivedKeys != len(keys) {
		t.Errorf("Wanted %d keys, received %d", keys, numReceivedKeys)
	}
	for i := 0; i < numReceivedKeys; i++ {
		wantedKey := keys[i].Marshal()
		receivedKey := res.ValidatingPublicKeys[i]
		if !bytes.Equal(wantedKey, receivedKey) {
			t.Errorf("Wanted %#x, received %#x", wantedKey, receivedKey)
		}
	}
}
