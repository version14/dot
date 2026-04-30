package expressnodetsconfig

import (
	"github.com/version14/dot/internal/state"
	"github.com/version14/dot/pkg/dotapi"
)

type Generator struct{}

func New() *Generator { return &Generator{} }

func (g *Generator) Name() string    { return Manifest.Name }
func (g *Generator) Version() string { return Manifest.Version }

func (g *Generator) Generate(ctx *dotapi.Context) error {
	// Override tsconfig with Node.js/CommonJS settings and remove "type": "module"
	// from package.json (ESNext module set by typescript_base is incompatible with Express ecosystem)
	if err := ctx.State.UpdateJSON("tsconfig.json", func(d *state.JSONDoc) error {
		d.Merge(map[string]interface{}{
			"compilerOptions": map[string]interface{}{
				"target":            "ES2022",
				"module":            "CommonJS",
				"moduleResolution":  "Node",
				"strict":            true,
				"esModuleInterop":   true,
				"skipLibCheck":      true,
				"outDir":            "dist",
				"rootDir":           "src",
				"resolveJsonModule": true,
			},
			"include": []interface{}{"src/**/*"},
			"exclude": []interface{}{"node_modules", "dist"},
		})
		return nil
	}); err != nil {
		return err
	}

	// CommonJS projects must not have "type": "module" in package.json
	return ctx.State.UpdateJSON("package.json", func(d *state.JSONDoc) error {
		d.Merge(map[string]interface{}{
			"type": "commonjs",
		})
		return nil
	})
}
