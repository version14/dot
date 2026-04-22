// Package genfs provides shared helpers for generators that embed file trees.
package genfs

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"strings"
	"text/template"

	"github.com/version14/dot/internal/generator"
)

const defaultPriority = 100

// WalkDir walks every file under root in fsys and returns one Create FileOp
// per file. Output paths are relative to root (root prefix is stripped).
// Use this for static file trees that need no rendering.
func WalkDir(fsys embed.FS, root, generatorName string) ([]generator.FileOp, error) {
	var ops []generator.FileOp

	err := fs.WalkDir(fsys, root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		content, err := fsys.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %s: %w", path, err)
		}
		outPath := strings.TrimPrefix(path, root+"/")
		ops = append(ops, makeOp(outPath, string(content), generatorName))
		return nil
	})

	return ops, err
}

// RenderDir walks every file under root in fsys. Files ending in ".tmpl" are
// rendered with data via text/template and the extension is stripped from the
// output path. All other files are copied verbatim.
func RenderDir(fsys embed.FS, root, generatorName string, data any) ([]generator.FileOp, error) {
	var ops []generator.FileOp

	err := fs.WalkDir(fsys, root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		content, err := fsys.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %s: %w", path, err)
		}

		outPath := strings.TrimPrefix(path, root+"/")

		if strings.HasSuffix(outPath, ".tmpl") {
			rendered, err := RenderTemplate(outPath, string(content), data)
			if err != nil {
				return fmt.Errorf("render %s: %w", path, err)
			}
			outPath = strings.TrimSuffix(outPath, ".tmpl")
			ops = append(ops, makeOp(outPath, rendered, generatorName))
		} else {
			ops = append(ops, makeOp(outPath, string(content), generatorName))
		}

		return nil
	})

	return ops, err
}

// RenderTemplate executes a text/template string against data.
func RenderTemplate(name, tmpl string, data any) (string, error) {
	t, err := template.New(name).Parse(tmpl)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func makeOp(path, content, generatorName string) generator.FileOp {
	return generator.FileOp{
		Kind:      generator.Create,
		Path:      path,
		Content:   content,
		Generator: generatorName,
		Priority:  defaultPriority,
	}
}
