package rules

import (
	"testing"

	"github.com/olddognewflex/graphql-painkiller/internal/config"
	"github.com/olddognewflex/graphql-painkiller/internal/extractors"
	"github.com/olddognewflex/graphql-painkiller/internal/models"
	"github.com/olddognewflex/graphql-painkiller/internal/severity"
)

func TestMaxDepth(t *testing.T) {
	cfg := config.Config{
		Rules: config.Rules{
			MaxDepth: 5,
		},
	}

	doc := extractors.Document{
		FilePath:  "test.graphql",
		StartLine: 1,
	}

	tests := []struct {
		name     string
		fields   []models.FieldInfo
		wantLen  int
		wantRule string
		wantPath string
		wantSev  severity.Severity
	}{
		{
			name: "depth within limit returns no findings",
			fields: []models.FieldInfo{
				{Name: "user", Depth: 1, Path: "user", Children: []models.FieldInfo{
					{Name: "posts", Depth: 2, Path: "user.posts", Children: []models.FieldInfo{
						{Name: "title", Depth: 3, Path: "user.posts.title"},
					}},
				}},
			},
			wantLen: 0,
		},
		{
			name: "depth exactly at limit returns no findings",
			fields: []models.FieldInfo{
				{Name: "a", Depth: 1, Path: "a", Children: []models.FieldInfo{
					{Name: "b", Depth: 2, Path: "a.b", Children: []models.FieldInfo{
						{Name: "c", Depth: 3, Path: "a.b.c", Children: []models.FieldInfo{
							{Name: "d", Depth: 4, Path: "a.b.c.d", Children: []models.FieldInfo{
								{Name: "e", Depth: 5, Path: "a.b.c.d.e"},
							}},
						}},
					}},
				}},
			},
			wantLen: 0,
		},
		{
			name: "depth exceeds limit by one returns finding",
			fields: []models.FieldInfo{
				{Name: "a", Depth: 1, Path: "a", Children: []models.FieldInfo{
					{Name: "b", Depth: 2, Path: "a.b", Children: []models.FieldInfo{
						{Name: "c", Depth: 3, Path: "a.b.c", Children: []models.FieldInfo{
							{Name: "d", Depth: 4, Path: "a.b.c.d", Children: []models.FieldInfo{
								{Name: "e", Depth: 5, Path: "a.b.c.d.e", Children: []models.FieldInfo{
									{Name: "f", Depth: 6, Path: "a.b.c.d.e.f"},
								}},
							}},
						}},
					}},
				}},
			},
			wantLen:  1,
			wantRule: "MAX_DEPTH",
			wantPath: "a.b.c.d.e.f",
			wantSev:  severity.Warning,
		},
		{
			name: "multiple branches reports deepest only",
			fields: []models.FieldInfo{
				{Name: "user", Depth: 1, Path: "user", Children: []models.FieldInfo{
					{Name: "posts", Depth: 2, Path: "user.posts", Children: []models.FieldInfo{
						{Name: "title", Depth: 3, Path: "user.posts.title"},
					}},
					{Name: "comments", Depth: 2, Path: "user.comments", Children: []models.FieldInfo{
						{Name: "author", Depth: 3, Path: "user.comments.author", Children: []models.FieldInfo{
							{Name: "profile", Depth: 4, Path: "user.comments.author.profile", Children: []models.FieldInfo{
								{Name: "settings", Depth: 5, Path: "user.comments.author.profile.settings", Children: []models.FieldInfo{
									{Name: "privacy", Depth: 6, Path: "user.comments.author.profile.settings.privacy"},
								}},
							}},
						}},
					}},
				}},
			},
			wantLen:  1,
			wantRule: "MAX_DEPTH",
			wantPath: "user.comments.author.profile.settings.privacy",
			wantSev:  severity.Warning,
		},
		{
			name:     "empty fields returns no findings",
			fields:   []models.FieldInfo{},
			wantLen:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MaxDepth(tt.fields, doc, cfg)
			if len(got) != tt.wantLen {
				t.Fatalf("MaxDepth() returned %d findings, want %d", len(got), tt.wantLen)
			}
			if tt.wantLen > 0 {
				if got[0].RuleID != tt.wantRule {
					t.Errorf("RuleID = %q, want %q", got[0].RuleID, tt.wantRule)
				}
				if got[0].Path != tt.wantPath {
					t.Errorf("Path = %q, want %q", got[0].Path, tt.wantPath)
				}
				if got[0].Severity != tt.wantSev {
					t.Errorf("Severity = %q, want %q", got[0].Severity, tt.wantSev)
				}
			}
		})
	}
}

func TestMaxDepthAdjustedLine(t *testing.T) {
	cfg := config.Config{
		Rules: config.Rules{
			MaxDepth: 3,
		},
	}

	doc := extractors.Document{
		FilePath:  "test.ts",
		StartLine: 10,
	}

	fields := []models.FieldInfo{
		{Name: "a", Depth: 1, Path: "a", Line: 3, Children: []models.FieldInfo{
			{Name: "b", Depth: 2, Path: "a.b", Line: 4, Children: []models.FieldInfo{
				{Name: "c", Depth: 3, Path: "a.b.c", Line: 5, Children: []models.FieldInfo{
					{Name: "d", Depth: 4, Path: "a.b.c.d", Line: 6},
				}},
			}},
		}},
	}

	got := MaxDepth(fields, doc, cfg)
	if len(got) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(got))
	}

	if got[0].Line != 15 {
		t.Errorf("Line = %d, want %d (start line + field line - 1)", got[0].Line, 15)
	}
}
