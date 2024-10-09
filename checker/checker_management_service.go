package checker

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/3rd_rec/air_api_tool/helper"
	"github.com/3rd_rec/air_api_tool/utils"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"io/ioutil"
	"os"
)

type Service struct {
	args                 *Args
	checkerCfg           []*Item
	fieldCheckerWrappers []*fieldCheckerWrapper
	checkerSavePath      string
	updated              bool
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
	if len(s.args.IgnoreFields) == 0 && len(s.args.RecoverFields) == 0 && !s.args.Show && !s.args.Clear {
		return errors.New("no valid args. Use `api_tool checker --help` to view the usage")
	}
	return nil
}

func (s *Service) init() error {
	var err error

	baseDir := helper.BuildModuleDir("checker")
	err = utils.EnsureDirExists(baseDir)
	if err != nil {
		return errors.New(fmt.Sprintf("Tmp dir %s create fail, err:%s", baseDir, err))
	}
	s.checkerSavePath = buildCheckerSavePath(s.args.Namespace, s.args.Industry, s.args.Table)
	s.checkerCfg, err = GetCheckerCfg(s.args.Namespace, s.args.Industry, s.args.Table)
	if err != nil {
		return err
	}
	if err = s.initChecker(); err != nil {
		return err
	}
	return nil
}

func (s *Service) executeCmd() error {
	defer func() {
		s.show()
	}()

	if err := s.ignoreFields(); err != nil {
		return err
	}
	if err := s.recoverFields(); err != nil {
		return err
	}
	if err := s.saveChecker(); err != nil {
		return err
	}
	return nil
}

func (s *Service) clear() error {
	if !s.args.Clear {
		return nil
	}
	s.checkerSavePath = buildCheckerSavePath(s.args.Namespace, s.args.Industry, s.args.Table)
	if !utils.FileExists(s.checkerSavePath) {
		fmt.Println("cleanup checker successful. Recover to default configuration.")
		return nil
	}
	err := os.Remove(s.checkerSavePath)
	if err != nil {
		return errors.New(fmt.Sprintf("failed while clear local checker configuration: %s, err:%s", s.checkerSavePath, err))
	}
	fmt.Println("cleanup checker successful. Recover to default configuration.")
	return nil
}

func (s *Service) ignoreFields() error {
	if len(s.args.IgnoreFields) == 0 {
		return nil
	}

	for _, item := range s.checkerCfg {
		for _, needIgnoredFieldName := range s.args.IgnoreFields {
			if item.FieldName == needIgnoredFieldName {
				item.Ignored = true
				break
			}
		}
	}
	s.updated = true
	return nil
}

func (s *Service) recoverFields() error {
	if len(s.args.RecoverFields) == 0 {
		return nil
	}

	for _, item := range s.checkerCfg {
		for _, needRecoveredFieldName := range s.args.RecoverFields {
			if item.FieldName == needRecoveredFieldName {
				item.Ignored = false
				break
			}
		}
	}
	s.updated = true
	return nil
}

func (s *Service) saveChecker() error {
	if !s.updated {
		return nil
	}
	checkerBytes, _ := json.Marshal(s.checkerCfg)
	err := ioutil.WriteFile(s.checkerSavePath, checkerBytes, 0644)
	if err != nil {
		return errors.New(fmt.Sprintf("failed while saving checker, err:%s", err))
	}

	if len(s.args.IgnoreFields) > 0 {
		addFieldsMsg := fmt.Sprintf("%d fields were ignored successfully. "+
			"Use 'api_tool checker --industry [industry] --table [table] --show' "+
			"to view the checker configuration", len(s.args.IgnoreFields))
		fmt.Println(addFieldsMsg)
	}

	if len(s.args.RecoverFields) > 0 {
		deleteFieldsMsg := fmt.Sprintf("%d fields were recovered successfully. "+
			"Use 'api_tool checker --industry [industry] --table [table] --show' "+
			"to view the checker configuration", len(s.args.RecoverFields))
		fmt.Println(deleteFieldsMsg)
	}
	return nil
}

func (s *Service) show() {
	if !s.args.Show {
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Field Name", "Checker", "Ignored", "Checker Details", "Description"})
	table.SetBorder(true)
	table.SetAutoMergeCellsByColumnIndex([]int{0})
	table.SetRowLine(true)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	table.AppendBulk(s.toTableRows(s.fieldCheckerWrappers))
	table.Render()
}

func (s *Service) toTableRows(fieldCheckerWrappers []*fieldCheckerWrapper) [][]string {
	rows := make([][]string, 0, len(fieldCheckerWrappers))
	for _, wrapper := range fieldCheckerWrappers {
		ignored := fmt.Sprintf("%t", wrapper.item.Ignored)
		if wrapper.item.Ignored {
			ignored = color.YellowString(ignored)
		}
		checkerDetails := wrapper.fieldChecker.CheckerDetails()
		if wrapper.item.DependCfg != nil {
			checkerDetails = fmt.Sprintf("%s, %s", checkerDetails, wrapper.item.DependCfg)
		}
		rows = append(rows, []string{wrapper.item.FieldName, wrapper.item.CheckerName, ignored,
			checkerDetails, helper.WrapMultiLineTextByWord(wrapper.fieldChecker.Description())})
	}
	return rows
}

func (s *Service) initChecker() error {
	s.fieldCheckerWrappers = make([]*fieldCheckerWrapper, 0, len(s.checkerCfg))
	for _, fieldCheckerCfgItem := range s.checkerCfg {
		fieldCheckerFactory := GetFieldCheckerFactory(fieldCheckerCfgItem.CheckerName)
		fieldChecker, err := fieldCheckerFactory.Create(fieldCheckerCfgItem.CheckerCfg)
		if err != nil {
			return err
		}
		s.fieldCheckerWrappers = append(s.fieldCheckerWrappers, &fieldCheckerWrapper{
			item:         fieldCheckerCfgItem,
			fieldChecker: fieldChecker,
		})
	}
	return nil
}
