// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2019-2021 IOTech Ltd
// Copyright (C) 2021 Schneider Electric
//
// SPDX-License-Identifier: Apache-2.0

package result

import (
	"fmt"
	"strings"
	"testing"

	"github.com/edgexfoundry/device-sdk-go/v4/pkg/models"
	"github.com/edgexfoundry/go-mod-core-contracts/v4/common"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// All TestNewResult* functions used from
// https://github.com/edgexfoundry/device-mqtt-go/blob/2edd794ead44ee1cbb795ccb8ebf3dc377aa3945/internal/driver/driver_test.go

func TestNewResult_bool(t *testing.T) {
	var reading any = true
	req := models.CommandRequest{
		DeviceResourceName: "light",
		Type:               common.ValueTypeBool,
	}

	cmdVal, err := NewResult(req, reading)
	if err != nil {
		t.Fatalf("Fail to create new ReadCommand result, %v", err)
	}
	val, err := cmdVal.BoolValue()
	if val != true || err != nil {
		t.Errorf("Convert new result(%v) failed, error: %v", val, err)
	}
}

func TestNewResult_uint8(t *testing.T) {
	var reading any = uint8(123)
	req := models.CommandRequest{
		DeviceResourceName: "temperature",
		Type:               common.ValueTypeUint8,
	}

	cmdVal, err := NewResult(req, reading)
	if err != nil {
		t.Fatalf("Fail to create new ReadCommand result, %v", err)
	}
	val, err := cmdVal.Uint8Value()
	if val != uint8(123) || err != nil {
		t.Errorf("Convert new result(%v) failed, error: %v", val, err)
	}
}

func TestNewResult_int8(t *testing.T) {
	var reading any = int8(123)
	req := models.CommandRequest{
		DeviceResourceName: "temperature",
		Type:               common.ValueTypeInt8,
	}

	cmdVal, err := NewResult(req, reading)
	if err != nil {
		t.Fatalf("Fail to create new ReadCommand result, %v", err)
	}
	val, err := cmdVal.Int8Value()
	if val != int8(123) || err != nil {
		t.Errorf("Convert new result(%v) failed, error: %v", val, err)
	}
}

func TestNewResultFailed_int8(t *testing.T) {
	var reading any = int16(256)
	req := models.CommandRequest{
		DeviceResourceName: "temperature",
		Type:               common.ValueTypeInt8,
	}

	_, err := NewResult(req, reading)
	fmt.Println(err)
	if err == nil || !strings.Contains(err.Error(), "Reading 256 is out of the value type(Int8)'s range") {
		t.Errorf("Convert new result should be failed")
	}
}

func TestNewResult_uint16(t *testing.T) {
	var reading any = uint16(123)
	req := models.CommandRequest{
		DeviceResourceName: "temperature",
		Type:               common.ValueTypeUint16,
	}

	cmdVal, err := NewResult(req, reading)
	if err != nil {
		t.Fatalf("Fail to create new ReadCommand result, %v", err)
	}
	val, err := cmdVal.Uint16Value()
	if val != uint16(123) || err != nil {
		t.Errorf("Convert new result(%v) failed, error: %v", val, err)
	}
}

func TestNewResult_int16(t *testing.T) {
	var reading any = int16(123)
	req := models.CommandRequest{
		DeviceResourceName: "temperature",
		Type:               common.ValueTypeInt16,
	}

	cmdVal, err := NewResult(req, reading)
	if err != nil {
		t.Fatalf("Fail to create new ReadCommand result, %v", err)
	}
	val, err := cmdVal.Int16Value()
	if val != int16(123) || err != nil {
		t.Errorf("Convert new result(%v) failed, error: %v", val, err)
	}
}

func TestNewResult_uint32(t *testing.T) {
	var reading any = uint32(123)
	req := models.CommandRequest{
		DeviceResourceName: "temperature",
		Type:               common.ValueTypeUint32,
	}

	cmdVal, err := NewResult(req, reading)
	if err != nil {
		t.Fatalf("Fail to create new ReadCommand result, %v", err)
	}
	val, err := cmdVal.Uint32Value()
	if val != uint32(123) || err != nil {
		t.Errorf("Convert new result(%v) failed, error: %v", val, err)
	}
}

func TestNewResult_int32(t *testing.T) {
	var reading any = int32(123)
	req := models.CommandRequest{
		DeviceResourceName: "temperature",
		Type:               common.ValueTypeInt32,
	}

	cmdVal, err := NewResult(req, reading)
	if err != nil {
		t.Fatalf("Fail to create new ReadCommand result, %v", err)
	}
	val, err := cmdVal.Int32Value()
	if val != int32(123) || err != nil {
		t.Errorf("Convert new result(%v) failed, error: %v", val, err)
	}
}

func TestNewResult_uint64(t *testing.T) {
	var reading any = uint64(123)
	req := models.CommandRequest{
		DeviceResourceName: "temperature",
		Type:               common.ValueTypeUint64,
	}

	cmdVal, err := NewResult(req, reading)
	if err != nil {
		t.Fatalf("Fail to create new ReadCommand result, %v", err)
	}
	val, err := cmdVal.Uint64Value()
	if val != uint64(123) || err != nil {
		t.Errorf("Convert new result(%v) failed, error: %v", val, err)
	}
}

func TestNewResult_int64(t *testing.T) {
	var reading any = int64(123)
	req := models.CommandRequest{
		DeviceResourceName: "temperature",
		Type:               common.ValueTypeInt64,
	}

	cmdVal, err := NewResult(req, reading)
	if err != nil {
		t.Fatalf("Fail to create new ReadCommand result, %v", err)
	}
	val, err := cmdVal.Int64Value()
	if val != int64(123) || err != nil {
		t.Errorf("Convert new result(%v) failed, error: %v", val, err)
	}
}

func TestNewResult_float32(t *testing.T) {
	req := models.CommandRequest{
		DeviceResourceName: "temperature",
		Type:               common.ValueTypeFloat32,
	}

	tests := []struct {
		name     string
		req      models.CommandRequest
		reading  any
		expected any
	}{
		{"Valid string 0", req, "0", float32(0)},
		{"Valid string 123.32", req, "123.32", float32(123.32)},
		{"Valid float32 0", req, 0, float32(0)},
		{"Valid float32 123.32", req, 123.32, float32(123.32)},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			cmdVal, err := NewResult(req, testCase.reading)
			require.NoError(t, err)
			result, err := cmdVal.Float32Value()
			require.NoError(t, err)

			assert.Equal(t, testCase.expected, result)
		})
	}
}

