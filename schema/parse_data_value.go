package schema

import (
	"encoding/json"
	"strings"

	"github.com/spf13/cast"
)

func isEmptyString(value interface{}) bool {
	valueStr, ok := value.(string)
	if !ok {
		return false
	}
	return valueStr == ""
}

func parseAsFloat64(value interface{}) (float64, bool) {
	if isEmptyString(value) {
		return 0, true
	}

	var (
		float64Val float64
		err        error
	)
	jsonNum, ok := parseAsJsonNumber(value)
	if ok {
		float64Val, err = jsonNum.Float64()
		if err != nil {
			return 0, false
		}
		return float64Val, true
	}
	float64Val, err = cast.ToFloat64E(value)
	if err != nil {
		return float64Val, false
	}
	return float64Val, true
}

func parseAsInt64(value interface{}) (int64, bool) {
	// 空字符串解析为0值
	if isEmptyString(value) {
		return 0, true
	}

	var (
		int64Val int64
		err      error
	)
	jsonNum, ok := parseAsJsonNumber(value)
	if ok {
		int64Val, err = jsonNum.Int64()
		if err != nil {
			return 0, false
		}
		return int64Val, true
	}
	int64Val, err = cast.ToInt64E(value)
	if err != nil {
		return int64Val, false
	}
	return int64Val, true
}

func parseAsJsonNumber(value interface{}) (json.Number, bool) {
	jsonNum, ok := value.(json.Number)
	if !ok {
		jsonNumStr, ok := value.(string)
		if !ok {
			return jsonNum, false
		}
		jsonNum = json.Number(jsonNumStr)
	}
	jsonNum = json.Number(FormatScientificString(jsonNum.String()))
	return jsonNum, true
}

func parseAsBool(value interface{}) (bool, bool) {
	// 空字符串解析为0值
	if isEmptyString(value) {
		return false, true
	}

	boolVal, ok := value.(bool)
	if ok {
		return boolVal, true
	}
	// adapt "true"/"false", 0/1
	boolValString, ok := value.(string)
	if ok {
		if strings.EqualFold(boolValString, "true") ||
			strings.EqualFold(boolValString, "1") {
			return true, true
		}
		if strings.EqualFold(boolValString, "false") ||
			strings.EqualFold(boolValString, "0") {
			return false, true
		}
	}
	int64Value, ok := parseAsInt64(value)
	if ok {
		return int64Value == 1, true
	}
	return boolVal, false
}

func parseAsString(value interface{}) (string, bool) {
	strVal, ok := value.(string)
	if !ok {
		return "", false
	}
	return strVal, true
}

func parseJSONString(value interface{}) (string, bool) {
	strVal, ok := value.(string)
	if !ok {
		return "", false
	}
	// json array string
	if strings.HasPrefix(strVal, "[") {
		array := make([]interface{}, 0)
		err := json.Unmarshal([]byte(strVal), &array)
		if err != nil {
			return "", false
		}
		return strVal, true
	}
	// json object string
	if strings.HasPrefix(strVal, "{") {
		object := make(map[string]interface{})
		err := json.Unmarshal([]byte(strVal), &object)
		if err != nil {
			return "", false
		}
		return strVal, true
	}
	return "", false
}
