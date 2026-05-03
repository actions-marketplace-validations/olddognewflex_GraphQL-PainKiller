package analyzer

import (
	"fmt"

	"github.com/olddognewflex/graphql-painkiller/internal/config"
	"github.com/olddognewflex/graphql-painkiller/internal/extractors"
	"github.com/olddognewflex/graphql-painkiller/internal/models"
	"github.com/olddognewflex/graphql-painkiller/internal/rules"
	"github.com/olddognewflex/graphql-painkiller/internal/severity"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/parser"
)

func AnalyzeDocument(doc extractors.Document, cfg config.Config) ([]models.Report, error) {
	query, err := parser.ParseQuery(&ast.Source{
		Name:  doc.FilePath,
		Input: doc.Source,
	})
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", doc.FilePath, err)
	}

	var reports []models.Report

	for _, op := range query.Operations {
		fields := BuildFieldTree(op.SelectionSet, "", 1)

		var findings []models.Finding
		findings = append(findings, rules.MaxDepth(fields, doc, cfg)...)
		findings = append(findings, rules.MissingPagination(fields, doc, cfg)...)
		findings = append(findings, rules.NestedCollection(fields, doc, cfg)...)
		findings = append(findings, rules.ExpensiveFields(fields, doc, cfg)...)
		findings = append(findings, rules.LargeCollectionSelection(fields, doc, cfg)...)
		findings = append(findings, rules.KnownResolvers(fields, doc, cfg)...)

		score := scoreRisk(findings)

		reports = append(reports, models.Report{
			FilePath:       doc.FilePath,
			OperationName: operationName(op),
			RiskScore:     score,
			Severity:      severity.FromScore(score),
			Findings:      findings,
		})
	}

	return reports, nil
}

func BuildFieldTree(selectionSet ast.SelectionSet, parentPath string, depth int) []models.FieldInfo {
	var fields []models.FieldInfo

	for _, selection := range selectionSet {
		field, ok := selection.(*ast.Field)
		if !ok {
			continue
		}

		path := field.Name
		if parentPath != "" {
			path = parentPath + "." + field.Name
		}

		info := models.FieldInfo{
			Name:      field.Name,
			Path:      path,
			Depth:     depth,
			Line:      lineFromPosition(field.Position),
			Arguments: argumentNames(field.Arguments),
			Children:  BuildFieldTree(field.SelectionSet, path, depth+1),
		}

		fields = append(fields, info)
	}

	return fields
}

func Flatten(fields []models.FieldInfo) []models.FieldInfo {
	var out []models.FieldInfo
	for _, field := range fields {
		out = append(out, field)
		out = append(out, Flatten(field.Children)...)
	}
	return out
}

func scoreRisk(findings []models.Finding) int {
	score := 0
	for _, finding := range findings {
		score += finding.ScoreImpact
	}

	if score > 10 {
		return 10
	}

	return score
}

func operationName(op *ast.OperationDefinition) string {
	if op.Name != "" {
		return op.Name
	}
	return "AnonymousOperation"
}

func argumentNames(args ast.ArgumentList) []string {
	names := make([]string, 0, len(args))
	for _, arg := range args {
		names = append(names, arg.Name)
	}
	return names
}

func lineFromPosition(pos *ast.Position) int {
	if pos == nil {
		return 0
	}
	return pos.Line
}
