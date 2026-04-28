// Package pluginreposkeleton scaffolds a publishable DOT plugin repository.
//
// The output is a complete Go module (go.mod + plugin.go + plugin.json +
// README + LICENSE + .gitignore) that someone else can install with
//
//	dot plugin install github.com/<your-author>/<plugin-id>
//
// Optional sample injection / sample generator code is included based on
// the user's answers, so the generated repo is either a minimal skeleton
// or a fully working example.
package backendArchitectureCleanArchitecture

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
