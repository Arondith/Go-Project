package github

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAnalyze(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/repos/acme/widget", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"name":"widget",
			"full_name":"acme/widget",
			"description":"Example repository",
			"html_url":"https://github.com/acme/widget",
			"stargazers_count":42,
			"forks_count":7,
			"open_issues_count":3,
			"default_branch":"main",
			"license":{"name":"MIT License"},
			"updated_at":"2026-09-19T00:00:00Z"
		}`))
	})
	mux.HandleFunc("/repos/acme/widget/languages", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"Go":900,"JavaScript":100}`))
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClient("")
	client.BaseURL = server.URL
	client.HTTPClient = server.Client()

	got, err := client.Analyze(context.Background(), "acme", "widget")
	if err != nil {
		t.Fatalf("Analyze returned error: %v", err)
	}

	if got.FullName != "acme/widget" {
		t.Fatalf("FullName = %q, want acme/widget", got.FullName)
	}
	if got.TotalBytes != 1000 {
		t.Fatalf("TotalBytes = %d, want 1000", got.TotalBytes)
	}
	if got.Languages["Go"] != 900 {
		t.Fatalf("Go bytes = %d, want 900", got.Languages["Go"])
	}
}

func TestAnalyzeRequiresOwnerAndRepo(t *testing.T) {
	client := NewClient("")
	if _, err := client.Analyze(context.Background(), "", "widget"); err == nil {
		t.Fatal("expected validation error")
	}
}
