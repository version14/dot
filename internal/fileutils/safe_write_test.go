package fileutils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSafeWrite_CreatesAndOverwrites(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nested", "file.txt")

	if err := SafeWrite(path, []byte("hello"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if string(got) != "hello" {
		t.Errorf("content = %q", got)
	}

	if err := SafeWrite(path, []byte("world"), 0o644); err != nil {
		t.Fatalf("overwrite: %v", err)
	}
	got, _ = os.ReadFile(path)
	if string(got) != "world" {
		t.Errorf("overwrite content = %q", got)
	}
}

func TestNormalizePath(t *testing.T) {
	cases := map[string]string{
		"./foo":      "foo",
		"foo/bar":    "foo/bar",
		"foo/../bar": "bar",
		"foo/./bar/": "foo/bar",
	}
	for in, want := range cases {
		if got := Normalize(in); got != want {
			t.Errorf("Normalize(%q) = %q, want %q", in, got, want)
		}
	}
}
