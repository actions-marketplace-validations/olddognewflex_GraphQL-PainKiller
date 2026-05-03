package github

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadPRMetadataFromFlags(t *testing.T) {
	meta := LoadPRMetadataFromFlags("owner", "repo", 42, "abc123")

	if meta.Owner != "owner" {
		t.Errorf("Owner = %q, want %q", meta.Owner, "owner")
	}
	if meta.Repo != "repo" {
		t.Errorf("Repo = %q, want %q", meta.Repo, "repo")
	}
	if meta.Number != 42 {
		t.Errorf("Number = %d, want %d", meta.Number, 42)
	}
	if meta.CommitSHA != "abc123" {
		t.Errorf("CommitSHA = %q, want %q", meta.CommitSHA, "abc123")
	}
}

func TestParsePRNumber(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int
		wantErr bool
	}{
		{name: "valid number", input: "42", want: 42, wantErr: false},
		{name: "zero", input: "0", want: 0, wantErr: false},
		{name: "large number", input: "999999", want: 999999, wantErr: false},
		{name: "negative number", input: "-1", want: -1, wantErr: false},
		{name: "not a number", input: "abc", want: 0, wantErr: true},
		{name: "empty string", input: "", want: 0, wantErr: true},
		{name: "decimal", input: "3.14", want: 0, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParsePRNumber(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParsePRNumber(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParsePRNumber(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

func TestLoadPRMetadataFromEnv(t *testing.T) {
	tests := []struct {
		name      string
		envRepo   string
		eventJSON string
		want      *PRMetadata
		wantErr   bool
	}{
		{
			name:    "valid event payload",
			envRepo: "owner/repo",
			eventJSON: `{
				"pull_request": {
					"number": 42,
					"head": {"sha": "abc123def456"}
				}
			}`,
			want: &PRMetadata{
				Owner:     "owner",
				Repo:      "repo",
				Number:    42,
				CommitSHA: "abc123def456",
			},
			wantErr: false,
		},
		{
			name:      "missing repository",
			envRepo:   "",
			eventJSON: `{}`,
			want:      nil,
			wantErr:   true,
		},
		{
			name:    "invalid repository format",
			envRepo: "invalid-repo-format",
			eventJSON: `{
				"pull_request": {
					"number": 1,
					"head": {"sha": "sha"}
				}
			}`,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "missing pull request number",
			envRepo: "owner/repo",
			eventJSON: `{
				"pull_request": {}
			}`,
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("GITHUB_REPOSITORY", tt.envRepo)

			tmpDir := t.TempDir()
			eventPath := filepath.Join(tmpDir, "event.json")
			if err := os.WriteFile(eventPath, []byte(tt.eventJSON), 0644); err != nil {
				t.Fatalf("failed to write event file: %v", err)
			}
			t.Setenv("GITHUB_EVENT_PATH", eventPath)

			got, err := LoadPRMetadataFromEnv()
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadPRMetadataFromEnv() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if got.Owner != tt.want.Owner {
				t.Errorf("Owner = %q, want %q", got.Owner, tt.want.Owner)
			}
			if got.Repo != tt.want.Repo {
				t.Errorf("Repo = %q, want %q", got.Repo, tt.want.Repo)
			}
			if got.Number != tt.want.Number {
				t.Errorf("Number = %d, want %d", got.Number, tt.want.Number)
			}
			if got.CommitSHA != tt.want.CommitSHA {
				t.Errorf("CommitSHA = %q, want %q", got.CommitSHA, tt.want.CommitSHA)
			}
		})
	}
}
