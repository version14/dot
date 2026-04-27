package fileutils

import (
	"path/filepath"
	"strings"
)

// Normalize returns a clean, forward-slash path suitable for use as a
// VirtualProjectState key. Use Resolve when interacting with the OS.
func Normalize(p string) string {
	p = filepath.ToSlash(p)
	p = filepath.Clean(p)
	p = strings.TrimPrefix(p, "./")
	return p
}

// Join joins segments using the OS separator after normalizing each.
func Join(segments ...string) string {
	return filepath.Join(segments...)
}

// Resolve joins root with a (possibly forward-slash) relative path and
// returns an OS-native absolute path under root. Returns an empty string
// when the resulting path escapes root.
func Resolve(root, rel string) string {
	rel = Normalize(rel)
	full := filepath.Join(root, filepath.FromSlash(rel))
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return ""
	}
	absFull, err := filepath.Abs(full)
	if err != nil {
		return ""
	}
	if !strings.HasPrefix(absFull+string(filepath.Separator), absRoot+string(filepath.Separator)) && absFull != absRoot {
		return ""
	}
	return full
}
