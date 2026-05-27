package arkui

import (
	"github.com/version14/dot/internal/state"
	"github.com/version14/dot/pkg/dotapi"
)

type Generator struct{}

func New() *Generator { return &Generator{} }

func (g *Generator) Name() string    { return Manifest.Name }
func (g *Generator) Version() string { return Manifest.Version }

func (g *Generator) Generate(ctx *dotapi.Context) error {
	if err := ctx.State.UpdateJSON("package.json", func(d *state.JSONDoc) error {
		d.Merge(map[string]interface{}{
			"dependencies": map[string]interface{}{
				"@ark-ui/react": "^4.0.0",
			},
		})
		return nil
	}); err != nil {
		return err
	}

	ctx.State.WriteFile("src/components/ui/button.tsx", []byte(buttonTSX), state.ContentRaw)
	ctx.State.WriteFile("src/components/ui/index.ts", []byte(indexTS), state.ContentRaw)

	return nil
}

const buttonTSX = `import { ark } from "@ark-ui/react/factory";

export const Button = ark.button;
`

const indexTS = `export * from "./button";
`
