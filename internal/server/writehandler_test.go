// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2021 Schneider Electric
//
// SPDX-License-Identifier: Apache-2.0

package server

import (
	"fmt"
	"testing"

	"github.com/edgexfoundry/device-opcua-go/internal/test"
	"github.com/edgexfoundry/device-opcua-go/pkg/gopcua"
	gopcuaMocks "github.com/edgexfoundry/device-opcua-go/pkg/gopcua/mocks"
	sdkModel "github.com/edgexfoundry/device-sdk-go/v4/pkg/models"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/common"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/models"
	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/ua"
	"github.com/stretchr/testify/mock"
)

func TestDriver_ProcessWriteCommands(t *testing.T) {
	origGetEndpoints := gopcua.GetEndpoints
	defer func() { gopcua.GetEndpoints = origGetEndpoints }()
	test.MockGetEndpoints()

	origNewOpcuaClient := gopcua.NewClient
	defer func() { gopcua.NewClient = origNewOpcuaClient }()

	t.Run("NOK - no endpoint defined", func(t *testing.T) {
		reqs := []sdkModel.CommandRequest{{DeviceResourceName: "TestVar1"}}
		params := []*sdkModel.CommandValue{{}}
		dsMock := test.NewDSMock(t)
		s := NewServer("Test", dsMock)
		s.config = &Config{Endpoint: ""}

		if err := s.ProcessWriteCommands(reqs, params); err == nil {
			t.Error("expected error but got none")
		}
	})

	t.Run("NOK - invalid node id", func(t *testing.T) {
		reqs := []sdkModel.CommandRequest{{
			DeviceResourceName: "TestResource1",
			Attributes:         map[string]any{NODE: "ns=2;i=3;x=42"},
			Type:               common.ValueTypeInt32,
		}}
		params := []*sdkModel.CommandValue{{
			DeviceResourceName: "TestResource1",
			Type:               common.ValueTypeInt32,
			Value:              int32(42),
		}}
		dsMock := test.NewDSMock(t)
		s := NewServer("Test", dsMock)
		s.config = &Config{Endpoint: test.Address}

		if err := s.ProcessWriteCommands(reqs, params); err == nil {
			t.Error("expected error but got none")
		}
	})

	t.Run("NOK - invalid value", func(t *testing.T) {
		reqs := []sdkModel.CommandRequest{{
			DeviceResourceName: "TestResource1",
			Attributes:         map[string]any{NODE: "ns=2;s=rw_int32"},
			Type:               common.ValueTypeInt32,
		}}
		params := []*sdkModel.CommandValue{{
			DeviceResourceName: "TestResource1",
			Type:               common.ValueTypeString,
			Value:              "foobar",
		}}
		dsMock := test.NewDSMock(t)
		s := NewServer("Test", dsMock)
		s.config = &Config{Endpoint: test.Address}

		if err := s.ProcessWriteCommands(reqs, params); err == nil {
			t.Error("expected error but got none")
		}
	})

	t.Run("NOK - client is closed", func(t *testing.T) {
		reqs := []sdkModel.CommandRequest{{
			DeviceResourceName: "TestResource1",
			Attributes:         map[string]any{NODE: "ns=2;s=rw_int32"},
			Type:               common.ValueTypeInt32,
		}}
		params := []*sdkModel.CommandValue{{
			DeviceResourceName: "TestResource1",
			Type:               common.ValueTypeInt32,
			Value:              int32(42),
		}}
		dsMock := test.NewDSMock(t)
		dsMock.On("GetDeviceByName", "Test").Return(models.Device{Protocols: map[string]models.ProtocolProperties{Protocol: {Endpoint: test.Address}}}, nil)

		s := NewServer("Test", dsMock)
		clientMock := gopcuaMocks.NewMockClient(t)
		s.client = &Client{clientMock, s.context.ctx}

		clientMock.On("State").Return(opcua.Closed)
		clientMock.On("Connect", mock.Anything).Return(fmt.Errorf("error"))

		gopcua.NewClient = func(endpoint string, opts ...opcua.Option) (gopcua.Client, error) {
			return clientMock, nil
		}

		if err := s.ProcessWriteCommands(reqs, params); err == nil {
			t.Error("expected error but got none")
		}
	})

	t.Run("OK - command request with one parameter", func(t *testing.T) {
		reqs := []sdkModel.CommandRequest{{
			DeviceResourceName: "TestResource1",
			Attributes:         map[string]any{NODE: "ns=2;s=rw_int32"},
			Type:               common.ValueTypeInt32,
		}}
		params := []*sdkModel.CommandValue{{
			DeviceResourceName: "TestResource1",
			Type:               common.ValueTypeInt32,
			Value:              int32(42),
		}}

		dsMock := test.NewDSMock(t)
		s := NewServer("Test", dsMock)

		clientMock := gopcuaMocks.NewMockClient(t)
		s.client = &Client{clientMock, s.context.ctx}

		clientMock.On("State").Return(opcua.Connected)
		clientMock.On("Write", s.context.ctx, mock.Anything).Return(&ua.WriteResponse{Results: []ua.StatusCode{ua.StatusOK}}, nil)

		gopcua.NewClient = func(endpoint string, opts ...opcua.Option) (gopcua.Client, error) {
			return clientMock, nil
		}

		if err := s.ProcessWriteCommands(reqs, params); err != nil {
			t.Errorf("Driver.HandleWriteCommands() unexpected error = %v", err)
		}
	})
}