func TestNewResult_float64(t *testing.T) {
	req := models.CommandRequest{
		DeviceResourceName: "temperature",
		Type:               common.ValueTypeFloat64,
	}

	tests := []struct {
		name     string
		req      models.CommandRequest
		reading  any
		expected any
	}{
		{"Valid string 0", req, "0", float64(0)},
		{"Valid string 0.123", req, "0.123", 0.123},
		{"Valid float64 0", req, 0, float64(0)},
		{"Valid float64 0.123", req, 0.123, 0.123},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			cmdVal, err := NewResult(req, testCase.reading)
			require.NoError(t, err)
			result, err := cmdVal.Float64Value()
			require.NoError(t, err)

			assert.Equal(t, testCase.expected, result)
		})
	}
}

func TestNewResult_float64ToInt8(t *testing.T) {
	var reading any = float64(123.11)
	req := models.CommandRequest{
		DeviceResourceName: "temperature",
		Type:               common.ValueTypeInt8,
	}

	cmdVal, err := NewResult(req, reading)
	if err != nil {
		t.Fatalf("Fail to create new ReadCommand result, %v", err)
	}
	val, err := cmdVal.Int8Value()
	if val != int8(reading.(float64)) || err != nil {
		t.Errorf("Convert new result(%v) failed, error: %v", val, err)
	}
}

