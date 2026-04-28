// Package backendArchitectureMVC scaffolds a publishable DOT mvc architecture.
package backendArchitectureMVC

import (
	"embed"

	"github.com/version14/dot/internal/render"
	"github.com/version14/dot/pkg/dotapi"
)

type Generator struct{}

func New() *Generator { return &Generator{} }

func (g *Generator) Name() string    { return Manifest.Name }
func (g *Generator) Version() string { return Manifest.Version }

//go:embed all:files
var fs embed.FS

// Generate builds the plugin repo skeleton from the populated answers.
func (g *Generator) Generate(ctx *dotapi.Context) error {
	renderer := render.NewLocalFolderRenderer(ctx.State)
	return renderer.Render(fs, nil)
}
