// Code generated by mockery v2.22.1. DO NOT EDIT.

package mocks

import (
	context "context"

	common "github.com/ethereum/go-ethereum/common"

	logpoller "github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"

	mock "github.com/stretchr/testify/mock"

	pg "github.com/smartcontractkit/chainlink/v2/core/services/pg"

	time "time"
)

// LogPoller is an autogenerated mock type for the LogPoller type
type LogPoller struct {
	mock.Mock
}

// Close provides a mock function with given fields:
func (_m *LogPoller) Close() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetBlocksRange provides a mock function with given fields: ctx, numbers, qopts
func (_m *LogPoller) GetBlocksRange(ctx context.Context, numbers []uint64, qopts ...pg.QOpt) ([]logpoller.LogPollerBlock, error) {
	_va := make([]interface{}, len(qopts))
	for _i := range qopts {
		_va[_i] = qopts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, numbers)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 []logpoller.LogPollerBlock
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, []uint64, ...pg.QOpt) ([]logpoller.LogPollerBlock, error)); ok {
		return rf(ctx, numbers, qopts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, []uint64, ...pg.QOpt) []logpoller.LogPollerBlock); ok {
		r0 = rf(ctx, numbers, qopts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]logpoller.LogPollerBlock)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, []uint64, ...pg.QOpt) error); ok {
		r1 = rf(ctx, numbers, qopts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// HealthReport provides a mock function with given fields:
func (_m *LogPoller) HealthReport() map[string]error {
	ret := _m.Called()

	var r0 map[string]error
	if rf, ok := ret.Get(0).(func() map[string]error); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]error)
		}
	}

	return r0
}

