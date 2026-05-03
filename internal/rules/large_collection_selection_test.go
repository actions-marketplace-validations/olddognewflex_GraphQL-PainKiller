package rules

import (
	"testing"

	"github.com/olddognewflex/graphql-painkiller/internal/config"
	"github.com/olddognewflex/graphql-painkiller/internal/extractors"
	"github.com/olddognewflex/graphql-painkiller/internal/models"
	"github.com/olddognewflex/graphql-painkiller/internal/severity"
)

func TestLargeCollectionSelection(t *testing.T) {
	cfg := config.Config{
		Rules: config.Rules{
			MaxCollectionSelectionFields: 3,
		},
		CollectionFieldPatterns: []string{"items", "nodes", "edges"},
	}

	doc := extractors.Document{
		FilePath:  "test.graphql",
		StartLine: 1,
	}

	tests := []struct {
		name     string
		fields   []models.FieldInfo
		wantLen  int
		wantPath string
	}{
		{
			name: "collection under limit produces no finding",
			fields: []models.FieldInfo{
				{Name: "posts", Depth: 1, Path: "posts", Children: []models.FieldInfo{
					{Name: "id", Depth: 2, Path: "posts.id"},
					{Name: "title", Depth: 2, Path: "posts.title"},
				}},
			},
			wantLen: 0,
		},
		{
			name: "collection at exact limit produces no finding",
			fields: []models.FieldInfo{
				{Name: "posts", Depth: 1, Path: "posts", Children: []models.FieldInfo{
					{Name: "id", Depth: 2, Path: "posts.id"},
					{Name: "title", Depth: 2, Path: "posts.title"},
					{Name: "body", Depth: 2, Path: "posts.body"},
				}},
			},
			wantLen: 0,
		},
		{
			name: "collection over limit produces finding",
			fields: []models.FieldInfo{
				{Name: "posts", Depth: 1, Path: "posts", Children: []models.FieldInfo{
					{Name: "id", Depth: 2, Path: "posts.id"},
					{Name: "title", Depth: 2, Path: "posts.title"},
					{Name: "body", Depth: 2, Path: "posts.body"},
					{Name: "status", Depth: 2, Path: "posts.status"},
				}},
			},
			wantLen:  1,
			wantPath: "posts",
		},
		{
			name: "non-collection field over limit produces no finding",
			fields: []models.FieldInfo{
				{Name: "user", Depth: 1, Path: "user", Children: []models.FieldInfo{
					{Name: "id", Depth: 2, Path: "user.id"},
					{Name: "name", Depth: 2, Path: "user.name"},
					{Name: "email", Depth: 2, Path: "user.email"},
					{Name: "bio", Depth: 2, Path: "user.bio"},
				}},
			},
			wantLen: 0,
		},
		{
			name: "configured pattern over limit produces finding",
			fields: []models.FieldInfo{
				{Name: "items", Depth: 1, Path: "items", Children: []models.FieldInfo{
					{Name: "id", Depth: 2, Path: "items.id"},
					{Name: "sku", Depth: 2, Path: "items.sku"},
					{Name: "price", Depth: 2, Path: "items.price"},
					{Name: "qty", Depth: 2, Path: "items.qty"},
				}},
			},
			wantLen:  1,
			wantPath: "items",
		},
		{
			name: "nested collection over limit",
			fields: []models.FieldInfo{
				{Name: "user", Depth: 1, Path: "user", Children: []models.FieldInfo{
					{Name: "orders", Depth: 2, Path: "user.orders", Children: []models.FieldInfo{
						{Name: "id", Depth: 3, Path: "user.orders.id"},
						{Name: "total", Depth: 3, Path: "user.orders.total"},
						{Name: "status", Depth: 3, Path: "user.orders.status"},
						{Name: "date", Depth: 3, Path: "user.orders.date"},
					}},
				}},
			},
			wantLen:  1,
			wantPath: "user.orders",
		},
		{
			name:    "empty fields produces no findings",
			fields:  []models.FieldInfo{},
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LargeCollectionSelection(tt.fields, doc, cfg)
			if len(got) != tt.wantLen {
				t.Fatalf("LargeCollectionSelection() returned %d findings, want %d. Findings: %+v", len(got), tt.wantLen, got)
			}
			if tt.wantLen > 0 {
				if got[0].Path != tt.wantPath {
					t.Errorf("Path = %q, want %q", got[0].Path, tt.wantPath)
				}
				if got[0].RuleID != "LARGE_COLLECTION_SELECTION" {
					t.Errorf("RuleID = %q, want %q", got[0].RuleID, "LARGE_COLLECTION_SELECTION")
				}
				if got[0].Severity != severity.Warning {
					t.Errorf("Severity = %q, want %q", got[0].Severity, severity.Warning)
				}
				if got[0].ScoreImpact != 2 {
					t.Errorf("ScoreImpact = %d, want 2", got[0].ScoreImpact)
				}
			}
		})
	}
}

func TestLargeCollectionSelectionAdjustedLine(t *testing.T) {
	cfg := config.Config{
		Rules: config.Rules{
			MaxCollectionSelectionFields: 2,
		},
		CollectionFieldPatterns: []string{"items"},
	}

	doc := extractors.Document{
		FilePath:  "test.ts",
		StartLine: 10,
	}

	fields := []models.FieldInfo{
		{Name: "posts", Depth: 1, Path: "posts", Line: 3, Children: []models.FieldInfo{
			{Name: "id", Depth: 2, Path: "posts.id", Line: 4},
			{Name: "title", Depth: 2, Path: "posts.title", Line: 5},
			{Name: "body", Depth: 2, Path: "posts.body", Line: 6},
		}},
	}

	got := LargeCollectionSelection(fields, doc, cfg)
	if len(got) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(got))
	}

	wantLine := 12
	if got[0].Line != wantLine {
		t.Errorf("Line = %d, want %d", got[0].Line, wantLine)
	}
}
