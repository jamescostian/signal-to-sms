// Code generated by mockery v2.9.4. DO NOT EDIT.

package frameio

import (
	mock "github.com/stretchr/testify/mock"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
)

// MockEncoder is an autogenerated mock type for the Encoder type
type MockEncoder struct {
	mock.Mock
}

// Marshal provides a mock function with given fields: m
func (_m *MockEncoder) Marshal(m protoreflect.ProtoMessage) ([]byte, error) {
	ret := _m.Called(m)

	var r0 []byte
	if rf, ok := ret.Get(0).(func(protoreflect.ProtoMessage) []byte); ok {
		r0 = rf(m)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(protoreflect.ProtoMessage) error); ok {
		r1 = rf(m)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Unmarshal provides a mock function with given fields: b, m
func (_m *MockEncoder) Unmarshal(b []byte, m protoreflect.ProtoMessage) error {
	ret := _m.Called(b, m)

	var r0 error
	if rf, ok := ret.Get(0).(func([]byte, protoreflect.ProtoMessage) error); ok {
		r0 = rf(b, m)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
