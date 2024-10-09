package shared

type DataExtension struct {
	LineNumber int

	LineString string

	OriginData map[string]interface{}

	StandardizedData map[string]interface{}

	// if ParseResult != nil, means parsing the data failed.
	ParseResult *ParseResult

	// Storage DataType Validate error(by schema) and Data Value Validate error(by checker)
	FieldValidateResult map[string]*FieldValidateItem
}

func NewDataExtension() *DataExtension {
	return &DataExtension{
		FieldValidateResult: make(map[string]*FieldValidateItem),
	}
}

type ParseResult struct {
	// Format Error、Invalid Data Length
	ErrType string

	ErrorDetails string
}

type FieldValidateItem struct {
	FieldName string
	*DataTypeErrorItem
	*DataValueErrorItem
}

type DataTypeErrorItem struct {
	FieldType     string
	FailureReason string
}

type DataValueErrorItem struct {
	CheckerName    string
	CheckerDetails string
	FailureReason  string
}
