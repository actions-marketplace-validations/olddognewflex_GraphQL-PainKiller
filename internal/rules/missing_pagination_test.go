package rules

import (
	"testing"

	"github.com/olddognewflex/graphql-painkiller/internal/config"
	"github.com/olddognewflex/graphql-painkiller/internal/extractors"
	"github.com/olddognewflex/graphql-painkiller/internal/models"
	"github.com/olddognewflex/graphql-painkiller/internal/severity"
)

func TestMissingPagination(t *testing.T) {
	cfg := config.Config{
		Rules: config.Rules{
			RequirePagination: true,
		},
		PaginationArgs:          []string{"first", "last", "limit", "take", "pageSize", "after", "before", "offset"},
		CollectionFieldPatterns: []string{"items", "nodes", "edges"},
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
			name: "collection without pagination produces finding",
			fields: []models.FieldInfo{
				{Name: "posts", Depth: 1, Path: "posts", Children: []models.FieldInfo{
					{Name: "id", Depth: 2, Path: "posts.id"},
					{Name: "title", Depth: 2, Path: "posts.title"},
				}},
			},
			wantLen:   1,
			wantPaths: []string{"posts"},
		},
		{
			name: "collection with first arg produces no finding",
			fields: []models.FieldInfo{
				{Name: "posts", Depth: 1, Path: "posts", Arguments: []string{"first"}, Children: []models.FieldInfo{
					{Name: "id", Depth: 2, Path: "posts.id"},
				}},
			},
			wantLen: 0,
		},
		{
			name: "collection with limit arg produces no finding",
			fields: []models.FieldInfo{
				{Name: "posts", Depth: 1, Path: "posts", Arguments: []string{"limit"}, Children: []models.FieldInfo{
					{Name: "id", Depth: 2, Path: "posts.id"},
				}},
			},
			wantLen: 0,
		},
		{
			name: "collection with pageSize arg produces no finding",
			fields: []models.FieldInfo{
				{Name: "posts", Depth: 1, Path: "posts", Arguments: []string{"pageSize"}, Children: []models.FieldInfo{
					{Name: "id", Depth: 2, Path: "posts.id"},
				}},
			},
			wantLen: 0,
		},
		{
			name: "non-collection field produces no finding",
			fields: []models.FieldInfo{
				{Name: "user", Depth: 1, Path: "user", Children: []models.FieldInfo{
					{Name: "name", Depth: 2, Path: "user.name"},
				}},
			},
			wantLen: 0,
		},
		{
			name: "leaf field produces no finding",
			fields: []models.FieldInfo{
				{Name: "posts", Depth: 1, Path: "posts"},
			},
			wantLen: 0,
		},
		{
			name: "configured pattern without pagination produces finding",
			fields: []models.FieldInfo{
				{Name: "items", Depth: 1, Path: "items", Children: []models.FieldInfo{
					{Name: "id", Depth: 2, Path: "items.id"},
				}},
			},
			wantLen:   1,
			wantPaths: []string{"items"},
		},
		{
			name: "multiple collections mixed pagination",
			fields: []models.FieldInfo{
				{Name: "posts", Depth: 1, Path: "posts", Arguments: []string{"first"}, Children: []models.FieldInfo{
					{Name: "comments", Depth: 2, Path: "posts.comments", Children: []models.FieldInfo{
						{Name: "body", Depth: 3, Path: "posts.comments.body"},
					}},
				}},
			},
			wantLen:   1,
			wantPaths: []string{"posts.comments"},
		},
		{
			name: "pagination arg matching is case insensitive",
			fields: []models.FieldInfo{
				{Name: "posts", Depth: 1, Path: "posts", Arguments: []string{"First"}, Children: []models.FieldInfo{
					{Name: "id", Depth: 2, Path: "posts.id"},
				}},
			},
			wantLen: 0,
		},
		{
			name:    "empty fields produces no findings",
			fields:  []models.FieldInfo{},
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MissingPagination(tt.fields, doc, cfg)
			if len(got) != tt.wantLen {
				t.Fatalf("MissingPagination() returned %d findings, want %d. Findings: %+v", len(got), tt.wantLen, got)
			}
			for i, path := range tt.wantPaths {
				if got[i].Path != path {
					t.Errorf("finding[%d].Path = %q, want %q", i, got[i].Path, path)
				}
				if got[i].RuleID != "MISSING_PAGINATION" {
					t.Errorf("finding[%d].RuleID = %q, want %q", i, got[i].RuleID, "MISSING_PAGINATION")
				}
				if got[i].Severity != severity.High {
					t.Errorf("finding[%d].Severity = %q, want %q", i, got[i].Severity, severity.High)
				}
			}
		})
	}
}

func TestMissingPaginationDisabled(t *testing.T) {
	cfg := config.Config{
		Rules: config.Rules{
			RequirePagination: false,
		},
		PaginationArgs:          []string{"first", "last", "limit"},
		CollectionFieldPatterns: []string{"items"},
	}

	doc := extractors.Document{
		FilePath:  "test.graphql",
		StartLine: 1,
	}

	fields := []models.FieldInfo{
		{Name: "posts", Depth: 1, Path: "posts", Children: []models.FieldInfo{
			{Name: "id", Depth: 2, Path: "posts.id"},
		}},
	}

	got := MissingPagination(fields, doc, cfg)
	if len(got) != 0 {
		t.Fatalf("MissingPagination() with requirePagination=false returned %d findings, want 0", len(got))
	}
}

func TestMissingPaginationAdjustedLine(t *testing.T) {
	cfg := config.Config{
		Rules: config.Rules{
			RequirePagination: true,
		},
		PaginationArgs:          []string{"first"},
		CollectionFieldPatterns: []string{"items"},
	}

	doc := extractors.Document{
		FilePath:  "test.ts",
		StartLine: 10,
	}

	fields := []models.FieldInfo{
		{Name: "posts", Depth: 1, Path: "posts", Line: 5, Children: []models.FieldInfo{
			{Name: "id", Depth: 2, Path: "posts.id", Line: 6},
		}},
	}

	got := MissingPagination(fields, doc, cfg)
	if len(got) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(got))
	}

	wantLine := 14
	if got[0].Line != wantLine {
		t.Errorf("Line = %d, want %d", got[0].Line, wantLine)
	}
}
