package github

import (
	"strings"
	"testing"

	"github.com/olddognewflex/graphql-painkiller/internal/models"
	"github.com/olddognewflex/graphql-painkiller/internal/severity"
)

func TestBuildReviewComments(t *testing.T) {
	tests := []struct {
		name   string
		reports []models.Report
		want   int
		check  func(t *testing.T, comments []ReviewComment)
	}{
		{
			name:   "empty reports",
			reports: []models.Report{},
			want:   0,
		},
		{
			name: "report with no findings",
			reports: []models.Report{
				{FilePath: "test.graphql", Findings: []models.Finding{}},
			},
			want: 0,
		},
		{
			name: "finding with valid line",
			reports: []models.Report{
				{
					FilePath: "test.graphql",
					Findings: []models.Finding{
						{
							RuleID:     "MAX_DEPTH",
							Message:    "Query exceeds max depth",
							FilePath:   "test.graphql",
							Line:       10,
							Path:       "user.posts.comments",
							Suggestion: "Reduce nesting",
							Severity:   severity.Warning,
						},
					},
				},
			},
			want: 1,
			check: func(t *testing.T, comments []ReviewComment) {
				c := comments[0]
				if c.Path != "test.graphql" {
					t.Errorf("Path = %q, want %q", c.Path, "test.graphql")
				}
				if c.Line != 10 {
					t.Errorf("Line = %d, want %d", c.Line, 10)
				}
				if !strings.Contains(c.Body, "MAX DEPTH") {
					t.Errorf("Body should contain formatted rule ID, got: %q", c.Body)
				}
				if !strings.Contains(c.Body, "Query exceeds max depth") {
					t.Errorf("Body should contain message, got: %q", c.Body)
				}
				if !strings.Contains(c.Body, "user.posts.comments") {
					t.Errorf("Body should contain path, got: %q", c.Body)
				}
				if !strings.Contains(c.Body, "Reduce nesting") {
					t.Errorf("Body should contain suggestion, got: %q", c.Body)
				}
			},
		},
		{
			name: "finding with docs url includes it in body",
			reports: []models.Report{
				{
					FilePath: "test.graphql",
					Findings: []models.Finding{
						{
							RuleID:     "MISSING_PAGINATION",
							Message:    "Missing pagination",
							FilePath:   "test.graphql",
							Line:       10,
							Path:       "posts",
							Suggestion: "Add pagination",
							Severity:   severity.High,
							DocsURL:    "https://graphql.org/learn/pagination/",
						},
					},
				},
			},
			want: 1,
			check: func(t *testing.T, comments []ReviewComment) {
				c := comments[0]
				if !strings.Contains(c.Body, "[Why this matters](https://graphql.org/learn/pagination/)") {
					t.Errorf("Body should contain docs link, got: %q", c.Body)
				}
			},
		},
		{
			name: "finding with line 0 is skipped",
			reports: []models.Report{
				{
					FilePath: "test.graphql",
					Findings: []models.Finding{
						{RuleID: "MAX_DEPTH", FilePath: "test.graphql", Line: 0, Message: "ignored"},
					},
				},
			},
			want: 0,
		},
		{
			name: "finding with negative line is skipped",
			reports: []models.Report{
				{
					FilePath: "test.graphql",
					Findings: []models.Finding{
						{RuleID: "MAX_DEPTH", FilePath: "test.graphql", Line: -1, Message: "ignored"},
					},
				},
			},
			want: 0,
		},
		{
			name: "multiple reports with multiple findings",
			reports: []models.Report{
				{
					FilePath: "a.graphql",
					Findings: []models.Finding{
						{RuleID: "MAX_DEPTH", FilePath: "a.graphql", Line: 5, Message: "deep"},
						{RuleID: "MISSING_PAGINATION", FilePath: "a.graphql", Line: 10, Message: "no pagination"},
					},
				},
				{
					FilePath: "b.graphql",
					Findings: []models.Finding{
						{RuleID: "EXPENSIVE_FIELD", FilePath: "b.graphql", Line: 3, Message: "expensive"},
					},
				},
			},
			want: 3,
			check: func(t *testing.T, comments []ReviewComment) {
				paths := make(map[string]int)
				for _, c := range comments {
					paths[c.Path]++
				}
				if paths["a.graphql"] != 2 {
					t.Errorf("Expected 2 comments for a.graphql, got %d", paths["a.graphql"])
				}
				if paths["b.graphql"] != 1 {
					t.Errorf("Expected 1 comment for b.graphql, got %d", paths["b.graphql"])
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildReviewComments(tt.reports)
			if len(got) != tt.want {
				t.Errorf("BuildReviewComments() returned %d comments, want %d", len(got), tt.want)
			}
			if tt.check != nil {
				tt.check(t, got)
			}
		})
	}
}
