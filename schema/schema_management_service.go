package schema

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/3rd_rec/air_api_tool/helper"
	"github.com/3rd_rec/air_api_tool/utils"
	"io/ioutil"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

type Service struct {
	args           *Args
	schemaSavePath string
	schema         []*Field
	schemaMap      map[string]*Field
	updated        bool
}

func newService(args *Args) *Service {
	return &Service{
		args: args,
	}
}

func (s *Service) Run() error {
	if err := s.checkArgs(); err != nil {
		return err
	}
	if s.args.Clear {
		if err := s.clear(); err != nil {
			return err
		}
		return nil
	}
	if err := s.init(); err != nil {
		return err
	}
	return s.executeCmd()
}

func (s *Service) checkArgs() error {
	if err := helper.CheckCommonArgs(s.args.Industry, s.args.Table); err != nil {
		return err
	}
	if len(s.args.AddFields) > 0 {
		for _, fieldPair := range s.args.AddFields {
			result := strings.Split(fieldPair, ",")
			if len(result) != 2 {
				return errors.New("invalid --add-field args. Use `api_tool schema --help` to view the usage")
			}
			fieldType := result[1]
			if _, exist := supportedFieldType[strings.ToLower(fieldType)]; !exist {
				msg := fmt.Sprintf("unsupported field type:%s exists in --add-field args(%s)."+
					" Use `api_tool schema --help` to view the usage", fieldType, fieldPair)
				return errors.New(msg)
			}
		}
	}
	if len(s.args.AddFields) == 0 && len(s.args.DeleteFields) == 0 && !s.args.Show && !s.args.Clear {
		return errors.New("no valid args. Use `api_tool schema --help` to view the usage")
	}
	return nil
}

func (s *Service) init() error {
	var err error

	baseDir := helper.BuildModuleDir("schema")
	err = utils.EnsureDirExists(baseDir)
	if err != nil {
		return errors.New(fmt.Sprintf("schema dir %s create fail, err:%s", baseDir, err))
	}
	s.schemaSavePath = buildSchemaSavePath(s.args.Namespace, s.args.Industry, s.args.Table)
	s.schema, err = GetSchema(s.args.Namespace, s.args.Industry, s.args.Table)
	if err != nil {
		return err
	}
	s.schemaMap = make(map[string]*Field, len(s.schema))
	for _, field := range s.schema {
		s.schemaMap[field.Name] = field
	}
	return nil
}

func (s *Service) executeCmd() error {
	defer func() {
		s.show()
	}()

	if err := s.addFields(); err != nil {
		return err
	}
	if err := s.deleteFields(); err != nil {
		return err
	}
	if err := s.saveSchema(); err != nil {
		return err
	}
	return nil
}

func (s *Service) clear() error {
	if !s.args.Clear {
		return nil
	}
	s.schemaSavePath = buildSchemaSavePath(s.args.Namespace, s.args.Industry, s.args.Table)
	if !utils.FileExists(s.schemaSavePath) {
		fmt.Println("cleanup schema successful. Recover to default configuration.")
		return nil
	}
	err := os.Remove(s.schemaSavePath)
	if err != nil {
		return errors.New(fmt.Sprintf("failed while clear local schema:%s, err:%s", s.schemaSavePath, err))
	}
	fmt.Println("cleanup schema successful. Recover to default configuration.")
	return nil
}

func (s *Service) addFields() error {
	if len(s.args.AddFields) == 0 {
		return nil
	}
	needAddFields := s.buildFields(s.args.AddFields)
	for _, needAddField := range needAddFields {
		_, exist := s.schemaMap[needAddField.Name]
		if exist {
			warningMsg := fmt.Sprintf("[Warning] field(%s) already exists. "+
				"Use 'api_tool schema --industry [industry] --table [table] --show' "+
				"to view the existing schema field", needAddField.Name)
			fmt.Println(warningMsg)
			continue
		}
		s.schema = append(s.schema, needAddField)
	}
	s.updated = true
	return nil
}

func (s *Service) buildFields(addFields []string) []*Field {
	needAddFields := make([]*Field, 0, len(addFields))
	for _, fieldPair := range addFields {
		result := strings.Split(fieldPair, ",")
		fieldName := result[0]
		fieldType := strings.ToLower(result[1])
		needAddFields = append(needAddFields, &Field{
			Name:        fieldName,
			Type:        fieldType,
			Description: descriptionCustomField,
			Custom:      true,
		})
	}
	return needAddFields
}

// Need to ensure order after deletion.
func (s *Service) deleteFields() error {
	if len(s.args.DeleteFields) == 0 {
		return nil
	}

	for _, needDeletedFieldName := range s.args.DeleteFields {
		existField, exist := s.schemaMap[needDeletedFieldName]
		if !exist {
			warningMsg := fmt.Sprintf("[Warning] delete field name:%s does not exist in schema. "+
				"Use 'api_tool schema --industry [industry] --table [table] --show' "+
				"to view the existing schema fields", needDeletedFieldName)
			fmt.Println(warningMsg)
			continue
		}
		if !existField.Custom {
			return errors.New(fmt.Sprintf("preset fields: %s cannot be deleted, only customized fields can be deleted", existField.Name))
		}
		delete(s.schemaMap, needDeletedFieldName)
	}

	schemaAfterDeleted := make([]*Field, 0, len(s.schema))
	for _, field := range s.schema {
		if _, exist := s.schemaMap[field.Name]; exist {
			schemaAfterDeleted = append(schemaAfterDeleted, field)
		}
	}

	s.schema = schemaAfterDeleted
	s.updated = true
	return nil
}

func (s *Service) show() {
	if !s.args.Show {
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Field Name", "Data Type", "Description", "Example Value"})
	table.SetBorder(true)
	table.SetRowLine(true)
	table.SetAutoWrapText(false)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.AppendBulk(s.toTableRows(s.schema))
	table.Render()
}

func (s *Service) toTableRows(schema []*Field) [][]string {
	rows := make([][]string, 0, len(schema))
	for _, field := range schema {
		description := field.Description
		if field.Custom {
			description = color.YellowString(descriptionCustomField)
		}
		rows = append(rows, []string{field.Name, field.Type,
			helper.WrapMultiLineTextByWord(description), helper.WrapMultiLineTextByWord(field.ExampleValue)})
	}
	return rows
}

func (s *Service) saveSchema() error {
	if !s.updated {
		return nil
	}
	schemaBytes, _ := json.Marshal(s.schema)
	err := ioutil.WriteFile(s.schemaSavePath, schemaBytes, 0644)
	if err != nil {
		return errors.New(fmt.Sprintf("failed while saving schema, err:%s", err))
	}

	if len(s.args.AddFields) > 0 {
		addFieldsMsg := fmt.Sprintf("%d fields were added successfully. "+
			"Use 'api_tool schema --industry [industry] --table [table] --show' "+
			"to view the schema", len(s.args.AddFields))
		fmt.Println(addFieldsMsg)
	}

	if len(s.args.DeleteFields) > 0 {
		deleteFieldsMsg := fmt.Sprintf("%d fields were deleted successfully. "+
			"Use 'api_tool schema --industry [industry] --table [table] --show' "+
			"to view the schema", len(s.args.DeleteFields))
		fmt.Println(deleteFieldsMsg)
	}
	return nil
}
