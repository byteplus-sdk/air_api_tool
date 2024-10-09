package checker

import (
	"github.com/3rd_rec/air_api_tool/consts"
	"github.com/3rd_rec/air_api_tool/shared"
	"github.com/3rd_rec/air_api_tool/utils"
)

type fieldCheckerWrapper struct {
	item         *Item
	fieldChecker FieldChecker
}

type Item struct {
	FieldName   string                   `json:"field_name,omitempty"`
	CheckerName string                   `json:"checker_name,omitempty"`
	CheckerCfg  *shared.CommonCheckerCfg `json:"checker_cfg,omitempty"`
	DependCfg   *shared.DependCfg        `json:"depend_cfg,omitempty"`
	Ignored     bool                     `json:"ignored,omitempty"`
}

var (
	// industry -> table -> checker_cfg
	defaultIndustryTableChecker = map[string]map[string][]*Item{
		consts.IndustrySaasRetail: {
			consts.TableNameUser:      defaultSaasRetailUserChecker,
			consts.TableNameProduct:   defaultSaasRetailProductChecker,
			consts.TableNameUserEvent: defaultSaasRetailUserEventChecker,
		},
		consts.IndustrySaasContent: {
			consts.TableNameUser:      defaultSaasContentUserChecker,
			consts.TableNameContent:   defaultSaasContentContentChecker,
			consts.TableNameUserEvent: defaultSaasContentUserEventChecker,
		},
	}
)