func TestNewResult_float64ToInt16(t *testing.T) {
	var reading any = float64(123.11)
	req := models.CommandRequest{
		DeviceResourceName: "temperature",
		Type:               common.ValueTypeInt16,
	}

	cmdVal, err := NewResult(req, reading)
	if err != nil {
		t.Fatalf("Fail to create new ReadCommand result, %v", err)
	}
	val, err := cmdVal.Int16Value()
	if val != int16(reading.(float64)) || err != nil {
		t.Errorf("Convert new result(%v) failed, error: %v", val, err)
	}
}

func TestNewResult_float64ToInt32(t *testing.T) {
	var reading any = float64(123.11)
	req := models.CommandRequest{
		DeviceResourceName: "temperature",
		Type:               common.ValueTypeInt32,
	}

	cmdVal, err := NewResult(req, reading)
	if err != nil {
		t.Fatalf("Fail to create new ReadCommand result, %v", err)
	}
	val, err := cmdVal.Int32Value()
	if val != int32(reading.(float64)) || err != nil {
		t.Errorf("Convert new result(%v) failed, error: %v", val, err)
	}
}

func TestNewResult_float64ToInt64(t *testing.T) {
	var reading any = float64(123.11)
	req := models.CommandRequest{
		DeviceResourceName: "temperature",
		Type:               common.ValueTypeInt64,
	}

	cmdVal, err := NewResult(req, reading)
	if err != nil {
		t.Fatalf("Fail to create new ReadCommand result, %v", err)
	}
	val, err := cmdVal.Int64Value()
	if val != int64(reading.(float64)) || err != nil {
		t.Errorf("Convert new result(%v) failed, error: %v", val, err)
	}
}

func TestNewResult_float64ToUint8(t *testing.T) {
	var reading any = float64(123.11)
	req := models.CommandRequest{
		DeviceResourceName: "temperature",
		Type:               common.ValueTypeUint8,
	}

	cmdVal, err := NewResult(req, reading)
	if err != nil {
		t.Fatalf("Fail to create new ReadCommand result, %v", err)
	}
	val, err := cmdVal.Uint8Value()
	if val != uint8(reading.(float64)) || err != nil {
		t.Errorf("Convert new result(%v) failed, error: %v", val, err)
	}
}

func TestNewResult_float64ToUint16(t *testing.T) {
	var reading any = float64(123.11)
	req := models.CommandRequest{
		DeviceResourceName: "temperature",
		Type:               common.ValueTypeUint16,
	}

	cmdVal, err := NewResult(req, reading)
	if err != nil {
		t.Fatalf("Fail to create new ReadCommand result, %v", err)
	}
	val, err := cmdVal.Uint16Value()
	if val != uint16(reading.(float64)) || err != nil {
		t.Errorf("Convert new result(%v) failed, error: %v", val, err)
	}
}

func TestNewResult_float64ToUint32(t *testing.T) {
	var reading any = float64(123.11)
	req := models.CommandRequest{
		DeviceResourceName: "temperature",
		Type:               common.ValueTypeUint32,
	}

	cmdVal, err := NewResult(req, reading)
	if err != nil {
		t.Fatalf("Fail to create new ReadCommand result, %v", err)
	}
	val, err := cmdVal.Uint32Value()
	if val != uint32(reading.(float64)) || err != nil {
		t.Errorf("Convert new result(%v) failed, error: %v", val, err)
	}
}

func TestNewResult_float64ToUint64(t *testing.T) {
	var reading any = float64(123.11)
	req := models.CommandRequest{
		DeviceResourceName: "temperature",
		Type:               common.ValueTypeUint64,
	}

	cmdVal, err := NewResult(req, reading)
	if err != nil {
		t.Fatalf("Fail to create new ReadCommand result, %v", err)
	}
	val, err := cmdVal.Uint64Value()
	if val != uint64(reading.(float64)) || err != nil {
		t.Errorf("Convert new result(%v) failed, error: %v", val, err)
	}
}

