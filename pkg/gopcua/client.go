// -*- Mode: Go; indent-tabs-mode: t -*-
//
// # Copyright (C) 2025 Schneider Electric
//
// SPDX-License-Identifier: Apache-2.0
package gopcua

import (
	"context"

	"github.com/gopcua/opcua"
	"github.com/gopcua/opcua/ua"
)

// Client is an interface for opcua.Client to allow for mocking
type Client interface {
	Connect(ctx context.Context) error
	Close(ctx context.Context) error
	Call(ctx context.Context, req *ua.CallMethodRequest) (*ua.CallMethodResult, error)
	Read(ctx context.Context, req *ua.ReadRequest) (*ua.ReadResponse, error)
	Write(ctx context.Context, req *ua.WriteRequest) (*ua.WriteResponse, error)
	Subscribe(ctx context.Context, params *opcua.SubscriptionParameters, notifyCh chan<- *opcua.PublishNotificationData) (*opcua.Subscription, error)
	State() opcua.ConnState
}

// NewClient returns a real OPC UA client, but can be overwritten for testing
var NewClient = func(endpoint string, opts ...opcua.Option) (Client, error) {
	return opcua.NewClient(endpoint, opts...)
}
