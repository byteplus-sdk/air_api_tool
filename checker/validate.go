package checker

import (
	"errors"
	"fmt"
	"github.com/3rd_rec/air_api_tool/shared"
)

type FieldChecker interface {
	Description() string

	CheckerDetails() string

	// Check 对字段值进行校验
	Check(fieldValueItr interface{}) error
}

type FieldCheckerFactory interface {
	// Create 解析传入的JSON字符串，构造得到Checker实例
	Create(checkerCfg *shared.CommonCheckerCfg) (FieldChecker, error)
}

var baseCheckerFactoryMap = make(map[string]FieldCheckerFactory)

func RegisterBaseCheckerFactory(name string, factory FieldCheckerFactory) {
	baseCheckerFactoryMap[name] = factory
}

func GetFieldCheckerFactory(name string) FieldCheckerFactory {
	return baseCheckerFactoryMap[name]
}

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
	CheckerName   string
	CheckerDetail string
	FailureReason string
}

// Validate Returns the verification result and whether the verification is successful
func (s *Service) Validate(data *shared.DataExtension) {
	// field_name -> item
	fieldValidateResult := data.FieldValidateResult
	dataMap := data.StandardizedData
	for _, wrapper := range s.fieldCheckerWrappers {
		if _, exist := fieldValidateResult[wrapper.item.FieldName]; exist {
			continue
		}
		if wrapper.item.Ignored {
			continue
		}
		if !s.matchDependCfg(dataMap, wrapper.item.DependCfg) {
			continue
		}
		err := wrapper.fieldChecker.Check(dataMap[wrapper.item.FieldName])
		if err != nil {
			err = s.reFormatWithDependCfg(err, wrapper.item.DependCfg)
			s.addFieldValidateResult(fieldValidateResult, wrapper, err)
		}
	}
}

func (s *Service) matchDependCfg(data map[string]interface{}, dependCfg *shared.DependCfg) bool {
	if dependCfg == nil || len(dependCfg.DependField) == 0 {
		return true
	}
	dependField := dependCfg.DependField
	dependFieldValue, exist := data[dependField]
	if !exist {
		return false
	}
	dependFieldValueStr, ok := dependFieldValue.(string)
	if !ok {
		return false
	}
	for _, value := range dependCfg.DependFieldValues {
		if dependFieldValueStr == value {
			return true
		}
	}

	return false
}

func (s *Service) reFormatWithDependCfg(originErr error, dependCfg *shared.DependCfg) error {
	if dependCfg == nil {
		return originErr
	}
	newErrMsg := fmt.Sprintf("%s, %s", originErr, dependCfg)
	return errors.New(newErrMsg)
}

func (s *Service) addFieldValidateResult(fieldValidateResult map[string]*shared.FieldValidateItem,
	wrapper *fieldCheckerWrapper, err error) {
	fieldValidateResult[wrapper.item.FieldName] = &shared.FieldValidateItem{
		FieldName: wrapper.item.FieldName,
		DataValueErrorItem: &shared.DataValueErrorItem{
			CheckerName:    wrapper.item.CheckerName,
			CheckerDetails: wrapper.fieldChecker.CheckerDetails(),
			FailureReason:  err.Error(),
		},
	}
}
