// Code generated by MockGen. DO NOT EDIT.
// Source: email.go

// Package mock_email is a generated GoMock package.
package mock_email

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockEmailSender is a mock of EmailSender interface.
type MockEmailSender struct {
	ctrl     *gomock.Controller
	recorder *MockEmailSenderMockRecorder
}

// MockEmailSenderMockRecorder is the mock recorder for MockEmailSender.
type MockEmailSenderMockRecorder struct {
	mock *MockEmailSender
}

// NewMockEmailSender creates a new mock instance.
func NewMockEmailSender(ctrl *gomock.Controller) *MockEmailSender {
	mock := &MockEmailSender{ctrl: ctrl}
	mock.recorder = &MockEmailSenderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEmailSender) EXPECT() *MockEmailSenderMockRecorder {
	return m.recorder
}

// Send mocks base method.
func (m *MockEmailSender) Send(templatePath, emailTo, subject string, data interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", templatePath, emailTo, subject, data)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send.
func (mr *MockEmailSenderMockRecorder) Send(templatePath, emailTo, subject, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockEmailSender)(nil).Send), templatePath, emailTo, subject, data)
}
