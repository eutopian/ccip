// Code generated by mockery v2.22.1. DO NOT EDIT.

package mocks

import (
	common "github.com/ethereum/go-ethereum/common"
	assets "github.com/smartcontractkit/chainlink/core/assets"

	mock "github.com/stretchr/testify/mock"
)

// Config is an autogenerated mock type for the Config type
type Config struct {
	mock.Mock
}

// EvmFinalityDepth provides a mock function with given fields:
func (_m *Config) EvmFinalityDepth() uint32 {
	ret := _m.Called()

	var r0 uint32
	if rf, ok := ret.Get(0).(func() uint32); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint32)
	}

	return r0
}

// EvmGasLimitDefault provides a mock function with given fields:
func (_m *Config) EvmGasLimitDefault() uint32 {
	ret := _m.Called()

	var r0 uint32
	if rf, ok := ret.Get(0).(func() uint32); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint32)
	}

	return r0
}

// EvmGasLimitVRFJobType provides a mock function with given fields:
func (_m *Config) EvmGasLimitVRFJobType() *uint32 {
	ret := _m.Called()

	var r0 *uint32
	if rf, ok := ret.Get(0).(func() *uint32); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*uint32)
		}
	}

	return r0
}

// KeySpecificMaxGasPriceWei provides a mock function with given fields: addr
func (_m *Config) KeySpecificMaxGasPriceWei(addr common.Address) *assets.Wei {
	ret := _m.Called(addr)

	var r0 *assets.Wei
	if rf, ok := ret.Get(0).(func(common.Address) *assets.Wei); ok {
		r0 = rf(addr)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*assets.Wei)
		}
	}

	return r0
}

// MinIncomingConfirmations provides a mock function with given fields:
func (_m *Config) MinIncomingConfirmations() uint32 {
	ret := _m.Called()

	var r0 uint32
	if rf, ok := ret.Get(0).(func() uint32); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(uint32)
	}

	return r0
}

type mockConstructorTestingTNewConfig interface {
	mock.TestingT
	Cleanup(func())
}

// NewConfig creates a new instance of Config. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewConfig(t mockConstructorTestingTNewConfig) *Config {
	mock := &Config{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
