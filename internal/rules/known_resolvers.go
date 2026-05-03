package rules

import (
	"fmt"

	"github.com/olddognewflex/graphql-painkiller/internal/config"
	"github.com/olddognewflex/graphql-painkiller/internal/extractors"
	"github.com/olddognewflex/graphql-painkiller/internal/models"
	"github.com/olddognewflex/graphql-painkiller/internal/severity"
)

func KnownResolvers(fields []models.FieldInfo, doc extractors.Document, cfg config.Config) []models.Finding {
	var findings []models.Finding

	for _, field := range Flatten(fields) {
		known, ok := cfg.KnownResolvers[field.Path]
		if !ok {
			continue
		}

		impact := 2
		if known.Risk == severity.High {
			impact = 3
		}
		if known.Risk == severity.Critical {
			impact = 4
		}

		service := ""
		if known.Service != "" {
			service = fmt.Sprintf(" Service: %s.", known.Service)
		}

		findings = append(findings, models.Finding{
			RuleID:      "KNOWN_RESOLVER_RISK",
			Severity:    known.Risk,
			Message:     fmt.Sprintf("Known resolver risk at %s.%s Reason: %s", field.Path, service, known.Reason),
			FilePath:    doc.FilePath,
			Line:        AdjustedLine(doc.StartLine, field.Line),
			Path:        field.Path,
			ScoreImpact: impact,
			Suggestion:  "Use team-defined remediation guidance for this resolver path.",
		})
	}

	return findings
}
