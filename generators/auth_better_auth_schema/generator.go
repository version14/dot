package authbetterauthschema

import (
	"embed"
	"strings"

	"github.com/version14/dot/internal/render"
	"github.com/version14/dot/internal/state"
	"github.com/version14/dot/pkg/dotapi"
)

type Generator struct{}

func New() *Generator { return &Generator{} }

func (g *Generator) Name() string    { return Manifest.Name }
func (g *Generator) Version() string { return Manifest.Version }

//go:embed all:files
var fs embed.FS

func (g *Generator) Generate(ctx *dotapi.Context) error {
	renderer := render.NewLocalFolderRenderer(ctx.State)
	if err := renderer.Render(fs, nil); err != nil {
		return err
	}

	existing := "export {};\n"
	if f, ok := ctx.State.GetFile("src/db/schema/index.ts"); ok {
		existing = string(f.Content)
	}
	if strings.TrimSpace(existing) == "export {};" {
		existing = ""
	}
	updated := existing + "export * from './auth.schema';\n"
	ctx.State.WriteFile("src/db/schema/index.ts", []byte(updated), state.ContentRaw)
	return nil
}
