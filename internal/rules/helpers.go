package rules

import (
	"strings"

	"github.com/olddognewflex/graphql-painkiller/internal/config"
	"github.com/olddognewflex/graphql-painkiller/internal/models"
)

func Flatten(fields []models.FieldInfo) []models.FieldInfo {
	var out []models.FieldInfo
	for _, field := range fields {
		out = append(out, field)
		out = append(out, Flatten(field.Children)...)
	}
	return out
}

// collectionSuffixes are name suffixes that strongly indicate a collection field.
var collectionSuffixes = []string{
	"list",
	"collection",
	"array",
	"set",
	"entries",
}

// nonPluralSuffixes are word endings where a trailing "s" does not indicate plural.
var nonPluralSuffixes = []string{
	"ss",
	"us",
	"is",
	"ous",
}

// nonPluralExact are specific field names ending in "s" that are not collections.
var nonPluralExact = map[string]bool{
	"alias":     true,
	"canvas":    true,
	"bias":      true,
	"atlas":     true,
	"gas":       true,
	"has":       true,
	"was":       true,
	"metadata":  true,
	"data":      true,
	"graphqlws": true,
}

func LooksCollectionLike(field models.FieldInfo, cfg config.Config) bool {
	lower := strings.ToLower(field.Name)

	// 1. Exact match against configured collection patterns.
	for _, pattern := range cfg.CollectionFieldPatterns {
		if lower == strings.ToLower(pattern) {
			return true
		}
	}

	// 2. Suffix match: names ending with "List", "Collection", etc.
	for _, suffix := range collectionSuffixes {
		if strings.HasSuffix(lower, suffix) && lower != suffix {
			return true
		}
	}

	// 3. Plural heuristic: ends with "s" but not a known non-plural pattern.
	if strings.HasSuffix(lower, "s") {
		if nonPluralExact[lower] {
			return false
		}
		for _, suffix := range nonPluralSuffixes {
			if strings.HasSuffix(lower, suffix) {
				return false
			}
		}
		return true
	}

	return false
}

func HasPagination(field models.FieldInfo, cfg config.Config) bool {
	args := map[string]bool{}
	for _, arg := range field.Arguments {
		args[strings.ToLower(arg)] = true
	}

	for _, arg := range cfg.PaginationArgs {
		if args[strings.ToLower(arg)] {
			return true
		}
	}

	return false
}

func AdjustedLine(docStartLine int, fieldLine int) int {
	if fieldLine <= 0 {
		return docStartLine
	}
	return docStartLine + fieldLine - 1
}
