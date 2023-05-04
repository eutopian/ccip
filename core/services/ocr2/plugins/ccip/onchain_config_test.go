package ccip

import (
	"math/big"
	"math/rand"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func randomAddress() common.Address {
	return common.BigToAddress(big.NewInt(rand.Int63()))
}

func TestCommitOnchainConfig(t *testing.T) {
	tests := []struct {
		name      string
		want      CommitOnchainConfig
		expectErr bool
	}{
		{
			name: "encodes and decodes config with all fields set",
			want: CommitOnchainConfig{
				PriceRegistry: randomAddress(),
				Afn:           randomAddress(),
			},
			expectErr: false,
		},
		{
			name: "encodes and fails decoding config with missing fields",
			want: CommitOnchainConfig{
				PriceRegistry: randomAddress(),
			},
			expectErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded, err := EncodeAbiStruct(tt.want)
			require.NoError(t, err)

			decoded, err := DecodeAbiStruct[CommitOnchainConfig](encoded)
			if tt.expectErr {
				require.ErrorContains(t, err, "must set")
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, decoded)
			}
		})
	}
}

func TestExecOnchainConfig(t *testing.T) {
	tests := []struct {
		name      string
		want      ExecOnchainConfig
		expectErr bool
	}{
		{
			name: "encodes and decodes config with all fields set",
			want: ExecOnchainConfig{
				PermissionLessExecutionThresholdSeconds: rand.Uint32(),
				Router:                                  randomAddress(),
				PriceRegistry:                           randomAddress(),
				Afn:                                     randomAddress(),
				MaxTokensLength:                         uint16(rand.Uint32()),
				MaxDataSize:                             rand.Uint32(),
			},
		},
		{
			name: "encodes and fails decoding config with missing fields",
			want: ExecOnchainConfig{
				PermissionLessExecutionThresholdSeconds: rand.Uint32(),
				MaxDataSize:                             rand.Uint32(),
			},
			expectErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded, err := EncodeAbiStruct(tt.want)
			require.NoError(t, err)

			decoded, err := DecodeAbiStruct[ExecOnchainConfig](encoded)
			if tt.expectErr {
				require.ErrorContains(t, err, "must set")
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.want, decoded)
			}
		})
	}
}
