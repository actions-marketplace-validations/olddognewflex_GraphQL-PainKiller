package rules

import (
	"testing"

	"github.com/olddognewflex/graphql-painkiller/internal/config"
	"github.com/olddognewflex/graphql-painkiller/internal/extractors"
	"github.com/olddognewflex/graphql-painkiller/internal/models"
	"github.com/olddognewflex/graphql-painkiller/internal/severity"
)

func TestNestedCollection(t *testing.T) {
	cfg := config.Config{
		CollectionFieldPatterns: []string{"items", "nodes", "edges"},
	}

	doc := extractors.Document{
		FilePath:  "test.graphql",
		StartLine: 1,
	}

	tests := []struct {
		name         string
		fields       []models.FieldInfo
		wantLen      int
		wantRuleIDs  []string
		wantPaths    []string
		wantSevs     []severity.Severity
		wantImpacts  []int
	}{
		{
			name: "collection nested under collection produces HIGH finding",
			fields: []models.FieldInfo{
				{Name: "posts", Depth: 1, Path: "posts", Children: []models.FieldInfo{
					{Name: "comments", Depth: 2, Path: "posts.comments", Children: []models.FieldInfo{
						{Name: "body", Depth: 3, Path: "posts.comments.body"},
					}},
				}},
			},
			wantLen:     1,
			wantRuleIDs: []string{"NESTED_COLLECTION_N_PLUS_ONE"},
			wantPaths:   []string{"posts.comments"},
			wantSevs:    []severity.Severity{severity.High},
			wantImpacts: []int{3},
		},
		{
			name: "object nested under collection produces WARNING finding",
			fields: []models.FieldInfo{
				{Name: "posts", Depth: 1, Path: "posts", Children: []models.FieldInfo{
					{Name: "author", Depth: 2, Path: "posts.author", Children: []models.FieldInfo{
						{Name: "name", Depth: 3, Path: "posts.author.name"},
					}},
				}},
			},
			wantLen:     1,
			wantRuleIDs: []string{"NESTED_OBJECT_UNDER_COLLECTION"},
			wantPaths:   []string{"posts.author"},
			wantSevs:    []severity.Severity{severity.Warning},
			wantImpacts: []int{1},
		},
		{
			name: "leaf fields under collection produce no findings",
			fields: []models.FieldInfo{
				{Name: "posts", Depth: 1, Path: "posts", Children: []models.FieldInfo{
					{Name: "id", Depth: 2, Path: "posts.id"},
					{Name: "title", Depth: 2, Path: "posts.title"},
				}},
			},
			wantLen: 0,
		},
		{
			name: "non-collection parent produces no findings",
			fields: []models.FieldInfo{
				{Name: "user", Depth: 1, Path: "user", Children: []models.FieldInfo{
					{Name: "profile", Depth: 2, Path: "user.profile", Children: []models.FieldInfo{
						{Name: "bio", Depth: 3, Path: "user.profile.bio"},
					}},
				}},
			},
			wantLen: 0,
		},
		{
			name: "deeply nested collection chain produces multiple findings",
			fields: []models.FieldInfo{
				{Name: "posts", Depth: 1, Path: "posts", Children: []models.FieldInfo{
					{Name: "comments", Depth: 2, Path: "posts.comments", Children: []models.FieldInfo{
						{Name: "author", Depth: 3, Path: "posts.comments.author", Children: []models.FieldInfo{
							{Name: "name", Depth: 4, Path: "posts.comments.author.name"},
						}},
					}},
				}},
			},
			wantLen:     2,
			wantRuleIDs: []string{"NESTED_COLLECTION_N_PLUS_ONE", "NESTED_OBJECT_UNDER_COLLECTION"},
			wantPaths:   []string{"posts.comments", "posts.comments.author"},
			wantSevs:    []severity.Severity{severity.High, severity.Warning},
			wantImpacts: []int{3, 1},
		},
		{
			name: "configured pattern 'items' as child triggers HIGH",
			fields: []models.FieldInfo{
				{Name: "orders", Depth: 1, Path: "orders", Children: []models.FieldInfo{
					{Name: "items", Depth: 2, Path: "orders.items", Children: []models.FieldInfo{
						{Name: "sku", Depth: 3, Path: "orders.items.sku"},
					}},
				}},
			},
			wantLen:     1,
			wantRuleIDs: []string{"NESTED_COLLECTION_N_PLUS_ONE"},
			wantPaths:   []string{"orders.items"},
			wantSevs:    []severity.Severity{severity.High},
			wantImpacts: []int{3},
		},
		{
			name:    "empty fields produce no findings",
			fields:  []models.FieldInfo{},
			wantLen: 0,
		},
		{
			name: "collection with mixed children",
			fields: []models.FieldInfo{
				{Name: "users", Depth: 1, Path: "users", Children: []models.FieldInfo{
					{Name: "id", Depth: 2, Path: "users.id"},
					{Name: "posts", Depth: 2, Path: "users.posts", Children: []models.FieldInfo{
						{Name: "title", Depth: 3, Path: "users.posts.title"},
					}},
					{Name: "profile", Depth: 2, Path: "users.profile", Children: []models.FieldInfo{
						{Name: "avatar", Depth: 3, Path: "users.profile.avatar"},
					}},
				}},
			},
			wantLen:     2,
			wantRuleIDs: []string{"NESTED_COLLECTION_N_PLUS_ONE", "NESTED_OBJECT_UNDER_COLLECTION"},
			wantPaths:   []string{"users.posts", "users.profile"},
			wantSevs:    []severity.Severity{severity.High, severity.Warning},
			wantImpacts: []int{3, 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NestedCollection(tt.fields, doc, cfg)
			if len(got) != tt.wantLen {
				t.Fatalf("NestedCollection() returned %d findings, want %d. Findings: %+v", len(got), tt.wantLen, got)
			}

			for i := 0; i < tt.wantLen; i++ {
				if got[i].RuleID != tt.wantRuleIDs[i] {
					t.Errorf("finding[%d].RuleID = %q, want %q", i, got[i].RuleID, tt.wantRuleIDs[i])
				}
				if got[i].Path != tt.wantPaths[i] {
					t.Errorf("finding[%d].Path = %q, want %q", i, got[i].Path, tt.wantPaths[i])
				}
				if got[i].Severity != tt.wantSevs[i] {
					t.Errorf("finding[%d].Severity = %q, want %q", i, got[i].Severity, tt.wantSevs[i])
				}
				if got[i].ScoreImpact != tt.wantImpacts[i] {
					t.Errorf("finding[%d].ScoreImpact = %d, want %d", i, got[i].ScoreImpact, tt.wantImpacts[i])
				}
			}
		})
	}
}

func TestNestedCollectionAdjustedLine(t *testing.T) {
	cfg := config.Config{
		CollectionFieldPatterns: []string{"items", "nodes", "edges"},
	}

	doc := extractors.Document{
		FilePath:  "test.ts",
		StartLine: 10,
	}

	fields := []models.FieldInfo{
		{Name: "posts", Depth: 1, Path: "posts", Line: 2, Children: []models.FieldInfo{
			{Name: "comments", Depth: 2, Path: "posts.comments", Line: 3, Children: []models.FieldInfo{
				{Name: "body", Depth: 3, Path: "posts.comments.body", Line: 4},
			}},
		}},
	}

	got := NestedCollection(fields, doc, cfg)
	if len(got) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(got))
	}

	wantLine := 12
	if got[0].Line != wantLine {
		t.Errorf("Line = %d, want %d (startLine %d + fieldLine %d - 1)", got[0].Line, wantLine, doc.StartLine, 3)
	}
}