// IndexedLogs provides a mock function with given fields: eventSig, address, topicIndex, topicValues, confs, qopts
func (_m *LogPoller) IndexedLogs(eventSig common.Hash, address common.Address, topicIndex int, topicValues []common.Hash, confs int, qopts ...pg.QOpt) ([]logpoller.Log, error) {
	_va := make([]interface{}, len(qopts))
	for _i := range qopts {
		_va[_i] = qopts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, eventSig, address, topicIndex, topicValues, confs)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 []logpoller.Log
	var r1 error
	if rf, ok := ret.Get(0).(func(common.Hash, common.Address, int, []common.Hash, int, ...pg.QOpt) ([]logpoller.Log, error)); ok {
		return rf(eventSig, address, topicIndex, topicValues, confs, qopts...)
	}
	if rf, ok := ret.Get(0).(func(common.Hash, common.Address, int, []common.Hash, int, ...pg.QOpt) []logpoller.Log); ok {
		r0 = rf(eventSig, address, topicIndex, topicValues, confs, qopts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]logpoller.Log)
		}
	}

	if rf, ok := ret.Get(1).(func(common.Hash, common.Address, int, []common.Hash, int, ...pg.QOpt) error); ok {
		r1 = rf(eventSig, address, topicIndex, topicValues, confs, qopts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IndexedLogsByBlockRange provides a mock function with given fields: start, end, eventSig, address, topicIndex, topicValues, qopts
func (_m *LogPoller) IndexedLogsByBlockRange(start int64, end int64, eventSig common.Hash, address common.Address, topicIndex int, topicValues []common.Hash, qopts ...pg.QOpt) ([]logpoller.Log, error) {
	_va := make([]interface{}, len(qopts))
	for _i := range qopts {
		_va[_i] = qopts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, start, end, eventSig, address, topicIndex, topicValues)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 []logpoller.Log
	var r1 error
	if rf, ok := ret.Get(0).(func(int64, int64, common.Hash, common.Address, int, []common.Hash, ...pg.QOpt) ([]logpoller.Log, error)); ok {
		return rf(start, end, eventSig, address, topicIndex, topicValues, qopts...)
	}
	if rf, ok := ret.Get(0).(func(int64, int64, common.Hash, common.Address, int, []common.Hash, ...pg.QOpt) []logpoller.Log); ok {
		r0 = rf(start, end, eventSig, address, topicIndex, topicValues, qopts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]logpoller.Log)
		}
	}

	if rf, ok := ret.Get(1).(func(int64, int64, common.Hash, common.Address, int, []common.Hash, ...pg.QOpt) error); ok {
		r1 = rf(start, end, eventSig, address, topicIndex, topicValues, qopts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IndexedLogsCreatedAfter provides a mock function with given fields: eventSig, address, topicIndex, topicValues, after, qopts
func (_m *LogPoller) IndexedLogsCreatedAfter(eventSig common.Hash, address common.Address, topicIndex int, topicValues []common.Hash, after time.Time, qopts ...pg.QOpt) ([]logpoller.Log, error) {
	_va := make([]interface{}, len(qopts))
	for _i := range qopts {
		_va[_i] = qopts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, eventSig, address, topicIndex, topicValues, after)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 []logpoller.Log
	var r1 error
	if rf, ok := ret.Get(0).(func(common.Hash, common.Address, int, []common.Hash, time.Time, ...pg.QOpt) ([]logpoller.Log, error)); ok {
		return rf(eventSig, address, topicIndex, topicValues, after, qopts...)
	}
	if rf, ok := ret.Get(0).(func(common.Hash, common.Address, int, []common.Hash, time.Time, ...pg.QOpt) []logpoller.Log); ok {
		r0 = rf(eventSig, address, topicIndex, topicValues, after, qopts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]logpoller.Log)
		}
	}

	if rf, ok := ret.Get(1).(func(common.Hash, common.Address, int, []common.Hash, time.Time, ...pg.QOpt) error); ok {
		r1 = rf(eventSig, address, topicIndex, topicValues, after, qopts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IndexedLogsTopicGreaterThan provides a mock function with given fields: eventSig, address, topicIndex, topicValueMin, confs, qopts
func (_m *LogPoller) IndexedLogsTopicGreaterThan(eventSig common.Hash, address common.Address, topicIndex int, topicValueMin common.Hash, confs int, qopts ...pg.QOpt) ([]logpoller.Log, error) {
	_va := make([]interface{}, len(qopts))
	for _i := range qopts {
		_va[_i] = qopts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, eventSig, address, topicIndex, topicValueMin, confs)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 []logpoller.Log
	var r1 error
	if rf, ok := ret.Get(0).(func(common.Hash, common.Address, int, common.Hash, int, ...pg.QOpt) ([]logpoller.Log, error)); ok {
		return rf(eventSig, address, topicIndex, topicValueMin, confs, qopts...)
	}
	if rf, ok := ret.Get(0).(func(common.Hash, common.Address, int, common.Hash, int, ...pg.QOpt) []logpoller.Log); ok {
		r0 = rf(eventSig, address, topicIndex, topicValueMin, confs, qopts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]logpoller.Log)
		}
	}

	if rf, ok := ret.Get(1).(func(common.Hash, common.Address, int, common.Hash, int, ...pg.QOpt) error); ok {
		r1 = rf(eventSig, address, topicIndex, topicValueMin, confs, qopts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IndexedLogsTopicRange provides a mock function with given fields: eventSig, address, topicIndex, topicValueMin, topicValueMax, confs, qopts
func (_m *LogPoller) IndexedLogsTopicRange(eventSig common.Hash, address common.Address, topicIndex int, topicValueMin common.Hash, topicValueMax common.Hash, confs int, qopts ...pg.QOpt) ([]logpoller.Log, error) {
	_va := make([]interface{}, len(qopts))
	for _i := range qopts {
		_va[_i] = qopts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, eventSig, address, topicIndex, topicValueMin, topicValueMax, confs)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 []logpoller.Log
	var r1 error
	if rf, ok := ret.Get(0).(func(common.Hash, common.Address, int, common.Hash, common.Hash, int, ...pg.QOpt) ([]logpoller.Log, error)); ok {
		return rf(eventSig, address, topicIndex, topicValueMin, topicValueMax, confs, qopts...)
	}
	if rf, ok := ret.Get(0).(func(common.Hash, common.Address, int, common.Hash, common.Hash, int, ...pg.QOpt) []logpoller.Log); ok {
		r0 = rf(eventSig, address, topicIndex, topicValueMin, topicValueMax, confs, qopts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]logpoller.Log)
		}
	}

	if rf, ok := ret.Get(1).(func(common.Hash, common.Address, int, common.Hash, common.Hash, int, ...pg.QOpt) error); ok {
		r1 = rf(eventSig, address, topicIndex, topicValueMin, topicValueMax, confs, qopts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// LatestBlock provides a mock function with given fields: qopts
func (_m *LogPoller) LatestBlock(qopts ...pg.QOpt) (int64, error) {
	_va := make([]interface{}, len(qopts))
	for _i := range qopts {
		_va[_i] = qopts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(...pg.QOpt) (int64, error)); ok {
		return rf(qopts...)
	}
	if rf, ok := ret.Get(0).(func(...pg.QOpt) int64); ok {
		r0 = rf(qopts...)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(...pg.QOpt) error); ok {
		r1 = rf(qopts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// LatestLogByEventSigWithConfs provides a mock function with given fields: eventSig, address, confs, qopts
func (_m *LogPoller) LatestLogByEventSigWithConfs(eventSig common.Hash, address common.Address, confs int, qopts ...pg.QOpt) (*logpoller.Log, error) {
	_va := make([]interface{}, len(qopts))
	for _i := range qopts {
		_va[_i] = qopts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, eventSig, address, confs)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *logpoller.Log
	var r1 error
	if rf, ok := ret.Get(0).(func(common.Hash, common.Address, int, ...pg.QOpt) (*logpoller.Log, error)); ok {
		return rf(eventSig, address, confs, qopts...)
	}
	if rf, ok := ret.Get(0).(func(common.Hash, common.Address, int, ...pg.QOpt) *logpoller.Log); ok {
		r0 = rf(eventSig, address, confs, qopts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*logpoller.Log)
		}
	}

	if rf, ok := ret.Get(1).(func(common.Hash, common.Address, int, ...pg.QOpt) error); ok {
		r1 = rf(eventSig, address, confs, qopts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// LatestLogEventSigsAddrsWithConfs provides a mock function with given fields: fromBlock, eventSigs, addresses, confs, qopts
func (_m *LogPoller) LatestLogEventSigsAddrsWithConfs(fromBlock int64, eventSigs []common.Hash, addresses []common.Address, confs int, qopts ...pg.QOpt) ([]logpoller.Log, error) {
	_va := make([]interface{}, len(qopts))
	for _i := range qopts {
		_va[_i] = qopts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, fromBlock, eventSigs, addresses, confs)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 []logpoller.Log
	var r1 error
	if rf, ok := ret.Get(0).(func(int64, []common.Hash, []common.Address, int, ...pg.QOpt) ([]logpoller.Log, error)); ok {
		return rf(fromBlock, eventSigs, addresses, confs, qopts...)
	}
	if rf, ok := ret.Get(0).(func(int64, []common.Hash, []common.Address, int, ...pg.QOpt) []logpoller.Log); ok {
		r0 = rf(fromBlock, eventSigs, addresses, confs, qopts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]logpoller.Log)
		}
	}

	if rf, ok := ret.Get(1).(func(int64, []common.Hash, []common.Address, int, ...pg.QOpt) error); ok {
		r1 = rf(fromBlock, eventSigs, addresses, confs, qopts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Logs provides a mock function with given fields: start, end, eventSig, address, qopts
func (_m *LogPoller) Logs(start int64, end int64, eventSig common.Hash, address common.Address, qopts ...pg.QOpt) ([]logpoller.Log, error) {
	_va := make([]interface{}, len(qopts))
	for _i := range qopts {
		_va[_i] = qopts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, start, end, eventSig, address)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 []logpoller.Log
	var r1 error
	if rf, ok := ret.Get(0).(func(int64, int64, common.Hash, common.Address, ...pg.QOpt) ([]logpoller.Log, error)); ok {
		return rf(start, end, eventSig, address, qopts...)
	}
	if rf, ok := ret.Get(0).(func(int64, int64, common.Hash, common.Address, ...pg.QOpt) []logpoller.Log); ok {
		r0 = rf(start, end, eventSig, address, qopts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]logpoller.Log)
		}
	}

	if rf, ok := ret.Get(1).(func(int64, int64, common.Hash, common.Address, ...pg.QOpt) error); ok {
		r1 = rf(start, end, eventSig, address, qopts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// LogsCreatedAfter provides a mock function with given fields: eventSig, address, _a2, qopts
func (_m *LogPoller) LogsCreatedAfter(eventSig common.Hash, address common.Address, _a2 time.Time, qopts ...pg.QOpt) ([]logpoller.Log, error) {
	_va := make([]interface{}, len(qopts))
	for _i := range qopts {
		_va[_i] = qopts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, eventSig, address, _a2)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 []logpoller.Log
	var r1 error
	if rf, ok := ret.Get(0).(func(common.Hash, common.Address, time.Time, ...pg.QOpt) ([]logpoller.Log, error)); ok {
		return rf(eventSig, address, _a2, qopts...)
	}
	if rf, ok := ret.Get(0).(func(common.Hash, common.Address, time.Time, ...pg.QOpt) []logpoller.Log); ok {
		r0 = rf(eventSig, address, _a2, qopts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]logpoller.Log)
		}
	}

	if rf, ok := ret.Get(1).(func(common.Hash, common.Address, time.Time, ...pg.QOpt) error); ok {
		r1 = rf(eventSig, address, _a2, qopts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// LogsDataWordGreaterThan provides a mock function with given fields: eventSig, address, wordIndex, wordValueMin, confs, qopts
func (_m *LogPoller) LogsDataWordGreaterThan(eventSig common.Hash, address common.Address, wordIndex int, wordValueMin common.Hash, confs int, qopts ...pg.QOpt) ([]logpoller.Log, error) {
	_va := make([]interface{}, len(qopts))
	for _i := range qopts {
		_va[_i] = qopts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, eventSig, address, wordIndex, wordValueMin, confs)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 []logpoller.Log
	var r1 error
	if rf, ok := ret.Get(0).(func(common.Hash, common.Address, int, common.Hash, int, ...pg.QOpt) ([]logpoller.Log, error)); ok {
		return rf(eventSig, address, wordIndex, wordValueMin, confs, qopts...)
	}
	if rf, ok := ret.Get(0).(func(common.Hash, common.Address, int, common.Hash, int, ...pg.QOpt) []logpoller.Log); ok {
		r0 = rf(eventSig, address, wordIndex, wordValueMin, confs, qopts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]logpoller.Log)
		}
	}

	if rf, ok := ret.Get(1).(func(common.Hash, common.Address, int, common.Hash, int, ...pg.QOpt) error); ok {
		r1 = rf(eventSig, address, wordIndex, wordValueMin, confs, qopts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// LogsDataWordRange provides a mock function with given fields: eventSig, address, wordIndex, wordValueMin, wordValueMax, confs, qopts
func (_m *LogPoller) LogsDataWordRange(eventSig common.Hash, address common.Address, wordIndex int, wordValueMin common.Hash, wordValueMax common.Hash, confs int, qopts ...pg.QOpt) ([]logpoller.Log, error) {
	_va := make([]interface{}, len(qopts))
	for _i := range qopts {
		_va[_i] = qopts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, eventSig, address, wordIndex, wordValueMin, wordValueMax, confs)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 []logpoller.Log
	var r1 error
	if rf, ok := ret.Get(0).(func(common.Hash, common.Address, int, common.Hash, common.Hash, int, ...pg.QOpt) ([]logpoller.Log, error)); ok {
		return rf(eventSig, address, wordIndex, wordValueMin, wordValueMax, confs, qopts...)
	}
	if rf, ok := ret.Get(0).(func(common.Hash, common.Address, int, common.Hash, common.Hash, int, ...pg.QOpt) []logpoller.Log); ok {
		r0 = rf(eventSig, address, wordIndex, wordValueMin, wordValueMax, confs, qopts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]logpoller.Log)
		}
	}

	if rf, ok := ret.Get(1).(func(common.Hash, common.Address, int, common.Hash, common.Hash, int, ...pg.QOpt) error); ok {
		r1 = rf(eventSig, address, wordIndex, wordValueMin, wordValueMax, confs, qopts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// LogsWithSigs provides a mock function with given fields: start, end, eventSigs, address, qopts
func (_m *LogPoller) LogsWithSigs(start int64, end int64, eventSigs []common.Hash, address common.Address, qopts ...pg.QOpt) ([]logpoller.Log, error) {
	_va := make([]interface{}, len(qopts))
	for _i := range qopts {
		_va[_i] = qopts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, start, end, eventSigs, address)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 []logpoller.Log
	var r1 error
	if rf, ok := ret.Get(0).(func(int64, int64, []common.Hash, common.Address, ...pg.QOpt) ([]logpoller.Log, error)); ok {
		return rf(start, end, eventSigs, address, qopts...)
	}
	if rf, ok := ret.Get(0).(func(int64, int64, []common.Hash, common.Address, ...pg.QOpt) []logpoller.Log); ok {
		r0 = rf(start, end, eventSigs, address, qopts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]logpoller.Log)
		}
	}

	if rf, ok := ret.Get(1).(func(int64, int64, []common.Hash, common.Address, ...pg.QOpt) error); ok {
		r1 = rf(start, end, eventSigs, address, qopts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Name provides a mock function with given fields:
func (_m *LogPoller) Name() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Ready provides a mock function with given fields:
func (_m *LogPoller) Ready() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RegisterFilter provides a mock function with given fields: filter
func (_m *LogPoller) RegisterFilter(filter logpoller.Filter) error {
	ret := _m.Called(filter)

	var r0 error
	if rf, ok := ret.Get(0).(func(logpoller.Filter) error); ok {
		r0 = rf(filter)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Replay provides a mock function with given fields: ctx, fromBlock
func (_m *LogPoller) Replay(ctx context.Context, fromBlock int64) error {
	ret := _m.Called(ctx, fromBlock)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) error); ok {
		r0 = rf(ctx, fromBlock)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Start provides a mock function with given fields: _a0
func (_m *LogPoller) Start(_a0 context.Context) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UnregisterFilter provides a mock function with given fields: name, q
func (_m *LogPoller) UnregisterFilter(name string, q pg.Queryer) error {
	ret := _m.Called(name, q)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, pg.Queryer) error); ok {
		r0 = rf(name, q)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewLogPoller interface {
	mock.TestingT
	Cleanup(func())
}

// NewLogPoller creates a new instance of LogPoller. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewLogPoller(t mockConstructorTestingTNewLogPoller) *LogPoller {
	mock := &LogPoller{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
