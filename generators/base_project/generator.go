package baseproject

import (
	"fmt"

	"github.com/version14/dot/internal/render"
	"github.com/version14/dot/internal/state"
	"github.com/version14/dot/pkg/dotapi"
)

// Generator implements dotapi.Generator for base_project. It writes universal
// project files: README.md, .gitignore, and LICENSE. It must run first because
// later generators may modify these (e.g. typescript_base appends node_modules
// to .gitignore).
type Generator struct{}

// New constructs a Generator instance for registration.
func New() *Generator { return &Generator{} }

func (g *Generator) Name() string    { return Manifest.Name }
func (g *Generator) Version() string { return Manifest.Version }

func (g *Generator) Generate(ctx *dotapi.Context) error {
	projectName, _ := ctx.Answers["project_name"].(string)
	if projectName == "" {
		projectName = ctx.Spec.Metadata.ProjectName
	}
	if projectName == "" {
		projectName = "my-project"
	}

	data := map[string]interface{}{
		"ProjectName": projectName,
	}

	if err := writeReadme(ctx.State, data); err != nil {
		return err
	}
	if err := writeGitignore(ctx.State); err != nil {
		return err
	}
	if err := writeLicense(ctx.State, data); err != nil {
		return err
	}
	return nil
}

func writeReadme(s *state.VirtualProjectState, data map[string]interface{}) error {
	out, err := render.Render(readmeTmpl, data)
	if err != nil {
		return fmt.Errorf("base_project: render README: %w", err)
	}
	s.WriteFile("README.md", out, state.ContentRaw)
	return nil
}

func writeGitignore(s *state.VirtualProjectState) error {
	s.WriteFile(".gitignore", []byte(gitignoreContent), state.ContentRaw)
	return nil
}

func writeLicense(s *state.VirtualProjectState, data map[string]interface{}) error {
	out, err := render.Render(licenseTmpl, data)
	if err != nil {
		return fmt.Errorf("base_project: render LICENSE: %w", err)
	}
	s.WriteFile("LICENSE", out, state.ContentRaw)
	return nil
}

const readmeTmpl = `# {{.ProjectName}}

> Scaffolded with [DOT](https://github.com/version14/dot).

## Getting Started

See language-specific instructions in the relevant subdirectories.
`

const gitignoreContent = `# Editor / OS
.DS_Store
*.swp
.idea/
.vscode/

# Build artifacts
dist/
build/
out/
*.log

# Dependencies
node_modules/
vendor/
`

const licenseTmpl = `MIT License

Copyright (c) {{.ProjectName}}

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction.
`
