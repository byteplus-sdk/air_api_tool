package fieldchecker

import (
	"errors"
	"fmt"
	"github.com/3rd_rec/air_api_tool/checker"
	"github.com/3rd_rec/air_api_tool/consts"
	"github.com/3rd_rec/air_api_tool/shared"
	"strings"
)

func init() {
	checker.RegisterBaseCheckerFactory(consts.CheckerNameEnums, &EnumsChecker{})
}

type EnumsChecker struct {
	enumIndex    map[string]bool
	enumsStr     string
	errMsgFormat string
}

func (c *EnumsChecker) Create(checkerCfg *shared.CommonCheckerCfg) (checker.FieldChecker, error) {
	const originErrMsgFormat = "an unexpected value '{}' is included, only '#' can be accepted"
	if checkerCfg == nil {
		return nil, errors.New("invalid enums checker cfg")
	}
	enums := checkerCfg.Enums
	if len(enums) == 0 {
		return nil, errors.New("invalid enums checker cfg")
	}
	enumsStr := strings.Join(strings.Fields(fmt.Sprint(enums)), ",")
	enumIndex := make(map[string]bool, len(enums))
	for _, e := range enums {
		enumIndex[strings.ToLower(e)] = true
	}
	errMsgFormat := strings.ReplaceAll(originErrMsgFormat, "#", enumsStr)
	instance := &EnumsChecker{
		errMsgFormat: errMsgFormat,
		enumsStr:     enumsStr,
		enumIndex:    enumIndex,
	}
	return instance, nil
}

func (c *EnumsChecker) Description() string {
	return "Only the specified enumeration values are supported."
}

func (c *EnumsChecker) CheckerDetails() string {
	return fmt.Sprintf("enum: %s", c.enumsStr)
}

func (c *EnumsChecker) Check(fieldValueItr interface{}) error {
	if c == nil || fieldValueItr == nil {
		return nil
	}
	dataValue := fieldValueItr
	if boolValue, ok := fieldValueItr.(bool); ok {
		if boolValue {
			dataValue = 1
		} else {
			dataValue = 0
		}
	}
	dataString := fmt.Sprintf("%+v", dataValue)
	if !c.enumIndex[strings.ToLower(dataString)] {
		errMsg := strings.ReplaceAll(c.errMsgFormat, "{}", dataString)
		return errors.New(errMsg)
	}
	return nil
}
