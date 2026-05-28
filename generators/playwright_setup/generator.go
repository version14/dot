package playwrightsetup

import (
	"embed"

	"github.com/version14/dot/internal/render"
	"github.com/version14/dot/internal/state"
	"github.com/version14/dot/pkg/dotapi"
)

//go:embed all:files
var fs embed.FS

type playwrightData struct {
	DevServerURL     string
	DevServerCommand string
}

type Generator struct{}

func New() *Generator { return &Generator{} }

func (g *Generator) Name() string    { return Manifest.Name }
func (g *Generator) Version() string { return Manifest.Version }

func (g *Generator) Generate(ctx *dotapi.Context) error {
	framework, _ := ctx.Answers["framework"].(string)

	devServerURL := "http://localhost:5173"
	devServerCommand := "pnpm exec vite --host 127.0.0.1"
	if framework == "next" {
		devServerURL = "http://localhost:3000"
		devServerCommand = "pnpm exec next start"
	}

	if err := ctx.State.UpdateJSON("package.json", func(d *state.JSONDoc) error {
		d.Merge(map[string]interface{}{
			"scripts": map[string]interface{}{
				"test:e2e": "playwright test",
			},
			"devDependencies": map[string]interface{}{
				"@playwright/test": "^1.44.0",
			},
		})
		return nil
	}); err != nil {
		return err
	}

	return render.NewLocalFolderRenderer(ctx.State).Render(fs, playwrightData{
		DevServerURL:     devServerURL,
		DevServerCommand: devServerCommand,
	})
}
