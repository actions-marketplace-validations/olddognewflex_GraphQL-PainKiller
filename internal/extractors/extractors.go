package extractors

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type DocumentKind string

const (
	KindGraphQLFile             DocumentKind = "graphql-file"
	KindTaggedTemplate          DocumentKind = "tagged-template"
	KindGraphQLCommentTemplate  DocumentKind = "graphql-comment-template"
)

type Document struct {
	FilePath  string       `json:"filePath"`
	Source    string       `json:"source"`
	StartLine int          `json:"startLine"`
	EndLine   int          `json:"endLine"`
	Kind      DocumentKind `json:"kind"`
}

var gqlTaggedRegex = regexp.MustCompile(`(?s)\bgql\s*` + "`" + `(.*?)` + "`")
var graphqlTaggedRegex = regexp.MustCompile(`(?s)\bgraphql\s*` + "`" + `(.*?)` + "`")
var graphqlCommentRegex = regexp.MustCompile(`(?s)/\*\s*GraphQL\s*\*/\s*` + "`" + `(.*?)` + "`")

func Extract(target string) ([]Document, error) {
	var docs []Document

	info, err := os.Stat(target)
	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		return extractFile(target)
	}

	err = filepath.WalkDir(target, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}

		if d.IsDir() && shouldSkipDir(d.Name()) {
			return filepath.SkipDir
		}

		if d.IsDir() {
			return nil
		}

		fileDocs, err := extractFile(path)
		if err != nil {
			return err
		}

		docs = append(docs, fileDocs...)
		return nil
	})

	return docs, err
}

func extractFile(path string) ([]Document, error) {
	ext := strings.ToLower(filepath.Ext(path))

	switch ext {
	case ".graphql", ".gql":
		bytes, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}

		source := string(bytes)
		return []Document{{
			FilePath:  path,
			Source:    source,
			StartLine: 1,
			EndLine:   countLines(source),
			Kind:      KindGraphQLFile,
		}}, nil

	case ".ts", ".tsx", ".js", ".jsx":
		return extractTemplates(path)

	default:
		return nil, nil
	}
}

func extractTemplates(path string) ([]Document, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	content := string(bytes)
	var docs []Document

	patterns := []struct {
		regex *regexp.Regexp
		kind  DocumentKind
	}{
		{gqlTaggedRegex, KindTaggedTemplate},
		{graphqlTaggedRegex, KindTaggedTemplate},
		{graphqlCommentRegex, KindGraphQLCommentTemplate},
	}

	for _, pattern := range patterns {
		matches := pattern.regex.FindAllStringSubmatchIndex(content, -1)

		for _, match := range matches {
			if len(match) < 4 {
				continue
			}

			fullStart := match[0]
			fullEnd := match[1]
			sourceStart := match[2]
			sourceEnd := match[3]

			source := content[sourceStart:sourceEnd]

			// V1 intentionally skips interpolated queries. Static analysis gets weird fast.
			if strings.Contains(source, "${") {
				continue
			}

			startLine := countLines(content[:fullStart])
			full := content[fullStart:fullEnd]

			docs = append(docs, Document{
				FilePath:  path,
				Source:    source,
				StartLine: startLine,
				EndLine:   startLine + countLines(full) - 1,
				Kind:      pattern.kind,
			})
		}
	}

	return docs, nil
}

func shouldSkipDir(name string) bool {
	switch name {
	case ".git", "node_modules", "dist", "build", ".next", "coverage":
		return true
	default:
		return false
	}
}

func countLines(value string) int {
	if value == "" {
		return 1
	}
	return strings.Count(value, "\n") + 1
}
