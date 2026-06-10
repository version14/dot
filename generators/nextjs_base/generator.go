package nextjsbase

import (
	"embed"

	"github.com/version14/dot/internal/render"
	"github.com/version14/dot/internal/state"
	"github.com/version14/dot/pkg/dotapi"
)

//go:embed all:files
var fs embed.FS

type nextjsData struct {
	ProjectName string
}

type Generator struct{}

func New() *Generator { return &Generator{} }

func (g *Generator) Name() string    { return Manifest.Name }
func (g *Generator) Version() string { return Manifest.Version }

func (g *Generator) Generate(ctx *dotapi.Context) error {
	projectName, _ := ctx.Answers["project_name"].(string)
	if projectName == "" {
		projectName = ctx.Spec.Metadata.ProjectName
	}
	if projectName == "" {
		projectName = "app"
	}

	if err := ctx.State.UpdateJSON("package.json", func(d *state.JSONDoc) error {
		d.Merge(map[string]interface{}{
			"scripts": map[string]interface{}{
				"dev":   "next dev --turbopack",
				"build": "next build",
				"start": "next start",
				"lint":  "next lint",
			},
			"dependencies": map[string]interface{}{
				"next":      "^16.2.9",
				"react":     "^19.2.7",
				"react-dom": "^19.2.7",
			},
			"devDependencies": map[string]interface{}{
				"@types/node":      "^25.9.2",
				"@types/react":     "^19.2.17",
				"@types/react-dom": "^19.2.3",
			},
		})
		return nil
	}); err != nil {
		return err
	}

	if err := ctx.State.UpdateJSON("tsconfig.json", func(d *state.JSONDoc) error {
		_ = d.SetNested("compilerOptions.jsx", "preserve")
		_ = d.SetNested("compilerOptions.lib", []interface{}{"DOM", "DOM.Iterable", "ES2022"})
		_ = d.SetNested("compilerOptions.plugins", []interface{}{map[string]interface{}{"name": "next"}})
		_ = d.SetNested("compilerOptions.paths", map[string]interface{}{"@/*": []interface{}{"./src/*"}})
		_ = d.SetNested("compilerOptions.allowJs", true)
		_ = d.SetNested("compilerOptions.noEmit", true)
		_ = d.SetNested("compilerOptions.incremental", true)
		_ = d.SetNested("compilerOptions.resolveJsonModule", true)
		_ = d.SetNested("compilerOptions.isolatedModules", true)
		_ = d.SetNested("exclude", []interface{}{"node_modules"})
		_ = d.SetNested("include", []interface{}{"src", ".next/types/**/*.ts"})
		return nil
	}); err != nil {
		return err
	}

	return render.NewLocalFolderRenderer(ctx.State).Render(fs, nextjsData{ProjectName: projectName})
}
