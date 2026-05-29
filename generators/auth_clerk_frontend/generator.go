package authclerkfrontend

import (
	"embed"

	"github.com/version14/dot/internal/state"
	"github.com/version14/dot/pkg/dotapi"
)

//go:embed all:files
var filesFS embed.FS

//go:embed all:next
var nextFS embed.FS

type Generator struct{}

func New() *Generator { return &Generator{} }

func (g *Generator) Name() string    { return Manifest.Name }
func (g *Generator) Version() string { return Manifest.Version }

func mustRead(fs embed.FS, path string) []byte {
	data, err := fs.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return data
}

func (g *Generator) Generate(ctx *dotapi.Context) error {
	framework, _ := ctx.Answers["framework"].(string)

	pkg := "@clerk/clerk-react"
	version := "^5.61.3"
	src := filesFS
	if framework == "next" {
		pkg = "@clerk/nextjs"
		version = "^7.4.2"
		src = nextFS
	}

	dir := "files"
	if framework == "next" {
		dir = "next"
	}

	if err := ctx.State.UpdateJSON("package.json", func(d *state.JSONDoc) error {
		d.Merge(map[string]interface{}{
			"dependencies": map[string]interface{}{
				pkg: version,
			},
		})
		return nil
	}); err != nil {
		return err
	}

	ctx.State.WriteFile("src/providers/ClerkProvider.tsx", mustRead(src, dir+"/ClerkProvider.tsx"), state.ContentRaw)
	ctx.State.WriteFile("src/hooks/useAuth.ts", mustRead(src, dir+"/useAuth.ts"), state.ContentRaw)
	ctx.State.WriteFile(".env.example", mustRead(src, dir+"/.env.example"), state.ContentRaw)

	return nil
}