func TestNewResult_float64ToFloat32(t *testing.T) {
	var reading any = float64(123.11)
	req := models.CommandRequest{
		DeviceResourceName: "temperature",
		Type:               common.ValueTypeFloat32,
	}

	cmdVal, err := NewResult(req, reading)
	if err != nil {
		t.Fatalf("Fail to create new ReadCommand result, %v", err)
	}
	val, err := cmdVal.Float32Value()
	if val != float32(reading.(float64)) || err != nil {
		t.Errorf("Convert new result(%v) failed, error: %v", val, err)
	}
}

func TestNewResult_float64ToString(t *testing.T) {
	var reading any = float64(123.11)
	req := models.CommandRequest{
		DeviceResourceName: "temperature",
		Type:               common.ValueTypeString,
	}

	cmdVal, err := NewResult(req, reading)
	if err != nil {
		t.Fatalf("Fail to create new ReadCommand result, %v", err)
	}
	val, err := cmdVal.StringValue()
	if val != fmt.Sprintf("%v", reading) || err != nil {
		t.Errorf("Convert new result(%v) failed, error: %v", val, err)
	}
}

func TestNewResult_string(t *testing.T) {
	var reading any = "test string"
	req := models.CommandRequest{
		DeviceResourceName: "temperature",
		Type:               common.ValueTypeString,
	}

	cmdVal, err := NewResult(req, reading)
	if err != nil {
		t.Fatalf("Fail to create new ReadCommand result, %v", err)
	}
	val, err := cmdVal.StringValue()
	if val != "test string" || err != nil {
		t.Errorf("Convert new result(%v) failed, error: %v", val, err)
	}
}

func TestNewResult_stringToInt64(t *testing.T) {
	var reading any = "123"
	req := models.CommandRequest{
		DeviceResourceName: "temperature",
		Type:               common.ValueTypeInt64,
	}

	cmdVal, err := NewResult(req, reading)
	if err != nil {
		t.Fatalf("Fail to create new ReadCommand result, %v", err)
	}
	val, err := cmdVal.Int64Value()
	if val != int64(123) || err != nil {
		t.Errorf("Convert new result(%v) failed, error: %v", val, err)
	}
}

func TestNewResult_stringToInt8(t *testing.T) {
	var reading any = "123"
	req := models.CommandRequest{
		DeviceResourceName: "temperature",
		Type:               common.ValueTypeInt8,
	}

	cmdVal, err := NewResult(req, reading)
	if err != nil {
		t.Fatalf("Fail to create new ReadCommand result, %v", err)
	}
	val, err := cmdVal.Int8Value()
	if val != int8(123) || err != nil {
		t.Errorf("Convert new result(%v) failed, error: %v", val, err)
	}
}

func TestNewResult_stringToUint8(t *testing.T) {
	var reading any = "123"
	req := models.CommandRequest{
		DeviceResourceName: "temperature",
		Type:               common.ValueTypeUint8,
	}

	cmdVal, err := NewResult(req, reading)
	if err != nil {
		t.Fatalf("Fail to create new ReadCommand result, %v", err)
	}
	val, err := cmdVal.Uint8Value()
	if val != uint8(123) || err != nil {
		t.Errorf("Convert new result(%v) failed, error: %v", val, err)
	}
}

func TestNewResult_stringToUint32(t *testing.T) {
	var reading any = "123"
	req := models.CommandRequest{
		DeviceResourceName: "temperature",
		Type:               common.ValueTypeUint32,
	}

	cmdVal, err := NewResult(req, reading)
	if err != nil {
		t.Fatalf("Fail to create new ReadCommand result, %v", err)
	}
	val, err := cmdVal.Uint32Value()
	if val != uint32(123) || err != nil {
		t.Errorf("Convert new result(%v) failed, error: %v", val, err)
	}
}

func TestNewResult_stringToUint64(t *testing.T) {
	var reading any = "123"
	req := models.CommandRequest{
		DeviceResourceName: "temperature",
		Type:               common.ValueTypeUint64,
	}

	cmdVal, err := NewResult(req, reading)
	if err != nil {
		t.Fatalf("Fail to create new ReadCommand result, %v", err)
	}
	val, err := cmdVal.Uint64Value()
	if val != uint64(123) || err != nil {
		t.Errorf("Convert new result(%v) failed, error: %v", val, err)
	}
}

