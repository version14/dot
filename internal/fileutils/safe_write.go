package fileutils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// SafeWrite writes data to path atomically by writing to a sibling temp file
// then renaming. Falls back to copy+remove when rename crosses devices.
func SafeWrite(path string, data []byte, mode os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("fileutils: mkdir %s: %w", filepath.Dir(path), err)
	}

	tmp, err := os.CreateTemp(filepath.Dir(path), ".dot-tmp-*")
	if err != nil {
		return fmt.Errorf("fileutils: create temp: %w", err)
	}
	tmpPath := tmp.Name()

	cleanup := func() { _ = os.Remove(tmpPath) }

	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		cleanup()
		return fmt.Errorf("fileutils: write %s: %w", tmpPath, err)
	}
	if err := tmp.Chmod(mode); err != nil {
		_ = tmp.Close()
		cleanup()
		return fmt.Errorf("fileutils: chmod %s: %w", tmpPath, err)
	}
	if err := tmp.Close(); err != nil {
		cleanup()
		return fmt.Errorf("fileutils: close %s: %w", tmpPath, err)
	}

	if err := os.Rename(tmpPath, path); err != nil {
		if copyErr := copyFile(tmpPath, path, mode); copyErr != nil {
			cleanup()
			return fmt.Errorf("fileutils: rename %s -> %s: %w", tmpPath, path, err)
		}
		cleanup()
	}
	return nil
}

func copyFile(src, dst string, mode os.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		_ = out.Close()
		return err
	}
	return out.Close()
}
