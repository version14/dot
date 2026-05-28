package reactapp

import (
	"embed"

	"github.com/version14/dot/internal/render"
	"github.com/version14/dot/internal/state"
	"github.com/version14/dot/pkg/dotapi"
)

//go:embed all:files
var filesFS embed.FS

type reactAppData struct {
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
				"dev":     "vite",
				"build":   "tsc && vite build",
				"preview": "vite preview",
			},
			"dependencies": map[string]interface{}{
				"react":     "^19.2.6",
				"react-dom": "^19.2.6",
			},
			"devDependencies": map[string]interface{}{
				"@types/react":         "^19.2.6",
				"@types/react-dom":     "^19.2.3",
				"@vitejs/plugin-react": "^6.0.2",
				"vite":                 "^8.0.14",
			},
		})
		return nil
	}); err != nil {
		return err
	}

	if err := ctx.State.UpdateJSON("tsconfig.json", func(d *state.JSONDoc) error {
		_ = d.SetNested("compilerOptions.jsx", "react-jsx")
		_ = d.SetNested("compilerOptions.lib", []interface{}{"DOM", "DOM.Iterable", "ES2022"})
		return nil
	}); err != nil {
		return err
	}

	return render.NewLocalFolderRenderer(ctx.State).Render(filesFS, reactAppData{ProjectName: projectName})
}
