package consts

import "github.com/fatih/color"

const (
	IndustrySaasRetail  = "saas_retail"
	IndustrySaasContent = "saas_content"
)

const (
	TableNameUser      = "user"
	TableNameProduct   = "product"
	TableNameContent   = "content"
	TableNameUserEvent = "user_event"
)

var (
	// SupportedIndustryTableMap industry -> table list
	SupportedIndustryTableMap = map[string][]string{
		IndustrySaasRetail: {
			TableNameUser,
			TableNameProduct,
			TableNameUserEvent,
		},
		IndustrySaasContent: {
			TableNameUser,
			TableNameContent,
			TableNameUserEvent,
		},
	}
)

const (
	ContentTypeCSV  = "csv"
	ContentTypeJSON = "json"
)

const (
	CheckerNameDatetime   = "datetime"
	CheckerNameEnums      = "enums"
	CheckerNameRequired   = "required"
	CheckerNameRanges     = "ranges"
	CheckerNameCategories = "categories"
)

const (
	FieldTypeInt32      = "int32"
	FieldTypeInt64      = "int64"
	FieldTypeFloat      = "float"
	FieldTypeDouble     = "double"
	FieldTypeString     = "string"
	FieldTypeBool       = "bool"
	FieldTypeJSONString = "json_string"
)

const (
	FieldErrorTypeDataTypeValidateError  = "Data Type Validate Error"
	FieldErrorTypeDataValueValidateError = "Data Value Validate Error"
)

const (
	ParseErrorTypeFormatError       = "Format Error"
	ParseErrorTypeInvalidDataLength = "Invalid Data Length"
)

const (
	KB = 1024
	MB = 1024 * KB
	GB = 1024 * MB
)

var (
	VerificationStatusSuccess = color.GreenString("SUCCESS")
	VerificationStatusFailure = color.RedString("FAILURE")
)

const EnvNameDataDir = "AIR_API_TOOL_DATA_DIR"
const DefaultDataDirPrefix = "/share/air_api_tool"
