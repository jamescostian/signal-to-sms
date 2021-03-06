// Code generated by mockery v2.9.4. DO NOT EDIT.

package frameio

import (
	signal "github.com/jamescostian/signal-to-sms/utils/proto/signal"
	mock "github.com/stretchr/testify/mock"
)

// MockFrameReader is an autogenerated mock type for the FrameReader type
type MockFrameReader struct {
	mock.Mock
}

// Read provides a mock function with given fields: backupFrameDest, blobDest
func (_m *MockFrameReader) Read(backupFrameDest *signal.BackupFrame, blobDest *[]byte) error {
	ret := _m.Called(backupFrameDest, blobDest)

	var r0 error
	if rf, ok := ret.Get(0).(func(*signal.BackupFrame, *[]byte) error); ok {
		r0 = rf(backupFrameDest, blobDest)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
