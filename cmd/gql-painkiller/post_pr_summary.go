package main

import (
	"context"
	"fmt"

	"github.com/olddognewflex/graphql-painkiller/internal/analyzer"
	"github.com/olddognewflex/graphql-painkiller/internal/config"
	"github.com/olddognewflex/graphql-painkiller/internal/extractors"
	"github.com/olddognewflex/graphql-painkiller/internal/github"
	"github.com/olddognewflex/graphql-painkiller/internal/models"
	"github.com/olddognewflex/graphql-painkiller/internal/reporters"
	"github.com/olddognewflex/graphql-painkiller/internal/severity"
	"github.com/spf13/cobra"
)

var (
	summaryOwner  string
	summaryRepo   string
	summaryNumber int
	summaryToken  string
	summaryFailOn string
)

func postPRSummaryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "post-pr-summary [path]",
		Short: "Analyze GraphQL and post a summary comment to a PR",
		Long: `Analyze GraphQL operations and post a summary comment (not a review) to a PR.

When running in GitHub Actions, PR metadata is auto-detected from environment variables.
For local testing or other CI systems, use --owner, --repo, and --pr flags.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			target := "."
			if len(args) == 1 {
				target = args[0]
			}

			cfg, err := config.Load(configPath)
			if err != nil {
				return err
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

			body := reporters.SummaryMarkdown(reports)

			var meta *github.PRMetadata
			if summaryOwner != "" && summaryRepo != "" && summaryNumber != 0 {
				meta = github.LoadPRMetadataFromFlags(summaryOwner, summaryRepo, summaryNumber, "")
			} else {
				meta, err = github.LoadPRMetadataFromEnv()
				if err != nil {
					return fmt.Errorf("failed to detect PR metadata from environment; use --owner, --repo, --pr flags: %w", err)
				}
			}

			client := github.NewClient(summaryToken)
			ctx := context.Background()

			fmt.Fprintf(cmd.OutOrStdout(), "Posting summary comment to %s/%s#%d...\n", meta.Owner, meta.Repo, meta.Number)

			if err := client.PostSummaryComment(ctx, meta.Owner, meta.Repo, meta.Number, body); err != nil {
				return fmt.Errorf("failed to post summary comment: %w", err)
			}

			fmt.Fprintln(cmd.OutOrStdout(), "Done.")

			if summaryFailOn != "" {
				cfg.Rules.FailOnSeverity = severity.Severity(summaryFailOn)
			}

			if shouldFail(reports, cfg.Rules.FailOnSeverity) {
				cmd.SilenceErrors = true
				err := fmt.Errorf("GraphQL Painkiller found findings at or above severity %q", cfg.Rules.FailOnSeverity)
				fmt.Fprintf(cmd.ErrOrStderr(), "Error: %v\n", err)
				return err
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&summaryOwner, "owner", "", "repository owner (overrides env detection)")
	cmd.Flags().StringVar(&summaryRepo, "repo", "", "repository name (overrides env detection)")
	cmd.Flags().IntVar(&summaryNumber, "pr", 0, "pull request number (overrides env detection)")
	cmd.Flags().StringVar(&summaryToken, "token", "", "GitHub token (defaults to GITHUB_TOKEN env var)")
	cmd.Flags().StringVar(&summaryFailOn, "fail-on", "", "override fail severity: none, info, warning, high, critical")

	return cmd
}