func TestNewResult_stringToBool(t *testing.T) {
	var reading any = "true"
	req := models.CommandRequest{
		DeviceResourceName: "temperature",
		Type:               common.ValueTypeBool,
	}

	cmdVal, err := NewResult(req, reading)
	if err != nil {
		t.Fatalf("Fail to create new ReadCommand result, %v", err)
	}
	val, err := cmdVal.BoolValue()
	if val != true || err != nil {
		t.Errorf("Convert new result(%v) failed, error: %v", val, err)
	}
}

func TestNewResult_numberToUint64(t *testing.T) {
	var reading any = 123
	req := models.CommandRequest{
		DeviceResourceName: "temperature",
		Type:               common.ValueTypeUint64,
	}

	cmdVal, err := NewResult(req, reading)
	if err != nil {
		t.Fatalf("Fail to create new ReadCommand result, %v", err)
	}
	val, err := cmdVal.Uint64Value()
	if val != uint64(123) || err != nil {
		t.Errorf("Convert new result(%v) failed, error: %v", val, err)
	}
}

func TestNewResult_floatNumberToFloat32(t *testing.T) {
	var reading any = 123.0
	req := models.CommandRequest{
		DeviceResourceName: "temperature",
		Type:               common.ValueTypeFloat32,
	}

	cmdVal, err := NewResult(req, reading)
	if err != nil {
		t.Fatalf("Fail to create new ReadCommand result, %v", err)
	}
	val, err := cmdVal.Float32Value()
	if val != float32(123) || err != nil {
		t.Errorf("Convert new result(%v) failed, error: %v", val, err)
	}
}

func TestNewResult_numberToString(t *testing.T) {
	var reading any = 123
	req := models.CommandRequest{
		DeviceResourceName: "temperature",
		Type:               common.ValueTypeString,
	}

	cmdVal, err := NewResult(req, reading)
	if err != nil {
		t.Fatalf("Fail to create new ReadCommand result, %v", err)
	}
	val, err := cmdVal.StringValue()
	if val != "123" || err != nil {
		t.Errorf("Convert new result(%v) failed, error: %v", val, err)
	}
}

func TestNewResult_boolArray(t *testing.T) {
	var reading any = []bool{true, false, true}
	req := models.CommandRequest{
		DeviceResourceName: "light",
		Type:               common.ValueTypeBoolArray,
	}

	cmdVal, err := NewResult(req, reading)
	require.NoError(t, err)

	val, err := cmdVal.BoolArrayValue()
	assert.Equal(t, val, []bool{true, false, true})
	assert.NoError(t, err)
}

func TestNewResult_stringArray(t *testing.T) {
	var reading any = []string{"a", "b", "c"}
	req := models.CommandRequest{
		DeviceResourceName: "stringArray",
		Type:               common.ValueTypeStringArray,
	}

	cmdVal, err := NewResult(req, reading)
	require.NoError(t, err)

	val, err := cmdVal.StringArrayValue()
	assert.Equal(t, val, []string{"a", "b", "c"})
	assert.NoError(t, err)
}

func TestNewResult_uint8Array(t *testing.T) {
	var reading any = []uint8{1, 2, 3}
	req := models.CommandRequest{
		DeviceResourceName: "uint8Array",
		Type:               common.ValueTypeUint8Array,
	}

	cmdVal, err := NewResult(req, reading)
	require.NoError(t, err)

	val, err := cmdVal.Uint8ArrayValue()
	assert.Equal(t, val, []uint8{1, 2, 3})
	assert.NoError(t, err)
}

func TestNewResult_uint16Array(t *testing.T) {
	var reading any = []uint16{1, 2, 3}
	req := models.CommandRequest{
		DeviceResourceName: "uint16Array",
		Type:               common.ValueTypeUint16Array,
	}

	cmdVal, err := NewResult(req, reading)
	require.NoError(t, err)

	val, err := cmdVal.Uint16ArrayValue()
	assert.Equal(t, val, []uint16{1, 2, 3})
	assert.NoError(t, err)
}

