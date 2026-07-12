package featureflagsvercel

import (
	"embed"

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
			"dependencies": deps.NPM("@vercel/edge-config", "@vercel/flags"),
		})
		return nil
	}); err != nil {
		return err
	}

	if err := render.NewLocalFolderRenderer(ctx.State).Render(fs, nil); err != nil {
		return err
	}
	ctx.State.AppendFile(".env.example", []byte("# Vercel Edge Config\nEDGE_CONFIG=your_edge_config_connection_string\n"))
	return nil
}
