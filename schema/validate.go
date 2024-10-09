package schema

import (
	"fmt"
	"github.com/3rd_rec/air_api_tool/consts"
	"github.com/3rd_rec/air_api_tool/shared"
)

func NewService(namespace, industry, table string) (*Service, error) {
	srv := newService(&Args{
		Namespace: namespace,
		Industry:  industry,
		Table:     table,
	})
	if err := srv.init(); err != nil {
		return nil, err
	}
	return srv, nil
}

type FieldValidateItem struct {
	FieldName     string
	FieldType     string
	FailureReason string
}

func (s *Service) Validate(data *shared.DataExtension) {
	errMsgFormat := "data type of %s(value:%v) should be %s"
	fieldValidateResult := data.FieldValidateResult
	dataMap := RemoveNilValue(data.OriginData)
	for _, field := range s.schema {
		if _, exist := fieldValidateResult[field.Name]; exist {
			continue
		}
		fieldValue, ok := dataMap[field.Name]
		if !ok {
			continue
		}
		errMsg := fmt.Sprintf(errMsgFormat, field.Name, fieldValue, field.Type)
		switch field.Type {
		case consts.FieldTypeFloat, consts.FieldTypeDouble:
			floatNum, ok := parseAsFloat64(fieldValue)
			if !ok {
				s.addFieldValidateResult(fieldValidateResult, field, errMsg)
				continue
			}
			dataMap[field.Name] = floatNum
		case consts.FieldTypeInt32:
			intNum, ok := parseAsInt64(fieldValue)
			if !ok {
				s.addFieldValidateResult(fieldValidateResult, field, errMsg)
				continue
			}
			dataMap[field.Name] = int32(intNum)
		case consts.FieldTypeInt64:
			intNum, ok := parseAsInt64(fieldValue)
			if !ok {
				s.addFieldValidateResult(fieldValidateResult, field, errMsg)
				continue
			}
			dataMap[field.Name] = intNum
		case consts.FieldTypeBool:
			boolVal, ok := parseAsInt64(fieldValue)
			if !ok {
				s.addFieldValidateResult(fieldValidateResult, field, errMsg)
				continue
			}
			dataMap[field.Name] = boolVal
		case consts.FieldTypeString:
			strVal, ok := parseAsString(fieldValue)
			if !ok {
				s.addFieldValidateResult(fieldValidateResult, field, errMsg)
				continue
			}
			dataMap[field.Name] = strVal
		case consts.FieldTypeJSONString:
			errMsg = fmt.Sprintf("field %s(value:%v) should be a valid json_string type. For example: {\"field\": \"[\\\"v1\\\", \\\"v2\\\"]\"}, the values of the field need to be serialized as a string. Please note not to concatenate manually. Programming languages typically have libraries for handling JSON, such as Python's JSON library. You can use json.dumps to achieve serialization.", field.Name, fieldValue)
			jsonStrVal, ok := parseJSONString(fieldValue)
			if !ok {
				s.addFieldValidateResult(fieldValidateResult, field, errMsg)
				continue
			}
			dataMap[field.Name] = jsonStrVal
		}
	}
	data.StandardizedData = dataMap
}

func (s *Service) addFieldValidateResult(fieldValidateResult map[string]*shared.FieldValidateItem, field *Field, errMsg string) {
	fieldValidateResult[field.Name] = &shared.FieldValidateItem{
		FieldName: field.Name,
		DataTypeErrorItem: &shared.DataTypeErrorItem{
			FieldType:     field.Type,
			FailureReason: errMsg,
		},
	}
}
