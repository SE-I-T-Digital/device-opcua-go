// -*- Mode: Go; indent-tabs-mode: t -*-
//
// Copyright (C) 2018 Canonical Ltd
// Copyright (C) 2018 IOTech Ltd
// Copyright (C) 2021 Schneider Electric
//
// SPDX-License-Identifier: Apache-2.0

package result

import (
	"math"
	"strings"

	"github.com/edgexfoundry/go-mod-core-contracts/v4/common"
	"github.com/spf13/cast"
)

// checkValueInRange checks value range is valid
func checkValueInRange(valueType string, reading any) bool {
	isValid := false

	if valueType == common.ValueTypeString || valueType == common.ValueTypeBool ||
		valueType == common.ValueTypeBoolArray || valueType == common.ValueTypeStringArray {
		return true
	}

	if valueType == common.ValueTypeInt8 || valueType == common.ValueTypeInt16 ||
		valueType == common.ValueTypeInt32 || valueType == common.ValueTypeInt64 {
		val, err := cast.ToInt64E(reading)
		if err != nil {
			return false
		}
		isValid = checkIntValueRange(valueType, val)
	}

	if valueType == common.ValueTypeUint8 || valueType == common.ValueTypeUint16 ||
		valueType == common.ValueTypeUint32 || valueType == common.ValueTypeUint64 {
		val, err := cast.ToUint64E(reading)
		if err != nil {
			return false
		}
		isValid = checkUintValueRange(valueType, val)
	}

	if valueType == common.ValueTypeFloat32 || valueType == common.ValueTypeFloat64 {
		val, err := cast.ToFloat64E(reading)
		if err != nil {
			return false
		}
		isValid = checkFloatValueRange(valueType, val)
	}

	if valueType == common.ValueTypeInt8Array || valueType == common.ValueTypeInt16Array ||
		valueType == common.ValueTypeInt32Array || valueType == common.ValueTypeInt64Array {
		val, err := cast.ToInt64SliceE(reading)
		if err != nil {
			return false
		}
		isValid = checkIntSliceValueRange(valueType, val)
	}

	if valueType == common.ValueTypeUint8Array || valueType == common.ValueTypeUint16Array ||
		valueType == common.ValueTypeUint32Array || valueType == common.ValueTypeUint64Array {
		val, err := cast.ToUint64SliceE(reading)
		if err != nil {
			return false
		}
		isValid = checkUintSliceValueRange(valueType, val)
	}

	if valueType == common.ValueTypeFloat32Array || valueType == common.ValueTypeFloat64Array {
		val, err := cast.ToFloat64SliceE(reading)
		if err != nil {
			return false
		}
		isValid = checkFloatSliceValueRange(valueType, val)
	}

	return isValid
}

func checkUintValueRange(valueType string, val uint64) bool {
	var isValid bool
	switch valueType {
	case common.ValueTypeUint8:
		if val <= math.MaxUint8 {
			isValid = true
		}
	case common.ValueTypeUint16:
		if val <= math.MaxUint16 {
			isValid = true
		}
	case common.ValueTypeUint32:
		if val <= math.MaxUint32 {
			isValid = true
		}
	case common.ValueTypeUint64:
		maxiMum := uint64(math.MaxUint64)
		if val <= maxiMum {
			isValid = true
		}
	}
	return isValid
}

func checkIntValueRange(valueType string, val int64) bool {
	var isValid bool
	switch valueType {
	case common.ValueTypeInt8:
		if val >= math.MinInt8 && val <= math.MaxInt8 {
			isValid = true
		}
	case common.ValueTypeInt16:
		if val >= math.MinInt16 && val <= math.MaxInt16 {
			isValid = true
		}
	case common.ValueTypeInt32:
		if val >= math.MinInt32 && val <= math.MaxInt32 {
			isValid = true
		}
	case common.ValueTypeInt64:
		isValid = true
	}
	return isValid
}

func checkFloatValueRange(valueType string, val float64) bool {
	var isValid bool
	switch valueType {
	case common.ValueTypeFloat32:
		if !math.IsNaN(val) && math.Abs(val) <= math.MaxFloat32 {
			isValid = true
		}
	case common.ValueTypeFloat64:
		if !math.IsNaN(val) && !math.IsInf(val, 0) {
			isValid = true
		}
	}
	return isValid
}

const arrayTypeSuffix = "Array"

func checkUintSliceValueRange(valueType string, val []uint64) bool {
	var isValid bool
	baseType := strings.TrimSuffix(valueType, arrayTypeSuffix)
	for _, v := range val {
		if checkUintValueRange(baseType, v) {
			isValid = true
			break
		}
	}
	return isValid
}

func checkIntSliceValueRange(valueType string, val []int64) bool {
	var isValid bool
	baseType := strings.TrimSuffix(valueType, arrayTypeSuffix)
	for _, v := range val {
		if checkIntValueRange(baseType, v) {
			isValid = true
			break
		}
	}
	return isValid
}

func checkFloatSliceValueRange(valueType string, val []float64) bool {
	var isValid bool
	baseType := strings.TrimSuffix(valueType, arrayTypeSuffix)
	for _, v := range val {
		if checkFloatValueRange(baseType, v) {
			isValid = true
			break
		}
	}
	return isValid
}
