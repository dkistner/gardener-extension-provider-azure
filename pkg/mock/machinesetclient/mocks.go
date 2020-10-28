// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/gardener/gardener-extension-provider-azure/pkg/internal/machinesetclient (interfaces: MachineSetClient)

// Package machinesetclient is a generated GoMock package.
package machinesetclient

import (
	compute "github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2019-12-01/compute"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockMachineSetClient is a mock of MachineSetClient interface
type MockMachineSetClient struct {
	ctrl     *gomock.Controller
	recorder *MockMachineSetClientMockRecorder
}

// MockMachineSetClientMockRecorder is the mock recorder for MockMachineSetClient
type MockMachineSetClientMockRecorder struct {
	mock *MockMachineSetClient
}

// NewMockMachineSetClient creates a new mock instance
func NewMockMachineSetClient(ctrl *gomock.Controller) *MockMachineSetClient {
	mock := &MockMachineSetClient{ctrl: ctrl}
	mock.recorder = &MockMachineSetClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockMachineSetClient) EXPECT() *MockMachineSetClientMockRecorder {
	return m.recorder
}

// CreateVMO mocks base method
func (m *MockMachineSetClient) CreateVMO() (*compute.VirtualMachineScaleSet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateVMO")
	ret0, _ := ret[0].(*compute.VirtualMachineScaleSet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateVMO indicates an expected call of CreateVMO
func (mr *MockMachineSetClientMockRecorder) CreateVMO() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateVMO", reflect.TypeOf((*MockMachineSetClient)(nil).CreateVMO))
}

// DeleteVMO mocks base method
func (m *MockMachineSetClient) DeleteVMO() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteVMO")
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteVMO indicates an expected call of DeleteVMO
func (mr *MockMachineSetClientMockRecorder) DeleteVMO() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteVMO", reflect.TypeOf((*MockMachineSetClient)(nil).DeleteVMO))
}

// GetVMO mocks base method
func (m *MockMachineSetClient) GetVMO() (*compute.VirtualMachineScaleSet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVMO")
	ret0, _ := ret[0].(*compute.VirtualMachineScaleSet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetVMO indicates an expected call of GetVMO
func (mr *MockMachineSetClientMockRecorder) GetVMO() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVMO", reflect.TypeOf((*MockMachineSetClient)(nil).GetVMO))
}

// IsVMORequired mocks base method
func (m *MockMachineSetClient) IsVMORequired() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsVMORequired")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsVMORequired indicates an expected call of IsVMORequired
func (mr *MockMachineSetClientMockRecorder) IsVMORequired() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsVMORequired", reflect.TypeOf((*MockMachineSetClient)(nil).IsVMORequired))
}

// ListVMOs mocks base method
func (m *MockMachineSetClient) ListVMOs() ([]compute.VirtualMachineScaleSet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListVMOs")
	ret0, _ := ret[0].([]compute.VirtualMachineScaleSet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListVMOs indicates an expected call of ListVMOs
func (mr *MockMachineSetClientMockRecorder) ListVMOs() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListVMOs", reflect.TypeOf((*MockMachineSetClient)(nil).ListVMOs))
}