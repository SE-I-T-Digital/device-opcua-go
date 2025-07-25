// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2021 Schneider Electric
//
// SPDX-License-Identifier: Apache-2.0

package server

import (
	"context"
	"fmt"
	"testing"

	"github.com/edgexfoundry/device-opcua-go/internal/test"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/models"
	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/ua"
)

func Test_StartSubscriptionListener(t *testing.T) {
	t.Run("call and exit", func(t *testing.T) {
		dsMock := test.NewDSMock(t)
		dsMock.On("GetDeviceByName", "Test").Return(models.Device{}, fmt.Errorf("error"))

		s := NewServer("Test", dsMock)
		err := s.StartSubscriptionListener()
		if err == nil {
			t.Error("expected err to exist in test environment")
		}
	})
}

func Test_configureMonitoredItems(t *testing.T) {
	t.Run("configure multiple items", func(t *testing.T) {
		dsMock := test.NewDSMock(t)
		dsMock.On("DeviceResource", "Test", "a").Return(models.DeviceResource{}, false)
		dsMock.On("DeviceResource", "Test", "b").Return(models.DeviceResource{}, false)
		dsMock.On("DeviceResource", "Test", "c").Return(models.DeviceResource{}, false)

		s := NewServer("Test", dsMock)
		s.config = &Config{Resources: []string{"a", "b", "c"}}
		err := s.configureMonitoredItems(nil)
		if err != nil {
			t.Errorf("expected no error, got = %v", err)
		}
	})
}

func Test_onIncomingDataReceived(t *testing.T) {
	t.Run("device resource unknown", func(t *testing.T) {
		dsMock := test.NewDSMock(t)
		dsMock.On("DeviceResource", "Test", "TestResource").Return(models.DeviceResource{}, false)

		s := NewServer("Test", dsMock)
		err := s.onIncomingDataReceived("42", "TestResource")
		if err == nil {
			t.Error("expected err to exist in test environment")
		}
	})
}

func TestDriver_initClient(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name:    "NOK - no endpoint configured",
			config:  &Config{},
			wantErr: true,
		},
		{
			name: "OK",
			config: &Config{
				Endpoint: test.Address,
				Policy:   "None",
				Mode:     "None",
			},
		},
	}

	origGetEndpoints := getEndpoints
	defer func() { getEndpoints = origGetEndpoints }()
	getEndpoints = func(ctx context.Context, endpointURL string, opts ...opcua.Option) ([]*ua.EndpointDescription, error) {
		if endpointURL == test.Address {
			return []*ua.EndpointDescription{{SecurityPolicyURI: ua.SecurityPolicyURINone, SecurityMode: ua.MessageSecurityModeNone}}, nil
		}
		return nil, fmt.Errorf("bad endpoint")
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dsMock := test.NewDSMock(t)
			s := NewServer("Test", dsMock)

			s.config = tt.config
			err := s.initClient()
			if (err != nil) != tt.wantErr {
				t.Errorf("Driver.getClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestDriver_handleDataChange(t *testing.T) {
	t.Run("OK - no monitored items", func(t *testing.T) {
		dsMock := test.NewDSMock(t)
		s := NewServer("Test", dsMock)
		s.handleDataChange(&ua.DataChangeNotification{MonitoredItems: make([]*ua.MonitoredItemNotification, 0)})
	})

	t.Run("OK - call onIncomingDataReceived", func(t *testing.T) {
		dsMock := test.NewDSMock(t)
		dsMock.On("DeviceResource", "Test", "").Return(models.DeviceResource{Name: "TestResource"}, true)
		s := NewServer("Test", dsMock)

		s.handleDataChange(&ua.DataChangeNotification{
			MonitoredItems: []*ua.MonitoredItemNotification{
				{ClientHandle: 123456, Value: &ua.DataValue{Value: ua.MustVariant("42")}},
			},
		})
	})
}
