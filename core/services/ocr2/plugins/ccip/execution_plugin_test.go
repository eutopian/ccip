package ccip

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/evm_2_evm_offramp"
	mock_contracts "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/mocks"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
)

func TestGetExecutionPluginFilterNamesFromSpec(t *testing.T) {
	testCases := []struct {
		description   string
		spec          *job.OCR2OracleSpec
		expectedNames []string
		expectingErr  bool
	}{
		{
			description:   "should not panic with nil spec",
			spec:          nil,
			expectedNames: nil,
			expectingErr:  true,
		},
		{
			description: "invalid config",
			spec: &job.OCR2OracleSpec{
				PluginConfig: map[string]interface{}{},
			},
			expectingErr: true,
		},
		{
			description: "invalid off ramp address",
			spec: &job.OCR2OracleSpec{
				PluginConfig: map[string]interface{}{"offRamp": "123"},
			},
			expectingErr: true,
		},
		{
			description: "invalid contract id",
			spec: &job.OCR2OracleSpec{
				ContractID: "whatever...",
			},
			expectingErr: true,
		},
	}

	for _, tc := range testCases {
		chainSet := &mocks.ChainSet{}
		t.Run(tc.description, func(t *testing.T) {
			names, err := GetExecutionPluginFilterNamesFromSpec(context.Background(), tc.spec, chainSet)
			assert.Equal(t, tc.expectedNames, names)
			if tc.expectingErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetExecutionPluginFilterNames(t *testing.T) {
	specContractID := common.HexToAddress("0xdafea492d9c6733ae3d56b7ed1adb60692c98bc1") // off-ramp addr
	onRampAddr := common.HexToAddress("0xdafea492d9c6733ae3d56b7ed1adb60692c98bc2")
	commitStoreAddr := common.HexToAddress("0xdafea492d9c6733ae3d56b7ed1adb60692c98bc3")

	mockOffRamp := mock_contracts.NewEVM2EVMOffRampInterface(t)
	mockOffRamp.On("GetStaticConfig", mock.Anything).Return(
		evm_2_evm_offramp.EVM2EVMOffRampStaticConfig{
			CommitStore: commitStoreAddr,
			OnRamp:      onRampAddr,
		}, nil)
	mockOffRamp.On("Address").Return(specContractID)

	filterNames, err := getExecutionPluginFilterNames(context.Background(), mockOffRamp)
	assert.NoError(t, err)
	assert.Equal(t, []string{
		"Exec ccip sends - 0xdafea492D9c6733aE3d56B7ED1aDb60692C98bc2",
		"Exec report accepts - 0xdafEa492d9C6733aE3D56b7eD1aDb60692c98bc3",
		"Exec execution state changes - 0xdafeA492d9c6733Ae3d56B7ed1AdB60692C98bC1",
		"Token pool added - 0xdafeA492d9c6733Ae3d56B7ed1AdB60692C98bC1",
		"Token pool removed - 0xdafeA492d9c6733Ae3d56B7ed1AdB60692C98bC1",
	}, filterNames)
}
