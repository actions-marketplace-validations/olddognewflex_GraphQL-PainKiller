package analyzer

import (
	"testing"

	"github.com/olddognewflex/graphql-painkiller/internal/config"
	"github.com/olddognewflex/graphql-painkiller/internal/extractors"
	"github.com/olddognewflex/graphql-painkiller/internal/severity"
)

func TestAnalyzeDocument_MaxDepth(t *testing.T) {
	cfg := config.Config{
		Rules: config.Rules{
			MaxDepth:                     3,
			MaxCollectionSelectionFields: 100,
			RequirePagination:            false,
		},
	}

	tests := []struct {
		name      string
		source    string
		wantRule  string
		wantPath  string
		wantSev   severity.Severity
		wantScore int
	}{
		{
			name: "query within max depth has no max depth finding",
			source: `query GetUser {
				user {
					name
				}
			}`,
		},
		{
			name: "query exceeding max depth produces finding",
			source: `query DeepQuery {
				a {
					b {
						c {
							d
						}
					}
				}
			}`,
			wantRule:  "MAX_DEPTH",
			wantPath:  "a.b.c.d",
			wantSev:   severity.Warning,
			wantScore: 2,
		},
		{
			name: "query at exact max depth boundary has no finding",
			source: `query BoundaryQuery {
				a {
					b {
						c
					}
				}
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := extractors.Document{
				FilePath:  "test.graphql",
				Source:    tt.source,
				StartLine: 1,
			}

			reports, err := AnalyzeDocument(doc, cfg)
			if err != nil {
				t.Fatalf("AnalyzeDocument() error = %v", err)
			}

			if len(reports) != 1 {
				t.Fatalf("expected 1 report, got %d", len(reports))
			}

			var found bool
			for _, finding := range reports[0].Findings {
				if finding.RuleID == "MAX_DEPTH" {
					found = true
					if finding.Path != tt.wantPath {
						t.Errorf("Path = %q, want %q", finding.Path, tt.wantPath)
					}
					if finding.Severity != tt.wantSev {
						t.Errorf("Severity = %q, want %q", finding.Severity, tt.wantSev)
					}
				}
			}

			if tt.wantRule != "" && !found {
				t.Errorf("expected MAX_DEPTH finding, got none. Findings: %+v", reports[0].Findings)
			}
			if tt.wantRule == "" && found {
				t.Errorf("expected no MAX_DEPTH finding, but found one")
			}

			if tt.wantScore > 0 && reports[0].RiskScore != tt.wantScore {
				t.Errorf("RiskScore = %d, want %d", reports[0].RiskScore, tt.wantScore)
			}
		})
	}
}
