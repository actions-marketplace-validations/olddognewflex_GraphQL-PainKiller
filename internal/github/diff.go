package github

import (
	"regexp"
	"strconv"
	"strings"
)

var hunkHeaderRe = regexp.MustCompile(`^@@\s+-\d+(?:,\d+)?\s+\+(\d+)(?:,\d+)?\s+@@`)

// ParsePatchLines parses a unified diff patch and returns the set of
// line numbers (on the new/right side) that appear within the diff hunks.
// These are the only lines where GitHub allows inline review comments.
func ParsePatchLines(patch string) map[int]bool {
	lines := make(map[int]bool)
	if patch == "" {
		return lines
	}

	var lineNum int
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
			lines[lineNum] = true
			lineNum++
		}
	}

	return lines
}
