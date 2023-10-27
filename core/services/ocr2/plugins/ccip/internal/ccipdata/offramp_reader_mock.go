// Code generated by mockery v2.35.4. DO NOT EDIT.

package ccipdata

import (
	bind "github.com/ethereum/go-ethereum/accounts/abi/bind"
	common "github.com/ethereum/go-ethereum/common"

	context "context"

	evm_2_evm_offramp "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/evm_2_evm_offramp"

	mock "github.com/stretchr/testify/mock"

	pg "github.com/smartcontractkit/chainlink/v2/core/services/pg"

	prices "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/prices"
)

// MockOffRampReader is an autogenerated mock type for the OffRampReader type
type MockOffRampReader struct {
	mock.Mock
}

// Address provides a mock function with given fields:
func (_m *MockOffRampReader) Address() common.Address {
	ret := _m.Called()

	var r0 common.Address
	if rf, ok := ret.Get(0).(func() common.Address); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(common.Address)
		}
	}

	return r0
}

// ChangeConfig provides a mock function with given fields: onchainConfig, offchainConfig
func (_m *MockOffRampReader) ChangeConfig(onchainConfig []byte, offchainConfig []byte) (common.Address, common.Address, error) {
	ret := _m.Called(onchainConfig, offchainConfig)

	var r0 common.Address
	var r1 common.Address
	var r2 error
	if rf, ok := ret.Get(0).(func([]byte, []byte) (common.Address, common.Address, error)); ok {
		return rf(onchainConfig, offchainConfig)
	}
	if rf, ok := ret.Get(0).(func([]byte, []byte) common.Address); ok {
		r0 = rf(onchainConfig, offchainConfig)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(common.Address)
		}
	}

	if rf, ok := ret.Get(1).(func([]byte, []byte) common.Address); ok {
		r1 = rf(onchainConfig, offchainConfig)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(common.Address)
		}
	}

	if rf, ok := ret.Get(2).(func([]byte, []byte) error); ok {
		r2 = rf(onchainConfig, offchainConfig)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// Close provides a mock function with given fields: qopts
func (_m *MockOffRampReader) Close(qopts ...pg.QOpt) error {
	_va := make([]interface{}, len(qopts))
	for _i := range qopts {
		_va[_i] = qopts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(...pg.QOpt) error); ok {
		r0 = rf(qopts...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// CurrentRateLimiterState provides a mock function with given fields: opts
func (_m *MockOffRampReader) CurrentRateLimiterState(opts *bind.CallOpts) (evm_2_evm_offramp.RateLimiterTokenBucket, error) {
	ret := _m.Called(opts)

	var r0 evm_2_evm_offramp.RateLimiterTokenBucket
	var r1 error
	if rf, ok := ret.Get(0).(func(*bind.CallOpts) (evm_2_evm_offramp.RateLimiterTokenBucket, error)); ok {
		return rf(opts)
	}
	if rf, ok := ret.Get(0).(func(*bind.CallOpts) evm_2_evm_offramp.RateLimiterTokenBucket); ok {
		r0 = rf(opts)
	} else {
		r0 = ret.Get(0).(evm_2_evm_offramp.RateLimiterTokenBucket)
	}

	if rf, ok := ret.Get(1).(func(*bind.CallOpts) error); ok {
		r1 = rf(opts)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DecodeExecutionReport provides a mock function with given fields: report
func (_m *MockOffRampReader) DecodeExecutionReport(report []byte) (ExecReport, error) {
	ret := _m.Called(report)

	var r0 ExecReport
	var r1 error
	if rf, ok := ret.Get(0).(func([]byte) (ExecReport, error)); ok {
		return rf(report)
	}
	if rf, ok := ret.Get(0).(func([]byte) ExecReport); ok {
		r0 = rf(report)
	} else {
		r0 = ret.Get(0).(ExecReport)
	}

	if rf, ok := ret.Get(1).(func([]byte) error); ok {
		r1 = rf(report)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// EncodeExecutionReport provides a mock function with given fields: report
func (_m *MockOffRampReader) EncodeExecutionReport(report ExecReport) ([]byte, error) {
	ret := _m.Called(report)

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(ExecReport) ([]byte, error)); ok {
		return rf(report)
	}
	if rf, ok := ret.Get(0).(func(ExecReport) []byte); ok {
		r0 = rf(report)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(ExecReport) error); ok {
		r1 = rf(report)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GasPriceEstimator provides a mock function with given fields:
func (_m *MockOffRampReader) GasPriceEstimator() prices.GasPriceEstimatorExec {
	ret := _m.Called()

	var r0 prices.GasPriceEstimatorExec
	if rf, ok := ret.Get(0).(func() prices.GasPriceEstimatorExec); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(prices.GasPriceEstimatorExec)
		}
	}

	return r0
}

// GetDestinationToken provides a mock function with given fields: ctx, address
func (_m *MockOffRampReader) GetDestinationToken(ctx context.Context, address common.Address) (common.Address, error) {
	ret := _m.Called(ctx, address)

	var r0 common.Address
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, common.Address) (common.Address, error)); ok {
		return rf(ctx, address)
	}
	if rf, ok := ret.Get(0).(func(context.Context, common.Address) common.Address); ok {
		r0 = rf(ctx, address)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(common.Address)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, common.Address) error); ok {
		r1 = rf(ctx, address)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetDestinationTokens provides a mock function with given fields: ctx
func (_m *MockOffRampReader) GetDestinationTokens(ctx context.Context) ([]common.Address, error) {
	ret := _m.Called(ctx)

	var r0 []common.Address
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]common.Address, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []common.Address); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]common.Address)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetExecutionState provides a mock function with given fields: opts, sequenceNumber
func (_m *MockOffRampReader) GetExecutionState(opts *bind.CallOpts, sequenceNumber uint64) (uint8, error) {
	ret := _m.Called(opts, sequenceNumber)

	var r0 uint8
	var r1 error
	if rf, ok := ret.Get(0).(func(*bind.CallOpts, uint64) (uint8, error)); ok {
		return rf(opts, sequenceNumber)
	}
	if rf, ok := ret.Get(0).(func(*bind.CallOpts, uint64) uint8); ok {
		r0 = rf(opts, sequenceNumber)
	} else {
		r0 = ret.Get(0).(uint8)
	}

	if rf, ok := ret.Get(1).(func(*bind.CallOpts, uint64) error); ok {
		r1 = rf(opts, sequenceNumber)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetExecutionStateChangesBetweenSeqNums provides a mock function with given fields: ctx, seqNumMin, seqNumMax, confs
func (_m *MockOffRampReader) GetExecutionStateChangesBetweenSeqNums(ctx context.Context, seqNumMin uint64, seqNumMax uint64, confs int) ([]Event[ExecutionStateChanged], error) {
	ret := _m.Called(ctx, seqNumMin, seqNumMax, confs)

	var r0 []Event[ExecutionStateChanged]
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint64, uint64, int) ([]Event[ExecutionStateChanged], error)); ok {
		return rf(ctx, seqNumMin, seqNumMax, confs)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint64, uint64, int) []Event[ExecutionStateChanged]); ok {
		r0 = rf(ctx, seqNumMin, seqNumMax, confs)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]Event[ExecutionStateChanged])
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint64, uint64, int) error); ok {
		r1 = rf(ctx, seqNumMin, seqNumMax, confs)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetOffRampAddress provides a mock function with given fields:
func (_m *MockOffRampReader) GetOffRampAddress() common.Address {
	ret := _m.Called()

	var r0 common.Address
	if rf, ok := ret.Get(0).(func() common.Address); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(common.Address)
		}
	}

	return r0
}

