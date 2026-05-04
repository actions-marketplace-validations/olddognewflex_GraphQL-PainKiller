package reporters

import (
	"fmt"
	"strings"

	"github.com/olddognewflex/graphql-painkiller/internal/models"
)

func SummaryMarkdown(reports []models.Report) string {
	if len(reports) == 0 {
		return "## GraphQL Painkiller Summary\n\nNo GraphQL operations found."
	}

	var b strings.Builder
	fmt.Fprintf(&b, "## GraphQL Painkiller Summary\n\n")
	fmt.Fprintf(&b, "**Operations scanned:** %d\n\n", len(reports))

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

	fmt.Fprintf(&b, "**Files with operations:** %d\n\n", len(uniqueFiles))
	fmt.Fprintf(&b, "**Total findings:** %d\n\n", totalFindings)

	if totalFindings > 0 {
		b.WriteString("**Findings by severity:**\n\n")
		for _, sev := range []string{"CRITICAL", "HIGH", "WARNING", "INFO"} {
			if count := severityCounts[sev]; count > 0 {
				fmt.Fprintf(&b, "- %s: %d\n", sev, count)
			}
		}
		b.WriteString("\n")
	}

	if len(reports) > 0 {
		fmt.Fprintf(&b, "**Average risk score:** %d/10\n\n", totalRisk/len(reports))
		fmt.Fprintf(&b, "**Max risk score:** %d/10\n", maxRisk)
	}

	return b.String()
}