// saas retail default checker configuration
var (
	defaultSaasRetailUserChecker = []*Item{
		{
			FieldName:   "user_id",
			CheckerName: "required",
			CheckerCfg: &shared.CommonCheckerCfg{
				Required: true,
			},
		},
		{
			FieldName:   "registration_timestamp",
			CheckerName: "datetime",
			CheckerCfg: &shared.CommonCheckerCfg{
				DatetimeCfg: &shared.DatetimeCfg{
					IsInt64: true,
				},
			},
		},
	}
	defaultSaasRetailProductChecker = []*Item{
		{
			FieldName:   "product_id",
			CheckerName: consts.CheckerNameRequired,
			CheckerCfg: &shared.CommonCheckerCfg{
				Required: true,
			},
		},
		{
			FieldName:   "is_recommendable",
			CheckerName: consts.CheckerNameRequired,
			CheckerCfg: &shared.CommonCheckerCfg{
				Required: true,
			},
		},
		{
			FieldName:   "is_recommendable",
			CheckerName: consts.CheckerNameEnums,
			CheckerCfg: &shared.CommonCheckerCfg{
				Enums: []string{"0", "1"},
			},
		},
		{
			FieldName:   "categories",
			CheckerName: consts.CheckerNameCategories,
			CheckerCfg: &shared.CommonCheckerCfg{
				CategoryCfg: &shared.CategoryCfg{
					Required:   true,
					StrictMode: true,
				},
			},
		},
		{
			FieldName:   "current_price",
			CheckerName: consts.CheckerNameRanges,
			CheckerCfg: &shared.CommonCheckerCfg{
				RangeCfgList: []*shared.RangeCfg{
					{
						Gt: utils.Float64Ptr(0),
					},
				},
			},
		},
		{
			FieldName:   "original_price",
			CheckerName: consts.CheckerNameRanges,
			CheckerCfg: &shared.CommonCheckerCfg{
				RangeCfgList: []*shared.RangeCfg{
					{
						Gte: utils.Float64Ptr(0),
					},
				},
			},
		},
		{
			FieldName:   "publish_timestamp",
			CheckerName: consts.CheckerNameDatetime,
			CheckerCfg: &shared.CommonCheckerCfg{
				DatetimeCfg: &shared.DatetimeCfg{
					IsInt64: true,
				},
			},
		},
		{
			FieldName:   "user_rating",
			CheckerName: consts.CheckerNameRanges,
			CheckerCfg: &shared.CommonCheckerCfg{
				RangeCfgList: []*shared.RangeCfg{
					{
						Gte: utils.Float64Ptr(0),
					},
				},
			},
		},
		{
			FieldName:   "seller_rating",
			CheckerName: consts.CheckerNameRanges,
			CheckerCfg: &shared.CommonCheckerCfg{
				RangeCfgList: []*shared.RangeCfg{
					{
						Gte: utils.Float64Ptr(0),
					},
				},
			},
		},
		{
			FieldName:   "sold_count",
			CheckerName: consts.CheckerNameRanges,
			CheckerCfg: &shared.CommonCheckerCfg{
				RangeCfgList: []*shared.RangeCfg{
					{
						Gte: utils.Float64Ptr(0),
					},
				},
			},
		},
		{
			FieldName:   "comment_count",
			CheckerName: consts.CheckerNameRanges,
			CheckerCfg: &shared.CommonCheckerCfg{
				RangeCfgList: []*shared.RangeCfg{
					{
						Gte: utils.Float64Ptr(0),
					},
				},
			},
		},
	}
	defaultSaasRetailUserEventChecker = []*Item{
		{
			FieldName:   "user_id",
			CheckerName: consts.CheckerNameRequired,
			CheckerCfg: &shared.CommonCheckerCfg{
				Required: true,
			},
		},
		{
			FieldName:   "event_timestamp",
			CheckerName: consts.CheckerNameRequired,
			CheckerCfg: &shared.CommonCheckerCfg{
				Required: true,
			},
		},
		{
			FieldName:   "event_timestamp",
			CheckerName: consts.CheckerNameDatetime,
			CheckerCfg: &shared.CommonCheckerCfg{
				DatetimeCfg: &shared.DatetimeCfg{
					IsInt64: true,
				},
			},
		},
		{
			FieldName:   "product_id",
			CheckerName: consts.CheckerNameRequired,
			CheckerCfg: &shared.CommonCheckerCfg{
				Required: true,
			},
			DependCfg: &shared.DependCfg{
				DependField: "event_type",
				DependFieldValues: []string{
					"impression",
					"click",
					"add-to-cart",
					"remove-from-cart",
					"add-to-favorites",
					"remove-from-favorites",
					"purchase",
					"checkout",
				},
			},
		},
		{
			FieldName:   "event_type",
			CheckerName: consts.CheckerNameRequired,
			CheckerCfg: &shared.CommonCheckerCfg{
				Required: true,
			},
		},
		{
			FieldName:   "scene_name",
			CheckerName: consts.CheckerNameRequired,
			CheckerCfg: &shared.CommonCheckerCfg{
				Required: true,
			},
			DependCfg: &shared.DependCfg{
				DependField: "event_type",
				DependFieldValues: []string{
					"impression",
					"click",
				},
			},
		},
		{
			FieldName:   "page_number",
			CheckerName: consts.CheckerNameRanges,
			CheckerCfg: &shared.CommonCheckerCfg{
				RangeCfgList: []*shared.RangeCfg{
					{
						Gte: utils.Float64Ptr(0),
					},
				},
			},
		},
		{
			FieldName:   "offset",
			CheckerName: consts.CheckerNameRanges,
			CheckerCfg: &shared.CommonCheckerCfg{
				RangeCfgList: []*shared.RangeCfg{
					{
						Gte: utils.Float64Ptr(0),
					},
				},
			},
		},
		{
			FieldName:   "purchase_count",
			CheckerName: consts.CheckerNameRequired,
			CheckerCfg: &shared.CommonCheckerCfg{
				Required: true,
			},
			DependCfg: &shared.DependCfg{
				DependField: "event_type",
				DependFieldValues: []string{
					"purchase",
				},
			},
		},
		{
			FieldName:   "purchase_count",
			CheckerName: consts.CheckerNameRanges,
			CheckerCfg: &shared.CommonCheckerCfg{
				RangeCfgList: []*shared.RangeCfg{
					{
						Gt: utils.Float64Ptr(0),
					},
				},
			},
			DependCfg: &shared.DependCfg{
				DependField: "event_type",
				DependFieldValues: []string{
					"purchase",
				},
			},
		},
		{
			FieldName:   "paid_price",
			CheckerName: consts.CheckerNameRequired,
			CheckerCfg: &shared.CommonCheckerCfg{
				Required: true,
			},
			DependCfg: &shared.DependCfg{
				DependField: "event_type",
				DependFieldValues: []string{
					"purchase",
				},
			},
		},
		{
			FieldName:   "paid_price",
			CheckerName: consts.CheckerNameRanges,
			CheckerCfg: &shared.CommonCheckerCfg{
				RangeCfgList: []*shared.RangeCfg{
					{
						Gt: utils.Float64Ptr(0),
					},
				},
			},
			DependCfg: &shared.DependCfg{
				DependField: "event_type",
				DependFieldValues: []string{
					"purchase",
				},
			},
		},
		{
			FieldName:   "traffic_source",
			CheckerName: consts.CheckerNameEnums,
			CheckerCfg: &shared.CommonCheckerCfg{
				Enums: []string{"self", "byteplus", "other"},
			},
		},
	}
)

