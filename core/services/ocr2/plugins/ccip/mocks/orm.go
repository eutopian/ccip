// Code generated by mockery v2.8.0. DO NOT EDIT.

package mocks

import (
	big "math/big"

	common "github.com/ethereum/go-ethereum/common"
	ccip "github.com/smartcontractkit/chainlink/core/services/ocr2/plugins/ccip"

	mock "github.com/stretchr/testify/mock"

	pg "github.com/smartcontractkit/chainlink/core/services/pg"
)

// ORM is an autogenerated mock type for the ORM type
type ORM struct {
	mock.Mock
}

// RelayReport provides a mock function with given fields: seqNum, qopts
func (_m *ORM) RelayReport(seqNum *big.Int, qopts ...pg.QOpt) (ccip.RelayReport, error) {
	_va := make([]interface{}, len(qopts))
	for _i := range qopts {
		_va[_i] = qopts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, seqNum)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 ccip.RelayReport
	if rf, ok := ret.Get(0).(func(*big.Int, ...pg.QOpt) ccip.RelayReport); ok {
		r0 = rf(seqNum, qopts...)
	} else {
		r0 = ret.Get(0).(ccip.RelayReport)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*big.Int, ...pg.QOpt) error); ok {
		r1 = rf(seqNum, qopts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Requests provides a mock function with given fields: sourceChainId, destChainId, onRamp, offRamp, minSeqNum, maxSeqNum, status, executor, options, opt
func (_m *ORM) Requests(sourceChainId *big.Int, destChainId *big.Int, onRamp common.Address, offRamp common.Address, minSeqNum *big.Int, maxSeqNum *big.Int, status ccip.RequestStatus, executor *common.Address, options []byte, opt ...pg.QOpt) ([]*ccip.Request, error) {
	_va := make([]interface{}, len(opt))
	for _i := range opt {
		_va[_i] = opt[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, sourceChainId, destChainId, onRamp, offRamp, minSeqNum, maxSeqNum, status, executor, options)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 []*ccip.Request
	if rf, ok := ret.Get(0).(func(*big.Int, *big.Int, common.Address, common.Address, *big.Int, *big.Int, ccip.RequestStatus, *common.Address, []byte, ...pg.QOpt) []*ccip.Request); ok {
		r0 = rf(sourceChainId, destChainId, onRamp, offRamp, minSeqNum, maxSeqNum, status, executor, options, opt...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*ccip.Request)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*big.Int, *big.Int, common.Address, common.Address, *big.Int, *big.Int, ccip.RequestStatus, *common.Address, []byte, ...pg.QOpt) error); ok {
		r1 = rf(sourceChainId, destChainId, onRamp, offRamp, minSeqNum, maxSeqNum, status, executor, options, opt...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ResetExpiredRequests provides a mock function with given fields: sourceChainId, destChainId, onRamp, offRamp, expiryTimeoutSeconds, fromStatus, toStatus, qopts
func (_m *ORM) ResetExpiredRequests(sourceChainId *big.Int, destChainId *big.Int, onRamp common.Address, offRamp common.Address, expiryTimeoutSeconds int, fromStatus ccip.RequestStatus, toStatus ccip.RequestStatus, qopts ...pg.QOpt) error {
	_va := make([]interface{}, len(qopts))
	for _i := range qopts {
		_va[_i] = qopts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, sourceChainId, destChainId, onRamp, offRamp, expiryTimeoutSeconds, fromStatus, toStatus)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(*big.Int, *big.Int, common.Address, common.Address, int, ccip.RequestStatus, ccip.RequestStatus, ...pg.QOpt) error); ok {
		r0 = rf(sourceChainId, destChainId, onRamp, offRamp, expiryTimeoutSeconds, fromStatus, toStatus, qopts...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SaveRelayReport provides a mock function with given fields: report, qopts
func (_m *ORM) SaveRelayReport(report ccip.RelayReport, qopts ...pg.QOpt) error {
	_va := make([]interface{}, len(qopts))
	for _i := range qopts {
		_va[_i] = qopts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, report)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(ccip.RelayReport, ...pg.QOpt) error); ok {
		r0 = rf(report, qopts...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SaveRequest provides a mock function with given fields: request, qopts
func (_m *ORM) SaveRequest(request *ccip.Request, qopts ...pg.QOpt) error {
	_va := make([]interface{}, len(qopts))
	for _i := range qopts {
		_va[_i] = qopts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, request)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(*ccip.Request, ...pg.QOpt) error); ok {
		r0 = rf(request, qopts...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateRequestSetStatus provides a mock function with given fields: sourceChainId, destChainId, onRamp, offRamp, seqNums, status, qopts
func (_m *ORM) UpdateRequestSetStatus(sourceChainId *big.Int, destChainId *big.Int, onRamp common.Address, offRamp common.Address, seqNums []*big.Int, status ccip.RequestStatus, qopts ...pg.QOpt) error {
	_va := make([]interface{}, len(qopts))
	for _i := range qopts {
		_va[_i] = qopts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, sourceChainId, destChainId, onRamp, offRamp, seqNums, status)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(*big.Int, *big.Int, common.Address, common.Address, []*big.Int, ccip.RequestStatus, ...pg.QOpt) error); ok {
		r0 = rf(sourceChainId, destChainId, onRamp, offRamp, seqNums, status, qopts...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateRequestStatus provides a mock function with given fields: sourceChainId, destChainId, onRamp, offRamp, minSeqNum, maxSeqNum, status, qopts
func (_m *ORM) UpdateRequestStatus(sourceChainId *big.Int, destChainId *big.Int, onRamp common.Address, offRamp common.Address, minSeqNum *big.Int, maxSeqNum *big.Int, status ccip.RequestStatus, qopts ...pg.QOpt) error {
	_va := make([]interface{}, len(qopts))
	for _i := range qopts {
		_va[_i] = qopts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, sourceChainId, destChainId, onRamp, offRamp, minSeqNum, maxSeqNum, status)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 error
	if rf, ok := ret.Get(0).(func(*big.Int, *big.Int, common.Address, common.Address, *big.Int, *big.Int, ccip.RequestStatus, ...pg.QOpt) error); ok {
		r0 = rf(sourceChainId, destChainId, onRamp, offRamp, minSeqNum, maxSeqNum, status, qopts...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
