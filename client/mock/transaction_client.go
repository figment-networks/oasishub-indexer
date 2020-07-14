// Code generated by MockGen. DO NOT EDIT.
// Source: transaction_client.go

// Package mock_client is a generated GoMock package.
package mock_client

import (
	transactionpb "github.com/figment-networks/oasis-rpc-proxy/grpc/transaction/transactionpb"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockTransactionClient is a mock of TransactionClient interface
type MockTransactionClient struct {
	ctrl     *gomock.Controller
	recorder *MockTransactionClientMockRecorder
}

// MockTransactionClientMockRecorder is the mock recorder for MockTransactionClient
type MockTransactionClientMockRecorder struct {
	mock *MockTransactionClient
}

// NewMockTransactionClient creates a new mock instance
func NewMockTransactionClient(ctrl *gomock.Controller) *MockTransactionClient {
	mock := &MockTransactionClient{ctrl: ctrl}
	mock.recorder = &MockTransactionClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockTransactionClient) EXPECT() *MockTransactionClientMockRecorder {
	return m.recorder
}

// GetByHeight mocks base method
func (m *MockTransactionClient) GetByHeight(arg0 int64) (*transactionpb.GetByHeightResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByHeight", arg0)
	ret0, _ := ret[0].(*transactionpb.GetByHeightResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByHeight indicates an expected call of GetByHeight
func (mr *MockTransactionClientMockRecorder) GetByHeight(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByHeight", reflect.TypeOf((*MockTransactionClient)(nil).GetByHeight), arg0)
}
