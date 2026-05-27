package version14ui

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
				"@version14/ui": "latest",
			},
		})
		return nil
	}); err != nil {
		return err
	}

	ctx.State.WriteFile("src/components/V14Example.tsx", []byte(exampleTSX), state.ContentRaw)

	return nil
}

const exampleTSX = `// Example usage of @version14/ui — replace with your own components.
export default function V14Example() {
  return (
    <div>
      @version14/ui is installed. Import components from "@version14/ui".
    </div>
  );
}
`
