package render

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/version14/dot/internal/fileutils"
)

// RepoSnapshot represents a snapshot of a GitHub repository's files.
type RepoSnapshot struct {
	Files       map[string][]byte // Map of file path to content (preserves folder structure)
	Ref         string            // The git ref (branch, tag, or commit)
	SourceURL   string            // The original repository URL
	TopLevelDir string            // The top-level directory name from GitHub tarball (e.g., "owner-repo-hash")
}

// GitHubRepoFetcher retrieves a repository snapshot from GitHub.
type GitHubRepoFetcher interface {
	FetchRepo(ctx context.Context, repoURL string, opts FetchOptions) (*RepoSnapshot, error)
}

// FetchOptions configures repository fetching behavior.
type FetchOptions struct {
	Ref       string // Git ref (branch, tag, commit). Defaults to default branch.
	AuthToken string // GitHub personal access token for private repos (optional).
}

// GitHubArchiveFetcher implements GitHubRepoFetcher using GitHub's tarball API.
type GitHubArchiveFetcher struct {
	Client  *http.Client
	Timeout time.Duration
}

// NewGitHubArchiveFetcher constructs a GitHubArchiveFetcher with sensible defaults.
func NewGitHubArchiveFetcher() *GitHubArchiveFetcher {
	return &GitHubArchiveFetcher{
		Client:  &http.Client{Timeout: 60 * time.Second},
		Timeout: 60 * time.Second,
	}
}

// FetchRepo downloads a GitHub repository as a tarball and extracts it into a map.
// Supports URLs like:
//   - https://github.com/owner/repo
//   - https://github.com/owner/repo.git
//   - github://owner/repo (assumes default branch)
func (f *GitHubArchiveFetcher) FetchRepo(ctx context.Context, repoURL string, opts FetchOptions) (*RepoSnapshot, error) {
	owner, repo, err := parseGitHubURL(repoURL)
	if err != nil {
		return nil, fmt.Errorf("render: parse GitHub URL: %w", err)
	}

	ref := opts.Ref
	if ref == "" {
		ref = "HEAD"
	}

	// Construct the GxitHub API tarball URL
	archiveURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/tarball/%s", owner, repo, ref)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, archiveURL, nil)
	if err != nil {
		return nil, fmt.Errorf("render: build request: %w", err)
	}

	if opts.AuthToken != "" {
		req.Header.Set("Authorization", "token "+opts.AuthToken)
	}

	client := f.Client
	if client == nil {
		client = http.DefaultClient
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("render: fetch tarball from %s: %w", archiveURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("render: fetch tarball: status %d from %s", resp.StatusCode, archiveURL)
	}

	files, topLevelDir, err := extractTarballToMap(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("render: extract tarball: %w", err)
	}

	return &RepoSnapshot{
		Files:       files,
		Ref:         ref,
		SourceURL:   repoURL,
		TopLevelDir: topLevelDir,
	}, nil
}

// parseGitHubURL extracts owner and repo from various GitHub URL formats.
func parseGitHubURL(urlStr string) (owner, repo string, err error) {
	if rest, ok := strings.CutPrefix(urlStr, "github://"); ok {
		parts := strings.SplitN(rest, "/", 2)
		if len(parts) != 2 {
			return "", "", fmt.Errorf("invalid github:// URL format, expected owner/repo")
		}
		return parts[0], parts[1], nil
	}

	if rest, ok := strings.CutPrefix(urlStr, "https://github.com/"); ok {
		parts := strings.SplitN(rest, "/", 2)
		if len(parts) != 2 {
			return "", "", fmt.Errorf("invalid GitHub URL format")
		}
		owner = parts[0]
		repo = strings.TrimSuffix(parts[1], ".git")
		return owner, repo, nil
	}

	if rest, ok := strings.CutPrefix(urlStr, "http://github.com/"); ok {
		parts := strings.SplitN(rest, "/", 2)
		if len(parts) != 2 {
			return "", "", fmt.Errorf("invalid GitHub URL format")
		}
		owner = parts[0]
		repo = strings.TrimSuffix(parts[1], ".git")
		return owner, repo, nil
	}

	return "", "", fmt.Errorf("unsupported URL format, expected https://github.com/owner/repo or github://owner/repo")
}

// extractTarballToMap extracts a tar.gz stream into a map of file paths to contents.
// It strips the top-level directory prefix from GitHub tarballs (e.g., "owner-repo-hash/file.txt" → "file.txt").
// Folder structure is preserved: files in subdirectories will have paths like "subdir/file.txt".
func extractTarballToMap(r io.Reader) (map[string][]byte, string, error) {
	files := make(map[string][]byte)
	topLevelDir := ""

	gzr, err := gzip.NewReader(r)
	if err != nil {
		return nil, "", fmt.Errorf("gzip read: %w", err)
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, "", fmt.Errorf("tar read: %w", err)
		}

		// Detect top-level directory, skipping metadata files
		if topLevelDir == "" && header.Name != "pax_global_header" {
			parts := strings.SplitN(header.Name, "/", 2)
			if len(parts) > 0 {
				topLevelDir = parts[0]
			}
		}

		if header.Typeflag == tar.TypeDir {
			continue
		}

		// Skip tar metadata files and .git directories
		if header.Name == "pax_global_header" || strings.Contains(header.Name, "/.git/") || strings.HasPrefix(header.Name, ".git/") {
			continue
		}

		content, err := io.ReadAll(tr)
		if err != nil {
			return nil, "", fmt.Errorf("read file %q: %w", header.Name, err)
		}

		// Strip the top-level directory prefix (e.g., "owner-repo-hash/file.txt" → "file.txt")
		filePath := header.Name
		if topLevelDir != "" && strings.HasPrefix(filePath, topLevelDir+"/") {
			filePath = strings.TrimPrefix(filePath, topLevelDir+"/")
		}

		normalizedPath := fileutils.Normalize(filePath)

		files[normalizedPath] = content
	}

	return files, topLevelDir, nil
}

// PopulateStateFromSnapshot adds files from a RepoSnapshot to a VirtualProjectState.
// It creates the directory structure as found in the repository.
func PopulateStateFromSnapshot(state interface {
	CreateFile(path string, content []byte) error
}, snapshot *RepoSnapshot) error {
	if snapshot == nil || len(snapshot.Files) == 0 {
		return fmt.Errorf("empty snapshot")
	}

	for path, content := range snapshot.Files {
		if err := state.CreateFile(path, content); err != nil {
			return fmt.Errorf("populate state: create file %q: %w", path, err)
		}
	}

	return nil
}
