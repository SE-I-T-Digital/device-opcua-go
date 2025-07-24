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
	sdkModel "github.com/edgexfoundry/device-sdk-go/v4/pkg/models"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/common"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/models"
	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/ua"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/mock"
)

func TestDriver_ProcessWriteCommands(t *testing.T) {
	type args struct {
		deviceName string
		protocols  map[string]models.ProtocolProperties
		reqs       []sdkModel.CommandRequest
		params     []*sdkModel.CommandValue
	}
	tests := []struct {
		name        string
		args        args
		wantErr     bool
		endpointErr bool
		nilClient   bool
	}{
		{
			name: "NOK - no endpoint defined",
			args: args{
				deviceName: "Test",
				protocols:  map[string]models.ProtocolProperties{Protocol: {}},
				reqs:       []sdkModel.CommandRequest{{DeviceResourceName: "TestVar1"}},
				params:     []*sdkModel.CommandValue{{}},
			},
			wantErr:     true,
			endpointErr: true,
		},
		{
			name: "NOK - invalid node id",
			args: args{
				deviceName: "Test",
				protocols:  map[string]models.ProtocolProperties{Protocol: {Endpoint: test.Address}},
				reqs: []sdkModel.CommandRequest{{
					DeviceResourceName: "TestResource1",
					Attributes:         map[string]any{NODE: "ns=2;i=3;x=42"},
					Type:               common.ValueTypeInt32,
				}},
				params: []*sdkModel.CommandValue{{
					DeviceResourceName: "TestResource1",
					Type:               common.ValueTypeInt32,
					Value:              int32(42),
				}},
			},
			wantErr: true,
		},
		{
			name: "NOK - invalid value",
			args: args{
				deviceName: "Test",
				protocols:  map[string]models.ProtocolProperties{Protocol: {Endpoint: test.Address}},
				reqs: []sdkModel.CommandRequest{{
					DeviceResourceName: "TestResource1",
					Attributes:         map[string]any{NODE: "ns=2;s=rw_int32"},
					Type:               common.ValueTypeInt32,
				}},
				params: []*sdkModel.CommandValue{{
					DeviceResourceName: "TestResource1",
					Type:               common.ValueTypeString,
					Value:              "foobar",
				}},
			},
			wantErr: true,
		},
		{
			name: "NOK - client is nil",
			args: args{
				deviceName: "Test",
				protocols:  map[string]models.ProtocolProperties{Protocol: {Endpoint: test.Address}},
				reqs: []sdkModel.CommandRequest{{
					DeviceResourceName: "TestResource1",
					Attributes:         map[string]any{NODE: "ns=2;s=rw_int32"},
					Type:               common.ValueTypeInt32,
				}},
				params: []*sdkModel.CommandValue{{
					DeviceResourceName: "TestResource1",
					Type:               common.ValueTypeInt32,
					Value:              int32(42),
				}},
			},
			nilClient: true,
			wantErr:   true,
		},
		{
			name: "OK - command request with one parameter",
			args: args{
				deviceName: "Test",
				protocols:  map[string]models.ProtocolProperties{Protocol: {Endpoint: test.Address}},
				reqs: []sdkModel.CommandRequest{{
					DeviceResourceName: "TestResource1",
					Attributes:         map[string]any{NODE: "ns=2;s=rw_int32"},
					Type:               common.ValueTypeInt32,
				}},
				params: []*sdkModel.CommandValue{{
					DeviceResourceName: "TestResource1",
					Type:               common.ValueTypeInt32,
					Value:              int32(42),
				}},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			endpoint := cast.ToString(tt.args.protocols[Protocol][Endpoint])

			dsMock := test.NewDSMock(t)
			s := NewServer(tt.args.deviceName, dsMock)
			clientMock := new(test.MockOpcuaClient)
			clientMock.On("Write", s.context.ctx, mock.Anything).Return(&ua.WriteResponse{Results: []ua.StatusCode{ua.StatusOK}}, nil)
			clientMock.On("State").Return(opcua.Connected)
			clientMock.On("Connect", s.context.ctx).Return(nil)
			s.config = &Config{Endpoint: endpoint}

			if tt.nilClient {
				s.client = nil
				clientMock.On("State").Return(opcua.Closed)
				clientMock.On("Connect", mock.Anything).Return(fmt.Errorf("error"))
				dsMock.On("GetDeviceByName", tt.args.deviceName).Return(models.Device{}, fmt.Errorf("error"))
			} else if !tt.endpointErr {
				s.client = &Client{clientMock, s.context.ctx}
			}
			if err := s.ProcessWriteCommands(tt.args.reqs, tt.args.params); (err != nil) != tt.wantErr {
				t.Errorf("Driver.HandleWriteCommands() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
