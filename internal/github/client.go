package github

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/v63/github"
)

type Client struct {
	client *github.Client
}

func NewClient(token string) *Client {
	if token == "" {
		token = os.Getenv("GITHUB_TOKEN")
	}
	return &Client{
		client: github.NewClient(nil).WithAuthToken(token),
	}
}

// PostSummaryComment posts a regular comment (not a review) on a PR.
func (c *Client) PostSummaryComment(ctx context.Context, owner, repo string, prNumber int, body string) error {
	issueComment := &github.IssueComment{
		Body: github.String(body),
	}
	_, _, err := c.client.Issues.CreateComment(ctx, owner, repo, prNumber, issueComment)
	if err != nil {
		return fmt.Errorf("failed to create PR comment: %w", err)
	}
	return nil
}

func (c *Client) getChangedFiles(ctx context.Context, owner, repo string, prNumber int) (map[string]map[int]int, error) {
	result := make(map[string]map[int]int)

	opts := &github.ListOptions{PerPage: 100}
	for {
		files, resp, err := c.client.PullRequests.ListFiles(ctx, owner, repo, prNumber, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to list PR files: %w", err)
		}

		for _, f := range files {
			name := f.GetFilename()
			if name == "" {
				continue
			}
			result[name] = ParsePatchLines(f.GetPatch())
		}

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}

	return result, nil
}

// PostReview posts all findings as a single PR review. Comments on lines
// within the diff are posted inline; comments on changed files but outside
// diff hunks become file-level comments; findings on unchanged files are
// listed in the review body.
func (c *Client) PostReview(ctx context.Context, owner, repo string, prNumber int, comments []ReviewComment, commitSHA string) error {
	diffLines, err := c.getChangedFiles(ctx, owner, repo, prNumber)
	if err != nil {
		return err
	}

	var inline []*github.DraftReviewComment
	var bodyFindings []ReviewComment

	for _, comment := range comments {
		fileLines, fileChanged := diffLines[comment.Path]
		position, lineInDiff := fileLines[comment.Line]
		switch {
		case fileChanged && lineInDiff:
			inline = append(inline, &github.DraftReviewComment{
				Path:     github.String(comment.Path),
				Position: github.Int(position),
				Body:     github.String(comment.Body),
			})

		case fileChanged:
			bodyFindings = append(bodyFindings, comment)

		default:
			bodyFindings = append(bodyFindings, comment)
		}
	}

	body := buildReviewBody(len(comments), len(inline), bodyFindings)

	review := &github.PullRequestReviewRequest{
		CommitID: github.String(commitSHA),
		Body:     github.String(body),
		Event:    github.String("COMMENT"),
		Comments: inline,
	}

	_, _, err = c.client.PullRequests.CreateReview(ctx, owner, repo, prNumber, review)
	if err != nil {
		return fmt.Errorf("failed to create review: %w", err)
	}

	return nil
}

func buildReviewBody(total, inlineCount int, bodyFindings []ReviewComment) string {
	var b strings.Builder
	b.WriteString("## GraphQL Painkiller\n\n")

	if len(bodyFindings) == 0 {
		fmt.Fprintf(&b, "Found **%d** finding(s) — all commented inline.\n", total)
		return b.String()
	}

	fmt.Fprintf(&b, "Found **%d** finding(s) — **%d** inline, **%d** in review body.\n\n",
		total, inlineCount, len(bodyFindings))

	b.WriteString("<details>\n<summary>Findings not shown inline</summary>\n\n")
	for _, f := range bodyFindings {
		fmt.Fprintf(&b, "---\n\n📍 `%s` line %d\n\n%s\n\n", f.Path, f.Line, f.Body)
	}
	b.WriteString("</details>\n")

	return b.String()
}
