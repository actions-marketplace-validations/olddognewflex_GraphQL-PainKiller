package reporters

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/olddognewflex/graphql-painkiller/internal/models"
)

func Text(w io.Writer, reports []models.Report) {
	if len(reports) == 0 {
		fmt.Fprintln(w, "GraphQL Painkiller Report")
		fmt.Fprintln(w, "")
		fmt.Fprintln(w, "No GraphQL operations found.")
		return
	}

	fmt.Fprintln(w, "GraphQL Painkiller Report")
	fmt.Fprintln(w, strings.Repeat("=", 26))
	fmt.Fprintln(w)

	for _, report := range reports {
		fmt.Fprintf(w, "Operation: %s\n", report.OperationName)
		fmt.Fprintf(w, "File: %s\n", filepath.Clean(report.FilePath))
		fmt.Fprintf(w, "Risk Score: %d/10 — %s\n", report.RiskScore, strings.ToUpper(string(report.Severity)))

		if len(report.Findings) == 0 {
			fmt.Fprintln(w, "Findings: none")
			fmt.Fprintln(w)
			continue
		}

		fmt.Fprintln(w, "Findings:")
		for _, finding := range report.Findings {
			line := ""
			if finding.Line > 0 {
				line = fmt.Sprintf(":%d", finding.Line)
			}

			fmt.Fprintf(w, "  - [%s] %s%s\n", strings.ToUpper(string(finding.Severity)), filepath.Clean(finding.FilePath), line)
			fmt.Fprintf(w, "    Rule: %s\n", finding.RuleID)
			if finding.Path != "" {
				fmt.Fprintf(w, "    Path: %s\n", finding.Path)
			}
			fmt.Fprintf(w, "    %s\n", finding.Message)
			if finding.Suggestion != "" {
				fmt.Fprintf(w, "    Suggestion: %s\n", finding.Suggestion)
			}
		}

		fmt.Fprintln(w)
	}
}
