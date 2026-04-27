package render

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// RemoteFetcher retrieves raw template files from a remote URL. The default
// implementation supports plain HTTPS and `github://owner/repo/path@ref` URIs;
// alternative fetchers can be plugged in for tests or air-gapped installs.
type RemoteFetcher interface {
	Fetch(ctx context.Context, uri string) ([]byte, error)
}

// HTTPFetcher is the default RemoteFetcher. It rewrites github:// URIs to
// raw.githubusercontent.com and performs a single HTTP GET with a timeout.
type HTTPFetcher struct {
	Client  *http.Client
	Timeout time.Duration
}

// NewHTTPFetcher constructs an HTTPFetcher with sensible defaults.
func NewHTTPFetcher() *HTTPFetcher {
	return &HTTPFetcher{
		Client:  &http.Client{Timeout: 30 * time.Second},
		Timeout: 30 * time.Second,
	}
}

// Fetch retrieves uri and returns its raw body. github:// URIs are normalized
// to raw.githubusercontent.com URLs.
func (f *HTTPFetcher) Fetch(ctx context.Context, uri string) ([]byte, error) {
	url, err := normalizeURI(uri)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("render: build request: %w", err)
	}

	client := f.Client
	if client == nil {
		client = http.DefaultClient
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("render: fetch %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("render: fetch %s: status %d", url, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("render: read body %s: %w", url, err)
	}
	return body, nil
}

// normalizeURI converts a friendly URI into an HTTP URL.
//
//	github://owner/repo/path/to/file@ref  → https://raw.githubusercontent.com/owner/repo/ref/path/to/file
//	https://...                            → unchanged
//	http://...                             → unchanged
func normalizeURI(uri string) (string, error) {
	switch {
	case strings.HasPrefix(uri, "github://"):
		rest := strings.TrimPrefix(uri, "github://")
		atIdx := strings.LastIndex(rest, "@")
		if atIdx < 0 {
			return "", fmt.Errorf("render: github URI %q missing @ref", uri)
		}
		path := rest[:atIdx]
		ref := rest[atIdx+1:]

		parts := strings.SplitN(path, "/", 3)
		if len(parts) < 3 {
			return "", fmt.Errorf("render: github URI %q must be owner/repo/path@ref", uri)
		}
		owner, repo, file := parts[0], parts[1], parts[2]
		return fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/%s/%s", owner, repo, ref, file), nil

	case strings.HasPrefix(uri, "http://"), strings.HasPrefix(uri, "https://"):
		return uri, nil

	default:
		return "", fmt.Errorf("render: unsupported URI scheme: %q", uri)
	}
}
