// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2021 Schneider Electric
//
// SPDX-License-Identifier: Apache-2.0

package test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/edgexfoundry/device-sdk-go/v4/pkg/interfaces/mocks"
	"github.com/edgexfoundry/device-sdk-go/v4/pkg/models"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/clients/logger"
	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/ua"
	"github.com/stretchr/testify/mock"
)

func NewDSMock(t *testing.T) *mocks.DeviceServiceSDK {
	dsMock := mocks.NewDeviceServiceSDK(t)
	logMock := logger.NewMockClient()
	dsMock.On("LoggingClient").Return(logMock).Maybe()
	dsMock.On("AsyncValuesChannel").Return(make(chan *models.AsyncValues, 1)).Maybe()

	return dsMock
}

type ResponseWriterMock struct {
	mock.Mock
}

func (r *ResponseWriterMock) Write([]byte) (int, error) {
	return 0, nil
}

func (r *ResponseWriterMock) Header() http.Header {
	return make(http.Header)
}

func (r *ResponseWriterMock) WriteHeader(statusCode int) {
	fmt.Printf("StatusCode=%d", statusCode)
}

// MockOpcuaClient is a mock of the opcuaClient interface
type MockOpcuaClient struct {
	mock.Mock
}

func (m *MockOpcuaClient) Connect(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockOpcuaClient) Close(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockOpcuaClient) Call(ctx context.Context, req *ua.CallMethodRequest) (*ua.CallMethodResult, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ua.CallMethodResult), args.Error(1)
}

func (m *MockOpcuaClient) Read(ctx context.Context, req *ua.ReadRequest) (*ua.ReadResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ua.ReadResponse), args.Error(1)
}

func (m *MockOpcuaClient) Write(ctx context.Context, req *ua.WriteRequest) (*ua.WriteResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ua.WriteResponse), args.Error(1)
}

func (m *MockOpcuaClient) Subscribe(ctx context.Context, params *opcua.SubscriptionParameters, notifyCh chan<- *opcua.PublishNotificationData) (*opcua.Subscription, error) {
	args := m.Called(ctx, params, notifyCh)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*opcua.Subscription), args.Error(1)
}

func (m *MockOpcuaClient) State() opcua.ConnState {
	args := m.Called()
	return args.Get(0).(opcua.ConnState)
}
