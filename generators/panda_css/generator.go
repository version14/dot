package pandacss

import (
	"embed"

	"github.com/version14/dot/internal/render"
	"github.com/version14/dot/internal/state"
	"github.com/version14/dot/pkg/dotapi"
)

//go:embed all:files
var filesFS embed.FS

type Generator struct{}

func New() *Generator { return &Generator{} }

func (g *Generator) Name() string    { return Manifest.Name }
func (g *Generator) Version() string { return Manifest.Version }

func (g *Generator) Generate(ctx *dotapi.Context) error {
	linter, _ := ctx.Spec.Answers["linter"].(string)
	if linter == "" {
		linter, _ = ctx.Spec.Answers["frontend-linter"].(string)
	}

	if err := ctx.State.UpdateJSON("package.json", func(d *state.JSONDoc) error {
		d.Merge(map[string]interface{}{
			"devDependencies": map[string]interface{}{
				"@pandacss/dev": "^0.45.0",
			},
		})
		prepare := "panda codegen"
		if existing, ok := d.GetNested("scripts.prepare"); ok {
			if s, isStr := existing.(string); isStr && s != "" {
				prepare = s + " && panda codegen"
			}
		}
		return d.SetNested("scripts.prepare", prepare)
	}); err != nil {
		return err
	}

	existing := ""
	if f, ok := ctx.State.GetFile(".prettierignore"); ok {
		existing = string(f.Content)
		if len(existing) > 0 && existing[len(existing)-1] != '\n' {
			existing += "\n"
		}
	}
	ctx.State.WriteFile(".prettierignore", []byte(existing+"styled-system/\n"), state.ContentRaw)

	if linter == "biome" {
		if err := ctx.State.UpdateJSON("biome.json", func(d *state.JSONDoc) error {
			return d.AppendStringSet("files.includes", "!**/styled-system")
		}); err != nil {
			return err
		}
	}

	return render.NewLocalFolderRenderer(ctx.State).Render(filesFS, nil)
}
