package rules

import (
	"testing"

	"github.com/olddognewflex/graphql-painkiller/internal/config"
	"github.com/olddognewflex/graphql-painkiller/internal/extractors"
	"github.com/olddognewflex/graphql-painkiller/internal/models"
	"github.com/olddognewflex/graphql-painkiller/internal/severity"
)

func TestExpensiveFields(t *testing.T) {
	cfg := config.Config{
		ExpensiveFieldPatterns: []string{"comments", "history", "events", "logs", "charges", "payments", "inspections", "accounts"},
	}

	doc := extractors.Document{
		FilePath:  "test.graphql",
		StartLine: 1,
	}

	tests := []struct {
		name      string
		fields    []models.FieldInfo
		wantLen   int
		wantPaths []string
	}{
		{
			name: "exact pattern match produces finding",
			fields: []models.FieldInfo{
				{Name: "comments", Depth: 1, Path: "comments"},
			},
			wantLen:   1,
			wantPaths: []string{"comments"},
		},
		{
			name: "substring match produces finding",
			fields: []models.FieldInfo{
				{Name: "dealCharges", Depth: 1, Path: "dealCharges"},
			},
			wantLen:   1,
			wantPaths: []string{"dealCharges"},
		},
		{
			name: "case insensitive match",
			fields: []models.FieldInfo{
				{Name: "Comments", Depth: 1, Path: "Comments"},
			},
			wantLen:   1,
			wantPaths: []string{"Comments"},
		},
		{
			name: "non-matching field produces no finding",
			fields: []models.FieldInfo{
				{Name: "title", Depth: 1, Path: "title"},
			},
			wantLen: 0,
		},
		{
			name: "nested field matching pattern",
			fields: []models.FieldInfo{
				{Name: "posts", Depth: 1, Path: "posts", Children: []models.FieldInfo{
					{Name: "comments", Depth: 2, Path: "posts.comments"},
				}},
			},
			wantLen:   1,
			wantPaths: []string{"posts.comments"},
		},
		{
			name: "multiple matches across different fields",
			fields: []models.FieldInfo{
				{Name: "comments", Depth: 1, Path: "comments"},
				{Name: "payments", Depth: 1, Path: "payments"},
			},
			wantLen:   2,
			wantPaths: []string{"comments", "payments"},
		},
		{
			name: "field matches only first pattern then breaks",
			fields: []models.FieldInfo{
				{Name: "accountsHistory", Depth: 1, Path: "accountsHistory"},
			},
			wantLen:   1,
			wantPaths: []string{"accountsHistory"},
		},
		{
			name:    "empty fields produces no findings",
			fields:  []models.FieldInfo{},
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExpensiveFields(tt.fields, doc, cfg)
			if len(got) != tt.wantLen {
				t.Fatalf("ExpensiveFields() returned %d findings, want %d. Findings: %+v", len(got), tt.wantLen, got)
			}
			for i, path := range tt.wantPaths {
				if got[i].Path != path {
					t.Errorf("finding[%d].Path = %q, want %q", i, got[i].Path, path)
				}
				if got[i].RuleID != "EXPENSIVE_FIELD_PATTERN" {
					t.Errorf("finding[%d].RuleID = %q, want %q", i, got[i].RuleID, "EXPENSIVE_FIELD_PATTERN")
				}
				if got[i].Severity != severity.Warning {
					t.Errorf("finding[%d].Severity = %q, want %q", i, got[i].Severity, severity.Warning)
				}
				if got[i].ScoreImpact != 1 {
					t.Errorf("finding[%d].ScoreImpact = %d, want 1", i, got[i].ScoreImpact)
				}
			}
		})
	}
}

func TestExpensiveFieldsEmptyConfig(t *testing.T) {
	cfg := config.Config{
		ExpensiveFieldPatterns: []string{},
	}

	doc := extractors.Document{
		FilePath:  "test.graphql",
		StartLine: 1,
	}

	fields := []models.FieldInfo{
		{Name: "comments", Depth: 1, Path: "comments"},
	}

	got := ExpensiveFields(fields, doc, cfg)
	if len(got) != 0 {
		t.Fatalf("ExpensiveFields() with empty config returned %d findings, want 0", len(got))
	}
}

func TestExpensiveFieldsAdjustedLine(t *testing.T) {
	cfg := config.Config{
		ExpensiveFieldPatterns: []string{"comments"},
	}

	doc := extractors.Document{
		FilePath:  "test.ts",
		StartLine: 10,
	}

	fields := []models.FieldInfo{
		{Name: "comments", Depth: 1, Path: "comments", Line: 4},
	}

	got := ExpensiveFields(fields, doc, cfg)
	if len(got) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(got))
	}

	wantLine := 13
	if got[0].Line != wantLine {
		t.Errorf("Line = %d, want %d", got[0].Line, wantLine)
	}
}
