// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2021 Schneider Electric
//
// SPDX-License-Identifier: Apache-2.0

package server

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/edgexfoundry/device-opcua-go/internal/test"
	"github.com/edgexfoundry/device-opcua-go/pkg/gopcua"
	gopcuaMocks "github.com/edgexfoundry/device-opcua-go/pkg/gopcua/mocks"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/models"
	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/ua"
	"github.com/stretchr/testify/mock"
)

func TestDriver_ProcessMethodCall(t *testing.T) {
	okDevice := models.Device{
		Name:           "TestDevice",
		AdminState:     models.Unlocked,
		OperatingState: models.Up,
		Protocols:      map[string]models.ProtocolProperties{Protocol: {Endpoint: test.Address}},
	}

	origGetEndpoints := gopcua.GetEndpoints
	defer func() { gopcua.GetEndpoints = origGetEndpoints }()
	test.MockGetEndpoints()

	origNewOpcuaClient := gopcua.NewClient
	defer func() { gopcua.NewClient = origNewOpcuaClient }()

	t.Run("NOK - device not found", func(t *testing.T) {
		dsMock := test.NewDSMock(t)
		s := NewServer("test", dsMock)
		dsMock.On("GetDeviceByName", mock.Anything).Return(models.Device{}, fmt.Errorf("device not found"))
		_, err := s.ProcessMethodCall("", nil)
		if err == nil {
			t.Error("expected error but got none")
		}
	})

	t.Run("NOK - device locked", func(t *testing.T) {
		dsMock := test.NewDSMock(t)
		s := NewServer("test", dsMock)
		device := models.Device{AdminState: models.Locked}
		dsMock.On("GetDeviceByName", mock.Anything).Return(device, nil)
		_, err := s.ProcessMethodCall("", nil)
		if err == nil {
			t.Error("expected error but got none")
		}
	})

	t.Run("NOK - device down", func(t *testing.T) {
		dsMock := test.NewDSMock(t)
		s := NewServer("test", dsMock)
		device := models.Device{OperatingState: models.Down}
		dsMock.On("GetDeviceByName", mock.Anything).Return(device, nil)
		_, err := s.ProcessMethodCall("", nil)
		if err == nil {
			t.Error("expected error but got none")
		}
	})

	t.Run("NOK - method call - method not found", func(t *testing.T) {
		dsMock := test.NewDSMock(t)
		s := NewServer("test", dsMock)
		dsMock.On("GetDeviceByName", mock.Anything).Return(okDevice, nil)
		dsMock.On("DeviceResource", mock.Anything, "TestResource0").Return(models.DeviceResource{}, false)
		_, err := s.ProcessMethodCall("TestResource0", nil)
		if err == nil {
			t.Error("expected error but got none")
		}
	})

	t.Run("NOK - method call - method is hidden", func(t *testing.T) {
		dsMock := test.NewDSMock(t)
		s := NewServer("test", dsMock)
		resource := models.DeviceResource{Name: "TestResource1", IsHidden: true}
		dsMock.On("GetDeviceByName", mock.Anything).Return(okDevice, nil)
		dsMock.On("DeviceResource", mock.Anything, "").Return(resource, true)
		_, err := s.ProcessMethodCall("", nil)
		if err == nil {
			t.Error("expected error but got none")
		}
	})

	t.Run("NOK - method call - invalid object node id", func(t *testing.T) {
		dsMock := test.NewDSMock(t)
		s := NewServer("test", dsMock)
		resource := models.DeviceResource{
			Name:       "TestResource1",
			Attributes: map[string]any{METHOD: "ns=2;s=test"},
		}
		dsMock.On("GetDeviceByName", mock.Anything).Return(okDevice, nil)
		dsMock.On("DeviceResource", mock.Anything, "").Return(resource, true)
		_, err := s.ProcessMethodCall("", nil)
		if err == nil {
			t.Error("expected error but got none")
		}
	})

	t.Run("NOK - method call - invalid method node id", func(t *testing.T) {
		dsMock := test.NewDSMock(t)
		s := NewServer("test", dsMock)
		resource := models.DeviceResource{
			Name:       "TestResource1",
			Attributes: map[string]any{OBJECT: "ns=2;s=main"},
		}
		dsMock.On("GetDeviceByName", mock.Anything).Return(okDevice, nil)
		dsMock.On("DeviceResource", mock.Anything, "").Return(resource, true)
		_, err := s.ProcessMethodCall("", nil)
		if err == nil {
			t.Error("expected error but got none")
		}
	})

	t.Run("NOK - method call - method does not exist", func(t *testing.T) {
		dsMock := test.NewDSMock(t)
		s := NewServer("test", dsMock)
		clientMock := gopcuaMocks.NewMockClient(t)
		s.client = &Client{clientMock, s.context.ctx}
		resource := models.DeviceResource{
			Name: "TestResource1",
			Attributes: map[string]any{
				METHOD: "NONE",
				OBJECT: "ns=2;s=main",
			},
		}
		dsMock.On("GetDeviceByName", mock.Anything).Return(okDevice, nil)
		dsMock.On("DeviceResource", mock.Anything, "").Return(resource, true)
		clientMock.On("State").Return(opcua.Connected)
		clientMock.On("Call", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("no method with that name"))
		_, err := s.ProcessMethodCall("", nil)
		if err == nil {
			t.Error("expected error but got none")
		}
	})

	t.Run("NOK - client connection is closed", func(t *testing.T) {
		dsMock := test.NewDSMock(t)
		s := NewServer("test", dsMock)
		clientMock := gopcuaMocks.NewMockClient(t)
		gopcua.NewClient = func(endpoint string, opts ...opcua.Option) (gopcua.Client, error) {
			return clientMock, nil
		}
		resource := models.DeviceResource{
			Name: "TestResource1",
			Attributes: map[string]any{
				METHOD: "ns=2;s=square",
				OBJECT: "ns=2;s=main",
			},
		}
		s.client = &Client{clientMock, s.context.ctx}
		dsMock.On("GetDeviceByName", mock.Anything).Return(okDevice, nil)
		dsMock.On("DeviceResource", mock.Anything, "").Return(resource, true)
		clientMock.On("State").Return(opcua.Closed)
		clientMock.On("Connect", mock.Anything).Return(fmt.Errorf("error"))
		_, err := s.ProcessMethodCall("", []string{"2"})
		if err == nil {
			t.Error("expected error but got none")
		}
	})

	t.Run("OK - call method from mock server", func(t *testing.T) {
		dsMock := test.NewDSMock(t)
		s := NewServer("test", dsMock)
		clientMock := gopcuaMocks.NewMockClient(t)
		s.client = &Client{clientMock, s.context.ctx}
		resource := models.DeviceResource{
			Name: "TestResource1",
			Attributes: map[string]any{
				METHOD: "ns=2;s=square",
				OBJECT: "ns=2;s=main",
			},
		}
		var want any = "4"
		dsMock.On("GetDeviceByName", mock.Anything).Return(okDevice, nil)
		dsMock.On("DeviceResource", mock.Anything, "").Return(resource, true)
		clientMock.On("State").Return(opcua.Connected)
		clientMock.On("Call", mock.Anything, mock.Anything).Return(&ua.CallMethodResult{
			StatusCode:      ua.StatusOK,
			OutputArguments: []*ua.Variant{ua.MustVariant("4")},
		}, nil)
		got, err := s.ProcessMethodCall("", []string{"2"})
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Driver.HandleReadCommands() = %v, want %v", got, want)
		}
	})
}
