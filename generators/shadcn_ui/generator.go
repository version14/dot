package shadcnui

import (
	"embed"

	"github.com/version14/dot/internal/deps"
	"github.com/version14/dot/internal/render"
	"github.com/version14/dot/internal/state"
	"github.com/version14/dot/pkg/dotapi"
)

//go:embed all:files
var filesFS embed.FS

//go:embed extra/globals.css
var globalCSSBytes []byte

type shadcnData struct {
	RSC     string
	CSSPath string
}

type Generator struct{}

func New() *Generator { return &Generator{} }

func (g *Generator) Name() string    { return Manifest.Name }
func (g *Generator) Version() string { return Manifest.Version }

func (g *Generator) Generate(ctx *dotapi.Context) error {
	framework, _ := ctx.Answers["framework"].(string)

	cssPath := "src/styles/globals.css"
	rsc := "false"
	if framework == "next" {
		cssPath = "src/app/globals.css"
		rsc = "true"
	}

	if err := ctx.State.UpdateJSON("package.json", func(d *state.JSONDoc) error {
		d.Merge(map[string]interface{}{
			"dependencies":    deps.NPM("tailwindcss", "class-variance-authority", "clsx", "tailwind-merge", "lucide-react"),
			"devDependencies": deps.NPM("@tailwindcss/vite"),
		})
		return nil
	}); err != nil {
		return err
	}

	if err := ctx.State.UpdateJSON("tsconfig.json", func(d *state.JSONDoc) error {
		_ = d.SetNested("compilerOptions.paths", map[string]interface{}{
			"@/*": []interface{}{"./src/*"},
		})
		return nil
	}); err != nil {
		return err
	}

	if err := render.NewLocalFolderRenderer(ctx.State).Render(filesFS, shadcnData{RSC: rsc, CSSPath: cssPath}); err != nil {
		return err
	}

	ctx.State.WriteFile(cssPath, globalCSSBytes, state.ContentRaw)

	return nil
}
