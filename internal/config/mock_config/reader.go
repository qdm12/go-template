// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/qdm12/go-template/internal/config (interfaces: Reader)

// Package mock_config is a generated GoMock package.
package mock_config

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	config "github.com/qdm12/go-template/internal/config"
)

// MockReader is a mock of Reader interface.
type MockReader struct {
	ctrl     *gomock.Controller
	recorder *MockReaderMockRecorder
}

// MockReaderMockRecorder is the mock recorder for MockReader.
type MockReaderMockRecorder struct {
	mock *MockReader
}

// NewMockReader creates a new mock instance.
func NewMockReader(ctrl *gomock.Controller) *MockReader {
	mock := &MockReader{ctrl: ctrl}
	mock.recorder = &MockReaderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockReader) EXPECT() *MockReaderMockRecorder {
	return m.recorder
}

// ReadConfig mocks base method.
func (m *MockReader) ReadConfig() (config.Config, []string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadConfig")
	ret0, _ := ret[0].(config.Config)
	ret1, _ := ret[1].([]string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ReadConfig indicates an expected call of ReadConfig.
func (mr *MockReaderMockRecorder) ReadConfig() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadConfig", reflect.TypeOf((*MockReader)(nil).ReadConfig))
}