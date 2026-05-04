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
			if finding.DocsURL != "" {
				fmt.Fprintf(w, "    Docs: %s\n", finding.DocsURL)
			}
		}

		fmt.Fprintln(w)
	}

	fmt.Fprintln(w, "Summary")
	fmt.Fprintln(w, strings.Repeat("=", 26))
	fmt.Fprintf(w, "Operations scanned: %d\n", len(reports))

	totalFindings := 0
	severityCounts := make(map[string]int)
	uniqueFiles := make(map[string]struct{})
	maxRisk := 0
	totalRisk := 0

	for _, report := range reports {
		uniqueFiles[report.FilePath] = struct{}{}
		totalFindings += len(report.Findings)
		for _, finding := range report.Findings {
			severityCounts[strings.ToUpper(string(finding.Severity))]++
		}
		if report.RiskScore > maxRisk {
			maxRisk = report.RiskScore
		}
		totalRisk += report.RiskScore
	}

	fmt.Fprintf(w, "Files with operations: %d\n", len(uniqueFiles))
	fmt.Fprintf(w, "Total findings: %d\n", totalFindings)
	if totalFindings > 0 {
		fmt.Fprintln(w, "Findings by severity:")
		for _, sev := range []string{"CRITICAL", "HIGH", "WARNING", "INFO"} {
			if count := severityCounts[sev]; count > 0 {
				fmt.Fprintf(w, "  - %s: %d\n", sev, count)
			}
		}
	}
	if len(reports) > 0 {
		fmt.Fprintf(w, "Average risk score: %d/10\n", totalRisk/len(reports))
		fmt.Fprintf(w, "Max risk score: %d/10\n", maxRisk)
	}
	fmt.Fprintln(w)
}
