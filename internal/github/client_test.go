package github

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPostReviewComment(t *testing.T) {
	tests := []struct {
		name       string
		comment    ReviewComment
		serverCode int
		wantErr    bool
	}{
		{
			name: "successful post",
			comment: ReviewComment{
				Path: "test.graphql",
				Line: 42,
				Body: "Test comment",
			},
			serverCode: http.StatusCreated,
			wantErr:    false,
		},
		{
			name: "server error",
			comment: ReviewComment{
				Path: "test.graphql",
				Line: 1,
				Body: "Test",
			},
			serverCode: http.StatusInternalServerError,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("Method = %q, want %q", r.Method, http.MethodPost)
				}

				expectedPath := "/repos/owner/repo/pulls/42/comments"
				if !strings.HasSuffix(r.URL.Path, expectedPath) {
					t.Errorf("Path = %q, want suffix %q", r.URL.Path, expectedPath)
				}

				var body map[string]interface{}
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Fatalf("failed to decode request body: %v", err)
				}

				if body["body"] != tt.comment.Body {
					t.Errorf("body = %q, want %q", body["body"], tt.comment.Body)
				}
				if body["path"] != tt.comment.Path {
					t.Errorf("path = %q, want %q", body["path"], tt.comment.Path)
				}
				if line, ok := body["line"].(float64); !ok || int(line) != tt.comment.Line {
					t.Errorf("line = %v, want %d", body["line"], tt.comment.Line)
				}

				w.WriteHeader(tt.serverCode)
				if tt.serverCode == http.StatusCreated {
					json.NewEncoder(w).Encode(map[string]interface{}{
						"id": 123,
					})
				}
			}))
			defer server.Close()

			client := NewClient("test-token")
			client.client.BaseURL, _ = client.client.BaseURL.Parse(server.URL + "/")

			err := client.PostReviewComment(t.Context(), "owner", "repo", 42, tt.comment, "sha123")
			if (err != nil) != tt.wantErr {
				t.Errorf("PostReviewComment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPostReviewComments(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]interface{}{"id": callCount})
	}))
	defer server.Close()

	client := NewClient("test-token")
	client.client.BaseURL, _ = client.client.BaseURL.Parse(server.URL + "/")

	comments := []ReviewComment{
		{Path: "a.graphql", Line: 1, Body: "First"},
		{Path: "b.graphql", Line: 2, Body: "Second"},
	}

	err := client.PostReviewComments(t.Context(), "owner", "repo", 1, comments, "sha")
	if err != nil {
		t.Errorf("PostReviewComments() error = %v", err)
	}
	if callCount != 2 {
		t.Errorf("Expected 2 API calls, got %d", callCount)
	}
}
