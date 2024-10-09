package fieldchecker

import (
	"errors"
	"github.com/3rd_rec/air_api_tool/checker"
	"github.com/3rd_rec/air_api_tool/consts"
	"github.com/3rd_rec/air_api_tool/shared"
)

func init() {
	checker.RegisterBaseCheckerFactory(consts.CheckerNameRequired, &required{})
}

type required struct {
	required bool
	err      error
}

func (c *required) Create(checkerCfg *shared.CommonCheckerCfg) (checker.FieldChecker, error) {
	if checkerCfg == nil {
		return nil, errors.New("invalid required checker cfg")
	}
	instance := &required{
		required: checkerCfg.Required,
		err:      errors.New("required field is empty"),
	}
	return instance, nil
}

func (c *required) Description() string {
	return "Field values are not allowed to be empty."
}

func (c *required) CheckerDetails() string {
	return "required: true"
}

func (c *required) Check(fieldValue interface{}) error {
	if c == nil || !c.required {
		return nil
	}
	if fieldValue == nil {
		return c.err
	}
	dataStr, ok := fieldValue.(string)
	if !ok {
		return nil
	}
	if dataStr == "" {
		return c.err
	}
	return nil
}
