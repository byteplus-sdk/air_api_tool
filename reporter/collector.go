package reporter

import (
	"encoding/json"
	"fmt"
	"github.com/3rd_rec/air_api_tool/consts"
	"github.com/3rd_rec/air_api_tool/helper"
	"github.com/3rd_rec/air_api_tool/shared"
	"github.com/3rd_rec/air_api_tool/utils"
	"github.com/fatih/color"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
)

type Collector struct {
	Metadata *Metadata
	// ErrorType -> DataParseErrorStats
	DataParseErrorSummary map[string]*DataParseErrorStats
	// FieldName -> DataTypeCheckErrorStats
	DataTypeCheckErrorSummary map[string]*DataTypeCheckErrorStats
	// FieldName -> DataValueCheckErrorStats
	DataValueCheckErrorSummary map[string]*DataValueCheckErrorStats
}

type Metadata struct {
	Filepath                 string
	SampleLines              []string
	SampleDatas              []string
	LineCount                int
	TotalCount               int
	SuccessCount             int
	DataParseErrorCount      int
	DataTypeCheckErrorCount  int
	DataValueCheckErrorCount int
}

type DataParseErrorStats struct {
	ErrorType        string
	ErrorDetails     string
	FailureCount     int
	SampleLineNumber int
	SampleLine       string
}

type DataTypeCheckErrorStats struct {
	FieldName     string
	FieldType     string
	FailureReason string
	FailureCount  int
	SampleData    string
}

type DataValueCheckErrorStats struct {
	FieldName      string
	CheckerDetails string
	FailureReason  string
	FailureCount   int
	SampleData     string
}

func NewReporter(filepath string) *Collector {
	return &Collector{
		Metadata: &Metadata{
			Filepath:    filepath,
			SampleLines: make([]string, 0),
			SampleDatas: make([]string, 0),
		},
		DataParseErrorSummary:      make(map[string]*DataParseErrorStats),
		DataTypeCheckErrorSummary:  make(map[string]*DataTypeCheckErrorStats),
		DataValueCheckErrorSummary: make(map[string]*DataValueCheckErrorStats),
	}
}

func (r *Collector) Report(data *shared.DataExtension) {
	r.Metadata.LineCount += 1
	r.addSampleLines(data)
	// csv header does not count
	if data.OriginData == nil && data.ParseResult == nil {
		return
	}

	r.Metadata.TotalCount += 1
	success1 := r.collectDataParseResult(data)
	success2 := r.collectDataTypeCheckError(data)
	success3 := r.collectDataValueCheckError(data)

	// no errors
	if !success1 && !success2 && !success3 {
		r.Metadata.SuccessCount += 1
	}
}

const sampleLineCount = 5

func (r *Collector) addSampleLines(data *shared.DataExtension) {
	if len(r.Metadata.SampleLines) < sampleLineCount {
		formattedSampleLine := fmt.Sprintf("%d %s", data.LineNumber, data.LineString)
		r.Metadata.SampleLines = append(r.Metadata.SampleLines, formattedSampleLine)
	}
	if len(r.Metadata.SampleDatas) < sampleLineCount {
		if data.StandardizedData != nil {
			dataBytes, _ := json.Marshal(data.StandardizedData)
			formattedSampleData := fmt.Sprintf("%d %s", data.LineNumber, dataBytes)
			r.Metadata.SampleDatas = append(r.Metadata.SampleDatas, formattedSampleData)
		}
	}
}

func (r *Collector) collectDataParseResult(data *shared.DataExtension) bool {
	parseResult := data.ParseResult
	if parseResult == nil {
		return false
	}
	r.Metadata.DataParseErrorCount += 1
	if _, exist := r.DataParseErrorSummary[parseResult.ErrType]; !exist {
		r.DataParseErrorSummary[parseResult.ErrType] = &DataParseErrorStats{
			FailureCount:     0,
			ErrorType:        parseResult.ErrType,
			ErrorDetails:     parseResult.ErrorDetails,
			SampleLineNumber: data.LineNumber,
			SampleLine:       data.LineString,
		}
	}
	r.DataParseErrorSummary[parseResult.ErrType].FailureCount += 1
	return true
}

