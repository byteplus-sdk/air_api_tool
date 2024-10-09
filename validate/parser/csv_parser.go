package parser

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/3rd_rec/air_api_tool/consts"
	"github.com/3rd_rec/air_api_tool/shared"
)

type CSVParser struct {
	columns []string
}

func (p *CSVParser) Parse(lineBytes []byte) (map[string]interface{}, *shared.ParseResult, error) {
	var err error

	// 1. parse columns first
	if p.columns == nil {
		r := csv.NewReader(bytes.NewReader(lineBytes))
		p.columns, err = r.Read()
		if err != nil {
			return nil, nil, errors.New(fmt.Sprintf("failed to read csv header, err:%s", err))
		}
		return nil, nil, nil
	}

	r := csv.NewReader(bytes.NewReader(lineBytes))
	record, err := r.Read()
	if err != nil {
		return nil, &shared.ParseResult{
			ErrType:      consts.ParseErrorTypeFormatError,
			ErrorDetails: fmt.Sprintf("error reading CSV: %s", err),
		}, nil
	}
	if len(record) != len(p.columns) {
		return nil, &shared.ParseResult{
			ErrType:      consts.ParseErrorTypeInvalidDataLength,
			ErrorDetails: fmt.Sprintf("data len: %d not equal to column(header) len: %d", len(record), len(p.columns)),
		}, nil
	}
	target := make(map[string]interface{}, len(record))
	for idx := range record {
		target[p.columns[idx]] = record[idx]
	}
	return target, nil, nil
}
