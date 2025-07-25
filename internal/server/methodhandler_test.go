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
	}
	type args struct {
		device     models.Device
		resource   models.DeviceResource
		method     string
		parameters []string
	}
	tests := []struct {
		name      string
		args      args
		deviceErr error
		want      any
		wantErr   bool
		nilClient bool
	}{
		{
			name:      "NOK - device not found",
			deviceErr: fmt.Errorf("device not found"),
			wantErr:   true,
		},
		{
			name: "NOK - device locked",
			args: args{
				device: models.Device{AdminState: models.Locked},
			},
			wantErr: true,
		},
		{
			name: "NOK - device down",
			args: args{
				device: models.Device{OperatingState: models.Down},
			},
			wantErr: true,
		},
		{
			name: "NOK - method call - method not found",
			args: args{
				method: "TestResource0",
				device: okDevice,
			},
			wantErr: true,
		},
		{
			name: "NOK - method call - method is hidden",
			args: args{
				resource: models.DeviceResource{Name: "TestResource1", IsHidden: true},
				device:   okDevice,
			},
			wantErr: true,
		},
		{
			name: "NOK - method call - invalid object node id",
			args: args{
				resource: models.DeviceResource{
					Name:       "TestResource1",
					Attributes: map[string]any{METHOD: "ns=2;s=test"},
				},
				device: okDevice,
			},
			wantErr: true,
		},
		{
			name: "NOK - method call - invalid method node id",
			args: args{
				resource: models.DeviceResource{
					Name:       "TestResource1",
					Attributes: map[string]any{OBJECT: "ns=2;s=main"},
				},
				device: okDevice,
			},
			wantErr: true,
		},
		{
			name: "NOK - method call - method does not exist",
			args: args{
				resource: models.DeviceResource{
					Name: "TestResource1",
					Attributes: map[string]any{
						METHOD: "NONE",
						OBJECT: "ns=2;s=main",
					},
				},
				device: okDevice,
			},
			wantErr: true,
		},
		{
			name: "NOK - client is nil",
			args: args{
				resource: models.DeviceResource{
					Name: "TestResource1",
					Attributes: map[string]any{
						METHOD: "ns=2;s=square",
						OBJECT: "ns=2;s=main",
					},
				},
				parameters: []string{"2"},
				device:     okDevice,
			},
			nilClient: true,
			wantErr:   true,
		},
		{
			name: "OK - call method from mock server",
			args: args{
				resource: models.DeviceResource{
					Name: "TestResource1",
					Attributes: map[string]any{
						METHOD: "ns=2;s=square",
						OBJECT: "ns=2;s=main",
					},
				},
				parameters: []string{"2"},
				device:     okDevice,
			},
			want: "4",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dsMock := test.NewDSMock(t)
			s := NewServer("test", dsMock)
			clientMock := new(test.MockOpcuaClient)

			dsMock.On("GetDeviceByName", mock.Anything).Return(tt.args.device, tt.deviceErr).Times(1)
			if tt.deviceErr == nil && tt.args.device.AdminState != models.Locked && tt.args.device.OperatingState != models.Down {
				dsMock.On("DeviceResource", mock.Anything, tt.args.method).Return(tt.args.resource, tt.args.resource.Name != "")
			}
			if tt.nilClient {
				s.client = nil
				clientMock.On("State").Return(opcua.Closed)
				clientMock.On("Connect", mock.Anything).Return(fmt.Errorf("error"))
				dsMock.On("GetDeviceByName", mock.Anything).Return(models.Device{}, fmt.Errorf("error")).Times(1)
			} else {
				s.client = &Client{clientMock, s.context.ctx}
				clientMock.On("State").Return(opcua.Connected)
				if tt.args.resource.Attributes[METHOD] != "NONE" {
					clientMock.On("Call", mock.Anything, mock.Anything).Return(&ua.CallMethodResult{
						StatusCode:      ua.StatusOK,
						OutputArguments: []*ua.Variant{ua.MustVariant("4")},
					}, nil)
				} else {
					clientMock.On("Call", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("no method with that name"))
				}
			}
			got, err := s.ProcessMethodCall(tt.args.method, tt.args.parameters)
			if (err != nil) != tt.wantErr {
				t.Errorf("Driver.HandleReadCommands() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Driver.HandleReadCommands() = %v, want %v", got, tt.want)
			}
		})
	}
}
