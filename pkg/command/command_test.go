// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2018 Canonical Ltd
// Copyright (C) 2018 IOTech Ltd
// Copyright (C) 2021 Schneider Electric
//
// SPDX-License-Identifier: Apache-2.0

package command

import (
	"reflect"
	"testing"

	sdkModel "github.com/edgexfoundry/device-sdk-go/v4/pkg/models"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/common"
)

func Test_newCommandValue(t *testing.T) {
	type args struct {
		valueType string
		param     *sdkModel.CommandValue
	}
	tests := []struct {
		name    string
		args    args
		want    any
		wantErr bool
	}{
		{
			name:    "NOK - unknown type",
			args:    args{valueType: "uknown"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "NOK - bool value - mismatching types",
			args:    args{valueType: common.ValueTypeBool, param: &sdkModel.CommandValue{Value: "42", Type: common.ValueTypeString}},
			want:    false,
			wantErr: true,
		},
		{
			name:    "OK - bool value - matching types",
			args:    args{valueType: common.ValueTypeBool, param: &sdkModel.CommandValue{Value: true, Type: common.ValueTypeBool}},
			want:    true,
			wantErr: false,
		},
		{
			name:    "OK - string value",
			args:    args{valueType: common.ValueTypeString, param: &sdkModel.CommandValue{Value: "test", Type: common.ValueTypeString}},
			want:    "test",
			wantErr: false,
		},
		{
			name:    "OK - uint8 value",
			args:    args{valueType: common.ValueTypeUint8, param: &sdkModel.CommandValue{Value: uint8(5), Type: common.ValueTypeUint8}},
			want:    uint8(5),
			wantErr: false,
		},
		{
			name:    "OK - uint16 value",
			args:    args{valueType: common.ValueTypeUint16, param: &sdkModel.CommandValue{Value: uint16(5), Type: common.ValueTypeUint16}},
			want:    uint16(5),
			wantErr: false,
		},
		{
			name:    "OK - uint32 value",
			args:    args{valueType: common.ValueTypeUint32, param: &sdkModel.CommandValue{Value: uint32(5), Type: common.ValueTypeUint32}},
			want:    uint32(5),
			wantErr: false,
		},
		{
			name:    "OK - uint64 value",
			args:    args{valueType: common.ValueTypeUint64, param: &sdkModel.CommandValue{Value: uint64(5), Type: common.ValueTypeUint64}},
			want:    uint64(5),
			wantErr: false,
		},
		{
			name:    "OK - int8 value",
			args:    args{valueType: common.ValueTypeInt8, param: &sdkModel.CommandValue{Value: int8(5), Type: common.ValueTypeInt8}},
			want:    int8(5),
			wantErr: false,
		},
		{
			name:    "OK - int16 value",
			args:    args{valueType: common.ValueTypeInt16, param: &sdkModel.CommandValue{Value: int16(5), Type: common.ValueTypeInt16}},
			want:    int16(5),
			wantErr: false,
		},
		{
			name:    "OK - int32 value",
			args:    args{valueType: common.ValueTypeInt32, param: &sdkModel.CommandValue{Value: int32(5), Type: common.ValueTypeInt32}},
			want:    int32(5),
			wantErr: false,
		},
		{
			name:    "OK - int64 value",
			args:    args{valueType: common.ValueTypeInt64, param: &sdkModel.CommandValue{Value: int64(5), Type: common.ValueTypeInt64}},
			want:    int64(5),
			wantErr: false,
		},
		{
			name:    "OK - float32 value",
			args:    args{valueType: common.ValueTypeFloat32, param: &sdkModel.CommandValue{Value: float32(5), Type: common.ValueTypeFloat32}},
			want:    float32(5),
			wantErr: false,
		},
		{
			name:    "OK - float64 value",
			args:    args{valueType: common.ValueTypeFloat64, param: &sdkModel.CommandValue{Value: float64(5), Type: common.ValueTypeFloat64}},
			want:    float64(5),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewValue(tt.args.valueType, tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("newCommandValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newCommandValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
