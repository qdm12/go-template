// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/qdm12/go-template/internal/metrics (interfaces: Metrics)

// Package mock_metrics is a generated GoMock package.
package mock_metrics

import (
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
)

// MockMetrics is a mock of Metrics interface.
type MockMetrics struct {
	ctrl     *gomock.Controller
	recorder *MockMetricsMockRecorder
}

// MockMetricsMockRecorder is the mock recorder for MockMetrics.
type MockMetricsMockRecorder struct {
	mock *MockMetrics
}

// NewMockMetrics creates a new mock instance.
func NewMockMetrics(ctrl *gomock.Controller) *MockMetrics {
	mock := &MockMetrics{ctrl: ctrl}
	mock.recorder = &MockMetricsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMetrics) EXPECT() *MockMetricsMockRecorder {
	return m.recorder
}

// InflightRequestsGaugeAdd mocks base method.
func (m *MockMetrics) InflightRequestsGaugeAdd(arg0 int) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "InflightRequestsGaugeAdd", arg0)
}

// InflightRequestsGaugeAdd indicates an expected call of InflightRequestsGaugeAdd.
func (mr *MockMetricsMockRecorder) InflightRequestsGaugeAdd(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InflightRequestsGaugeAdd", reflect.TypeOf((*MockMetrics)(nil).InflightRequestsGaugeAdd), arg0)
}

// RequestCountInc mocks base method.
func (m *MockMetrics) RequestCountInc(arg0 string, arg1 int) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RequestCountInc", arg0, arg1)
}

// RequestCountInc indicates an expected call of RequestCountInc.
func (mr *MockMetricsMockRecorder) RequestCountInc(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RequestCountInc", reflect.TypeOf((*MockMetrics)(nil).RequestCountInc), arg0, arg1)
}

// ResponseBytesCountAdd mocks base method.
func (m *MockMetrics) ResponseBytesCountAdd(arg0 string, arg1, arg2 int) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ResponseBytesCountAdd", arg0, arg1, arg2)
}

// ResponseBytesCountAdd indicates an expected call of ResponseBytesCountAdd.
func (mr *MockMetricsMockRecorder) ResponseBytesCountAdd(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResponseBytesCountAdd", reflect.TypeOf((*MockMetrics)(nil).ResponseBytesCountAdd), arg0, arg1, arg2)
}

// ResponseTimeHistogramObserve mocks base method.
func (m *MockMetrics) ResponseTimeHistogramObserve(arg0 string, arg1 int, arg2 time.Duration) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ResponseTimeHistogramObserve", arg0, arg1, arg2)
}

// ResponseTimeHistogramObserve indicates an expected call of ResponseTimeHistogramObserve.
func (mr *MockMetricsMockRecorder) ResponseTimeHistogramObserve(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResponseTimeHistogramObserve", reflect.TypeOf((*MockMetrics)(nil).ResponseTimeHistogramObserve), arg0, arg1, arg2)
}
