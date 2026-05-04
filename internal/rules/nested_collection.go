package rules

import (
	"fmt"

	"github.com/olddognewflex/graphql-painkiller/internal/config"
	"github.com/olddognewflex/graphql-painkiller/internal/extractors"
	"github.com/olddognewflex/graphql-painkiller/internal/models"
	"github.com/olddognewflex/graphql-painkiller/internal/severity"
)

func NestedCollection(fields []models.FieldInfo, doc extractors.Document, cfg config.Config) []models.Finding {
	var findings []models.Finding

	for _, field := range Flatten(fields) {
		if !LooksCollectionLike(field, cfg) {
			continue
		}

		for _, child := range field.Children {
			if len(child.Children) == 0 {
				continue
			}

			path := field.Path + "." + child.Name

			if LooksCollectionLike(child, cfg) {
			findings = append(findings, models.Finding{
				RuleID:      "NESTED_COLLECTION_N_PLUS_ONE",
				Severity:    severity.High,
				Message:     fmt.Sprintf("%s is a collection nested under collection %s — likely N+1 resolver fan-out.", path, field.Name),
				FilePath:    doc.FilePath,
				Line:        AdjustedLine(doc.StartLine, child.Line),
				Path:        path,
				ScoreImpact: 3,
				Suggestion:  "Add pagination, use DataLoader/batching, or avoid nested collection selection.",
				DocsURL:     "https://www.graphql-js.org/docs/n1-dataloader/",
			})
			} else {
			findings = append(findings, models.Finding{
				RuleID:      "NESTED_OBJECT_UNDER_COLLECTION",
				Severity:    severity.Warning,
				Message:     fmt.Sprintf("%s resolves per record in %s — potential fan-out if not batched.", path, field.Name),
				FilePath:    doc.FilePath,
				Line:        AdjustedLine(doc.StartLine, child.Line),
				Path:        path,
				ScoreImpact: 1,
				Suggestion:  "Confirm this resolver uses batching or is inexpensive per-record.",
				DocsURL:     "https://graphql.org/learn/performance/#the-n1-problem",
			})
			}
		}
	}

	return findings
}
