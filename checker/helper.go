package checker

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/3rd_rec/air_api_tool/helper"
	"github.com/3rd_rec/air_api_tool/utils"
	"io/ioutil"
	"path/filepath"
)

func GetCheckerCfg(namespace, industry, table string) ([]*Item, error) {
	checkerPath := buildCheckerSavePath(namespace, industry, table)
	if !utils.FileExists(checkerPath) {
		return copyDefaultChecker(defaultIndustryTableChecker[industry][table]), nil
	}
	// load from local file
	currentCheckerBytes, err := ioutil.ReadFile(checkerPath)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed while loading local checker configuration, path:%s, err:%s",
			checkerPath, err))
	}
	checkerItems := make([]*Item, 0)
	err = json.Unmarshal(currentCheckerBytes, &checkerItems)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed while loading local checker configuration, path:%s, err:%s",
			checkerPath, err))
	}
	return checkerItems, nil
}

func buildCheckerSavePath(namespace, industry, table string) string {
	baseDir := helper.BuildModuleDir("checker")
	filename := fmt.Sprintf("%s_%s_%s_checker.json", namespace, industry, table)
	return filepath.Join(baseDir, filename)
}

func copyDefaultChecker(defaultChecker []*Item) []*Item {
	result := make([]*Item, 0, len(defaultChecker))
	for _, checkerItem := range defaultChecker {
		result = append(result, &Item{
			FieldName:   checkerItem.FieldName,
			CheckerName: checkerItem.CheckerName,
			CheckerCfg:  checkerItem.CheckerCfg, // Shallow Copy.
			DependCfg:   checkerItem.DependCfg,  // Shallow Copy.
			Ignored:     checkerItem.Ignored,
		})
	}
	return result
}