func (r *Collector) collectDataTypeCheckError(data *shared.DataExtension) bool {
	fieldValidateResult := data.FieldValidateResult
	if len(fieldValidateResult) == 0 {
		return false
	}
	success := false
	summary := r.DataTypeCheckErrorSummary
	for fieldName, validateItem := range fieldValidateResult {
		if validateItem.DataTypeErrorItem == nil {
			continue
		}
		success = true
		if _, exist := summary[fieldName]; !exist {
			sampleDataBytes, _ := json.Marshal(data.OriginData)
			summary[fieldName] = &DataTypeCheckErrorStats{
				FieldName:     fieldName,
				FieldType:     validateItem.DataTypeErrorItem.FieldType,
				FailureReason: validateItem.DataTypeErrorItem.FailureReason,
				FailureCount:  0,
				SampleData:    string(sampleDataBytes),
			}
		}
		summary[fieldName].FailureCount += 1
	}
	if success {
		r.Metadata.DataTypeCheckErrorCount += 1
	}
	return success
}

func (r *Collector) collectDataValueCheckError(data *shared.DataExtension) bool {
	fieldValidateResult := data.FieldValidateResult
	if len(fieldValidateResult) == 0 {
		return false
	}
	success := false
	summary := r.DataValueCheckErrorSummary
	for fieldName, validateItem := range fieldValidateResult {
		if validateItem.DataValueErrorItem == nil {
			continue
		}
		success = true
		if _, exist := summary[fieldName]; !exist {
			sampleDataBytes, _ := json.Marshal(data.OriginData)
			summary[fieldName] = &DataValueCheckErrorStats{
				FieldName:      fieldName,
				CheckerDetails: validateItem.DataValueErrorItem.CheckerDetails,
				FailureReason:  validateItem.DataValueErrorItem.FailureReason,
				FailureCount:   0,
				SampleData:     string(sampleDataBytes),
			}
		}
		summary[fieldName].FailureCount += 1
	}
	if success {
		r.Metadata.DataValueCheckErrorCount += 1
	}
	return success
}

func (r *Collector) Show() {
	r.showResult()
	r.showMetaInfo()
	r.showValidationSummary()
	r.showDataParseErrorSummary()
	r.showDataTypeCheckErrorSummary()
	r.showDataValueCheckErrorSummary()
}

func (r *Collector) showResult() {
	status := consts.VerificationStatusSuccess
	if r.Metadata.TotalCount > r.Metadata.SuccessCount {
		status = consts.VerificationStatusFailure
	}
	fmt.Println(color.New(color.Bold).Sprintf("[Verification Result] %s", status))
	r.newChapter()
}

func (r *Collector) newChapter() {
	r.newline()
	r.newline()
}

func (r *Collector) newline() {
	fmt.Println()
}

func (r *Collector) showMetaInfo() {
	fmt.Println(color.New(color.Bold).Sprint("[MetaInfo]"))

	// build sample lines
	sampleLineList := make([]string, 0, len(r.Metadata.SampleLines))
	for _, sampleLine := range r.Metadata.SampleLines {
		sampleLineList = append(sampleLineList, fmt.Sprintf(" * %s", sampleLine))
	}
	sampleLines := strings.Join(sampleLineList, "\n")

	// build sample datas
	sampleDataList := make([]string, 0, len(r.Metadata.SampleLines))
	for _, sampleData := range r.Metadata.SampleDatas {
		sampleDataList = append(sampleDataList, fmt.Sprintf(" * %s", sampleData))
	}
	sampleDatas := strings.Join(sampleDataList, "\n")

	fmt.Printf(`- filepath: %s
- total line cnt: %d
- sample lines:
%s
- sample datas(parsed):
%s
`, r.Metadata.Filepath, r.Metadata.LineCount, sampleLines, sampleDatas)
	r.newChapter()
}

