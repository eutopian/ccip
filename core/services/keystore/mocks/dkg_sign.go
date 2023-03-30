// Code generated by mockery v2.22.1. DO NOT EDIT.

package mocks

import (
	dkgsignkey "github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/dkgsignkey"

	mock "github.com/stretchr/testify/mock"
)

// DKGSign is an autogenerated mock type for the DKGSign type
type DKGSign struct {
	mock.Mock
}

// Add provides a mock function with given fields: key
func (_m *DKGSign) Add(key dkgsignkey.Key) error {
	ret := _m.Called(key)

	var r0 error
	if rf, ok := ret.Get(0).(func(dkgsignkey.Key) error); ok {
		r0 = rf(key)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Create provides a mock function with given fields:
func (_m *DKGSign) Create() (dkgsignkey.Key, error) {
	ret := _m.Called()

	var r0 dkgsignkey.Key
	var r1 error
	if rf, ok := ret.Get(0).(func() (dkgsignkey.Key, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() dkgsignkey.Key); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(dkgsignkey.Key)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: id
func (_m *DKGSign) Delete(id string) (dkgsignkey.Key, error) {
	ret := _m.Called(id)

	var r0 dkgsignkey.Key
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (dkgsignkey.Key, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(string) dkgsignkey.Key); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(dkgsignkey.Key)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// EnsureKey provides a mock function with given fields:
func (_m *DKGSign) EnsureKey() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Export provides a mock function with given fields: id, password
func (_m *DKGSign) Export(id string, password string) ([]byte, error) {
	ret := _m.Called(id, password)

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string) ([]byte, error)); ok {
		return rf(id, password)
	}
	if rf, ok := ret.Get(0).(func(string, string) []byte); ok {
		r0 = rf(id, password)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(id, password)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Get provides a mock function with given fields: id
func (_m *DKGSign) Get(id string) (dkgsignkey.Key, error) {
	ret := _m.Called(id)

	var r0 dkgsignkey.Key
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (dkgsignkey.Key, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(string) dkgsignkey.Key); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(dkgsignkey.Key)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAll provides a mock function with given fields:
func (_m *DKGSign) GetAll() ([]dkgsignkey.Key, error) {
	ret := _m.Called()

	var r0 []dkgsignkey.Key
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]dkgsignkey.Key, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []dkgsignkey.Key); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]dkgsignkey.Key)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Import provides a mock function with given fields: keyJSON, password
func (_m *DKGSign) Import(keyJSON []byte, password string) (dkgsignkey.Key, error) {
	ret := _m.Called(keyJSON, password)

	var r0 dkgsignkey.Key
	var r1 error
	if rf, ok := ret.Get(0).(func([]byte, string) (dkgsignkey.Key, error)); ok {
		return rf(keyJSON, password)
	}
	if rf, ok := ret.Get(0).(func([]byte, string) dkgsignkey.Key); ok {
		r0 = rf(keyJSON, password)
	} else {
		r0 = ret.Get(0).(dkgsignkey.Key)
	}

	if rf, ok := ret.Get(1).(func([]byte, string) error); ok {
		r1 = rf(keyJSON, password)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewDKGSign interface {
	mock.TestingT
	Cleanup(func())
}

// NewDKGSign creates a new instance of DKGSign. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewDKGSign(t mockConstructorTestingTNewDKGSign) *DKGSign {
	mock := &DKGSign{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
