package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

const (
	url = "http://rec-b-ap-singapore-1.byteplusapi.com/openapi/schema/spreadsheet/parse"
)

type LoadSchemaRequest struct {
	SpreadsheetLink     string `json:"spreadsheet_link,omitempty"`
	SubSheetName        string `json:"sub_sheet_name,omitempty"`
	FieldNameColumnName string `json:"field_name_column_name,omitempty"`
	FieldTypeColumnName string `json:"field_type_column_name,omitempty"`
}

type LoadSchemaResponse struct {
	Status *Status       `json:"status,omitempty"`
	Data   *ParsedSchema `json:"data,omitempty"`
}

type ParsedSchema struct {
	Fields []*Field `json:"fields,omitempty"`
}

type Status struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Success bool   `json:"success,omitempty"`
}

type Field struct {
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
}

func LoadSchema(request *LoadSchemaRequest) (*ParsedSchema, error) {
	reqBytes, _ := json.Marshal(request)
	bodyBuffer := bytes.NewBuffer(reqBytes)
	resp, err := http.Post(url, "application/json", bodyBuffer)
	if err != nil {
		return nil, errors.New("parse schema from lark spreadsheet fail, please try again later")
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	var respBytes []byte
	respBytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("parse schema from lark spreadsheet fail, please try again later")
	}
	parseResp := &LoadSchemaResponse{}
	err = json.Unmarshal(respBytes, parseResp)
	// if unmarshal fail, explanation is not a business exception.
	if err != nil {
		return nil, errors.New("parse schema from lark spreadsheet fail, please try again later")
	}
	if parseResp.Status.Code != 0 {
		return nil, errors.New(parseResp.Status.Message)
	}
	return parseResp.Data, nil
}
