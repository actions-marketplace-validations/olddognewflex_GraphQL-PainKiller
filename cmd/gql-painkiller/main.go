package main

import (
	"fmt"
	"os"

	"github.com/olddognewflex/graphql-painkiller/internal/analyzer"
	"github.com/olddognewflex/graphql-painkiller/internal/config"
	"github.com/olddognewflex/graphql-painkiller/internal/extractors"
	"github.com/olddognewflex/graphql-painkiller/internal/models"
	"github.com/olddognewflex/graphql-painkiller/internal/reporters"
	"github.com/olddognewflex/graphql-painkiller/internal/severity"
	"github.com/spf13/cobra"
)

var (
	configPath string
	jsonOutput bool
	failOn     string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "gql-painkiller",
		Short: "Static analysis for risky GraphQL query patterns",
	}

	rootCmd.PersistentFlags().StringVar(&configPath, "config", "gql-painkiller.config.yaml", "path to config file")

	rootCmd.AddCommand(initCmd())
	rootCmd.AddCommand(analyzeCmd())
	rootCmd.AddCommand(postPRCmd())
	rootCmd.AddCommand(postPRSummaryCmd())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func initCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Create a default gql-painkiller.config.yaml",
		RunE: func(cmd *cobra.Command, args []string) error {
			if _, err := os.Stat(configPath); err == nil {
				return fmt.Errorf("%s already exists", configPath)
			}

			if err := config.WriteDefault(configPath); err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Created %s\n", configPath)
			return nil
		},
	}
}

func analyzeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "analyze [path]",
		Short: "Analyze GraphQL operations",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.SilenceUsage = true
			target := "."
			if len(args) == 1 {
				target = args[0]
			}

			cfg, err := config.Load(configPath)
			if err != nil {
				return err
			}

			if failOn != "" {
				cfg.Rules.FailOnSeverity = severity.Severity(failOn)
			}

			docs, err := extractors.Extract(target)
			if err != nil {
				return err
			}

			var reports []models.Report

			for _, doc := range docs {
				docReports, err := analyzer.AnalyzeDocument(doc, cfg)
				if err != nil {
					fmt.Fprintf(cmd.ErrOrStderr(), "warning: %v\n", err)
					continue
				}

				reports = append(reports, docReports...)
			}

			if jsonOutput {
				if err := reporters.JSON(cmd.OutOrStdout(), reports); err != nil {
					return err
				}
			} else {
				reporters.Text(cmd.OutOrStdout(), reports)
			}

			if shouldFail(reports, cfg.Rules.FailOnSeverity) {
				return fmt.Errorf("GraphQL Painkiller found findings at or above severity %q", cfg.Rules.FailOnSeverity)
			}

			return nil
		},
	}

	cmd.Flags().BoolVar(&jsonOutput, "json", false, "emit JSON report")
	cmd.Flags().StringVar(&failOn, "fail-on", "", "override fail severity: none, info, warning, high, critical")

	return cmd
}

func shouldFail(reports []models.Report, threshold severity.Severity) bool {
	if threshold == severity.None {
		return false
	}
	for _, report := range reports {
		for _, finding := range report.Findings {
			if severity.GTE(finding.Severity, threshold) {
				return true
			}
		}
	}
	return false
}
