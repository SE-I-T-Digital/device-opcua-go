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
	sdkModel "github.com/edgexfoundry/device-sdk-go/v4/pkg/models"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/clients/logger"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/common"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/models"
	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/ua"
	"github.com/stretchr/testify/mock"
)

const (
	Protocol string = "opcua"
	Endpoint string = "Endpoint"
)

func TestDriver_ProcessReadCommands(t *testing.T) {
	origGetEndpoints := gopcua.GetEndpoints
	defer func() { gopcua.GetEndpoints = origGetEndpoints }()
	test.MockGetEndpoints()

	origNewOpcuaClient := gopcua.NewClient
	defer func() { gopcua.NewClient = origNewOpcuaClient }()

	t.Run("NOK - no endpoint defined", func(t *testing.T) {
		reqs := []sdkModel.CommandRequest{{
			DeviceResourceName: "TestVar1",
			Attributes:         map[string]any{NODE: "ns=2;s=fake"},
			Type:               common.ValueTypeInt32,
		}}
		dsMock := test.NewDSMock(t)
		dsMock.On("GetDeviceByName", "Test").Return(models.Device{Protocols: map[string]models.ProtocolProperties{Protocol: {Endpoint: ""}}}, nil)
		s := NewServer("Test", dsMock)

		_, err := s.ProcessReadCommands(reqs)
		if err == nil {
			t.Error("expected an error but got none")
		}
	})

	t.Run("OK - non-existent variable", func(t *testing.T) {
		reqs := []sdkModel.CommandRequest{{
			DeviceResourceName: "TestVar1",
			Attributes:         map[string]any{NODE: "ns=2;s=fake"},
			Type:               common.ValueTypeInt32,
		}}
		want := make([]*sdkModel.CommandValue, 1)

		dsMock := test.NewDSMock(t)
		s := NewServer("TestWithFakeVar", dsMock)
		s.config = &Config{Endpoint: test.Address}
		clientMock := gopcuaMocks.NewMockClient(t)

		readResponse := &ua.ReadResponse{}

		clientMock.On("Read", s.context.ctx, mock.Anything).Return(readResponse, nil)
		clientMock.On("State").Return(opcua.Connected)
		s.client = &Client{clientMock, s.context.ctx}

		got, err := s.ProcessReadCommands(reqs)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Driver.HandleReadCommands() = %v, want %v", got, want)
		}
	})

	t.Run("NOK - read command - invalid node id", func(t *testing.T) {
		reqs := []sdkModel.CommandRequest{{
			DeviceResourceName: "TestResource1",
			Attributes:         map[string]any{NODE: "ns=2;i=22;z=43"},
			Type:               common.ValueTypeInt32,
		}}
		dsMock := test.NewDSMock(t)
		s := NewServer("Test", dsMock)
		s.config = &Config{Endpoint: test.Address}

		_, err := s.ProcessReadCommands(reqs)
		if err == nil {
			t.Error("expected an error but got none")
		}
	})

	t.Run("NOK - not allowed to call method with reader", func(t *testing.T) {
		reqs := []sdkModel.CommandRequest{{
			DeviceResourceName: "SquareResource",
			Attributes:         map[string]any{METHOD: "ns=2;s=square", OBJECT: "ns=2;s=main"},
			Type:               common.ValueTypeInt64,
		}}
		dsMock := test.NewDSMock(t)
		s := NewServer("Test", dsMock)
		s.config = &Config{Endpoint: test.Address}

		_, err := s.ProcessReadCommands(reqs)
		if err == nil {
			t.Error("expected an error but got none")
		}
	})

	t.Run("OK - read value from mock server", func(t *testing.T) {
		reqs := []sdkModel.CommandRequest{{
			DeviceResourceName: "TestVar1",
			Attributes:         map[string]any{NODE: "ns=2;s=ro_int32"},
			Type:               common.ValueTypeInt32,
		}}
		want := []*sdkModel.CommandValue{{
			DeviceResourceName: "TestVar1",
			Type:               common.ValueTypeInt32,
			Value:              int32(5),
			Tags:               make(map[string]string),
		}}
		dsMock := test.NewDSMock(t)
		s := NewServer("Test", dsMock)
		s.config = &Config{Endpoint: test.Address}

		clientMock := gopcuaMocks.NewMockClient(t)
		readResponse := &ua.ReadResponse{
			Results: []*ua.DataValue{
				{Value: ua.MustVariant(int32(5))},
			},
		}
		clientMock.On("Read", s.context.ctx, mock.Anything).Return(readResponse, nil)
		clientMock.On("State").Return(opcua.Connected)
		s.client = &Client{clientMock, s.context.ctx}

		got, err := s.ProcessReadCommands(reqs)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		for i := range got {
			if got[i] != nil {
				got[i].Origin = 0
			}
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Driver.HandleReadCommands() = %v, want %v", got, want)
		}
	})

	t.Run("OK - read many values from mock server", func(t *testing.T) {
		reqs := []sdkModel.CommandRequest{{
			DeviceResourceName: "TestVar1",
			Attributes:         map[string]any{NODE: "ns=2;s=ro_int32"},
			Type:               common.ValueTypeInt32,
		}, {
			DeviceResourceName: "TestVar2",
			Attributes:         map[string]any{NODE: "ns=2;s=ro_bool"},
			Type:               common.ValueTypeBool,
		}}
		want := []*sdkModel.CommandValue{{
			DeviceResourceName: "TestVar1",
			Type:               common.ValueTypeInt32,
			Value:              int32(5),
			Tags:               make(map[string]string),
		}, {
			DeviceResourceName: "TestVar2",
			Type:               common.ValueTypeBool,
			Value:              true,
			Tags:               make(map[string]string),
		}}

		dsMock := test.NewDSMock(t)
		s := NewServer("Test", dsMock)
		s.config = &Config{Endpoint: test.Address}

		clientMock := gopcuaMocks.NewMockClient(t)
		readResponse := &ua.ReadResponse{
			Results: []*ua.DataValue{
				{Value: ua.MustVariant(int32(5))},
				{Value: ua.MustVariant(true)},
			},
		}
		clientMock.On("Read", s.context.ctx, mock.Anything).Return(readResponse, nil)
		clientMock.On("State").Return(opcua.Connected)
		s.client = &Client{clientMock, s.context.ctx}

		got, err := s.ProcessReadCommands(reqs)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		for i := range got {
			if got[i] != nil {
				got[i].Origin = 0
			}
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Driver.HandleReadCommands() = %v, want %v", got, want)
		}
	})

	t.Run("NOK - error from nil client", func(t *testing.T) {
		reqs := []sdkModel.CommandRequest{{
			DeviceResourceName: "TestVar1",
			Attributes:         map[string]any{NODE: "ns=2;s=ro_int32"},
			Type:               common.ValueTypeInt32,
		}}
		dsMock := test.NewDSMock(t)
		dsMock.On("GetDeviceByName", "Test").Return(models.Device{}, fmt.Errorf("error"))
		s := NewServer("Test", dsMock)
		s.client = nil

		_, err := s.ProcessReadCommands(reqs)
		if err == nil {
			t.Error("expected an error but got none")
		}
	})

	t.Run("NOK - error from disconnected server", func(t *testing.T) {
		reqs := []sdkModel.CommandRequest{{
			DeviceResourceName: "TestVar1",
			Attributes:         map[string]any{NODE: "ns=2;s=ro_int32"},
			Type:               common.ValueTypeInt32,
		}}
		dsMock := test.NewDSMock(t)
		dsMock.On("GetDeviceByName", "Test").Return(models.Device{Protocols: map[string]models.ProtocolProperties{Protocol: {Endpoint: test.Address}}}, nil)

		s := NewServer("Test", dsMock)

		clientMock := gopcuaMocks.NewMockClient(t)
		clientMock.On("State").Return(opcua.Disconnected)
		clientMock.On("Connect", mock.Anything).Return(fmt.Errorf("error"))
		s.client = &Client{clientMock, s.context.ctx}

		gopcua.NewClient = func(endpoint string, opts ...opcua.Option) (gopcua.Client, error) {
			return clientMock, nil
		}

		_, err := s.ProcessReadCommands(reqs)
		if err == nil {
			t.Error("expected an error but got none")
		}
	})
}

