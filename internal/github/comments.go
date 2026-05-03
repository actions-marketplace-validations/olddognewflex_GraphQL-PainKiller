package github

import (
	"fmt"
	"strings"

	"github.com/olddognewflex/graphql-painkiller/internal/models"
)

// ReviewComment is a local DTO for a future GitHub API integration.
// V1 intentionally does not call GitHub directly.
type ReviewComment struct {
	Path string `json:"path"`
	Line int    `json:"line"`
	Body string `json:"body"`
}

func BuildReviewComments(reports []models.Report) []ReviewComment {
	var comments []ReviewComment

	for _, report := range reports {
		for _, finding := range report.Findings {
			if finding.Line <= 0 {
				continue
			}

			body := fmt.Sprintf(
				"⚠️ **%s**\n\n%s\n\n**Path:** `%s`\n\n**Suggestion:** %s",
				strings.ReplaceAll(finding.RuleID, "_", " "),
				finding.Message,
				finding.Path,
				finding.Suggestion,
			)

			comments = append(comments, ReviewComment{
				Path: finding.FilePath,
				Line: finding.Line,
				Body: body,
			})
		}
	}

	return comments
}
