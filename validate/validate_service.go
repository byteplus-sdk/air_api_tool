package validate

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/3rd_rec/air_api_tool/checker"
	"github.com/3rd_rec/air_api_tool/consts"
	"github.com/3rd_rec/air_api_tool/helper"
	"github.com/3rd_rec/air_api_tool/reporter"
	"github.com/3rd_rec/air_api_tool/schema"
	"github.com/3rd_rec/air_api_tool/shared"
	"github.com/3rd_rec/air_api_tool/validate/parser"
	"io"
	"os"
	"strings"
)

type Service struct {
	args            *Args
	localFileReader io.ReadCloser
	lineScanner     *bufio.Scanner
	checker         *checker.Service
	schema          *schema.Service
	reporter        *reporter.Collector
}

func NewService(args *Args) *Service {
	return &Service{args: args}
}

func (s *Service) Run() error {
	defer func() {
		if s.localFileReader != nil {
			_ = s.localFileReader.Close()
		}
	}()
	if err := s.checkArgs(); err != nil {
		return err
	}
	if err := s.init(); err != nil {
		return err
	}
	if err := s.parseAndValidate(); err != nil {
		return err
	}
	return nil
}

func (s *Service) checkArgs() error {
	if err := helper.CheckCommonArgs(s.args.Industry, s.args.Table); err != nil {
		return err
	}
	if len(s.args.FilePath) == 0 {
		return errors.New("--file-path args is required")
	}
	if len(s.args.ContentType) == 0 {
		return errors.New("--content-type args is required")
	}
	if _, exist := parser.GetParser(s.args.ContentType); !exist {
		errMsg := fmt.Sprintf("invalid --content-type args(%s). The accepted values are: [%s]",
			s.args.ContentType, strings.Join(parser.GetSupportedContentType(), ","))
		return errors.New(errMsg)
	}
	return nil
}

func (s *Service) init() error {
	var err error

	s.localFileReader, err = os.Open(s.args.FilePath)
	if err != nil {
		errMsg := fmt.Sprintf("reading local file: %s failed, err:%s", s.args.FilePath, err)
		return errors.New(errMsg)
	}
	s.lineScanner, err = s.createLineScanner(s.localFileReader)
	if err != nil {
		return err
	}
	s.schema, err = schema.NewService(s.args.Namespace, s.args.Industry, s.args.Table)
	if err != nil {
		return err
	}
	s.checker, err = checker.NewService(s.args.Namespace, s.args.Industry, s.args.Table)
	if err != nil {
		return err
	}
	s.reporter = reporter.NewReporter(s.args.FilePath)
	return nil
}

const (
	defaultMaxBufferSize = 4 * consts.MB
	startBufferSize      = consts.MB
)

func (s *Service) createLineScanner(reader io.Reader) (*bufio.Scanner, error) {
	bufReader := bufio.NewReaderSize(reader, defaultMaxBufferSize)
	lineScanner := bufio.NewScanner(bufReader)
	lineScanner.Buffer(make([]byte, 0, startBufferSize), defaultMaxBufferSize)
	return lineScanner, nil
}

func (s *Service) parseAndValidate() error {
	lineNumber := 1

	dataParser, _ := parser.GetParser(s.args.ContentType)
	for s.lineScanner.Scan() {
		lineBytes := s.lineScanner.Bytes()
		if len(lineBytes) == 0 {
			continue
		}
		parsedData, parseResult, err := dataParser.Parse(lineBytes)
		if err != nil {
			return err
		}

		dataExtension := shared.NewDataExtension()
		dataExtension.LineNumber = lineNumber
		dataExtension.LineString = string(lineBytes)
		dataExtension.OriginData = parsedData
		dataExtension.ParseResult = parseResult

		s.doValidateAndReport(dataExtension)
		if err != nil {
			return err
		}

		lineNumber += 1
		if s.args.Lines > 0 && lineNumber > s.args.Lines {
			break
		}
	}
	s.reporter.Show()
	return nil
}

func (s *Service) doValidateAndReport(data *shared.DataExtension) {
	if data.ParseResult != nil {
		s.reporter.Report(data)
		return
	}
	// The first line of the CSV file is CSV Header, parsedData==nil
	if data.OriginData == nil {
		s.reporter.Report(data)
		return
	}
	s.schema.Validate(data)
	s.checker.Validate(data)
	s.reporter.Report(data)
	return
}