// GetPoolByDestToken provides a mock function with given fields: ctx, address
func (_m *MockOffRampReader) GetPoolByDestToken(ctx context.Context, address common.Address) (common.Address, error) {
	ret := _m.Called(ctx, address)

	var r0 common.Address
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, common.Address) (common.Address, error)); ok {
		return rf(ctx, address)
	}
	if rf, ok := ret.Get(0).(func(context.Context, common.Address) common.Address); ok {
		r0 = rf(ctx, address)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(common.Address)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, common.Address) error); ok {
		r1 = rf(ctx, address)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSenderNonce provides a mock function with given fields: opts, sender
func (_m *MockOffRampReader) GetSenderNonce(opts *bind.CallOpts, sender common.Address) (uint64, error) {
	ret := _m.Called(opts, sender)

	var r0 uint64
	var r1 error
	if rf, ok := ret.Get(0).(func(*bind.CallOpts, common.Address) (uint64, error)); ok {
		return rf(opts, sender)
	}
	if rf, ok := ret.Get(0).(func(*bind.CallOpts, common.Address) uint64); ok {
		r0 = rf(opts, sender)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	if rf, ok := ret.Get(1).(func(*bind.CallOpts, common.Address) error); ok {
		r1 = rf(opts, sender)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSupportedTokens provides a mock function with given fields: ctx
func (_m *MockOffRampReader) GetSupportedTokens(ctx context.Context) ([]common.Address, error) {
	ret := _m.Called(ctx)

	var r0 []common.Address
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]common.Address, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []common.Address); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]common.Address)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// OffchainConfig provides a mock function with given fields:
func (_m *MockOffRampReader) OffchainConfig() ExecOffchainConfig {
	ret := _m.Called()

	var r0 ExecOffchainConfig
	if rf, ok := ret.Get(0).(func() ExecOffchainConfig); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(ExecOffchainConfig)
	}

	return r0
}

// OnchainConfig provides a mock function with given fields:
func (_m *MockOffRampReader) OnchainConfig() ExecOnchainConfig {
	ret := _m.Called()

	var r0 ExecOnchainConfig
	if rf, ok := ret.Get(0).(func() ExecOnchainConfig); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(ExecOnchainConfig)
	}

	return r0
}

// TokenEvents provides a mock function with given fields:
func (_m *MockOffRampReader) TokenEvents() []common.Hash {
	ret := _m.Called()

	var r0 []common.Hash
	if rf, ok := ret.Get(0).(func() []common.Hash); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]common.Hash)
		}
	}

	return r0
}

// NewMockOffRampReader creates a new instance of MockOffRampReader. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockOffRampReader(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockOffRampReader {
	mock := &MockOffRampReader{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
