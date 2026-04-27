package render

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"testing"
)

func TestParseGitHubURL(t *testing.T) {
	tests := []struct {
		url         string
		wantOwner   string
		wantRepo    string
		wantErr     bool
		description string
	}{
		{
			url:         "https://github.com/mathieusouflis/github-template",
			wantOwner:   "mathieusouflis",
			wantRepo:    "github-template",
			wantErr:     false,
			description: "HTTPS URL without .git suffix",
		},
		{
			url:         "https://github.com/mathieusouflis/github-template.git",
			wantOwner:   "mathieusouflis",
			wantRepo:    "github-template",
			wantErr:     false,
			description: "HTTPS URL with .git suffix",
		},
		{
			url:         "http://github.com/owner/repo.git",
			wantOwner:   "owner",
			wantRepo:    "repo",
			wantErr:     false,
			description: "HTTP URL with .git suffix",
		},
		{
			url:         "github://owner/repo",
			wantOwner:   "owner",
			wantRepo:    "repo",
			wantErr:     false,
			description: "github:// URL format",
		},
		{
			url:         "invalid://owner/repo",
			wantOwner:   "",
			wantRepo:    "",
			wantErr:     true,
			description: "Unsupported URL scheme",
		},
		{
			url:         "https://github.com/onlyowner",
			wantOwner:   "",
			wantRepo:    "",
			wantErr:     true,
			description: "Invalid format - missing repo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			owner, repo, err := parseGitHubURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseGitHubURL(%q) error = %v, wantErr %v", tt.url, err, tt.wantErr)
			}
			if owner != tt.wantOwner {
				t.Errorf("parseGitHubURL(%q) owner = %q, want %q", tt.url, owner, tt.wantOwner)
			}
			if repo != tt.wantRepo {
				t.Errorf("parseGitHubURL(%q) repo = %q, want %q", tt.url, repo, tt.wantRepo)
			}
		})
	}
}

func TestExtractTarballToMap(t *testing.T) {
	testTarball := createTestTarball()

	files, topLevelDir, err := extractTarballToMap(bytes.NewReader(testTarball))
	if err != nil {
		t.Fatalf("extractTarballToMap failed: %v", err)
	}

	if topLevelDir != "owner-repo-abc123" {
		t.Errorf("topLevelDir = %q, want %q", topLevelDir, "owner-repo-abc123")
	}

	expectedFiles := map[string]string{
		"README.md":           "# Test Repository",
		"src/main.go":         "package main",
		"src/utils/helper.go": "package utils",
		".gitignore":          "*.o",
	}

	for path, expectedContent := range expectedFiles {
		content, ok := files[path]
		if !ok {
			t.Errorf("expected file %q not found in extracted files", path)
			continue
		}
		if string(content) != expectedContent {
			t.Errorf("file %q content = %q, want %q", path, string(content), expectedContent)
		}
	}

	if len(files) != len(expectedFiles) {
		t.Errorf("extractTarballToMap returned %d files, want %d", len(files), len(expectedFiles))
	}
}

func TestFolderStructurePreservation(t *testing.T) {
	testTarball := createTestTarball()

	files, _, err := extractTarballToMap(bytes.NewReader(testTarball))
	if err != nil {
		t.Fatalf("extractTarballToMap failed: %v", err)
	}

	if _, ok := files["src/main.go"]; !ok {
		t.Error("folder structure not preserved: src/main.go not found")
	}

	if _, ok := files["src/utils/helper.go"]; !ok {
		t.Error("nested folder structure not preserved: src/utils/helper.go not found")
	}
}

func TestPaxGlobalHeaderFiltered(t *testing.T) {
	testTarball := createTestTarballWithMetadata()

	files, _, err := extractTarballToMap(bytes.NewReader(testTarball))
	if err != nil {
		t.Fatalf("extractTarballToMap failed: %v", err)
	}

	if _, ok := files["pax_global_header"]; ok {
		t.Error("pax_global_header should be filtered out")
	}
}

func TestPopulateStateFromSnapshot(t *testing.T) {
	snapshot := &RepoSnapshot{
		Files: map[string][]byte{
			"file1.txt":     []byte("content1"),
			"dir/file2.txt": []byte("content2"),
		},
		Ref:       "main",
		SourceURL: "https://github.com/test/repo",
	}

	state := &mockState{files: make(map[string][]byte)}

	err := PopulateStateFromSnapshot(state, snapshot)
	if err != nil {
		t.Fatalf("PopulateStateFromSnapshot failed: %v", err)
	}

	if len(state.files) != 2 {
		t.Errorf("expected 2 files in state, got %d", len(state.files))
	}

	if content, ok := state.files["file1.txt"]; !ok || string(content) != "content1" {
		t.Errorf("file1.txt not correctly populated")
	}

	if content, ok := state.files["dir/file2.txt"]; !ok || string(content) != "content2" {
		t.Errorf("dir/file2.txt not correctly populated")
	}
}

type mockState struct {
	files map[string][]byte
}

func (m *mockState) CreateFile(path string, content []byte) error {
	if _, exists := m.files[path]; exists {
		return nil
	}
	m.files[path] = content
	return nil
}

func TestGitHubArchiveFetcherURLConstruction(t *testing.T) {
	fetcher := NewGitHubArchiveFetcher()
	if fetcher == nil {
		t.Error("NewGitHubArchiveFetcher returned nil")
	}
	if fetcher.Client == nil {
		t.Error("GitHubArchiveFetcher Client is nil")
	}
}

func createTestTarball() []byte {
	var buf bytes.Buffer

	gzipWriter, _ := gzip.NewWriterLevel(&buf, gzip.DefaultCompression)
	tarWriter := tar.NewWriter(gzipWriter)

	files := []struct {
		name    string
		content string
	}{
		{"owner-repo-abc123/README.md", "# Test Repository"},
		{"owner-repo-abc123/src/main.go", "package main"},
		{"owner-repo-abc123/src/utils/helper.go", "package utils"},
		{"owner-repo-abc123/.gitignore", "*.o"},
	}

	for _, f := range files {
		header := &tar.Header{
			Name: f.name,
			Size: int64(len(f.content)),
		}
		if f.name[len(f.name)-1] != '/' {
			header.Typeflag = tar.TypeReg
		}
		tarWriter.WriteHeader(header)
		tarWriter.Write([]byte(f.content))
	}

	tarWriter.Close()
	gzipWriter.Close()

	return buf.Bytes()
}

func createTestTarballWithMetadata() []byte {
	var buf bytes.Buffer

	gzipWriter, _ := gzip.NewWriterLevel(&buf, gzip.DefaultCompression)
	tarWriter := tar.NewWriter(gzipWriter)

	files := []struct {
		name    string
		content string
	}{
		{"pax_global_header", ""},
		{"owner-repo-abc123/README.md", "# Test Repository"},
		{"owner-repo-abc123/src/main.go", "package main"},
	}

	for _, f := range files {
		header := &tar.Header{
			Name: f.name,
			Size: int64(len(f.content)),
		}
		if f.name[len(f.name)-1] != '/' {
			header.Typeflag = tar.TypeReg
		}
		tarWriter.WriteHeader(header)
		tarWriter.Write([]byte(f.content))
	}

	tarWriter.Close()
	gzipWriter.Close()

	return buf.Bytes()
}
