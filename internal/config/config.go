package config

import (
	"os"

	"github.com/olddognewflex/graphql-painkiller/internal/severity"
	"gopkg.in/yaml.v3"
)

type KnownResolver struct {
	Risk    severity.Severity `yaml:"risk" json:"risk"`
	Reason  string            `yaml:"reason" json:"reason"`
	Service string            `yaml:"service,omitempty" json:"service,omitempty"`
}

type Rules struct {
	MaxDepth                     int               `yaml:"maxDepth" json:"maxDepth"`
	MaxCollectionSelectionFields int               `yaml:"maxCollectionSelectionFields" json:"maxCollectionSelectionFields"`
	RequirePagination           bool              `yaml:"requirePagination" json:"requirePagination"`
	FailOnSeverity              severity.Severity `yaml:"failOnSeverity" json:"failOnSeverity"`
}

type Config struct {
	Rules                   Rules                    `yaml:"rules" json:"rules"`
	PaginationArgs          []string                 `yaml:"paginationArgs" json:"paginationArgs"`
	CollectionFieldPatterns []string                 `yaml:"collectionFieldPatterns" json:"collectionFieldPatterns"`
	ExpensiveFieldPatterns  []string                 `yaml:"expensiveFieldPatterns" json:"expensiveFieldPatterns"`
	KnownResolvers          map[string]KnownResolver `yaml:"knownResolvers" json:"knownResolvers"`
}

func Default() Config {
	return Config{
		Rules: Rules{
			MaxDepth:                     5,
			MaxCollectionSelectionFields: 8,
			RequirePagination:           true,
			FailOnSeverity:              severity.High,
		},
		PaginationArgs:          []string{"first", "last", "limit", "take", "pageSize", "after", "before", "offset"},
		CollectionFieldPatterns: []string{"items", "nodes", "edges"},
		ExpensiveFieldPatterns:  []string{"comments", "history", "events", "logs", "charges", "payments", "inspections", "accounts", "permissions", "audit"},
		KnownResolvers:          map[string]KnownResolver{},
	}
}

func Load(path string) (Config, error) {
	cfg := Default()

	bytes, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return cfg, err
	}

	var fromFile Config
	if err := yaml.Unmarshal(bytes, &fromFile); err != nil {
		return cfg, err
	}

	merge(&cfg, fromFile)
	return cfg, nil
}

func WriteDefault(path string) error {
	cfg := Default()
	bytes, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(path, bytes, 0644)
}

func merge(base *Config, incoming Config) {
	if incoming.Rules.MaxDepth != 0 {
		base.Rules.MaxDepth = incoming.Rules.MaxDepth
	}
	if incoming.Rules.MaxCollectionSelectionFields != 0 {
		base.Rules.MaxCollectionSelectionFields = incoming.Rules.MaxCollectionSelectionFields
	}
	base.Rules.RequirePagination = incoming.Rules.RequirePagination || base.Rules.RequirePagination
	if incoming.Rules.FailOnSeverity != "" {
		base.Rules.FailOnSeverity = incoming.Rules.FailOnSeverity
	}
	if len(incoming.PaginationArgs) > 0 {
		base.PaginationArgs = incoming.PaginationArgs
	}
	if len(incoming.CollectionFieldPatterns) > 0 {
		base.CollectionFieldPatterns = incoming.CollectionFieldPatterns
	}
	if len(incoming.ExpensiveFieldPatterns) > 0 {
		base.ExpensiveFieldPatterns = incoming.ExpensiveFieldPatterns
	}
	if incoming.KnownResolvers != nil {
		base.KnownResolvers = incoming.KnownResolvers
	}
}
