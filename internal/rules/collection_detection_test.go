package rules

import (
	"testing"

	"github.com/olddognewflex/graphql-painkiller/internal/config"
	"github.com/olddognewflex/graphql-painkiller/internal/models"
)

func TestLooksCollectionLike(t *testing.T) {
	cfg := config.Config{
		CollectionFieldPatterns: []string{"items", "nodes", "edges"},
	}

	tests := []struct {
		name  string
		field string
		want  bool
	}{
		{name: "items is configured pattern", field: "items", want: true},
		{name: "nodes is configured pattern", field: "nodes", want: true},
		{name: "edges is configured pattern", field: "edges", want: true},
		{name: "Items is case insensitive", field: "Items", want: true},

		{name: "userList suffix", field: "userList", want: true},
		{name: "taskCollection suffix", field: "taskCollection", want: true},
		{name: "resultSet suffix", field: "resultSet", want: true},
		{name: "logEntries suffix", field: "logEntries", want: true},
		{name: "dataArray suffix", field: "dataArray", want: true},

		{name: "list alone is not collection", field: "list", want: false},
		{name: "set alone is not collection", field: "set", want: false},
		{name: "array alone is not collection", field: "array", want: false},
		{name: "collection alone is not collection", field: "collection", want: false},

		{name: "posts is plural", field: "posts", want: true},
		{name: "comments is plural", field: "comments", want: true},
		{name: "users is plural", field: "users", want: true},
		{name: "orders is plural", field: "orders", want: true},
		{name: "departments is plural", field: "departments", want: true},
		{name: "teams is plural", field: "teams", want: true},
		{name: "members is plural", field: "members", want: true},
		{name: "accounts is plural", field: "accounts", want: true},
		{name: "categories is plural", field: "categories", want: true},
		{name: "activities is plural", field: "activities", want: true},

		{name: "address is not plural", field: "address", want: false},
		{name: "access is not plural", field: "access", want: false},
		{name: "process is not plural", field: "process", want: false},
		{name: "progress is not plural", field: "progress", want: false},
		{name: "success is not plural", field: "success", want: false},
		{name: "business is not plural", field: "business", want: false},
		{name: "class is not plural", field: "class", want: false},

		{name: "status is not plural", field: "status", want: false},
		{name: "focus is not plural", field: "focus", want: false},
		{name: "bonus is not plural", field: "bonus", want: false},
		{name: "radius is not plural", field: "radius", want: false},
		{name: "campus is not plural", field: "campus", want: false},
		{name: "census is not plural", field: "census", want: false},
		{name: "nexus is not plural", field: "nexus", want: false},
		{name: "apparatus is not plural", field: "apparatus", want: false},

		{name: "basis is not plural", field: "basis", want: false},
		{name: "analysis is not plural", field: "analysis", want: false},
		{name: "diagnosis is not plural", field: "diagnosis", want: false},
		{name: "synopsis is not plural", field: "synopsis", want: false},
		{name: "thesis is not plural", field: "thesis", want: false},
		{name: "axis is not plural", field: "axis", want: false},
		{name: "crisis is not plural", field: "crisis", want: false},

		{name: "previous is not plural", field: "previous", want: false},
		{name: "various is not plural", field: "various", want: false},
		{name: "anonymous is not plural", field: "anonymous", want: false},

		{name: "alias is excluded", field: "alias", want: false},
		{name: "canvas is excluded", field: "canvas", want: false},
		{name: "metadata is excluded", field: "metadata", want: false},
		{name: "data is excluded", field: "data", want: false},
		{name: "bias is excluded", field: "bias", want: false},

		{name: "user is not collection", field: "user", want: false},
		{name: "author is not collection", field: "author", want: false},
		{name: "profile is not collection", field: "profile", want: false},
		{name: "id is not collection", field: "id", want: false},
		{name: "name is not collection", field: "name", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := models.FieldInfo{Name: tt.field}
			got := LooksCollectionLike(field, cfg)
			if got != tt.want {
				t.Errorf("LooksCollectionLike(%q) = %v, want %v", tt.field, got, tt.want)
			}
		})
	}
}