// saas content default checker configuration
var (
	defaultSaasContentUserChecker = []*Item{
		{
			FieldName:   "user_id",
			CheckerName: "required",
			CheckerCfg: &shared.CommonCheckerCfg{
				Required: true,
			},
		},
		{
			FieldName:   "registration_timestamp",
			CheckerName: "datetime",
			CheckerCfg: &shared.CommonCheckerCfg{
				DatetimeCfg: &shared.DatetimeCfg{
					IsInt64: true,
				},
			},
		},
	}
	defaultSaasContentContentChecker = []*Item{
		{
			FieldName:   "content_id",
			CheckerName: consts.CheckerNameRequired,
			CheckerCfg: &shared.CommonCheckerCfg{
				Required: true,
			},
		},
		{
			FieldName:   "is_recommendable",
			CheckerName: consts.CheckerNameRequired,
			CheckerCfg: &shared.CommonCheckerCfg{
				Required: true,
			},
		},
		{
			FieldName:   "is_recommendable",
			CheckerName: consts.CheckerNameEnums,
			CheckerCfg: &shared.CommonCheckerCfg{
				Enums: []string{"0", "1"},
			},
		},
		{
			FieldName:   "content_type",
			CheckerName: consts.CheckerNameRequired,
			CheckerCfg: &shared.CommonCheckerCfg{
				Required: true,
			},
		},
		{
			FieldName:   "categories",
			CheckerName: consts.CheckerNameCategories,
			CheckerCfg: &shared.CommonCheckerCfg{
				CategoryCfg: &shared.CategoryCfg{
					Required:   true,
					StrictMode: true,
				},
			},
		},
		{
			FieldName:   "publish_timestamp",
			CheckerName: consts.CheckerNameDatetime,
			CheckerCfg: &shared.CommonCheckerCfg{
				DatetimeCfg: &shared.DatetimeCfg{
					IsInt64: true,
				},
			},
		},
		{
			FieldName:   "current_price",
			CheckerName: consts.CheckerNameRanges,
			CheckerCfg: &shared.CommonCheckerCfg{
				RangeCfgList: []*shared.RangeCfg{
					{
						Gt: utils.Float64Ptr(0),
					},
				},
			},
		},
		{
			FieldName:   "original_price",
			CheckerName: consts.CheckerNameRanges,
			CheckerCfg: &shared.CommonCheckerCfg{
				RangeCfgList: []*shared.RangeCfg{
					{
						Gte: utils.Float64Ptr(0),
					},
				},
			},
		},
		{
			FieldName:   "video_duration",
			CheckerName: consts.CheckerNameRanges,
			CheckerCfg: &shared.CommonCheckerCfg{
				RangeCfgList: []*shared.RangeCfg{
					{
						Gte: utils.Float64Ptr(0),
					},
				},
			},
		},
		{
			FieldName:   "user_rating",
			CheckerName: consts.CheckerNameRanges,
			CheckerCfg: &shared.CommonCheckerCfg{
				RangeCfgList: []*shared.RangeCfg{
					{
						Gte: utils.Float64Ptr(0),
					},
				},
			},
		},
	}
	defaultSaasContentUserEventChecker = []*Item{
		{
			FieldName:   "user_id",
			CheckerName: consts.CheckerNameRequired,
			CheckerCfg: &shared.CommonCheckerCfg{
				Required: true,
			},
		},
		{
			FieldName:   "content_id",
			CheckerName: consts.CheckerNameRequired,
			CheckerCfg: &shared.CommonCheckerCfg{
				Required: true,
			},
			DependCfg: &shared.DependCfg{
				DependField: "event_type",
				DependFieldValues: []string{
					"impression",
					"click",
					"like",
					"share",
					"comment",
					"cart",
					"favorite",
					"stay",
					"checkout",
				},
			},
		},
		{
			FieldName:   "event_timestamp",
			CheckerName: consts.CheckerNameRequired,
			CheckerCfg: &shared.CommonCheckerCfg{
				Required: true,
			},
		},
		{
			FieldName:   "event_timestamp",
			CheckerName: consts.CheckerNameDatetime,
			CheckerCfg: &shared.CommonCheckerCfg{
				DatetimeCfg: &shared.DatetimeCfg{
					IsInt64: true,
				},
			},
		},
		{
			FieldName:   "event_type",
			CheckerName: consts.CheckerNameRequired,
			CheckerCfg: &shared.CommonCheckerCfg{
				Required: true,
			},
		},
		{
			FieldName:   "scene_name",
			CheckerName: consts.CheckerNameRequired,
			CheckerCfg: &shared.CommonCheckerCfg{
				Required: true,
			},
			DependCfg: &shared.DependCfg{
				DependField: "event_type",
				DependFieldValues: []string{
					"impression",
					"click",
					"stay",
				},
			},
		},
		{
			FieldName:   "traffic_source",
			CheckerName: consts.CheckerNameEnums,
			CheckerCfg: &shared.CommonCheckerCfg{
				Enums: []string{"self", "byteplus", "other"},
			},
		},
		{
			FieldName:   "page_number",
			CheckerName: consts.CheckerNameRanges,
			CheckerCfg: &shared.CommonCheckerCfg{
				RangeCfgList: []*shared.RangeCfg{
					{
						Gte: utils.Float64Ptr(0),
					},
				},
			},
		},
		{
			FieldName:   "offset",
			CheckerName: consts.CheckerNameRanges,
			CheckerCfg: &shared.CommonCheckerCfg{
				RangeCfgList: []*shared.RangeCfg{
					{
						Gte: utils.Float64Ptr(0),
					},
				},
			},
		},
		{
			FieldName:   "stay_duration",
			CheckerName: consts.CheckerNameRequired,
			CheckerCfg: &shared.CommonCheckerCfg{
				Required: true,
			},
			DependCfg: &shared.DependCfg{
				DependField: "event_type",
				DependFieldValues: []string{
					"stay",
				},
			},
		},
		{
			FieldName:   "stay_duration",
			CheckerName: consts.CheckerNameRanges,
			CheckerCfg: &shared.CommonCheckerCfg{
				RangeCfgList: []*shared.RangeCfg{
					{
						Gt: utils.Float64Ptr(0),
					},
				},
			},
			DependCfg: &shared.DependCfg{
				DependField: "event_type",
				DependFieldValues: []string{
					"stay",
				},
			},
		},
		{
			FieldName:   "purchase_count",
			CheckerName: consts.CheckerNameRequired,
			CheckerCfg: &shared.CommonCheckerCfg{
				Required: true,
			},
			DependCfg: &shared.DependCfg{
				DependField: "event_type",
				DependFieldValues: []string{
					"purchase",
				},
			},
		},
		{
			FieldName:   "purchase_count",
			CheckerName: consts.CheckerNameRanges,
			CheckerCfg: &shared.CommonCheckerCfg{
				RangeCfgList: []*shared.RangeCfg{
					{
						Gt: utils.Float64Ptr(0),
					},
				},
			},
			DependCfg: &shared.DependCfg{
				DependField: "event_type",
				DependFieldValues: []string{
					"purchase",
				},
			},
		},
		{
			FieldName:   "paid_price",
			CheckerName: consts.CheckerNameRequired,
			CheckerCfg: &shared.CommonCheckerCfg{
				Required: true,
			},
			DependCfg: &shared.DependCfg{
				DependField: "event_type",
				DependFieldValues: []string{
					"purchase",
				},
			},
		},
		{
			FieldName:   "paid_price",
			CheckerName: consts.CheckerNameRanges,
			CheckerCfg: &shared.CommonCheckerCfg{
				RangeCfgList: []*shared.RangeCfg{
					{
						Gt: utils.Float64Ptr(0),
					},
				},
			},
			DependCfg: &shared.DependCfg{
				DependField: "event_type",
				DependFieldValues: []string{
					"purchase",
				},
			},
		},
	}
)
