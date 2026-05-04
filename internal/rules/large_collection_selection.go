package rules

import (
	"fmt"

	"github.com/olddognewflex/graphql-painkiller/internal/config"
	"github.com/olddognewflex/graphql-painkiller/internal/extractors"
	"github.com/olddognewflex/graphql-painkiller/internal/models"
	"github.com/olddognewflex/graphql-painkiller/internal/severity"
)

func LargeCollectionSelection(fields []models.FieldInfo, doc extractors.Document, cfg config.Config) []models.Finding {
	var findings []models.Finding

	for _, field := range Flatten(fields) {
		if !LooksCollectionLike(field, cfg) {
			continue
		}

		if len(field.Children) <= cfg.Rules.MaxCollectionSelectionFields {
			continue
		}

		findings = append(findings, models.Finding{
			RuleID:      "LARGE_COLLECTION_SELECTION",
			Severity:    severity.Warning,
			Message:     fmt.Sprintf("%s has %d selected fields under a collection-like field.", field.Path, len(field.Children)),
			FilePath:    doc.FilePath,
			Line:        AdjustedLine(doc.StartLine, field.Line),
			Path:        field.Path,
			ScoreImpact: 2,
			Suggestion:  "Reduce selected fields or split heavy detail fields into a follow-up query.",
			DocsURL:     "https://graphql.org/learn/performance/#demand-control",
		})
	}

	return findings
}
