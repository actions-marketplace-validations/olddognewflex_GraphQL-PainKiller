package github

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// PRMetadata holds the information needed to post review comments on a PR.
type PRMetadata struct {
	Owner      string
	Repo       string
	Number     int
	CommitSHA  string
}

// LoadPRMetadataFromEnv extracts PR metadata from GitHub Actions environment variables.
func LoadPRMetadataFromEnv() (*PRMetadata, error) {
	repo := os.Getenv("GITHUB_REPOSITORY")
	if repo == "" {
		return nil, fmt.Errorf("GITHUB_REPOSITORY not set")
	}

	parts := strings.SplitN(repo, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid GITHUB_REPOSITORY format: %s", repo)
	}

	owner := parts[0]
	repoName := parts[1]

	eventPath := os.Getenv("GITHUB_EVENT_PATH")
	if eventPath == "" {
		return nil, fmt.Errorf("GITHUB_EVENT_PATH not set")
	}

	data, err := os.ReadFile(eventPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read GITHUB_EVENT_PATH: %w", err)
	}

	var payload struct {
		PullRequest struct {
			Number  int `json:"number"`
			Head    struct {
				SHA string `json:"sha"`
			} `json:"head"`
		} `json:"pull_request"`
	}

	if err := json.Unmarshal(data, &payload); err != nil {
		return nil, fmt.Errorf("failed to parse GitHub event payload: %w", err)
	}

	if payload.PullRequest.Number == 0 {
		return nil, fmt.Errorf("pull_request.number not found in event payload")
	}

	return &PRMetadata{
		Owner:     owner,
		Repo:      repoName,
		Number:    payload.PullRequest.Number,
		CommitSHA: payload.PullRequest.Head.SHA,
	}, nil
}

// LoadPRMetadataFromFlags allows manual override via CLI flags (for testing/local use).
func LoadPRMetadataFromFlags(owner, repo string, prNumber int, commitSHA string) *PRMetadata {
	return &PRMetadata{
		Owner:     owner,
		Repo:      repo,
		Number:    prNumber,
		CommitSHA: commitSHA,
	}
}

// ParsePRNumber parses a PR number from a string.
func ParsePRNumber(s string) (int, error) {
	return strconv.Atoi(s)
}
