package rules

import (
	"fmt"

	"github.com/olddognewflex/graphql-painkiller/internal/config"
	"github.com/olddognewflex/graphql-painkiller/internal/extractors"
	"github.com/olddognewflex/graphql-painkiller/internal/models"
	"github.com/olddognewflex/graphql-painkiller/internal/severity"
)

func MissingPagination(fields []models.FieldInfo, doc extractors.Document, cfg config.Config) []models.Finding {
	if !cfg.Rules.RequirePagination {
		return nil
	}

	var findings []models.Finding

	for _, field := range Flatten(fields) {
		if len(field.Children) == 0 {
			continue
		}
		if !LooksCollectionLike(field, cfg) {
			continue
		}
		if HasPagination(field, cfg) {
			continue
		}

		findings = append(findings, models.Finding{
			RuleID:      "MISSING_PAGINATION",
			Severity:    severity.High,
			Message:     fmt.Sprintf("%s looks like a collection but has no pagination argument.", field.Path),
			FilePath:    doc.FilePath,
			Line:        AdjustedLine(doc.StartLine, field.Line),
			Path:        field.Path,
			ScoreImpact: 2,
			Suggestion:  "Add pagination args such as first, limit, take, pageSize, after, or before.",
		})
	}

	return findings
}
