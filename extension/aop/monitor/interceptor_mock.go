// Code generated by mockery v2.12.2. DO NOT EDIT.

package monitor

import (
	mock "github.com/stretchr/testify/mock"

	aop "github.com/alibaba/ioc-golang/aop"

	testing "testing"
)

// mockInterceptorImplIOCInterface is an autogenerated mock type for the interceptorImplIOCInterface type
type mockInterceptorImplIOCInterface struct {
	mock.Mock
}

// AfterInvoke provides a mock function with given fields: ctx
func (_m *mockInterceptorImplIOCInterface) AfterInvoke(ctx *aop.InvocationContext) {
	_m.Called(ctx)
}

// BeforeInvoke provides a mock function with given fields: ctx
func (_m *mockInterceptorImplIOCInterface) BeforeInvoke(ctx *aop.InvocationContext) {
	_m.Called(ctx)
}

// Monitor provides a mock function with given fields: monitorCtx
func (_m *mockInterceptorImplIOCInterface) Monitor(monitorCtx contextIOCInterface) {
	_m.Called(monitorCtx)
}

// StopMonitor provides a mock function with given fields:
func (_m *mockInterceptorImplIOCInterface) StopMonitor() {
	_m.Called()
}

// newMockInterceptorImplIOCInterface creates a new instance of mockInterceptorImplIOCInterface. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func newMockInterceptorImplIOCInterface(t testing.TB) *mockInterceptorImplIOCInterface {
	mock := &mockInterceptorImplIOCInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
