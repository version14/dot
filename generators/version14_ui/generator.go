package version14ui

import (
	"embed"
	"slices"
	"strings"

	"github.com/version14/dot/internal/render"
	"github.com/version14/dot/internal/state"
	"github.com/version14/dot/pkg/dotapi"
)

//go:embed all:files
var fs embed.FS

type Generator struct{}

func New() *Generator { return &Generator{} }

func (g *Generator) Name() string    { return Manifest.Name }
func (g *Generator) Version() string { return Manifest.Version }

func (g *Generator) Generate(ctx *dotapi.Context) error {
	if err := ctx.State.UpdateJSON("package.json", func(d *state.JSONDoc) error {
		d.Merge(map[string]interface{}{
			"dependencies": map[string]interface{}{
				"@version14/ui": "0.8.0",
				"@ark-ui/react": "^5.37.2",
			},
		})
		return nil
	}); err != nil {
		return err
	}

	if err := render.NewLocalFolderRenderer(ctx.State).Render(fs, nil); err != nil {
		return err
	}

	// Component styles are pre-generated CSS shipped by the package itself —
	// required regardless of which styling engine (Tailwind, CSS Modules,
	// Panda) the project uses, so this is wired here rather than left to the
	// styling generator.
	injectStylesImport(ctx.State, ctx.PreviousGens)

	return nil
}

func injectStylesImport(s *state.VirtualProjectState, previousGens []string) {
	switch {
	case slices.Contains(previousGens, "react_app"):
		injectImport(s, "src/main.tsx",
			`import React from "react";`,
			"import React from \"react\";\nimport \"@version14/ui/styles.css\";")
	case slices.Contains(previousGens, "nextjs_base"):
		injectImport(s, "src/app/layout.tsx",
			`import "./globals.css";`,
			"import \"./globals.css\";\nimport \"@version14/ui/styles.css\";")
	}
}

func injectImport(s *state.VirtualProjectState, path, marker, replacement string) {
	f, ok := s.GetFile(path)
	if !ok {
		return
	}
	content := string(f.Content)
	if strings.Contains(content, "@version14/ui/styles.css") {
		return
	}
	content = strings.Replace(content, marker, replacement, 1)
	s.WriteFile(path, []byte(content), state.ContentRaw)
}
