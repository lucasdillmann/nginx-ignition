package logline

import (
	"regexp"
	"strings"
)

func Search(logLines []LogLine, query string, surroundingLines int) ([]LogLine, error) {
	if query == "" {
		return logLines, nil
	}

	re, err := compileRegex(query)
	if err != nil {
		return nil, err
	}

	surroundingLines = clampSurroundingLines(surroundingLines)
	matchedIndices := findMatchedIndices(logLines, re)
	if len(matchedIndices) == 0 {
		return []LogLine{}, nil
	}

	includeMap := calculateIncludedIndices(matchedIndices, surroundingLines, len(logLines))
	return buildResult(logLines, includeMap, re), nil
}

func compileRegex(query string) (*regexp.Regexp, error) {
	parts := strings.Split(query, " ")
	for i, part := range parts {
		parts[i] = regexp.QuoteMeta(part)
	}

	regexStr := "(?i)" + strings.Join(parts, ".*?")
	return regexp.Compile(regexStr)
}

func clampSurroundingLines(lines int) int {
	if lines < 0 {
		return 0
	}

	if lines > 10 {
		return 10
	}

	return lines
}

func findMatchedIndices(logLines []LogLine, re *regexp.Regexp) []int {
	var matchedIndices []int
	for index, line := range logLines {
		if re.MatchString(line.Contents) {
			matchedIndices = append(matchedIndices, index)
		}
	}

	return matchedIndices
}

func calculateIncludedIndices(matchedIndices []int, surroundingLines, totalLines int) map[int]bool {
	includeMap := make(map[int]bool)
	for _, index := range matchedIndices {
		start := index - surroundingLines
		if start < 0 {
			start = 0
		}

		end := index + surroundingLines
		if end >= totalLines {
			end = totalLines - 1
		}

		for idx := start; idx <= end; idx++ {
			includeMap[idx] = true
		}
	}

	return includeMap
}

func buildResult(logLines []LogLine, includeMap map[int]bool, re *regexp.Regexp) []LogLine {
	var result []LogLine

	for index := 0; index < len(logLines); index++ {
		if includeMap[index] {
			line := logLines[index]
			loc := re.FindStringIndex(line.Contents)

			if loc != nil {
				line.Highlight = &Highlight{
					Start: loc[0],
					End:   loc[1],
				}
			}

			result = append(result, line)
		}
	}

	return result
}
