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

			findings = append(findings, models.Finding{
				RuleID:      "POTENTIAL_N_PLUS_ONE",
				Severity:    severity.High,
				Message:     fmt.Sprintf("%s may cause resolver fan-out when %s returns many records.", path, field.Name),
				FilePath:    doc.FilePath,
				Line:        AdjustedLine(doc.StartLine, child.Line),
				Path:        path,
				ScoreImpact: 3,
				Suggestion:  "Confirm batching/DataLoader behavior or avoid nested selection when not needed.",
			})
		}
	}

	return findings
}
