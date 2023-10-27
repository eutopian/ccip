// Code generated by mockery v2.35.4. DO NOT EDIT.

package mocks

import (
	p2pkey "github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
	mock "github.com/stretchr/testify/mock"
)

// P2P is an autogenerated mock type for the P2P type
type P2P struct {
	mock.Mock
}

// Add provides a mock function with given fields: key
func (_m *P2P) Add(key p2pkey.KeyV2) error {
	ret := _m.Called(key)

	var r0 error
	if rf, ok := ret.Get(0).(func(p2pkey.KeyV2) error); ok {
		r0 = rf(key)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Create provides a mock function with given fields:
func (_m *P2P) Create() (p2pkey.KeyV2, error) {
	ret := _m.Called()

	var r0 p2pkey.KeyV2
	var r1 error
	if rf, ok := ret.Get(0).(func() (p2pkey.KeyV2, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() p2pkey.KeyV2); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(p2pkey.KeyV2)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: id
func (_m *P2P) Delete(id p2pkey.PeerID) (p2pkey.KeyV2, error) {
	ret := _m.Called(id)

	var r0 p2pkey.KeyV2
	var r1 error
	if rf, ok := ret.Get(0).(func(p2pkey.PeerID) (p2pkey.KeyV2, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(p2pkey.PeerID) p2pkey.KeyV2); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(p2pkey.KeyV2)
	}

	if rf, ok := ret.Get(1).(func(p2pkey.PeerID) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// EnsureKey provides a mock function with given fields:
func (_m *P2P) EnsureKey() error {
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
func (_m *P2P) Export(id p2pkey.PeerID, password string) ([]byte, error) {
	ret := _m.Called(id, password)

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(p2pkey.PeerID, string) ([]byte, error)); ok {
		return rf(id, password)
	}
	if rf, ok := ret.Get(0).(func(p2pkey.PeerID, string) []byte); ok {
		r0 = rf(id, password)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(p2pkey.PeerID, string) error); ok {
		r1 = rf(id, password)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Get provides a mock function with given fields: id
func (_m *P2P) Get(id p2pkey.PeerID) (p2pkey.KeyV2, error) {
	ret := _m.Called(id)

	var r0 p2pkey.KeyV2
	var r1 error
	if rf, ok := ret.Get(0).(func(p2pkey.PeerID) (p2pkey.KeyV2, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(p2pkey.PeerID) p2pkey.KeyV2); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(p2pkey.KeyV2)
	}

	if rf, ok := ret.Get(1).(func(p2pkey.PeerID) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAll provides a mock function with given fields:
func (_m *P2P) GetAll() ([]p2pkey.KeyV2, error) {
	ret := _m.Called()

	var r0 []p2pkey.KeyV2
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]p2pkey.KeyV2, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []p2pkey.KeyV2); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]p2pkey.KeyV2)
		}
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetOrFirst provides a mock function with given fields: id
func (_m *P2P) GetOrFirst(id p2pkey.PeerID) (p2pkey.KeyV2, error) {
	ret := _m.Called(id)

	var r0 p2pkey.KeyV2
	var r1 error
	if rf, ok := ret.Get(0).(func(p2pkey.PeerID) (p2pkey.KeyV2, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(p2pkey.PeerID) p2pkey.KeyV2); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(p2pkey.KeyV2)
	}

	if rf, ok := ret.Get(1).(func(p2pkey.PeerID) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetV1KeysAsV2 provides a mock function with given fields:
func (_m *P2P) GetV1KeysAsV2() ([]p2pkey.KeyV2, error) {
	ret := _m.Called()

	var r0 []p2pkey.KeyV2
	var r1 error
	if rf, ok := ret.Get(0).(func() ([]p2pkey.KeyV2, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() []p2pkey.KeyV2); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]p2pkey.KeyV2)
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
func (_m *P2P) Import(keyJSON []byte, password string) (p2pkey.KeyV2, error) {
	ret := _m.Called(keyJSON, password)

	var r0 p2pkey.KeyV2
	var r1 error
	if rf, ok := ret.Get(0).(func([]byte, string) (p2pkey.KeyV2, error)); ok {
		return rf(keyJSON, password)
	}
	if rf, ok := ret.Get(0).(func([]byte, string) p2pkey.KeyV2); ok {
		r0 = rf(keyJSON, password)
	} else {
		r0 = ret.Get(0).(p2pkey.KeyV2)
	}

	if rf, ok := ret.Get(1).(func([]byte, string) error); ok {
		r1 = rf(keyJSON, password)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewP2P creates a new instance of P2P. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewP2P(t interface {
	mock.TestingT
	Cleanup(func())
}) *P2P {
	mock := &P2P{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
