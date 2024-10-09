package fieldchecker

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/3rd_rec/air_api_tool/checker"
	"github.com/3rd_rec/air_api_tool/consts"
	"github.com/3rd_rec/air_api_tool/shared"
)

func init() {
	checker.RegisterBaseCheckerFactory(consts.CheckerNameCategories, &CategoriesChecker{})
}

type CategoriesChecker struct {
	isRequired   bool
	isStrictMode bool
}

type Category struct {
	CategoryDepth int32           `json:"category_depth"`
	CategoryNodes []*CategoryNode `json:"category_nodes"`
}

type CategoryNode struct {
	IdOrName string `json:"id_or_name"`
}

func (c *CategoriesChecker) Create(checkerCfg *shared.CommonCheckerCfg) (checker.FieldChecker, error) {
	if checkerCfg == nil {
		return nil, errors.New("invalid categories checker cfg")
	}
	cateCfg := checkerCfg.CategoryCfg
	if cateCfg == nil {
		return nil, errors.New("invalid categories checker cfg")
	}
	instance := &CategoriesChecker{
		isRequired:   cateCfg.Required,
		isStrictMode: cateCfg.StrictMode,
	}
	return instance, nil
}

func (c *CategoriesChecker) Description() string {
	var correctFormat = "correct format example: \"[{\\\"category_depth\\\":1,\\\"category_nodes\\\":[{\\\"id_or_name\\\":\\\"Movie\\\"}]},\n{\\\"category_depth\\\":2,\\\"category_nodes\\\":[{\\\"id_or_name\\\":\\\"Comedy\\\"}]}]\""
	return fmt.Sprintf("json array string of category element. categories#category_nodes[n]#id_or_name can't be empty, categories#category_depth should start with 1 with continuous integers. %s", correctFormat)
}

func (c *CategoriesChecker) CheckerDetails() string {
	return "categories: correct categories format"
}

func (c *CategoriesChecker) Check(fieldValue interface{}) error {
	if c == nil || !c.isRequired {
		return nil
	}
	if fieldValue == nil {
		return nil
	}
	categoriesStr, ok := fieldValue.(string)
	if !ok {
		return nil
	}
	if categoriesStr == "" {
		return nil
	}
	return c.doCheck(categoriesStr)
}

func (c *CategoriesChecker) doCheck(categoriesStr string) error {
	var categories []*Category
	err := json.Unmarshal([]byte(categoriesStr), &categories)
	if err != nil {
		return errors.New("the value of `categories` is format as json err")
	}
	if len(categories) == 0 {
		return nil
	}
	for idx, cate := range categories {
		if cate.CategoryDepth <= 0 {
			return errors.New("the value of `categories#category_depth` is invalid")
		}
		if len(cate.CategoryNodes) == 0 {
			return errors.New("the value of `categories#category_nodes` is empty")
		}
		for _, node := range cate.CategoryNodes {
			if len(node.IdOrName) == 0 {
				return errors.New("the value of `categories#category_nodes#id_or_name` is empty")
			}
		}
		if c.isStrictMode {
			// Category Depth should start with 1 with continuous integers
			if int(cate.CategoryDepth) != idx+1 {
				return errors.New("the value of `categories#category_depth` is invalid," +
					" category_depth should start with 1 with continuous integers")
			}
		}
	}
	return nil
}
