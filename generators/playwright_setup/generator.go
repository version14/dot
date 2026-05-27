package playwrightsetup

import (
	"fmt"

	"github.com/version14/dot/internal/state"
	"github.com/version14/dot/pkg/dotapi"
)

type Generator struct{}

func New() *Generator { return &Generator{} }

func (g *Generator) Name() string    { return Manifest.Name }
func (g *Generator) Version() string { return Manifest.Version }

func (g *Generator) Generate(ctx *dotapi.Context) error {
	framework, _ := ctx.Answers["framework"].(string)

	devServerURL := "http://localhost:5173"
	devServerCommand := "pnpm exec vite --host 127.0.0.1"
	if framework == "next" {
		devServerURL = "http://localhost:3000"
		devServerCommand = "pnpm exec next start"
	}

	if err := ctx.State.UpdateJSON("package.json", func(d *state.JSONDoc) error {
		d.Merge(map[string]interface{}{
			"scripts": map[string]interface{}{
				"test:e2e": "playwright test",
			},
			"devDependencies": map[string]interface{}{
				"@playwright/test": "^1.44.0",
			},
		})
		return nil
	}); err != nil {
		return err
	}

	playwrightConfig := fmt.Sprintf(playwrightConfigTmpl, devServerURL, devServerCommand, devServerURL)
	ctx.State.WriteFile("playwright.config.ts", []byte(playwrightConfig), state.ContentRaw)
	ctx.State.WriteFile("e2e/example.spec.ts", []byte(exampleSpec), state.ContentRaw)

	return nil
}

const playwrightConfigTmpl = `import { defineConfig, devices } from "@playwright/test";

export default defineConfig({
  testDir: "./e2e",
  fullyParallel: true,
  reporter: "html",
  use: { baseURL: "%s" },
  webServer: {
    command: "%s",
    url: "%s",
    reuseExistingServer: true,
    timeout: 120_000,
  },
  projects: [
    { name: "chromium", use: { ...devices["Desktop Chrome"] } },
  ],
});
`

const exampleSpec = `import { test, expect } from "@playwright/test";

test("homepage renders", async ({ page }) => {
  await page.goto("/");
  await expect(page).toHaveTitle(/.*/);
});
`
