package vitesttestinglibrary

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
			"scripts": map[string]interface{}{
				"test":          "vitest run",
				"test:watch":    "vitest",
				"test:coverage": "vitest run --coverage",
			},
			"devDependencies": map[string]interface{}{
				"vitest":                      "^2.0.0",
				"@vitest/coverage-v8":         "^2.0.0",
				"@testing-library/react":      "^16.0.0",
				"@testing-library/user-event": "^14.0.0",
				"@testing-library/jest-dom":   "^6.0.0",
				"jsdom":                       "^24.0.0",
			},
		})
		return nil
	}); err != nil {
		return err
	}

	ctx.State.WriteFile("vitest.config.ts", []byte(vitestConfig), state.ContentRaw)
	ctx.State.WriteFile("src/test/setup.ts", []byte(setupTS), state.ContentRaw)
	ctx.State.WriteFile("src/test/App.test.tsx", []byte(appTestTSX), state.ContentRaw)

	return nil
}

const vitestConfig = `import { defineConfig } from "vitest/config";
import react from "@vitejs/plugin-react";

export default defineConfig({
  plugins: [react()],
  test: {
    globals: true,
    environment: "jsdom",
    exclude: ["**/e2e/**", "**/node_modules/**"],
    setupFiles: ["./src/test/setup.ts"],
  },
});
`

const setupTS = `import "@testing-library/jest-dom";
`

const appTestTSX = `import { render } from "@testing-library/react";
import { describe, it, expect } from "vitest";
import App from "../App";

describe("App", () => {
  it("renders without crashing", () => {
    render(<App />);
    expect(document.body).toBeDefined();
  });
});
`
