package rules

import (
	"strings"
	"testing"

	"github.com/olddognewflex/graphql-painkiller/internal/config"
	"github.com/olddognewflex/graphql-painkiller/internal/extractors"
	"github.com/olddognewflex/graphql-painkiller/internal/models"
	"github.com/olddognewflex/graphql-painkiller/internal/severity"
)

func TestKnownResolvers(t *testing.T) {
	cfg := config.Config{
		KnownResolvers: map[string]config.KnownResolver{
			"posts.comments": {
				Risk:   severity.High,
				Reason: "comments resolve per post",
			},
			"posts.comments.author": {
				Risk:    severity.Critical,
				Reason:  "author fans out per comment",
				Service: "user-service",
			},
			"orders.items": {
				Risk:   severity.Warning,
				Reason: "items are embedded but still queried",
			},
		},
	}

	doc := extractors.Document{
		FilePath:  "test.graphql",
		StartLine: 1,
	}

	tests := []struct {
		name        string
		fields      []models.FieldInfo
		wantLen     int
		wantPaths   []string
		wantSevs    []severity.Severity
		wantImpacts []int
	}{
		{
			name: "matching high risk path",
			fields: []models.FieldInfo{
				{Name: "posts", Depth: 1, Path: "posts", Children: []models.FieldInfo{
					{Name: "comments", Depth: 2, Path: "posts.comments"},
				}},
			},
			wantLen:     1,
			wantPaths:   []string{"posts.comments"},
			wantSevs:    []severity.Severity{severity.High},
			wantImpacts: []int{3},
		},
		{
			name: "matching critical risk path",
			fields: []models.FieldInfo{
				{Name: "posts", Depth: 1, Path: "posts", Children: []models.FieldInfo{
					{Name: "comments", Depth: 2, Path: "posts.comments", Children: []models.FieldInfo{
						{Name: "author", Depth: 3, Path: "posts.comments.author"},
					}},
				}},
			},
			wantLen:     2,
			wantPaths:   []string{"posts.comments", "posts.comments.author"},
			wantSevs:    []severity.Severity{severity.High, severity.Critical},
			wantImpacts: []int{3, 4},
		},
		{
			name: "matching warning risk path",
			fields: []models.FieldInfo{
				{Name: "orders", Depth: 1, Path: "orders", Children: []models.FieldInfo{
					{Name: "items", Depth: 2, Path: "orders.items"},
				}},
			},
			wantLen:     1,
			wantPaths:   []string{"orders.items"},
			wantSevs:    []severity.Severity{severity.Warning},
			wantImpacts: []int{2},
		},
		{
			name: "non-matching path produces no finding",
			fields: []models.FieldInfo{
				{Name: "users", Depth: 1, Path: "users", Children: []models.FieldInfo{
					{Name: "profile", Depth: 2, Path: "users.profile"},
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
			got := KnownResolvers(tt.fields, doc, cfg)
			if len(got) != tt.wantLen {
				t.Fatalf("KnownResolvers() returned %d findings, want %d. Findings: %+v", len(got), tt.wantLen, got)
			}
			for i := 0; i < tt.wantLen; i++ {
				if got[i].Path != tt.wantPaths[i] {
					t.Errorf("finding[%d].Path = %q, want %q", i, got[i].Path, tt.wantPaths[i])
				}
				if got[i].RuleID != "KNOWN_RESOLVER_RISK" {
					t.Errorf("finding[%d].RuleID = %q, want %q", i, got[i].RuleID, "KNOWN_RESOLVER_RISK")
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

func TestKnownResolversServiceField(t *testing.T) {
	doc := extractors.Document{
		FilePath:  "test.graphql",
		StartLine: 1,
	}

	t.Run("service included in message when present", func(t *testing.T) {
		cfg := config.Config{
			KnownResolvers: map[string]config.KnownResolver{
				"posts.comments.author": {
					Risk:    severity.Critical,
					Reason:  "fans out per comment",
					Service: "user-service",
				},
			},
		}

		fields := []models.FieldInfo{
			{Name: "posts", Depth: 1, Path: "posts", Children: []models.FieldInfo{
				{Name: "comments", Depth: 2, Path: "posts.comments", Children: []models.FieldInfo{
					{Name: "author", Depth: 3, Path: "posts.comments.author"},
				}},
			}},
		}

		got := KnownResolvers(fields, doc, cfg)
		if len(got) != 1 {
			t.Fatalf("expected 1 finding, got %d", len(got))
		}
		if !strings.Contains(got[0].Message, "user-service") {
			t.Errorf("Message = %q, expected it to contain 'user-service'", got[0].Message)
		}
	})

	t.Run("service omitted from message when empty", func(t *testing.T) {
		cfg := config.Config{
			KnownResolvers: map[string]config.KnownResolver{
				"posts.comments": {
					Risk:   severity.High,
					Reason: "comments resolve per post",
				},
			},
		}

		fields := []models.FieldInfo{
			{Name: "posts", Depth: 1, Path: "posts", Children: []models.FieldInfo{
				{Name: "comments", Depth: 2, Path: "posts.comments"},
			}},
		}

		got := KnownResolvers(fields, doc, cfg)
		if len(got) != 1 {
			t.Fatalf("expected 1 finding, got %d", len(got))
		}
		if strings.Contains(got[0].Message, "Service:") {
			t.Errorf("Message = %q, expected it NOT to contain 'Service:'", got[0].Message)
		}
	})
}

func TestKnownResolversEmptyConfig(t *testing.T) {
	cfg := config.Config{
		KnownResolvers: map[string]config.KnownResolver{},
	}

	doc := extractors.Document{
		FilePath:  "test.graphql",
		StartLine: 1,
	}

	fields := []models.FieldInfo{
		{Name: "posts", Depth: 1, Path: "posts", Children: []models.FieldInfo{
			{Name: "comments", Depth: 2, Path: "posts.comments"},
		}},
	}

	got := KnownResolvers(fields, doc, cfg)
	if len(got) != 0 {
		t.Fatalf("KnownResolvers() with empty config returned %d findings, want 0", len(got))
	}
}

func TestKnownResolversAdjustedLine(t *testing.T) {
	cfg := config.Config{
		KnownResolvers: map[string]config.KnownResolver{
			"posts.comments": {
				Risk:   severity.High,
				Reason: "resolves per post",
			},
		},
	}

	doc := extractors.Document{
		FilePath:  "test.ts",
		StartLine: 10,
	}

	fields := []models.FieldInfo{
		{Name: "posts", Depth: 1, Path: "posts", Line: 2, Children: []models.FieldInfo{
			{Name: "comments", Depth: 2, Path: "posts.comments", Line: 3},
		}},
	}

	got := KnownResolvers(fields, doc, cfg)
	if len(got) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(got))
	}

	wantLine := 12
	if got[0].Line != wantLine {
		t.Errorf("Line = %d, want %d", got[0].Line, wantLine)
	}
}
