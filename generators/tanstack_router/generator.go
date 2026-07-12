package tanstackrouter

import (
	"embed"
	"slices"
	"strings"

	"github.com/version14/dot/internal/deps"
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
			"dependencies":    deps.NPM("@tanstack/react-router"),
			"devDependencies": deps.NPM("@tanstack/router-plugin"),
		})
		return nil
	}); err != nil {
		return err
	}

	if err := render.NewLocalFolderRenderer(ctx.State).Render(fs, nil); err != nil {
		return err
	}

	// Inject global CSS import when panda_css already ran (and wrote the CSS file).
	if slices.Contains(ctx.PreviousGens, "panda_css") {
		if f, ok := ctx.State.GetFile("src/main.tsx"); ok {
			content := string(f.Content)
			if !strings.Contains(content, "styles/global.css") {
				content = strings.Replace(content,
					`import React from "react";`,
					"import React from \"react\";\nimport \"./styles/global.css\";",
					1)
				ctx.State.WriteFile("src/main.tsx", []byte(content), state.ContentRaw)
			}
		}
	}

	return nil
}
