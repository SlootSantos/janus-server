// Code generated by MockGen. DO NOT EDIT.
// Source: cdn.go

// Package cdn is a generated GoMock package.
package cdn

import (
	acm "github.com/aws/aws-sdk-go/service/acm"
	cloudfront "github.com/aws/aws-sdk-go/service/cloudfront"
	route53 "github.com/aws/aws-sdk-go/service/route53"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// Mockcdnandler is a mock of cdnandler interface
type Mockcdnandler struct {
	ctrl     *gomock.Controller
	recorder *MockcdnandlerMockRecorder
}

// MockcdnandlerMockRecorder is the mock recorder for Mockcdnandler
type MockcdnandlerMockRecorder struct {
	mock *Mockcdnandler
}

// NewMockcdnandler creates a new mock instance
func NewMockcdnandler(ctrl *gomock.Controller) *Mockcdnandler {
	mock := &Mockcdnandler{ctrl: ctrl}
	mock.recorder = &MockcdnandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *Mockcdnandler) EXPECT() *MockcdnandlerMockRecorder {
	return m.recorder
}

// CreateDistribution mocks base method
func (m *Mockcdnandler) CreateDistribution(arg0 *cloudfront.CreateDistributionInput) (*cloudfront.CreateDistributionOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateDistribution", arg0)
	ret0, _ := ret[0].(*cloudfront.CreateDistributionOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateDistribution indicates an expected call of CreateDistribution
func (mr *MockcdnandlerMockRecorder) CreateDistribution(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateDistribution", reflect.TypeOf((*Mockcdnandler)(nil).CreateDistribution), arg0)
}

// DeleteDistribution mocks base method
func (m *Mockcdnandler) DeleteDistribution(arg0 *cloudfront.DeleteDistributionInput) (*cloudfront.DeleteDistributionOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteDistribution", arg0)
	ret0, _ := ret[0].(*cloudfront.DeleteDistributionOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteDistribution indicates an expected call of DeleteDistribution
func (mr *MockcdnandlerMockRecorder) DeleteDistribution(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteDistribution", reflect.TypeOf((*Mockcdnandler)(nil).DeleteDistribution), arg0)
}

// GetDistribution mocks base method
func (m *Mockcdnandler) GetDistribution(arg0 *cloudfront.GetDistributionInput) (*cloudfront.GetDistributionOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDistribution", arg0)
	ret0, _ := ret[0].(*cloudfront.GetDistributionOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDistribution indicates an expected call of GetDistribution
func (mr *MockcdnandlerMockRecorder) GetDistribution(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDistribution", reflect.TypeOf((*Mockcdnandler)(nil).GetDistribution), arg0)
}

// UpdateDistribution mocks base method
func (m *Mockcdnandler) UpdateDistribution(arg0 *cloudfront.UpdateDistributionInput) (*cloudfront.UpdateDistributionOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateDistribution", arg0)
	ret0, _ := ret[0].(*cloudfront.UpdateDistributionOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateDistribution indicates an expected call of UpdateDistribution
func (mr *MockcdnandlerMockRecorder) UpdateDistribution(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateDistribution", reflect.TypeOf((*Mockcdnandler)(nil).UpdateDistribution), arg0)
}

// GetCloudFrontOriginAccessIdentity mocks base method
func (m *Mockcdnandler) GetCloudFrontOriginAccessIdentity(arg0 *cloudfront.GetCloudFrontOriginAccessIdentityInput) (*cloudfront.GetCloudFrontOriginAccessIdentityOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCloudFrontOriginAccessIdentity", arg0)
	ret0, _ := ret[0].(*cloudfront.GetCloudFrontOriginAccessIdentityOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCloudFrontOriginAccessIdentity indicates an expected call of GetCloudFrontOriginAccessIdentity
func (mr *MockcdnandlerMockRecorder) GetCloudFrontOriginAccessIdentity(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCloudFrontOriginAccessIdentity", reflect.TypeOf((*Mockcdnandler)(nil).GetCloudFrontOriginAccessIdentity), arg0)
}

// CreateCloudFrontOriginAccessIdentity mocks base method
func (m *Mockcdnandler) CreateCloudFrontOriginAccessIdentity(arg0 *cloudfront.CreateCloudFrontOriginAccessIdentityInput) (*cloudfront.CreateCloudFrontOriginAccessIdentityOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCloudFrontOriginAccessIdentity", arg0)
	ret0, _ := ret[0].(*cloudfront.CreateCloudFrontOriginAccessIdentityOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateCloudFrontOriginAccessIdentity indicates an expected call of CreateCloudFrontOriginAccessIdentity
func (mr *MockcdnandlerMockRecorder) CreateCloudFrontOriginAccessIdentity(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCloudFrontOriginAccessIdentity", reflect.TypeOf((*Mockcdnandler)(nil).CreateCloudFrontOriginAccessIdentity), arg0)
}

// Mockdnshandler is a mock of dnshandler interface
type Mockdnshandler struct {
	ctrl     *gomock.Controller
	recorder *MockdnshandlerMockRecorder
}

// MockdnshandlerMockRecorder is the mock recorder for Mockdnshandler
type MockdnshandlerMockRecorder struct {
	mock *Mockdnshandler
}

// NewMockdnshandler creates a new mock instance
func NewMockdnshandler(ctrl *gomock.Controller) *Mockdnshandler {
	mock := &Mockdnshandler{ctrl: ctrl}
	mock.recorder = &MockdnshandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *Mockdnshandler) EXPECT() *MockdnshandlerMockRecorder {
	return m.recorder
}

// ChangeResourceRecordSets mocks base method
func (m *Mockdnshandler) ChangeResourceRecordSets(arg0 *route53.ChangeResourceRecordSetsInput) (*route53.ChangeResourceRecordSetsOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangeResourceRecordSets", arg0)
	ret0, _ := ret[0].(*route53.ChangeResourceRecordSetsOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ChangeResourceRecordSets indicates an expected call of ChangeResourceRecordSets
func (mr *MockdnshandlerMockRecorder) ChangeResourceRecordSets(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangeResourceRecordSets", reflect.TypeOf((*Mockdnshandler)(nil).ChangeResourceRecordSets), arg0)
}

// MockcertificateHandler is a mock of certificateHandler interface
type MockcertificateHandler struct {
	ctrl     *gomock.Controller
	recorder *MockcertificateHandlerMockRecorder
}

// MockcertificateHandlerMockRecorder is the mock recorder for MockcertificateHandler
type MockcertificateHandlerMockRecorder struct {
	mock *MockcertificateHandler
}

// NewMockcertificateHandler creates a new mock instance
func NewMockcertificateHandler(ctrl *gomock.Controller) *MockcertificateHandler {
	mock := &MockcertificateHandler{ctrl: ctrl}
	mock.recorder = &MockcertificateHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockcertificateHandler) EXPECT() *MockcertificateHandlerMockRecorder {
	return m.recorder
}

// RequestCertificate mocks base method
func (m *MockcertificateHandler) RequestCertificate(arg0 *acm.RequestCertificateInput) (*acm.RequestCertificateOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RequestCertificate", arg0)
	ret0, _ := ret[0].(*acm.RequestCertificateOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RequestCertificate indicates an expected call of RequestCertificate
func (mr *MockcertificateHandlerMockRecorder) RequestCertificate(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RequestCertificate", reflect.TypeOf((*MockcertificateHandler)(nil).RequestCertificate), arg0)
}

// GetCertificate mocks base method
func (m *MockcertificateHandler) GetCertificate(arg0 *acm.GetCertificateInput) (*acm.GetCertificateOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCertificate", arg0)
	ret0, _ := ret[0].(*acm.GetCertificateOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCertificate indicates an expected call of GetCertificate
func (mr *MockcertificateHandlerMockRecorder) GetCertificate(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCertificate", reflect.TypeOf((*MockcertificateHandler)(nil).GetCertificate), arg0)
}

// DescribeCertificate mocks base method
func (m *MockcertificateHandler) DescribeCertificate(arg0 *acm.DescribeCertificateInput) (*acm.DescribeCertificateOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DescribeCertificate", arg0)
	ret0, _ := ret[0].(*acm.DescribeCertificateOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DescribeCertificate indicates an expected call of DescribeCertificate
func (mr *MockcertificateHandlerMockRecorder) DescribeCertificate(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DescribeCertificate", reflect.TypeOf((*MockcertificateHandler)(nil).DescribeCertificate), arg0)
}

// DeleteCertificate mocks base method
func (m *MockcertificateHandler) DeleteCertificate(arg0 *acm.DeleteCertificateInput) (*acm.DeleteCertificateOutput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCertificate", arg0)
	ret0, _ := ret[0].(*acm.DeleteCertificateOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteCertificate indicates an expected call of DeleteCertificate
func (mr *MockcertificateHandlerMockRecorder) DeleteCertificate(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCertificate", reflect.TypeOf((*MockcertificateHandler)(nil).DeleteCertificate), arg0)
}
