package parser

import (
	"bufio"
	"github.com/3rd_rec/air_api_tool/consts"
	"io"
)

const (
	defaultMaxBufferSize = 4 * consts.MB
	startBufferSize      = consts.MB
)

func createLineScanner(reader io.ReadCloser) *bufio.Scanner {
	bufReader := bufio.NewReaderSize(reader, defaultMaxBufferSize)
	lineScanner := bufio.NewScanner(bufReader)
	lineScanner.Buffer(make([]byte, 0, startBufferSize), defaultMaxBufferSize)
	return lineScanner
}
