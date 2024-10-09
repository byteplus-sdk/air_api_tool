package fieldchecker

import (
	"errors"
	"fmt"
	"github.com/3rd_rec/air_api_tool/checker"
	"github.com/3rd_rec/air_api_tool/consts"
	"github.com/3rd_rec/air_api_tool/shared"
	"github.com/3rd_rec/air_api_tool/utils"
	"strings"
)

func init() {
	checker.RegisterBaseCheckerFactory(consts.CheckerNameRanges, &RangesChecker{})
}

type RangesChecker struct {
	ranges       []*RangeChecker
	rangesMsg    string
	errMsgFormat string
}

func (c *RangesChecker) Create(checkerCfg *shared.CommonCheckerCfg) (checker.FieldChecker, error) {
	const errMsgOriginFormat = "An unexpected value '{}' is included, only '%v' can be accepted"
	if checkerCfg == nil {
		return nil, errors.New("invalid range checker cfg")
	}
	rangeCfgList := checkerCfg.RangeCfgList
	if len(rangeCfgList) == 0 {
		return nil, errors.New("invalid range checker cfg")
	}
	checkers := make([]*RangeChecker, len(rangeCfgList))
	rangeMsgs := make([]string, len(rangeCfgList))
	for i, cfg := range rangeCfgList {
		checkers[i] = &RangeChecker{cfg: cfg}
		rangeMsgs[i] = "(" + cfg.String() + ")"
	}
	rangesMsg := strings.Join(rangeMsgs, " or ")
	errMsgFormat := fmt.Sprintf(errMsgOriginFormat, rangesMsg)
	instance := &RangesChecker{
		ranges:       checkers,
		rangesMsg:    rangesMsg,
		errMsgFormat: errMsgFormat,
	}
	return instance, nil
}

func (c *RangesChecker) Description() string {
	return "The field value must satisfy the range condition."
}

func (c *RangesChecker) CheckerDetails() string {
	return fmt.Sprintf("ranges: %s", c.rangesMsg)
}

func (c *RangesChecker) Check(fieldValue interface{}) error {
	if c == nil || fieldValue == nil {
		return nil
	}
	var num float64
	switch t := fieldValue.(type) {
	case uint64:
		num = float64(t)
	case int64:
		num = float64(t)
	case int32:
		num = float64(t)
	case float64:
		num = t
	default: // only check number type
		return nil
	}
	for _, rangeChecker := range c.ranges {
		// Passing check requires meeting only one range
		if rangeChecker.check(num) {
			return nil
		}
	}
	errMsg := strings.ReplaceAll(c.errMsgFormat, "{}", utils.Float64ToA(num))
	return errors.New(errMsg)
}

type RangeChecker struct {
	cfg *shared.RangeCfg
}

func (c *RangeChecker) check(num float64) bool {
	cfg := c.cfg
	if cfg.Lt != nil {
		if num >= *cfg.Lt {
			return false
		}
	}
	if cfg.Lte != nil {
		if num > *cfg.Lte {
			return false
		}
	}
	if cfg.Gt != nil {
		if num <= *cfg.Gt {
			return false
		}
	}
	if cfg.Gte != nil {
		if num < *cfg.Gte {
			return false
		}
	}
	return true
}

func (c *RangesChecker) MatchConditions(_ map[string]interface{}) bool {
	return true
}
