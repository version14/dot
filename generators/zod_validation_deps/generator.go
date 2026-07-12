package zodvalidationdeps

import (
	"github.com/version14/dot/internal/deps"
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
			"dependencies":    deps.NPM("zod", "@asteasolutions/zod-to-openapi", "reflect-metadata", "swagger-ui-express"),
			"devDependencies": deps.NPM("@types/swagger-ui-express"),
		})
		return nil
	}); err != nil {
		return err
	}

	return ctx.State.UpdateJSON("tsconfig.json", func(d *state.JSONDoc) error {
		d.Merge(map[string]interface{}{
			"compilerOptions": map[string]interface{}{
				"experimentalDecorators": true,
				"emitDecoratorMetadata":  true,
			},
		})
		return nil
	})
}
