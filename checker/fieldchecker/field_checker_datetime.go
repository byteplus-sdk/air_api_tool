package fieldchecker

import (
	"errors"
	"fmt"
	"github.com/3rd_rec/air_api_tool/checker"
	"github.com/3rd_rec/air_api_tool/consts"
	"github.com/3rd_rec/air_api_tool/shared"
	"strings"
	"time"

	"github.com/spf13/cast"
)

func init() {
	checker.RegisterBaseCheckerFactory(consts.CheckerNameDatetime, &DatetimeChecker{})
}

type DatetimeChecker struct {
	errMessageFormat string
	timeLayout       string
	isInt64          bool
}

func (c *DatetimeChecker) Create(checkerCfg *shared.CommonCheckerCfg) (checker.FieldChecker, error) {
	if checkerCfg == nil {
		return nil, errors.New("invalid datetime checker cfg")
	}
	cfg := checkerCfg.DatetimeCfg
	if cfg == nil {
		return nil, errors.New("invalid datetime checker cfg")
	}
	instance := &DatetimeChecker{
		timeLayout: cfg.TimeLayout,
		isInt64:    cfg.IsInt64,
	}
	if !instance.isInt64 && instance.timeLayout == "" {
		instance.timeLayout = time.RFC3339
	}
	if instance.isInt64 {
		instance.errMessageFormat = "the value(%+v) is invalid datetime. The field value must be a legal timestamp (second timestamp>1000000000 or millisecond timestamp>1000000000000). It supports second timestamp or millisecond timestamp. It is recommended to use second timestamp."
	} else {
		const originErrMsgFormat = "the value(%+v) is invalid datetime,  please format it like '#'"
		instance.errMessageFormat = strings.ReplaceAll(originErrMsgFormat, "#", instance.timeLayout)
	}
	return instance, nil
}

func (c *DatetimeChecker) Description() string {
	return "The field value must be a legal timestamp (second timestamp>1000000000 or millisecond timestamp>1000000000000). It supports second timestamp or millisecond timestamp. It is recommended to use second timestamp."
}

func (c *DatetimeChecker) CheckerDetails() string {
	return "datetime: valid timestamp"
}

func (c *DatetimeChecker) Check(fieldValueItr interface{}) (err error) {
	// don't check required.
	if c == nil || fieldValueItr == nil {
		return nil
	}
	defer func() {
		if err != nil {
			errMsg := fmt.Sprintf(c.errMessageFormat, fieldValueItr)
			err = errors.New(errMsg)
		}
	}()
	if c.isInt64 {
		int64Ts := cast.ToInt64(fieldValueItr)
		// 时间戳为0时不进行校验
		if int64Ts == 0 {
			return nil
		}
		err = checkTimestamp(int64Ts)
		return err
	}
	dataStr, ok := fieldValueItr.(string)
	if !ok {
		return nil
	}
	if len(dataStr) == 0 {
		return nil
	}
	_, err = time.Parse(c.timeLayout, dataStr)
	if err != nil {
		return err
	}
	return nil
}

const (
	SecTimestampThreshold int64 = 1000000000
)

const (
	microsecondsTimestampThreshold = 1000000000000000
)

func checkTimestamp(tsInt int64) error {
	// Microsecond timestamps are not supported
	if tsInt > microsecondsTimestampThreshold {
		return errors.New("invalid timestamp")
	}
	if tsInt < SecTimestampThreshold {
		return errors.New("invalid timestamp")
	}
	return nil
}
