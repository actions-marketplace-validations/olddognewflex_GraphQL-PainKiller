package models

import "github.com/olddognewflex/graphql-painkiller/internal/severity"

type FieldInfo struct {
	Name      string      `json:"name"`
	Path      string      `json:"path"`
	Depth     int         `json:"depth"`
	Line      int         `json:"line,omitempty"`
	Arguments []string    `json:"arguments"`
	Children  []FieldInfo `json:"children"`
}

type Finding struct {
	RuleID      string            `json:"ruleId"`
	Severity    severity.Severity `json:"severity"`
	Message     string            `json:"message"`
	FilePath    string            `json:"filePath"`
	Line        int               `json:"line,omitempty"`
	Path        string            `json:"path,omitempty"`
	ScoreImpact int               `json:"scoreImpact"`
	Suggestion  string            `json:"suggestion,omitempty"`
}

type Report struct {
	FilePath      string            `json:"filePath"`
	OperationName string           `json:"operationName"`
	RiskScore     int              `json:"riskScore"`
	Severity      severity.Severity `json:"severity"`
	Findings      []Finding        `json:"findings"`
}
