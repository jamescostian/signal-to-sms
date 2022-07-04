// Code generated by mockery v2.9.4. DO NOT EDIT.

package frameio

import (
	signal "github.com/jamescostian/signal-to-sms/utils/proto/signal"
	mock "github.com/stretchr/testify/mock"
)

// MockFrameWriteCloser is an autogenerated mock type for the FrameWriteCloser type
type MockFrameWriteCloser struct {
	mock.Mock
}

// Close provides a mock function with given fields:
func (_m *MockFrameWriteCloser) Close() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Write provides a mock function with given fields: backupFrame, blob
func (_m *MockFrameWriteCloser) Write(backupFrame *signal.BackupFrame, blob []byte) error {
	ret := _m.Called(backupFrame, blob)

	var r0 error
	if rf, ok := ret.Get(0).(func(*signal.BackupFrame, []byte) error); ok {
		r0 = rf(backupFrame, blob)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}