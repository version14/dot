package render

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/version14/dot/internal/state"
)

// Renderer defines the interface for rendering files.
type Renderer interface {
	Render(srcPath string, rootPath string, data interface{}) error
}

// LocalFolderRenderer renders a local folder to a VirtualProjectState.
type LocalFolderRenderer struct {
	State *state.VirtualProjectState
}

// NewLocalFolderRenderer creates a new LocalFolderRenderer.
func NewLocalFolderRenderer(s *state.VirtualProjectState) *LocalFolderRenderer {
	return &LocalFolderRenderer{State: s}
}

// Render processes a source directory and outputs to the virtual state.
func (r *LocalFolderRenderer) Render(embed embed.FS, data interface{}) error {
	return fs.WalkDir(embed, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip the root source directory itself.
		if path == "." {
			return nil
		}

		// Determine the destination path relative to the source directory.
		part1Path := strings.Split(path, "/")
		destPath := strings.Join(part1Path[1:], "/")

		fmt.Println(destPath)

		if d.IsDir() {
			return nil
		}

		// Read the file content.
		content, err := fs.ReadFile(embed, path)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", path, err)
		}

		// If it's a template file, render it.
		if strings.HasSuffix(path, ".tmpl") {
			tmpl, err := template.New(filepath.Base(path)).Parse(string(content))
			if err != nil {
				return fmt.Errorf("failed to parse template %s: %w", path, err)
			}

			var renderedContent bytes.Buffer
			if err := tmpl.Execute(&renderedContent, data); err != nil {
				return fmt.Errorf("failed to execute template %s: %w", path, err)
			}

			// Remove .tmpl extension from destination path.
			destPath = strings.TrimSuffix(destPath, ".tmpl")
			r.State.WriteFile(destPath, renderedContent.Bytes(), state.ContentRaw)
		} else {
			// Otherwise, copy the file as is.
			r.State.WriteFile(destPath, content, state.ContentRaw)
		}

		return nil
	})
}
