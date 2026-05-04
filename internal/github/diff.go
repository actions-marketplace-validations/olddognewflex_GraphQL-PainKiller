package github

import (
	"regexp"
	"strconv"
	"strings"
)

var hunkHeaderRe = regexp.MustCompile(`^@@\s+-\d+(?:,\d+)?\s+\+(\d+)(?:,\d+)?\s+@@`)

func ParsePatchLines(patch string) map[int]int {
	lines := make(map[int]int)
	if patch == "" {
		return lines
	}

	var lineNum int
	var position int
	for _, raw := range strings.Split(patch, "\n") {
		if m := hunkHeaderRe.FindStringSubmatch(raw); m != nil {
			n, err := strconv.Atoi(m[1])
			if err != nil {
				continue
			}
			lineNum = n
			continue
		}

		if strings.HasPrefix(raw, "-") {
			continue
		}

		if strings.HasPrefix(raw, "+") || strings.HasPrefix(raw, " ") {
			position++
			lines[lineNum] = position
			lineNum++
		}
	}

	return lines
}