func TestNewResult_uint32Array(t *testing.T) {
	var reading any = []uint32{1, 2, 3}
	req := models.CommandRequest{
		DeviceResourceName: "uint32Array",
		Type:               common.ValueTypeUint32Array,
	}

	cmdVal, err := NewResult(req, reading)
	require.NoError(t, err)

	val, err := cmdVal.Uint32ArrayValue()
	assert.Equal(t, val, []uint32{1, 2, 3})
	assert.NoError(t, err)
}

func TestNewResult_uint64Array(t *testing.T) {
	var reading any = []uint64{1, 2, 3}
	req := models.CommandRequest{
		DeviceResourceName: "uint64Array",
		Type:               common.ValueTypeUint64Array,
	}

	cmdVal, err := NewResult(req, reading)
	require.NoError(t, err)

	val, err := cmdVal.Uint64ArrayValue()
	assert.Equal(t, val, []uint64{1, 2, 3})
	assert.NoError(t, err)
}

func TestNewResult_int8Array(t *testing.T) {
	var reading any = []int8{1, 2, 3}
	req := models.CommandRequest{
		DeviceResourceName: "int8Array",
		Type:               common.ValueTypeInt8Array,
	}

	cmdVal, err := NewResult(req, reading)
	require.NoError(t, err)

	val, err := cmdVal.Int8ArrayValue()
	assert.Equal(t, val, []int8{1, 2, 3})
	assert.NoError(t, err)
}

func TestNewResult_int16Array(t *testing.T) {
	var reading any = []int16{1, 2, 3}
	req := models.CommandRequest{
		DeviceResourceName: "int16Array",
		Type:               common.ValueTypeInt16Array,
	}

	cmdVal, err := NewResult(req, reading)
	require.NoError(t, err)

	val, err := cmdVal.Int16ArrayValue()
	assert.Equal(t, val, []int16{1, 2, 3})
	assert.NoError(t, err)
}

func TestNewResult_int32Array(t *testing.T) {
	var reading any = []int32{1, 2, 3}
	req := models.CommandRequest{
		DeviceResourceName: "int32Array",
		Type:               common.ValueTypeInt32Array,
	}

	cmdVal, err := NewResult(req, reading)
	require.NoError(t, err)

	val, err := cmdVal.Int32ArrayValue()
	assert.Equal(t, val, []int32{1, 2, 3})
	assert.NoError(t, err)
}

func TestNewResult_int64Array(t *testing.T) {
	var reading any = []int64{1, 2, 3}
	req := models.CommandRequest{
		DeviceResourceName: "int64Array",
		Type:               common.ValueTypeInt64Array,
	}

	cmdVal, err := NewResult(req, reading)
	require.NoError(t, err)

	val, err := cmdVal.Int64ArrayValue()
	assert.Equal(t, val, []int64{1, 2, 3})
	assert.NoError(t, err)
}

func TestNewResult_float32Array(t *testing.T) {
	var reading any = []float32{1.1, 2.2, 3.3}
	req := models.CommandRequest{
		DeviceResourceName: "float32Array",
		Type:               common.ValueTypeFloat32Array,
	}

	cmdVal, err := NewResult(req, reading)
	require.NoError(t, err)

	val, err := cmdVal.Float32ArrayValue()
	assert.Equal(t, val, []float32{1.1, 2.2, 3.3})
	assert.NoError(t, err)
}

func TestNewResult_float64Array(t *testing.T) {
	var reading any = []float64{1.1, 2.2, 3.3}
	req := models.CommandRequest{
		DeviceResourceName: "float64Array",
		Type:               common.ValueTypeFloat64Array,
	}

	cmdVal, err := NewResult(req, reading)
	require.NoError(t, err)

	val, err := cmdVal.Float64ArrayValue()
	assert.Equal(t, val, []float64{1.1, 2.2, 3.3})
	assert.NoError(t, err)
}
