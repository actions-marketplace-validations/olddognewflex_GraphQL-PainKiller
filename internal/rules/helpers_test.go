package rules

import (
	"testing"

	"github.com/olddognewflex/graphql-painkiller/internal/config"
	"github.com/olddognewflex/graphql-painkiller/internal/models"
)

func TestHasPagination(t *testing.T) {
	cfg := config.Config{
		PaginationArgs: []string{"first", "last", "limit", "take", "pageSize", "after", "before", "offset"},
	}

	tests := []struct {
		name string
		args []string
		want bool
	}{
		{name: "first arg matches", args: []string{"first"}, want: true},
		{name: "limit arg matches", args: []string{"limit"}, want: true},
		{name: "pageSize arg matches", args: []string{"pageSize"}, want: true},
		{name: "after arg matches", args: []string{"after"}, want: true},
		{name: "case insensitive match", args: []string{"First"}, want: true},
		{name: "case insensitive LIMIT", args: []string{"LIMIT"}, want: true},
		{name: "non-pagination arg does not match", args: []string{"where"}, want: false},
		{name: "empty args does not match", args: []string{}, want: false},
		{name: "nil args does not match", args: nil, want: false},
		{name: "multiple args one matches", args: []string{"where", "orderBy", "first"}, want: true},
		{name: "multiple args none match", args: []string{"where", "orderBy", "filter"}, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := models.FieldInfo{Name: "test", Arguments: tt.args}
			got := HasPagination(field, cfg)
			if got != tt.want {
				t.Errorf("HasPagination(args=%v) = %v, want %v", tt.args, got, tt.want)
			}
		})
	}
}

func TestHasPaginationEmptyConfig(t *testing.T) {
	cfg := config.Config{
		PaginationArgs: []string{},
	}

	field := models.FieldInfo{Name: "test", Arguments: []string{"first"}}
	got := HasPagination(field, cfg)
	if got {
		t.Errorf("HasPagination() with empty config should return false")
	}
}

func TestFlatten(t *testing.T) {
	tests := []struct {
		name      string
		fields    []models.FieldInfo
		wantNames []string
	}{
		{
			name:      "empty input",
			fields:    []models.FieldInfo{},
			wantNames: nil,
		},
		{
			name: "flat list returns same order",
			fields: []models.FieldInfo{
				{Name: "a"},
				{Name: "b"},
				{Name: "c"},
			},
			wantNames: []string{"a", "b", "c"},
		},
		{
			name: "nested tree flattens depth first",
			fields: []models.FieldInfo{
				{Name: "a", Children: []models.FieldInfo{
					{Name: "b", Children: []models.FieldInfo{
						{Name: "c"},
					}},
					{Name: "d"},
				}},
			},
			wantNames: []string{"a", "b", "c", "d"},
		},
		{
			name: "multiple roots with children",
			fields: []models.FieldInfo{
				{Name: "x", Children: []models.FieldInfo{
					{Name: "y"},
				}},
				{Name: "z"},
			},
			wantNames: []string{"x", "y", "z"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Flatten(tt.fields)
			if len(got) != len(tt.wantNames) {
				t.Fatalf("Flatten() returned %d items, want %d", len(got), len(tt.wantNames))
			}
			for i, name := range tt.wantNames {
				if got[i].Name != name {
					t.Errorf("Flatten()[%d].Name = %q, want %q", i, got[i].Name, name)
				}
			}
		})
	}
}

func TestAdjustedLine(t *testing.T) {
	tests := []struct {
		name         string
		docStartLine int
		fieldLine    int
		want         int
	}{
		{name: "positive field line", docStartLine: 10, fieldLine: 5, want: 14},
		{name: "field line 1", docStartLine: 10, fieldLine: 1, want: 10},
		{name: "zero field line returns doc start", docStartLine: 10, fieldLine: 0, want: 10},
		{name: "negative field line returns doc start", docStartLine: 10, fieldLine: -1, want: 10},
		{name: "doc start 1 field line 1", docStartLine: 1, fieldLine: 1, want: 1},
		{name: "doc start 1 field line 3", docStartLine: 1, fieldLine: 3, want: 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AdjustedLine(tt.docStartLine, tt.fieldLine)
			if got != tt.want {
				t.Errorf("AdjustedLine(%d, %d) = %d, want %d", tt.docStartLine, tt.fieldLine, got, tt.want)
			}
		})
	}
}
