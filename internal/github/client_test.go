package github

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPostReview_InlineComments(t *testing.T) {
	mux := http.NewServeMux()

	mux.HandleFunc("/repos/owner/repo/pulls/1/files", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]map[string]interface{}{
			{
				"filename": "src/query.graphql",
				"patch":    "@@ -1,3 +1,5 @@\n line1\n+line2\n+line3\n line4\n line5",
			},
		})
	})

	var reviewBody map[string]interface{}
	mux.HandleFunc("/repos/owner/repo/pulls/1/reviews", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Method = %q, want POST", r.Method)
		}
		if err := json.NewDecoder(r.Body).Decode(&reviewBody); err != nil {
			t.Fatalf("failed to decode review body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"id": 1})
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient("test-token")
	client.client.BaseURL, _ = client.client.BaseURL.Parse(server.URL + "/")

	comments := []ReviewComment{
		{Path: "src/query.graphql", Line: 2, Body: "deep query"},
	}

	if err := client.PostReview(t.Context(), "owner", "repo", 1, comments, "sha123"); err != nil {
		t.Fatalf("PostReview() error = %v", err)
	}

	reviewComments := reviewBody["comments"].([]interface{})
	if len(reviewComments) != 1 {
		t.Fatalf("expected 1 inline comment, got %d", len(reviewComments))
	}

	first := reviewComments[0].(map[string]interface{})
	if first["path"] != "src/query.graphql" {
		t.Errorf("path = %q, want %q", first["path"], "src/query.graphql")
	}
	if int(first["line"].(float64)) != 2 {
		t.Errorf("line = %v, want 2", first["line"])
	}
}

func TestPostReview_FileLevelFallback(t *testing.T) {
	mux := http.NewServeMux()

	mux.HandleFunc("/repos/owner/repo/pulls/1/files", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]map[string]interface{}{
			{
				"filename": "src/query.graphql",
				"patch":    "@@ -1,2 +1,2 @@\n-old\n+new",
			},
		})
	})

	var reviewBody map[string]interface{}
	mux.HandleFunc("/repos/owner/repo/pulls/1/reviews", func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&reviewBody)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"id": 1})
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient("test-token")
	client.client.BaseURL, _ = client.client.BaseURL.Parse(server.URL + "/")

	comments := []ReviewComment{
		{Path: "src/query.graphql", Line: 99, Body: "finding on unchanged line"},
	}

	if err := client.PostReview(t.Context(), "owner", "repo", 1, comments, "sha123"); err != nil {
		t.Fatalf("PostReview() error = %v", err)
	}

	reviewComments := reviewBody["comments"].([]interface{})
	if len(reviewComments) != 1 {
		t.Fatalf("expected 1 file-level comment, got %d", len(reviewComments))
	}

	first := reviewComments[0].(map[string]interface{})
	if first["subject_type"] != "file" {
		t.Errorf("subject_type = %q, want %q", first["subject_type"], "file")
	}
	if !strings.Contains(first["body"].(string), "Line 99") {
		t.Errorf("body should reference line 99, got %q", first["body"])
	}
}

func TestPostReview_UnchangedFileFallsToBody(t *testing.T) {
	mux := http.NewServeMux()

	mux.HandleFunc("/repos/owner/repo/pulls/1/files", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]map[string]interface{}{
			{
				"filename": "src/other.graphql",
				"patch":    "@@ -1,1 +1,1 @@\n-old\n+new",
			},
		})
	})

	var reviewBody map[string]interface{}
	mux.HandleFunc("/repos/owner/repo/pulls/1/reviews", func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&reviewBody)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"id": 1})
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient("test-token")
	client.client.BaseURL, _ = client.client.BaseURL.Parse(server.URL + "/")

	comments := []ReviewComment{
		{Path: "src/untouched.graphql", Line: 10, Body: "finding in unchanged file"},
	}

	if err := client.PostReview(t.Context(), "owner", "repo", 1, comments, "sha123"); err != nil {
		t.Fatalf("PostReview() error = %v", err)
	}

	body := reviewBody["body"].(string)
	if !strings.Contains(body, "unchanged files") {
		t.Errorf("review body should mention unchanged files, got %q", body)
	}
	if !strings.Contains(body, "src/untouched.graphql") {
		t.Errorf("review body should reference the file, got %q", body)
	}

	if reviewBody["comments"] != nil {
		commentsList := reviewBody["comments"].([]interface{})
		if len(commentsList) != 0 {
			t.Errorf("expected 0 inline comments, got %d", len(commentsList))
		}
	}
}

func TestPostReview_MixedPlacement(t *testing.T) {
	mux := http.NewServeMux()

	mux.HandleFunc("/repos/owner/repo/pulls/1/files", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]map[string]interface{}{
			{
				"filename": "src/query.graphql",
				"patch":    "@@ -1,3 +1,4 @@\n line1\n+line2\n line3\n line4",
			},
		})
	})

	var reviewBody map[string]interface{}
	mux.HandleFunc("/repos/owner/repo/pulls/1/reviews", func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&reviewBody)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"id": 1})
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient("test-token")
	client.client.BaseURL, _ = client.client.BaseURL.Parse(server.URL + "/")

	comments := []ReviewComment{
		{Path: "src/query.graphql", Line: 2, Body: "inline finding"},
		{Path: "src/query.graphql", Line: 50, Body: "file-level finding"},
		{Path: "src/missing.graphql", Line: 1, Body: "body finding"},
	}

	if err := client.PostReview(t.Context(), "owner", "repo", 1, comments, "sha123"); err != nil {
		t.Fatalf("PostReview() error = %v", err)
	}

	reviewComments := reviewBody["comments"].([]interface{})
	if len(reviewComments) != 2 {
		t.Fatalf("expected 2 inline/file-level comments, got %d", len(reviewComments))
	}

	body := reviewBody["body"].(string)
	if !strings.Contains(body, "src/missing.graphql") {
		t.Errorf("review body should reference unchanged file, got %q", body)
	}
}
