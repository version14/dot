package storybooksetup

import (
	"embed"

	"github.com/version14/dot/internal/deps"
	"github.com/version14/dot/internal/render"
	"github.com/version14/dot/internal/state"
	"github.com/version14/dot/pkg/dotapi"
)

//go:embed all:files
var fs embed.FS

type storybookData struct {
	FrameworkPkg string
}

type Generator struct{}

func New() *Generator { return &Generator{} }

func (g *Generator) Name() string    { return Manifest.Name }
func (g *Generator) Version() string { return Manifest.Version }

func (g *Generator) Generate(ctx *dotapi.Context) error {
	framework, _ := ctx.Answers["framework"].(string)

	frameworkPkg := "@storybook/react-vite"
	if framework == "next" {
		frameworkPkg = "@storybook/nextjs"
	}

	if err := ctx.State.UpdateJSON("package.json", func(d *state.JSONDoc) error {
		d.Merge(map[string]interface{}{
			"scripts": map[string]interface{}{
				"storybook":       "storybook dev -p 6006",
				"build-storybook": "storybook build",
			},
			"devDependencies": deps.NPM("storybook", "@storybook/react", frameworkPkg),
		})
		return nil
	}); err != nil {
		return err
	}

	return render.NewLocalFolderRenderer(ctx.State).Render(fs, storybookData{FrameworkPkg: frameworkPkg})
}
