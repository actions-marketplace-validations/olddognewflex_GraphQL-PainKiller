package rules

import (
	"fmt"
	"strings"

	"github.com/olddognewflex/graphql-painkiller/internal/config"
	"github.com/olddognewflex/graphql-painkiller/internal/extractors"
	"github.com/olddognewflex/graphql-painkiller/internal/models"
	"github.com/olddognewflex/graphql-painkiller/internal/severity"
)

func ExpensiveFields(fields []models.FieldInfo, doc extractors.Document, cfg config.Config) []models.Finding {
	var findings []models.Finding

	for _, field := range Flatten(fields) {
		lower := strings.ToLower(field.Name)

		for _, pattern := range cfg.ExpensiveFieldPatterns {
			if strings.Contains(lower, strings.ToLower(pattern)) {
			findings = append(findings, models.Finding{
				RuleID:      "EXPENSIVE_FIELD_PATTERN",
				Severity:    severity.Warning,
				Message:     fmt.Sprintf("%s matches expensive field pattern %q.", field.Path, pattern),
				FilePath:    doc.FilePath,
				Line:        AdjustedLine(doc.StartLine, field.Line),
				Path:        field.Path,
				ScoreImpact: 1,
				Suggestion:  "Confirm this field is necessary for this operation.",
				DocsURL:     "https://graphql.org/learn/performance/",
			})
				break
			}
		}
	}

	return findings
}
