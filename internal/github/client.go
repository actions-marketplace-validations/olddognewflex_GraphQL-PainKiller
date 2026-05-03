package github

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/v63/github"
)

// Client wraps the go-github client for posting PR review comments.
type Client struct {
	client *github.Client
}

// NewClient creates a new GitHub client from a token.
func NewClient(token string) *Client {
	if token == "" {
		token = os.Getenv("GITHUB_TOKEN")
	}
	return &Client{
		client: github.NewClient(nil).WithAuthToken(token),
	}
}

// PostReviewComment posts a single review comment on a PR.
func (c *Client) PostReviewComment(ctx context.Context, owner, repo string, prNumber int, comment ReviewComment, commitSHA string) error {
	_, _, err := c.client.PullRequests.CreateComment(ctx, owner, repo, prNumber, &github.PullRequestComment{
		Body:     github.String(comment.Body),
		CommitID: github.String(commitSHA),
		Path:     github.String(comment.Path),
		Line:     github.Int(comment.Line),
		Side:     github.String("RIGHT"),
	})
	if err != nil {
		return fmt.Errorf("failed to post review comment on %s:%d: %w", comment.Path, comment.Line, err)
	}
	return nil
}

// PostReviewComments posts multiple review comments on a PR.
func (c *Client) PostReviewComments(ctx context.Context, owner, repo string, prNumber int, comments []ReviewComment, commitSHA string) error {
	for _, comment := range comments {
		if err := c.PostReviewComment(ctx, owner, repo, prNumber, comment, commitSHA); err != nil {
			return err
		}
	}
	return nil
}
