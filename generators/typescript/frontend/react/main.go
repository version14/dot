package frontend_react_generator

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"strings"
	"text/template"

	frontend_architecture_generator "github.com/version14/dot/generators/common/frontend/architecture"
	"github.com/version14/dot/internal/generator"
	newtest_generator "github.com/version14/dot/internal/scaffold"
	"github.com/version14/dot/internal/spec"
)

// all: is required to include dot-files (.gitkeep) in the embedded FS.
//
//go:embed all:files
var reactFiles embed.FS

// ReactTS is the generator for React + TypeScript + Vite projects.
// Wire it into a question with ReactTS.Func():
//
//	question.Select(...).ChoiceWithGen("React", "react", ReactTS.Func())
var ReactTS = &newtest_generator.Generator{
	Name:     "react-ts",
	Language: "typescript",

	ApplyFunction: func(s spec.Spec) ([]generator.FileOp, error) {
		arch, _ := s.Extensions["architecture"].(string)
		data := struct {
			Name         string
			Architecture string
		}{s.Project.Name, arch}

		ops, err := renderDir(reactFiles, "files", "", data)
		if err != nil {
			return nil, fmt.Errorf("react-ts base: %w", err)
		}

		archOps, err := frontend_architecture_generator.ArchitectureTS.Apply(s)
		if err != nil {
			return nil, fmt.Errorf("react-ts arch: %w", err)
		}

		return append(ops, archOps...), nil
	},

	PostApplyFunction: func(s spec.Spec) []generator.PostOp {
		return []generator.PostOp{
			{Command: "pnpm", Args: []string{"install"}, Dir: ".", Generator: "react-ts"},
		}
	},
}

// renderDir walks fsys under root, skips the skipDir subtree, renders .tmpl
// files with data, and returns a FileOp per file.
func renderDir(fsys embed.FS, root, skipDir string, data any) ([]generator.FileOp, error) {
	var ops []generator.FileOp

	err := fs.WalkDir(fsys, root, func(path string, dir fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if dir.IsDir() {
			if skipDir != "" && strings.HasSuffix(path, "/"+skipDir) {
				return fs.SkipDir
			}
			return nil
		}

		content, err := fsys.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %s: %w", path, err)
		}

		outPath := strings.TrimPrefix(path, root+"/")

		if strings.HasSuffix(outPath, ".tmpl") {
			rendered, err := renderTemplate(outPath, string(content), data)
			if err != nil {
				return fmt.Errorf("render %s: %w", path, err)
			}
			outPath = strings.TrimSuffix(outPath, ".tmpl")
			ops = append(ops, fileOp(outPath, rendered))
		} else {
			ops = append(ops, fileOp(outPath, string(content)))
		}

		return nil
	})

	return ops, err
}

func renderTemplate(name, templateString string, data any) (string, error) {
	template, err := template.New(name).Parse(templateString)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := template.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func fileOp(path, content string) generator.FileOp {
	return generator.FileOp{
		Kind:      generator.Create,
		Path:      path,
		Content:   content,
		Generator: "react-ts",
		Priority:  100,
	}
}
