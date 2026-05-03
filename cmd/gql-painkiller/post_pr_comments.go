package main

import (
	"context"
	"fmt"

	"github.com/olddognewflex/graphql-painkiller/internal/analyzer"
	"github.com/olddognewflex/graphql-painkiller/internal/config"
	"github.com/olddognewflex/graphql-painkiller/internal/extractors"
	"github.com/olddognewflex/graphql-painkiller/internal/github"
	"github.com/olddognewflex/graphql-painkiller/internal/models"
	"github.com/spf13/cobra"
)

var (
	prOwner     string
	prRepo      string
	prNumber    int
	prCommitSHA string
	prToken     string
)

func postPRCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "post-pr-comments [path]",
		Short: "Analyze GraphQL and post findings as PR review comments",
		Long: `Analyze GraphQL operations and post findings as inline PR review comments.

When running in GitHub Actions, PR metadata is auto-detected from environment variables.
For local testing or other CI systems, use --owner, --repo, --pr, and --commit flags.`,
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

			comments := github.BuildReviewComments(reports)
			if len(comments) == 0 {
				fmt.Fprintln(cmd.OutOrStdout(), "No findings to report.")
				return nil
			}

			var meta *github.PRMetadata
			if prOwner != "" && prRepo != "" && prNumber != 0 && prCommitSHA != "" {
				meta = github.LoadPRMetadataFromFlags(prOwner, prRepo, prNumber, prCommitSHA)
			} else {
				meta, err = github.LoadPRMetadataFromEnv()
				if err != nil {
					return fmt.Errorf("failed to detect PR metadata from environment; use --owner, --repo, --pr, --commit flags: %w", err)
				}
			}

			client := github.NewClient(prToken)
			ctx := context.Background()

			fmt.Fprintf(cmd.OutOrStdout(), "Posting %d review comments to %s/%s#%d...\n", len(comments), meta.Owner, meta.Repo, meta.Number)

			if err := client.PostReviewComments(ctx, meta.Owner, meta.Repo, meta.Number, comments, meta.CommitSHA); err != nil {
				return fmt.Errorf("failed to post review comments: %w", err)
			}

			fmt.Fprintln(cmd.OutOrStdout(), "Done.")
			return nil
		},
	}

	cmd.Flags().StringVar(&prOwner, "owner", "", "repository owner (overrides env detection)")
	cmd.Flags().StringVar(&prRepo, "repo", "", "repository name (overrides env detection)")
	cmd.Flags().IntVar(&prNumber, "pr", 0, "pull request number (overrides env detection)")
	cmd.Flags().StringVar(&prCommitSHA, "commit", "", "head commit SHA (overrides env detection)")
	cmd.Flags().StringVar(&prToken, "token", "", "GitHub token (defaults to GITHUB_TOKEN env var)")

	return cmd
}
