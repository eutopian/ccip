// Code generated by mockery v2.35.4. DO NOT EDIT.

package mocks

import (
	common "github.com/ethereum/go-ethereum/common"
	ccipdata "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"

	context "context"

	internal "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal"

	mock "github.com/stretchr/testify/mock"
)

// OnRampReader is an autogenerated mock type for the OnRampReader type
type OnRampReader struct {
	mock.Mock
}

// Address provides a mock function with given fields:
func (_m *OnRampReader) Address() (common.Address, error) {
	ret := _m.Called()

	var r0 common.Address
	var r1 error
	if rf, ok := ret.Get(0).(func() (common.Address, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() common.Address); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(common.Address)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Close provides a mock function with given fields:
func (_m *OnRampReader) Close() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetDynamicConfig provides a mock function with given fields:
func (_m *OnRampReader) GetDynamicConfig() (ccipdata.OnRampDynamicConfig, error) {
	ret := _m.Called()

	var r0 ccipdata.OnRampDynamicConfig
	var r1 error
	if rf, ok := ret.Get(0).(func() (ccipdata.OnRampDynamicConfig, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() ccipdata.OnRampDynamicConfig); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(ccipdata.OnRampDynamicConfig)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSendRequestsBetweenSeqNums provides a mock function with given fields: ctx, seqNumMin, seqNumMax
func (_m *OnRampReader) GetSendRequestsBetweenSeqNums(ctx context.Context, seqNumMin uint64, seqNumMax uint64) ([]ccipdata.Event[internal.EVM2EVMMessage], error) {
	ret := _m.Called(ctx, seqNumMin, seqNumMax)

	var r0 []ccipdata.Event[internal.EVM2EVMMessage]
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint64, uint64) ([]ccipdata.Event[internal.EVM2EVMMessage], error)); ok {
		return rf(ctx, seqNumMin, seqNumMax)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint64, uint64) []ccipdata.Event[internal.EVM2EVMMessage]); ok {
		r0 = rf(ctx, seqNumMin, seqNumMax)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]ccipdata.Event[internal.EVM2EVMMessage])
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint64, uint64) error); ok {
		r1 = rf(ctx, seqNumMin, seqNumMax)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RouterAddress provides a mock function with given fields:
func (_m *OnRampReader) RouterAddress() (common.Address, error) {
	ret := _m.Called()

	var r0 common.Address
	var r1 error
	if rf, ok := ret.Get(0).(func() (common.Address, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() common.Address); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(common.Address)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewOnRampReader creates a new instance of OnRampReader. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewOnRampReader(t interface {
	mock.TestingT
	Cleanup(func())
}) *OnRampReader {
	mock := &OnRampReader{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
