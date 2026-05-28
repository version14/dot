package tailwindv4

import (
	"github.com/version14/dot/internal/state"
	"github.com/version14/dot/pkg/dotapi"
)

type Generator struct{}

func New() *Generator { return &Generator{} }

func (g *Generator) Name() string    { return Manifest.Name }
func (g *Generator) Version() string { return Manifest.Version }

func (g *Generator) Generate(ctx *dotapi.Context) error {
	framework, _ := ctx.Answers["framework"].(string)

	deps := map[string]interface{}{
		"tailwindcss": "^4.0.0",
	}
	devDeps := map[string]interface{}{
		"@tailwindcss/vite":    "^4.0.0",
		"@tailwindcss/postcss": "^4.0.0",
	}

	var postcssContent string
	if framework == "next" {
		postcssContent = postcssNext
	} else {
		postcssContent = postcssVite
		ctx.State.WriteFile("vite.config.ts", []byte(viteConfigTailwind), state.ContentRaw)
	}

	if err := ctx.State.UpdateJSON("package.json", func(d *state.JSONDoc) error {
		d.Merge(map[string]interface{}{
			"dependencies":    deps,
			"devDependencies": devDeps,
		})
		return nil
	}); err != nil {
		return err
	}

	ctx.State.WriteFile("src/styles/globals.css", []byte(globalCSS), state.ContentRaw)
	ctx.State.WriteFile("postcss.config.mjs", []byte(postcssContent), state.ContentRaw)

	return nil
}

const globalCSS = `@import "tailwindcss";
`

const postcssVite = `export default {
  plugins: {
    "@tailwindcss/postcss": {},
  },
};
`

const postcssNext = `export default {
  plugins: {
    "@tailwindcss/postcss": {},
  },
};
`

const viteConfigTailwind = `import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import tailwindcss from "@tailwindcss/vite";

export default defineConfig({
  plugins: [tailwindcss(), react()],
});
`
