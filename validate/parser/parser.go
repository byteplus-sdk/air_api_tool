package parser

import (
	"github.com/3rd_rec/air_api_tool/consts"
	"github.com/3rd_rec/air_api_tool/shared"
)

type Parser interface {
	Parse(lineBytes []byte) (map[string]interface{}, *shared.ParseResult, error)
}

var parserMap = map[string]Parser{
	consts.ContentTypeCSV:  &CSVParser{},
	consts.ContentTypeJSON: &JSONParser{},
}

func GetSupportedContentType() []string {
	result := make([]string, 0, len(parserMap))
	for contentType := range parserMap {
		result = append(result, contentType)
	}
	return result
}

func GetParser(contentType string) (Parser, bool) {
	parser, exist := parserMap[contentType]
	return parser, exist
}
