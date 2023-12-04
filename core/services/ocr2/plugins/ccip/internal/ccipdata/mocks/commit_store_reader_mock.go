// Code generated by mockery v2.35.4. DO NOT EDIT.

package mocks

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/commit_store"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/internal/ccipdata/offramp"

	"github.com/stretchr/testify/mock"

	"github.com/smartcontractkit/chainlink/v2/core/services/pg"

	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ccip/prices"
)

// CommitStoreReader is an autogenerated mock type for the CommitStoreReader type
type CommitStoreReader struct {
	mock.Mock
}

// ChangeConfig provides a mock function with given fields: onchainConfig, offchainConfig
func (_m *CommitStoreReader) ChangeConfig(onchainConfig []byte, offchainConfig []byte) (common.Address, error) {
	ret := _m.Called(onchainConfig, offchainConfig)

	var r0 common.Address
	var r1 error
	if rf, ok := ret.Get(0).(func([]byte, []byte) (common.Address, error)); ok {
		return rf(onchainConfig, offchainConfig)
	}
	if rf, ok := ret.Get(0).(func([]byte, []byte) common.Address); ok {
		r0 = rf(onchainConfig, offchainConfig)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(common.Address)
		}
	}

	if rf, ok := ret.Get(1).(func([]byte, []byte) error); ok {
		r1 = rf(onchainConfig, offchainConfig)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Close provides a mock function with given fields: qopts
func (_m *CommitStoreReader) Close(qopts ...pg.QOpt) error {
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

// DecodeCommitReport provides a mock function with given fields: report
func (_m *CommitStoreReader) DecodeCommitReport(report []byte) (commit_store.CommitStoreReport, error) {
	ret := _m.Called(report)

	var r0 commit_store.CommitStoreReport
	var r1 error
	if rf, ok := ret.Get(0).(func([]byte) (commit_store.CommitStoreReport, error)); ok {
		return rf(report)
	}
	if rf, ok := ret.Get(0).(func([]byte) commit_store.CommitStoreReport); ok {
		r0 = rf(report)
	} else {
		r0 = ret.Get(0).(commit_store.CommitStoreReport)
	}

	if rf, ok := ret.Get(1).(func([]byte) error); ok {
		r1 = rf(report)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// EncodeCommitReport provides a mock function with given fields: report
func (_m *CommitStoreReader) EncodeCommitReport(report commit_store.CommitStoreReport) ([]byte, error) {
	ret := _m.Called(report)

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(commit_store.CommitStoreReport) ([]byte, error)); ok {
		return rf(report)
	}
	if rf, ok := ret.Get(0).(func(commit_store.CommitStoreReport) []byte); ok {
		r0 = rf(report)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(commit_store.CommitStoreReport) error); ok {
		r1 = rf(report)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GasPriceEstimator provides a mock function with given fields:
func (_m *CommitStoreReader) GasPriceEstimator() prices.GasPriceEstimatorCommit {
	ret := _m.Called()

	var r0 prices.GasPriceEstimatorCommit
	if rf, ok := ret.Get(0).(func() prices.GasPriceEstimatorCommit); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(prices.GasPriceEstimatorCommit)
		}
	}

	return r0
}

// GetAcceptedCommitReportsGteTimestamp provides a mock function with given fields: ctx, ts, confs
func (_m *CommitStoreReader) GetAcceptedCommitReportsGteTimestamp(ctx context.Context, ts time.Time, confs int) ([]ccipdata.Event[commit_store.CommitStoreReport], error) {
	ret := _m.Called(ctx, ts, confs)

	var r0 []ccipdata.Event[commit_store.CommitStoreReport]
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, time.Time, int) ([]ccipdata.Event[commit_store.CommitStoreReport], error)); ok {
		return rf(ctx, ts, confs)
	}
	if rf, ok := ret.Get(0).(func(context.Context, time.Time, int) []ccipdata.Event[commit_store.CommitStoreReport]); ok {
		r0 = rf(ctx, ts, confs)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]ccipdata.Event[commit_store.CommitStoreReport])
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, time.Time, int) error); ok {
		r1 = rf(ctx, ts, confs)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetCommitReportMatchingSeqNum provides a mock function with given fields: ctx, seqNum, confs
func (_m *CommitStoreReader) GetCommitReportMatchingSeqNum(ctx context.Context, seqNum uint64, confs int) ([]ccipdata.Event[commit_store.CommitStoreReport], error) {
	ret := _m.Called(ctx, seqNum, confs)

	var r0 []ccipdata.Event[commit_store.CommitStoreReport]
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint64, int) ([]ccipdata.Event[commit_store.CommitStoreReport], error)); ok {
		return rf(ctx, seqNum, confs)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint64, int) []ccipdata.Event[commit_store.CommitStoreReport]); ok {
		r0 = rf(ctx, seqNum, confs)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]ccipdata.Event[commit_store.CommitStoreReport])
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint64, int) error); ok {
		r1 = rf(ctx, seqNum, confs)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetCommitStoreStaticConfig provides a mock function with given fields: ctx
func (_m *CommitStoreReader) GetCommitStoreStaticConfig(ctx context.Context) (commit_store.CommitStoreStaticConfig, error) {
	ret := _m.Called(ctx)

	var r0 commit_store.CommitStoreStaticConfig
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (commit_store.CommitStoreStaticConfig, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) commit_store.CommitStoreStaticConfig); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(commit_store.CommitStoreStaticConfig)
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetExpectedNextSequenceNumber provides a mock function with given fields: _a0
func (_m *CommitStoreReader) GetExpectedNextSequenceNumber(_a0 context.Context) (uint64, error) {
	ret := _m.Called(_a0)

	var r0 uint64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (uint64, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(context.Context) uint64); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetLatestPriceEpochAndRound provides a mock function with given fields: _a0
func (_m *CommitStoreReader) GetLatestPriceEpochAndRound(_a0 context.Context) (uint64, error) {
	ret := _m.Called(_a0)

	var r0 uint64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (uint64, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(context.Context) uint64); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(uint64)
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsBlessed provides a mock function with given fields: ctx, root
func (_m *CommitStoreReader) IsBlessed(ctx context.Context, root [32]byte) (bool, error) {
	ret := _m.Called(ctx, root)

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, [32]byte) (bool, error)); ok {
		return rf(ctx, root)
	}
	if rf, ok := ret.Get(0).(func(context.Context, [32]byte) bool); ok {
		r0 = rf(ctx, root)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, [32]byte) error); ok {
		r1 = rf(ctx, root)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsDown provides a mock function with given fields: ctx
func (_m *CommitStoreReader) IsDown(ctx context.Context) (bool, error) {
	ret := _m.Called(ctx)

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (bool, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) bool); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// OffchainConfig provides a mock function with given fields:
func (_m *CommitStoreReader) OffchainConfig() commit_store.CommitOffchainConfig {
	ret := _m.Called()

	var r0 commit_store.CommitOffchainConfig
	if rf, ok := ret.Get(0).(func() commit_store.CommitOffchainConfig); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(commit_store.CommitOffchainConfig)
	}

	return r0
}

// RegisterFilters provides a mock function with given fields: qopts
func (_m *CommitStoreReader) RegisterFilters(qopts ...pg.QOpt) error {
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

// VerifyExecutionReport provides a mock function with given fields: ctx, report
func (_m *CommitStoreReader) VerifyExecutionReport(ctx context.Context, report offramp.ExecReport) (bool, error) {
	ret := _m.Called(ctx, report)

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, offramp.ExecReport) (bool, error)); ok {
		return rf(ctx, report)
	}
	if rf, ok := ret.Get(0).(func(context.Context, offramp.ExecReport) bool); ok {
		r0 = rf(ctx, report)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, offramp.ExecReport) error); ok {
		r1 = rf(ctx, report)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewCommitStoreReader creates a new instance of CommitStoreReader. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewCommitStoreReader(t interface {
	mock.TestingT
	Cleanup(func())
}) *CommitStoreReader {
	mock := &CommitStoreReader{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
