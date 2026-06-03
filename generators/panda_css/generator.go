package pandacss

import (
	"embed"
	"slices"
	"strings"

	"github.com/version14/dot/internal/render"
	"github.com/version14/dot/internal/state"
	"github.com/version14/dot/pkg/dotapi"
)

//go:embed all:files
var filesFS embed.FS

type Generator struct{}

func New() *Generator { return &Generator{} }

func (g *Generator) Name() string    { return Manifest.Name }
func (g *Generator) Version() string { return Manifest.Version }

func (g *Generator) Generate(ctx *dotapi.Context) error {
	linter, _ := ctx.Answers["linter"].(string)
	if linter == "" {
		linter, _ = ctx.Answers["frontend-linter"].(string)
	}

	if err := ctx.State.UpdateJSON("package.json", func(d *state.JSONDoc) error {
		d.Merge(map[string]interface{}{
			"devDependencies": map[string]interface{}{
				"@pandacss/dev": "^1.11.1",
			},
		})
		prepare := "panda codegen"
		if existing, ok := d.GetNested("scripts.prepare"); ok {
			if s, isStr := existing.(string); isStr && s != "" {
				prepare = s + " && panda codegen"
			}
		}
		return d.SetNested("scripts.prepare", prepare)
	}); err != nil {
		return err
	}

	existing := ""
	if f, ok := ctx.State.GetFile(".prettierignore"); ok {
		existing = string(f.Content)
		if len(existing) > 0 && existing[len(existing)-1] != '\n' {
			existing += "\n"
		}
	}
	ctx.State.WriteFile(".prettierignore", []byte(existing+"styled-system/\n"), state.ContentRaw)

	if linter == "biome" {
		if err := ctx.State.UpdateJSON("biome.json", func(d *state.JSONDoc) error {
			return d.AppendStringSet("files.includes", "!**/styled-system")
		}); err != nil {
			return err
		}
	}

	if err := ctx.State.UpdateJSON("tsconfig.json", func(d *state.JSONDoc) error {
		return d.SetNested("compilerOptions.paths.@styled-system/*", []interface{}{"./styled-system/*"})
	}); err != nil {
		return err
	}

	if err := render.NewLocalFolderRenderer(ctx.State).Render(filesFS, nil); err != nil {
		return err
	}

	// Inject global CSS import into main.tsx when the frontend generator already ran.
	hasFrontend := slices.Contains(ctx.PreviousGens, "tanstack_router") ||
		slices.Contains(ctx.PreviousGens, "react_app")
	if hasFrontend {
		injectGlobalCSS(ctx.State)
	}

	injectViteAlias(ctx.State)

	if slices.Contains(ctx.PreviousGens, "version14_ui") {
		injectV14Fonts(ctx.State)
		injectV14Preset(ctx.State)
	}

	return nil
}

func injectViteAlias(s *state.VirtualProjectState) {
	f, ok := s.GetFile("vite.config.ts")
	if !ok {
		return
	}
	content := string(f.Content)
	if strings.Contains(content, "@styled-system") {
		return
	}
	content = strings.Replace(content,
		`"@": path.resolve(__dirname, "src"),`,
		`"@": path.resolve(__dirname, "src"),
      "@styled-system": path.resolve(__dirname, "styled-system"),`,
		1)
	s.WriteFile("vite.config.ts", []byte(content), state.ContentRaw)
}

func injectV14Fonts(s *state.VirtualProjectState) {
	f, ok := s.GetFile("src/styles/global.css")
	if !ok {
		return
	}
	content := string(f.Content)
	if strings.Contains(content, "@version14/ui/fonts.css") {
		return
	}
	s.WriteFile("src/styles/global.css", []byte("@import \"@version14/ui/fonts.css\";\n"+content), state.ContentRaw)
}

func injectV14Preset(s *state.VirtualProjectState) {
	f, ok := s.GetFile("panda.config.ts")
	if !ok {
		return
	}
	content := string(f.Content)
	if strings.Contains(content, "v14Preset") {
		return
	}
	content = strings.Replace(content,
		`import { defineConfig } from "@pandacss/dev";`,
		`import { v14Preset } from "@version14/ui/preset";
import { defineConfig } from "@pandacss/dev";`,
		1)
	content = strings.Replace(content,
		"  jsxFactory: \"react\"\n});",
		"  jsxFactory: \"react\",\n  presets: [v14Preset]\n});",
		1)
	s.WriteFile("panda.config.ts", []byte(content), state.ContentRaw)
}

func injectGlobalCSS(s *state.VirtualProjectState) {
	f, ok := s.GetFile("src/main.tsx")
	if !ok {
		return
	}
	content := string(f.Content)
	if strings.Contains(content, "styles/global.css") {
		return
	}
	content = strings.Replace(content,
		`import React from "react";`,
		"import React from \"react\";\nimport \"./styles/global.css\";",
		1)
	s.WriteFile("src/main.tsx", []byte(content), state.ContentRaw)
}
