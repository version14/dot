package pandacss

import (
	"github.com/version14/dot/internal/state"
	"github.com/version14/dot/pkg/dotapi"
)

type Generator struct{}

func New() *Generator { return &Generator{} }

func (g *Generator) Name() string    { return Manifest.Name }
func (g *Generator) Version() string { return Manifest.Version }

func (g *Generator) Generate(ctx *dotapi.Context) error {
	// Determine linter — init flow uses "linter", frontend flow uses "frontend-linter".
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
		// Add panda codegen to prepare so it runs automatically on pnpm install.
		// Append to any existing prepare value rather than overwriting it.
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

	// Exclude the generated styled-system directory from prettier.
	// Read existing .prettierignore content (if any) so we append rather than overwrite.
	existing := ""
	if f, ok := ctx.State.GetFile(".prettierignore"); ok {
		existing = string(f.Content)
		if len(existing) > 0 && existing[len(existing)-1] != '\n' {
			existing += "\n"
		}
	}
	ctx.State.WriteFile(".prettierignore", []byte(existing+"styled-system/\n"), state.ContentRaw)

	// Exclude the generated styled-system directory from biome when it is the linter.
	// AppendStringSet is safe even if biome_config hasn't run yet — its own
	// AppendStringSet call will merge rather than overwrite this entry.
	if linter == "biome" {
		if err := ctx.State.UpdateJSON("biome.json", func(d *state.JSONDoc) error {
			return d.AppendStringSet("files.includes", "!styled-system/**")
		}); err != nil {
			return err
		}
	}

	ctx.State.WriteFile("panda.config.ts", []byte(pandaConfig), state.ContentRaw)
	ctx.State.WriteFile("postcss.config.mjs", []byte(postcssConfig), state.ContentRaw)

	return nil
}

const pandaConfig = `import { defineConfig } from "@pandacss/dev";

export default defineConfig({
  preflight: true,
  include: ["./src/**/*.{js,jsx,ts,tsx}"],
  exclude: [],
  outdir: "styled-system",
});
`

const postcssConfig = `export default {
  plugins: {
    "@pandacss/dev/postcss": {},
  },
};
`
