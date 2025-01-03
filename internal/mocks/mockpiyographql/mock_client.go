// Code generated by MockGen. DO NOT EDIT.
// Source: ../../adapters/secondary/piyographql/client.go

// Package mockpiyographql is a generated GoMock package.
package mockpiyographql

import (
	context "context"
	reflect "reflect"

	models "github.com/NishimuraTakuya-nt/go-rest-clean-plane-chi/internal/core/domain/models"
	gomock "github.com/golang/mock/gomock"
)

// MockClient is a mock of Client interface.
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient.
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance.
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// GetSample mocks base method.
func (m *MockClient) GetSample(ctx context.Context, id string) (*models.Sample, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSample", ctx, id)
	ret0, _ := ret[0].(*models.Sample)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSample indicates an expected call of GetSample.
func (mr *MockClientMockRecorder) GetSample(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSample", reflect.TypeOf((*MockClient)(nil).GetSample), ctx, id)
}

// ListSample mocks base method.
func (m *MockClient) ListSample(ctx context.Context, offset, limit *int) ([]models.Sample, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListSample", ctx, offset, limit)
	ret0, _ := ret[0].([]models.Sample)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListSample indicates an expected call of ListSample.
func (mr *MockClientMockRecorder) ListSample(ctx, offset, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListSample", reflect.TypeOf((*MockClient)(nil).ListSample), ctx, offset, limit)
}
