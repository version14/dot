package authjwtusersschema

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

	existing := ""
	if f, ok := ctx.State.GetFile("src/db/schema/index.ts"); ok {
		existing = string(f.Content)
	}
	existing = strings.ReplaceAll(existing, "export {};\n", "")
	existing = strings.ReplaceAll(existing, "export {};", "")
	updated := existing + "export * from './users.table';\n"
	ctx.State.WriteFile("src/db/schema/index.ts", []byte(updated), state.ContentRaw)
	return nil
}
