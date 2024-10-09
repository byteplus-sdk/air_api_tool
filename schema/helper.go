package schema

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/3rd_rec/air_api_tool/helper"
	"github.com/3rd_rec/air_api_tool/utils"
	"github.com/shopspring/decimal"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/modern-go/reflect2"
)

func GetSchema(namespace, industry, table string) ([]*Field, error) {
	schemaPath := buildSchemaSavePath(namespace, industry, table)
	if !utils.FileExists(schemaPath) {
		return copyDefaultSchema(defaultIndustryTableSchema[industry][table]), nil
	}
	// load from local file
	currentSchemaBytes, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed while loading local schema, path:%s, err:%s", schemaPath, err))
	}
	schema := make([]*Field, 0)
	err = json.Unmarshal(currentSchemaBytes, &schema)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed while loading local schema, path:%s, err:%s", schemaPath, err))
	}
	return schema, nil
}

func buildSchemaSavePath(namespace, industry, table string) string {
	baseDir := helper.BuildModuleDir("schema")
	filename := fmt.Sprintf("%s_%s_%s_schema.json", namespace, industry, table)
	return filepath.Join(baseDir, filename)
}

func copyDefaultSchema(defaultSchema []*Field) []*Field {
	result := make([]*Field, 0, len(defaultSchema))
	for _, field := range defaultSchema {
		result = append(result, &Field{
			Name:         field.Name,
			Type:         field.Type,
			Description:  field.Description,
			ExampleValue: field.ExampleValue,
			Custom:       field.Custom,
		})
	}
	return result
}

func RemoveNilValue(data map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{}, len(data))
	for key, value := range data {
		if reflect2.IsNil(value) {
			continue
		}
		result[key] = value
	}
	return result
}

// FormatScientificString Convert a string represented in scientific notation to a normal string
func FormatScientificString(scientificStr string) string {
	if expIndex := strings.IndexAny(scientificStr, "Ee"); expIndex == -1 {
		return scientificStr
	}
	realStr, err := decimal.NewFromString(scientificStr)
	// 解析不了不做处理
	if err != nil {
		return scientificStr
	}
	return realStr.String()
}
