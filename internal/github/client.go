package github

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const defaultBaseURL = "https://api.github.com"

// Client retrieves repository information from the GitHub REST API.
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	Token      string
}

// NewClient creates a GitHub API client with sensible defaults.
func NewClient(token string) *Client {
	return &Client{
		BaseURL: defaultBaseURL,
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		Token: strings.TrimSpace(token),
	}
}

// Analyze retrieves repository metadata and language data concurrently.
func (c *Client) Analyze(ctx context.Context, owner, repo string) (Analysis, error) {
	owner = strings.TrimSpace(owner)
	repo = strings.TrimSpace(repo)
	if owner == "" || repo == "" {
		return Analysis{}, errors.New("owner and repository are required")
	}

	var (
		repository Repository
		languages  map[string]int
		repoErr    error
		langErr    error
		wg         sync.WaitGroup
	)

	wg.Add(2)

	go func() {
		defer wg.Done()
		path := fmt.Sprintf("/repos/%s/%s", url.PathEscape(owner), url.PathEscape(repo))
		repoErr = c.getJSON(ctx, path, &repository)
	}()

	go func() {
		defer wg.Done()
		path := fmt.Sprintf("/repos/%s/%s/languages", url.PathEscape(owner), url.PathEscape(repo))
		langErr = c.getJSON(ctx, path, &languages)
	}()

	wg.Wait()

	if repoErr != nil {
		return Analysis{}, repoErr
	}
	if langErr != nil {
		return Analysis{}, langErr
	}

	total := 0
	for _, bytes := range languages {
		total += bytes
	}

	license := "Not specified"
	if repository.License != nil && repository.License.Name != "" {
		license = repository.License.Name
	}

	return Analysis{
		Name:          repository.Name,
		FullName:      repository.FullName,
		Description:   repository.Description,
		URL:           repository.HTMLURL,
		Stars:         repository.Stars,
		Forks:         repository.Forks,
		OpenIssues:    repository.OpenIssues,
		DefaultBranch: repository.DefaultBranch,
		License:       license,
		UpdatedAt:     repository.UpdatedAt,
		Languages:     languages,
		TotalBytes:    total,
	}, nil
}

func (c *Client) getJSON(ctx context.Context, path string, target any) error {
	base := strings.TrimRight(c.BaseURL, "/")
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, base+path, nil)
	if err != nil {
		return fmt.Errorf("create GitHub request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Set("User-Agent", "RepoRadar")
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("GitHub request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		var apiError struct {
			Message string `json:"message"`
		}
		_ = json.NewDecoder(res.Body).Decode(&apiError)
		if apiError.Message == "" {
			apiError.Message = res.Status
		}
		return fmt.Errorf("GitHub API: %s", apiError.Message)
	}

	if err := json.NewDecoder(res.Body).Decode(target); err != nil {
		return fmt.Errorf("decode GitHub response: %w", err)
	}
	return nil
}
