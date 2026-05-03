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

func LooksCollectionLike(field models.FieldInfo, cfg config.Config) bool {
	lower := strings.ToLower(field.Name)

	for _, pattern := range cfg.CollectionFieldPatterns {
		if lower == strings.ToLower(pattern) {
			return true
		}
	}

	if strings.HasSuffix(lower, "s") && lower != "status" {
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
