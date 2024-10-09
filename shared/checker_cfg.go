package shared

import (
	"fmt"
	"github.com/3rd_rec/air_api_tool/utils"
	"strings"
)

type DependCfg struct {
	DependField       string   `json:"depend_field,omitempty"`
	DependFieldValues []string `json:"depend_field_values,omitempty"`
}

func (d *DependCfg) String() string {
	return fmt.Sprintf("when `%s` field value is in [%s]", d.DependField, strings.Join(d.DependFieldValues, ", "))
}

type CommonCheckerCfg struct {
	Required     bool         `json:"required,omitempty"`
	Enums        []string     `json:"enums,omitempty"`
	DatetimeCfg  *DatetimeCfg `json:"datetime_cfg,omitempty"`
	RangeCfgList []*RangeCfg  `json:"range_cfg_list,omitempty"`
	CategoryCfg  *CategoryCfg `json:"category_cfg,omitempty"`
}

type DatetimeCfg struct {
	IsInt64    bool   `json:"is_int64,omitempty"`
	TimeLayout string `json:"time_layout,omitempty"`
}

type RangeCfg struct {
	Gt  *float64 `json:"gt,omitempty"`
	Gte *float64 `json:"gte,omitempty"`
	Lt  *float64 `json:"lt,omitempty"`
	Lte *float64 `json:"lte,omitempty"`
}

// CategoryCfg TODO: move it to shared package
// CategoryCfg Field Explanation
// Required: Whether to validate the category. If false, return directly. If true, the structure/content of the category field will be checked.
// StrictMode: Whether to enable strict validation mode. Once enabled, it will require the levels in the categories field provided by the client to start from 1 and increment consecutively.
type CategoryCfg struct {
	Required   bool `json:"required"`
	StrictMode bool `json:"strict_mode"`
}

func (cfg *RangeCfg) String() string {
	var msgBuff = &strings.Builder{}
	if cfg.Gte != nil {
		msgBuff.WriteString(utils.Float64ToA(*cfg.Gte))
		msgBuff.WriteString("=<")
	} else if cfg.Gt != nil {
		msgBuff.WriteString(utils.Float64ToA(*cfg.Gt))
		msgBuff.WriteString("<")
	}
	msgBuff.WriteString("X")
	if cfg.Lte != nil {
		msgBuff.WriteString("<=")
		msgBuff.WriteString(utils.Float64ToA(*cfg.Lte))
	} else if cfg.Lt != nil {
		msgBuff.WriteString("<")
		msgBuff.WriteString(utils.Float64ToA(*cfg.Lt))
	}
	return msgBuff.String()
}
