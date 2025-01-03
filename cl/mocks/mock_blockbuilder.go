// Code generated by MockGen. DO NOT EDIT.
// Source: redisapp/leaderfollower/leaderfollower.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockBlockBuilder is a mock of BlockBuilder interface.
type MockBlockBuilder struct {
	ctrl     *gomock.Controller
	recorder *MockBlockBuilderMockRecorder
}

// MockBlockBuilderMockRecorder is the mock recorder for MockBlockBuilder.
type MockBlockBuilderMockRecorder struct {
	mock *MockBlockBuilder
}

// NewMockBlockBuilder creates a new mock instance.
func NewMockBlockBuilder(ctrl *gomock.Controller) *MockBlockBuilder {
	mock := &MockBlockBuilder{ctrl: ctrl}
	mock.recorder = &MockBlockBuilderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBlockBuilder) EXPECT() *MockBlockBuilderMockRecorder {
	return m.recorder
}

// FinalizeBlock mocks base method.
func (m *MockBlockBuilder) FinalizeBlock(ctx context.Context, payloadIDStr, executionPayloadStr, msgID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FinalizeBlock", ctx, payloadIDStr, executionPayloadStr, msgID)
	ret0, _ := ret[0].(error)
	return ret0
}

// FinalizeBlock indicates an expected call of FinalizeBlock.
func (mr *MockBlockBuilderMockRecorder) FinalizeBlock(ctx, payloadIDStr, executionPayloadStr, msgID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FinalizeBlock", reflect.TypeOf((*MockBlockBuilder)(nil).FinalizeBlock), ctx, payloadIDStr, executionPayloadStr, msgID)
}

// GetPayload mocks base method.
func (m *MockBlockBuilder) GetPayload(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPayload", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetPayload indicates an expected call of GetPayload.
func (mr *MockBlockBuilderMockRecorder) GetPayload(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPayload", reflect.TypeOf((*MockBlockBuilder)(nil).GetPayload), ctx)
}

// ProcessLastPayload mocks base method.
func (m *MockBlockBuilder) ProcessLastPayload(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProcessLastPayload", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// ProcessLastPayload indicates an expected call of ProcessLastPayload.
func (mr *MockBlockBuilderMockRecorder) ProcessLastPayload(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProcessLastPayload", reflect.TypeOf((*MockBlockBuilder)(nil).ProcessLastPayload), ctx)
}

// SetLastCallTimeToZero mocks base method.
func (m *MockBlockBuilder) SetLastCallTimeToZero() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetLastCallTimeToZero")
}

// SetLastCallTimeToZero indicates an expected call of SetLastCallTimeToZero.
func (mr *MockBlockBuilderMockRecorder) SetLastCallTimeToZero() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetLastCallTimeToZero", reflect.TypeOf((*MockBlockBuilder)(nil).SetLastCallTimeToZero))
}