func TestBuildReadRequest(t *testing.T) {
	tests := []struct {
		name                    string
		reqNodeIds              []string
		expectedNodeIds         map[string]struct{}
		expectedResultToRequest ResultToRequest
	}{
		{
			name:                    "OK - Empty Request",
			reqNodeIds:              []string{},
			expectedNodeIds:         map[string]struct{}{},
			expectedResultToRequest: map[int][]int{},
		}, {
			name:       "OK - On Read Request",
			reqNodeIds: []string{"ns=1;i=1"},
			expectedNodeIds: map[string]struct{}{
				"ns=1;i=1": {},
			},
			expectedResultToRequest: map[int][]int{
				0: {0},
			},
		}, {
			name: "OK - Multi Read Request",
			reqNodeIds: []string{
				"ns=1;i=1",
				"ns=1;i=2",
				"ns=1;i=3",
				"ns=1;i=4",
				"ns=1;i=5",
			},
			expectedNodeIds: map[string]struct{}{
				"ns=1;i=1": {},
				"ns=1;i=2": {},
				"ns=1;i=3": {},
				"ns=1;i=4": {},
				"ns=1;i=5": {},
			},
			expectedResultToRequest: map[int][]int{
				0: {0},
				1: {1},
				2: {2},
				3: {3},
				4: {4},
			},
		}, {
			name: "OK - Two Overlapping Read Requests",
			reqNodeIds: []string{
				"ns=1;i=1",
				"ns=1;i=1",
			},
			expectedNodeIds: map[string]struct{}{
				"ns=1;i=1": {},
			},
			expectedResultToRequest: map[int][]int{
				0: {0, 1},
			},
		}, {
			name: "OK - Complex Read Requests",
			reqNodeIds: []string{
				"ns=1;i=1",
				"ns=1;i=1",
				"ns=1;i=2",
				"ns=1;i=1",
				"ns=1;i=3",
				"ns=1;i=2",
				"ns=1;i=1",
			},
			expectedNodeIds: map[string]struct{}{
				"ns=1;i=1": {},
				"ns=1;i=2": {},
				"ns=1;i=3": {},
			},
			expectedResultToRequest: map[int][]int{
				0: {0, 1, 3, 6},
				1: {2, 5},
				2: {4},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqs := make([]sdkModel.CommandRequest, 0, len(tt.reqNodeIds))
			for _, id := range tt.reqNodeIds {
				reqs = append(reqs, sdkModel.CommandRequest{Attributes: map[string]any{NODE: id}})
			}

			nodesToRead, resultToRequest, err := buildNodesToReadRequest(reqs)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			for _, nodes := range nodesToRead {
				id := nodes.NodeID.String()
				if _, ok := tt.expectedNodeIds[id]; !ok {
					t.Fatalf("Node %s not found in request", id)
				} else {
					delete(tt.expectedNodeIds, id)
				}
			}

			if len(tt.expectedNodeIds) > 0 {
				t.Fatalf("Unexpected node in request: %+v", nodesToRead)
			}

			if !reflect.DeepEqual(resultToRequest, tt.expectedResultToRequest) {
				t.Fatalf("Unexpected result to request: expected %+v; got %+v", tt.expectedResultToRequest, resultToRequest)
			}
		})
	}
}

