package helper

import (
	"errors"
	"fmt"
	"github.com/3rd_rec/air_api_tool/consts"
	"strings"
)

func CheckCommonArgs(industry string, table string) error {
	if len(industry) == 0 {
		return errors.New("--industry args is required")
	}
	tableList, exist := consts.SupportedIndustryTableMap[industry]
	if !exist {
		errMsg := fmt.Sprintf("invalid --industry args(%s). The accepted values are: [saas_retail, saas_content]", industry)
		return errors.New(errMsg)
	}
	if len(table) == 0 {
		return errors.New("--table args is required")
	}
	if _, exist = slice2Map(tableList)[table]; !exist {
		errMsg := fmt.Sprintf("invalid --table args(%s). The accepted values are: [%s]", table, strings.Join(tableList, ", "))
		return errors.New(errMsg)
	}
	return nil
}

func slice2Map(slice []string) map[string]bool {
	result := make(map[string]bool, len(slice))
	for _, value := range slice {
		result[value] = true
	}
	return result
}
