package parser

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/3rd_rec/air_api_tool/consts"
	"github.com/3rd_rec/air_api_tool/shared"
	"io"
)

type JSONParser struct {
	filepath        string
	localFileReader io.ReadCloser
	lineScanner     *bufio.Scanner
	currentLines    int
	done            bool
}

func (p *JSONParser) Parse(lineBytes []byte) (map[string]interface{}, *shared.ParseResult, error) {
	target := make(map[string]interface{})
	decoder := json.NewDecoder(bytes.NewReader(lineBytes))
	decoder.UseNumber()
	err := decoder.Decode(&target)
	if err != nil {
		return nil, &shared.ParseResult{
			ErrType:      consts.ParseErrorTypeFormatError,
			ErrorDetails: fmt.Sprintf("err reading JSON: %s", err),
		}, nil
	}
	return target, nil, nil
}

func (p *JSONParser) Next() ([]*shared.DataExtension, error) {
	dataExtension := shared.NewDataExtension()

	var lineBytes []byte
	originData := make(map[string]interface{})
	for p.lineScanner.Scan() {
		lineBytes = p.lineScanner.Bytes()
		if len(lineBytes) == 0 {
			continue
		}
		dataExtension.LineNumber = p.currentLines
		dataExtension.LineString = string(lineBytes)
		p.currentLines += 1
		err := json.Unmarshal(lineBytes, &originData)
		if err != nil {
			dataExtension.ParseResult = &shared.ParseResult{
				ErrType:      consts.ParseErrorTypeFormatError,
				ErrorDetails: fmt.Sprintf("json format err:%s", err),
			}
		}
		dataExtension.OriginData = originData
		return []*shared.DataExtension{dataExtension}, nil
	}
	if p.lineScanner.Err() != nil {
		errMsg := fmt.Sprintf("reading local file: %s failed, err:%s", p.filepath, p.lineScanner.Err())
		return nil, errors.New(errMsg)
	}
	p.done = true
	return nil, nil
}