func TestBuildReadRequestOnMethod(t *testing.T) {
	reqs := []sdkModel.CommandRequest{
		{
			Attributes: map[string]any{METHOD: true},
		},
	}
	_, _, err := buildNodesToReadRequest(reqs)

	if err == nil {
		t.Fatalf("Method request should not be allowed")
	}
}

func TestBuildReadRequestOnMissingNodeId(t *testing.T) {
	reqs := []sdkModel.CommandRequest{
		{
			Attributes: map[string]any{METHOD: false},
		},
	}
	_, _, err := buildNodesToReadRequest(reqs)

	if err == nil {
		t.Fatalf("Node Id is missing from properties; error expected")
	}
}

func TestBuildCommandValues(t *testing.T) {
	reqs := []sdkModel.CommandRequest{
		{
			DeviceResourceName: "Res1",
			Type:               common.ValueTypeInt32,
		},
		{
			DeviceResourceName: "Res2",
			Type:               common.ValueTypeInt32,
		},
		{
			DeviceResourceName: "Res3",
			Type:               common.ValueTypeInt32,
		},
	}

	lc := logger.NewMockClient()

	t.Run("Read on one node", func(t *testing.T) {
		var resultToRequest ResultToRequest = map[int][]int{0: {0, 1, 2}}

		uaResponse := &ua.ReadResponse{
			Results: []*ua.DataValue{
				{
					Value: ua.MustVariant(int32(1)),
				},
			},
		}

		commandValues := resultToRequest.buildCommandValues(reqs, uaResponse, lc)

		if len(commandValues) != 3 {
			t.Fatalf("Expected number of command values 3; got %d;", len(commandValues))
		}

		if commandValues[0].DeviceResourceName != "Res1" {
			t.Fatalf("Expected device resource name [0] Res1; got %s", commandValues[0].DeviceResourceName)
		}

		if commandValues[1].DeviceResourceName != "Res2" {
			t.Fatalf("Expected device resource name [1] Res2; got %s", commandValues[1].DeviceResourceName)
		}

		if commandValues[2].DeviceResourceName != "Res3" {
			t.Fatalf("Expected device resource name [2] Res3; got %s", commandValues[2].DeviceResourceName)
		}

		if commandValues[0].Value != int32(1) {
			t.Fatalf("Expected device resource value [0] 1; got %v", commandValues[0].Value)
		}

		if commandValues[1].Value != int32(1) {
			t.Fatalf("Expected device resource value [1] 1; got %v", commandValues[1].Value)
		}

		if commandValues[2].Value != int32(1) {
			t.Fatalf("Expected device resource value [2] 1; got %v", commandValues[2].Value)
		}

	})

	t.Run("Read on multiple nodes", func(t *testing.T) {
		var resultToRequest ResultToRequest = map[int][]int{0: {0}, 1: {1}, 2: {2}}

		uaResponse := &ua.ReadResponse{
			Results: []*ua.DataValue{
				{
					Value: ua.MustVariant(int32(1)),
				}, {
					Value: ua.MustVariant(int32(2)),
				}, {
					Value: ua.MustVariant(int32(3)),
				},
			},
		}

		commandValues := resultToRequest.buildCommandValues(reqs, uaResponse, lc)

		if len(commandValues) != 3 {
			t.Fatalf("Expected number of command values 3; got %d;", len(commandValues))
		}

		if commandValues[0].DeviceResourceName != "Res1" {
			t.Fatalf("Expected device resource name [0] Res1; got %s", commandValues[0].DeviceResourceName)
		}

		if commandValues[1].DeviceResourceName != "Res2" {
			t.Fatalf("Expected device resource name [1] Res2; got %s", commandValues[1].DeviceResourceName)
		}

		if commandValues[2].DeviceResourceName != "Res3" {
			t.Fatalf("Expected device resource name [2] Res3; got %s", commandValues[2].DeviceResourceName)
		}

		if commandValues[0].Value != int32(1) {
			t.Fatalf("Expected device resource value [0] 1; got %v", commandValues[0].Value)
		}

		if commandValues[1].Value != int32(2) {
			t.Fatalf("Expected device resource value [1] 2; got %v", commandValues[1].Value)
		}

		if commandValues[2].Value != int32(3) {
			t.Fatalf("Expected device resource value [2] 3; got %v", commandValues[2].Value)
		}

	})
}
