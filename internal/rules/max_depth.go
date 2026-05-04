package rules

import (
	"fmt"

	"github.com/olddognewflex/graphql-painkiller/internal/config"
	"github.com/olddognewflex/graphql-painkiller/internal/extractors"
	"github.com/olddognewflex/graphql-painkiller/internal/models"
	"github.com/olddognewflex/graphql-painkiller/internal/severity"
)

func MaxDepth(fields []models.FieldInfo, doc extractors.Document, cfg config.Config) []models.Finding {
	var deepest *models.FieldInfo

	for _, field := range Flatten(fields) {
		current := field
		if deepest == nil || current.Depth > deepest.Depth {
			deepest = &current
		}
	}

	if deepest == nil || deepest.Depth <= cfg.Rules.MaxDepth {
		return nil
	}

	return []models.Finding{{
		RuleID:      "MAX_DEPTH",
		Severity:    severity.Warning,
		Message:     fmt.Sprintf("Query depth is %d, which exceeds configured max depth %d.", deepest.Depth, cfg.Rules.MaxDepth),
		FilePath:    doc.FilePath,
		Line:        AdjustedLine(doc.StartLine, deepest.Line),
		Path:        deepest.Path,
		ScoreImpact: 2,
		Suggestion:  "Consider splitting the operation or reducing nested selections.",
		DocsURL:     "https://graphql.org/learn/performance/#demand-control",
	}}
}
