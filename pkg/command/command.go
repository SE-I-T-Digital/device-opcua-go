// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2018 Canonical Ltd
// Copyright (C) 2018 IOTech Ltd
// Copyright (C) 2021 Schneider Electric
//
// SPDX-License-Identifier: Apache-2.0

package command

import (
	"fmt"

	sdkModel "github.com/edgexfoundry/device-sdk-go/v4/pkg/models"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/common"
)

func NewValue(valueType string, param *sdkModel.CommandValue) (any, error) {
	var commandValue any
	var err error
	switch valueType {
	case common.ValueTypeBool:
		commandValue, err = param.BoolValue()
	case common.ValueTypeString:
		commandValue, err = param.StringValue()
	case common.ValueTypeUint8:
		commandValue, err = param.Uint8Value()
	case common.ValueTypeUint16:
		commandValue, err = param.Uint16Value()
	case common.ValueTypeUint32:
		commandValue, err = param.Uint32Value()
	case common.ValueTypeUint64:
		commandValue, err = param.Uint64Value()
	case common.ValueTypeInt8:
		commandValue, err = param.Int8Value()
	case common.ValueTypeInt16:
		commandValue, err = param.Int16Value()
	case common.ValueTypeInt32:
		commandValue, err = param.Int32Value()
	case common.ValueTypeInt64:
		commandValue, err = param.Int64Value()
	case common.ValueTypeFloat32:
		commandValue, err = param.Float32Value()
	case common.ValueTypeFloat64:
		commandValue, err = param.Float64Value()
	case common.ValueTypeBoolArray:
		commandValue, err = param.BoolArrayValue()
	case common.ValueTypeStringArray:
		commandValue, err = param.StringArrayValue()
	case common.ValueTypeUint8Array:
		commandValue, err = param.Uint8ArrayValue()
	case common.ValueTypeUint16Array:
		commandValue, err = param.Uint16ArrayValue()
	case common.ValueTypeUint32Array:
		commandValue, err = param.Uint32ArrayValue()
	case common.ValueTypeUint64Array:
		commandValue, err = param.Uint64ArrayValue()
	case common.ValueTypeInt8Array:
		commandValue, err = param.Int8ArrayValue()
	case common.ValueTypeInt16Array:
		commandValue, err = param.Int16ArrayValue()
	case common.ValueTypeInt32Array:
		commandValue, err = param.Int32ArrayValue()
	case common.ValueTypeInt64Array:
		commandValue, err = param.Int64ArrayValue()
	case common.ValueTypeFloat32Array:
		commandValue, err = param.Float32ArrayValue()
	case common.ValueTypeFloat64Array:
		commandValue, err = param.Float64ArrayValue()
	default:
		err = fmt.Errorf("fail to convert param, none supported value type: %v", valueType)
	}

	return commandValue, err
}
