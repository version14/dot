package sentryfrontend

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

	sentryPkg := "@sentry/react"
	src := filesFS
	dir := "files"
	if framework == "next" {
		sentryPkg = "@sentry/nextjs"
		src = nextFS
		dir = "next"
	}

	if err := ctx.State.UpdateJSON("package.json", func(d *state.JSONDoc) error {
		d.Merge(map[string]interface{}{
			"dependencies": map[string]interface{}{
				sentryPkg: "^10.55.0",
			},
		})
		return nil
	}); err != nil {
		return err
	}

	ctx.State.WriteFile("src/lib/sentry.ts", mustRead(src, dir+"/sentry.ts"), state.ContentRaw)
	ctx.State.AppendFile(".env.example", mustRead(src, dir+"/.env.example"))

	if framework == "next" {
		ctx.State.WriteFile("sentry.client.config.ts", mustRead(nextFS, "next/sentry.client.config.ts"), state.ContentRaw)
	}

	return nil
}
