// Code generated by mockery v2.35.4. DO NOT EDIT.

package mocks

import (
	common "github.com/ethereum/go-ethereum/common"
	fluxmonitorv2 "github.com/smartcontractkit/chainlink/v2/core/services/fluxmonitorv2"
	mock "github.com/stretchr/testify/mock"

	pg "github.com/smartcontractkit/chainlink/v2/core/services/pg"
)

// ORM is an autogenerated mock type for the ORM type
type ORM struct {
	mock.Mock
}

// CountFluxMonitorRoundStats provides a mock function with given fields:
func (_m *ORM) CountFluxMonitorRoundStats() (int, error) {
	ret := _m.Called()

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func() (int, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() int); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateEthTransaction provides a mock function with given fields: fromAddress, toAddress, payload, gasLimit, idempotencyKey
func (_m *ORM) CreateEthTransaction(fromAddress common.Address, toAddress common.Address, payload []byte, gasLimit uint32, idempotencyKey *string) error {
	ret := _m.Called(fromAddress, toAddress, payload, gasLimit, idempotencyKey)

	var r0 error
	if rf, ok := ret.Get(0).(func(common.Address, common.Address, []byte, uint32, *string) error); ok {
		r0 = rf(fromAddress, toAddress, payload, gasLimit, idempotencyKey)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteFluxMonitorRoundsBackThrough provides a mock function with given fields: aggregator, roundID
func (_m *ORM) DeleteFluxMonitorRoundsBackThrough(aggregator common.Address, roundID uint32) error {
	ret := _m.Called(aggregator, roundID)

	var r0 error
	if rf, ok := ret.Get(0).(func(common.Address, uint32) error); ok {
		r0 = rf(aggregator, roundID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindOrCreateFluxMonitorRoundStats provides a mock function with given fields: aggregator, roundID, newRoundLogs
func (_m *ORM) FindOrCreateFluxMonitorRoundStats(aggregator common.Address, roundID uint32, newRoundLogs uint) (fluxmonitorv2.FluxMonitorRoundStatsV2, error) {
	ret := _m.Called(aggregator, roundID, newRoundLogs)

	var r0 fluxmonitorv2.FluxMonitorRoundStatsV2
	var r1 error
	if rf, ok := ret.Get(0).(func(common.Address, uint32, uint) (fluxmonitorv2.FluxMonitorRoundStatsV2, error)); ok {
		return rf(aggregator, roundID, newRoundLogs)
	}
	if rf, ok := ret.Get(0).(func(common.Address, uint32, uint) fluxmonitorv2.FluxMonitorRoundStatsV2); ok {
		r0 = rf(aggregator, roundID, newRoundLogs)
	} else {
		r0 = ret.Get(0).(fluxmonitorv2.FluxMonitorRoundStatsV2)
	}

	if rf, ok := ret.Get(1).(func(common.Address, uint32, uint) error); ok {
		r1 = rf(aggregator, roundID, newRoundLogs)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MostRecentFluxMonitorRoundID provides a mock function with given fields: aggregator
func (_m *ORM) MostRecentFluxMonitorRoundID(aggregator common.Address) (uint32, error) {
	ret := _m.Called(aggregator)

	var r0 uint32
	var r1 error
	if rf, ok := ret.Get(0).(func(common.Address) (uint32, error)); ok {
		return rf(aggregator)
	}
	if rf, ok := ret.Get(0).(func(common.Address) uint32); ok {
		r0 = rf(aggregator)
	} else {
		r0 = ret.Get(0).(uint32)
	}

	if rf, ok := ret.Get(1).(func(common.Address) error); ok {
		r1 = rf(aggregator)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateFluxMonitorRoundStats provides a mock function with given fields: aggregator, roundID, runID, newRoundLogsAddition, qopts
func (_m *ORM) UpdateFluxMonitorRoundStats(aggregator common.Address, roundID uint32, runID int64, newRoundLogsAddition uint, qopts ...pg.QOpt) error {
	_va := make([]interface{}, len(qopts))
	for _i := range qopts {
		_va[_i] = qopts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, aggregator, roundID, runID, newRoundLogsAddition)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(common.Address, uint32, int64, uint, ...pg.QOpt) error); ok {
		r0 = rf(aggregator, roundID, runID, newRoundLogsAddition, qopts...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewORM creates a new instance of ORM. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewORM(t interface {
	mock.TestingT
	Cleanup(func())
}) *ORM {
	mock := &ORM{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