func (r *Collector) showValidationSummary() {
	fmt.Println(color.New(color.Bold).Sprint("[Validation Summary]"))
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"File Path", "Total Count", "Success Count", "Failure Count", "Error Distribution"})
	table.SetBorder(true)
	table.SetRowLine(true)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetAutoMergeCellsByColumnIndex([]int{0, 1, 2, 3})

	summary := [][]string{
		{
			helper.WrapMultiLineTextByWord(r.Metadata.Filepath), utils.ToA(r.Metadata.TotalCount),
			utils.ToA(r.Metadata.SuccessCount), utils.ToA(r.Metadata.TotalCount - r.Metadata.SuccessCount),
			fmt.Sprintf("Data Parse Error: %d", r.Metadata.DataParseErrorCount),
		},
		{
			helper.WrapMultiLineTextByWord(r.Metadata.Filepath), utils.ToA(r.Metadata.TotalCount),
			utils.ToA(r.Metadata.SuccessCount), utils.ToA(r.Metadata.TotalCount - r.Metadata.SuccessCount),
			fmt.Sprintf("Data Type Check Error: %d", r.Metadata.DataTypeCheckErrorCount),
		},
		{
			helper.WrapMultiLineTextByWord(r.Metadata.Filepath), utils.ToA(r.Metadata.TotalCount),
			utils.ToA(r.Metadata.SuccessCount), utils.ToA(r.Metadata.TotalCount - r.Metadata.SuccessCount),
			fmt.Sprintf("Data Value Check Error: %d", r.Metadata.DataValueCheckErrorCount),
		},
	}
	table.AppendBulk(summary)
	table.Render()
	r.newChapter()
}

func (r *Collector) showDataParseErrorSummary() {
	if len(r.DataParseErrorSummary) == 0 {
		return
	}
	fmt.Println(color.New(color.Bold).Sprint("[Data Parse Error Summary]"))
	fmt.Printf("Refer to %s for more guidance on data repair.\n",
		color.HiYellowString("https://bytedance.larkoffice.com/docx/FvFvdQysyorz9cxziSJc7UGfngc#GSlMd9IxFoYVxhxCZgtcKZJnnRd"))
	r.newline()
	for errorType, stats := range r.DataParseErrorSummary {
		fmt.Printf("- %s\n", errorType)
		fmt.Printf("  ∟ Failure Count: %d\n", stats.FailureCount)
		fmt.Printf("  ∟ Error Details: %s\n", stats.ErrorDetails)
		fmt.Printf("  ∟ Sample Line Number: %d\n", stats.SampleLineNumber)
		fmt.Printf("  ∟ Sample Line: %s\n", stats.SampleLine)
	}
	r.newChapter()
}

func (r *Collector) showDataTypeCheckErrorSummary() {
	if len(r.DataTypeCheckErrorSummary) == 0 {
		return
	}
	fmt.Println(color.New(color.Bold).Sprint("[Data Type Check Error Summary]"))
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Field Name", "Data Type", "Failure Count", "Failure Reason", "Sample Data"})
	table.SetBorder(true)
	table.SetRowLine(true)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	summary := make([][]string, 0)
	for _, stats := range r.DataTypeCheckErrorSummary {
		row := []string{
			stats.FieldName,
			stats.FieldType,
			utils.ToA(stats.FailureCount),
			helper.WrapMultiLineTextByWord(stats.FailureReason),
			helper.WrapMultiLineTextByWord(stats.SampleData),
		}
		summary = append(summary, row)
	}
	table.AppendBulk(summary)
	table.Render()
	r.newChapter()
}

func (r *Collector) showDataValueCheckErrorSummary() {
	if len(r.DataValueCheckErrorSummary) == 0 {
		return
	}
	fmt.Println(color.New(color.Bold).Sprint("[Data Value Check Error Summary]"))
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Field Name", "Checker Details", "Failure Count", "Failure Reason", "Sample Data"})
	table.SetBorder(true)
	table.SetRowLine(true)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	summary := make([][]string, 0)
	for _, stats := range r.DataValueCheckErrorSummary {
		row := []string{
			stats.FieldName,
			stats.CheckerDetails,
			utils.ToA(stats.FailureCount),
			helper.WrapMultiLineTextByWord(stats.FailureReason),
			helper.WrapMultiLineTextByWord(stats.SampleData),
		}
		summary = append(summary, row)
	}
	table.AppendBulk(summary)
	table.Render()
	r.newChapter()
}
