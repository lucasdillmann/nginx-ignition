package logline

import (
	"regexp"
	"strings"
)

func Search(logLines []LogLine, query string, surroundingLines int) (result []LogLine, err error) {
	if query == "" {
		return logLines, nil
	}

	parts := strings.Split(query, " ")
	for i, part := range parts {
		parts[i] = regexp.QuoteMeta(part)
	}

	regexStr := "(?i)" + strings.Join(parts, ".*")
	re, err := regexp.Compile(regexStr)
	if err != nil {
		return nil, err
	}

	if surroundingLines < 0 {
		surroundingLines = 0
	} else if surroundingLines > 10 {
		surroundingLines = 10
	}

	var matchedIndices []int
	for i, line := range logLines {
		if re.MatchString(line.Contents) {
			matchedIndices = append(matchedIndices, i)
		}
	}

	if len(matchedIndices) == 0 {
		return []LogLine{}, nil
	}

	includeMap := make(map[int]bool)
	for _, index := range matchedIndices {
		start := index - surroundingLines
		if start < 0 {
			start = 0
		}

		end := index + surroundingLines
		if end >= len(logLines) {
			end = len(logLines) - 1
		}

		for idx := start; idx <= end; idx++ {
			includeMap[idx] = true
		}
	}

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

	return result, nil
}
