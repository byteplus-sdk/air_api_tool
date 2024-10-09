package helper

import (
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const (
	defaultTerminalWide = 150
	defaultColumnCount  = 5
)

func WrapMultiLineTextByWord(text string, columnCount ...int) string {
	// 保留原始文本中的换行符
	lines := strings.Split(text, "\n")
	var wrappedLines []string

	for _, line := range lines {
		wrappedLines = append(wrappedLines, WrapTextByWord(line, columnCount...))
	}
	return strings.TrimPrefix(strings.Join(wrappedLines, "\n"), "\n")
}

func WrapTextByWord(text string, columnCount ...int) string {
	maxCharacters := getMaxCharacters(columnCount...)
	words := strings.Fields(text)
	var lines []string
	currentLine := ""
	for _, word := range words {
		if len(currentLine)+len(word) > maxCharacters {
			lines = append(lines, currentLine)
			currentLine = ""
		}

		if len(word) > maxCharacters {
			lines = append(lines, WrapText(word))
			currentLine = ""
		} else {
			if currentLine != "" {
				currentLine += " "
			}
			currentLine += word
		}
	}

	if currentLine != "" {
		lines = append(lines, currentLine)
	}

	return strings.TrimPrefix(strings.Join(lines, "\n"), "\n")
}

func WrapText(text string, columnCount ...int) string {
	var result string
	lineLength := 0
	maxCharacters := getMaxCharacters(columnCount...)

	for _, char := range text {
		result += string(char)
		lineLength++

		if lineLength == maxCharacters {
			result += "\n"
			lineLength = 0
		}
	}

	return strings.TrimPrefix(result, "\n")
}

func getMaxCharacters(columnCountList ...int) int {
	var (
		maxCharacters int
		terminalWide  int
		columnCount   = defaultColumnCount
	)

	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, _ := cmd.Output()
	if len(out) > 0 {
		outStr := strings.TrimSuffix(string(out), "\n")
		// length, wide
		terminalSize := strings.Split(outStr, " ")
		if len(terminalSize) == 2 {
			terminalWide, _ = strconv.Atoi(terminalSize[1])
		}
	}
	if terminalWide <= 0 {
		terminalWide = defaultTerminalWide
	}
	if len(columnCountList) > 0 {
		columnCount = columnCountList[0]
	}

	maxCharacters = terminalWide / columnCount
	return maxCharacters
}
