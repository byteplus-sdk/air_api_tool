package utils

import (
	"encoding/json"
	"strconv"
)

const base10 int = 10
const int64BitSize int = 64
const uint64BitSize int = 64
const int32BitSize int = 32
const uint32BitSize int = 32
const float32BitSize int = 32
const float64BitSize int = 64
const withoutExponent = -1

func ItoA(num int) string {
	return strconv.Itoa(num)
}

func Int32ToA(num int32) string {
	return strconv.FormatInt(int64(num), base10)
}

func Int64ToA(num int64) string {
	return strconv.FormatInt(int64(num), base10)
}

func Uint32ToA(num uint32) string {
	return strconv.FormatUint(uint64(num), base10)
}

func Uint64ToA(num uint64) string {
	return strconv.FormatUint(uint64(num), base10)
}

func Float32ToA(num float32) string {
	return strconv.FormatFloat(float64(num), 'f', withoutExponent, float32BitSize)
}

func Float64ToA(num float64) string {
	return strconv.FormatFloat(num, 'f', withoutExponent, float64BitSize)
}

func BoolToA(s bool) string {
	return strconv.FormatBool(s)
}

func BoolToNumberStr(s bool) string {
	if s {
		return "1"
	}
	return "0"
}

func ToA(source interface{}) string {
	switch value := source.(type) {
	case string:
		return value
	case float64:
		return Float64ToA(value)
	case int64:
		return Int64ToA(value)
	case bool:
		return BoolToA(value)
	case float32:
		return Float32ToA(value)
	case int:
		return ItoA(value)
	case int32:
		return Int32ToA(value)
	case int16:
		return Int32ToA(int32(value))
	case int8:
		return Int32ToA(int32(value))
	case uint64:
		return Uint64ToA(value)
	case uint:
		return Uint64ToA(uint64(value))
	case uint32:
		return Uint32ToA(value)
	case uint16:
		return Uint32ToA(uint32(value))
	case uint8:
		return Uint32ToA(uint32(value))
	case nil:
		return "nil"
	default:
		result, _ := json.Marshal(source)
		return string(result)
	}
}
