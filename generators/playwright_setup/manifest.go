package playwrightsetup

import "github.com/version14/dot/pkg/dotapi"

var Manifest = dotapi.Manifest{
	Name:        "playwright_setup",
	Version:     "0.1.0",
	Description: "Playwright end-to-end testing setup",
	DependsOn:   []string{"typescript_base"},
	Outputs: []string{
		"playwright.config.ts",
		"e2e/example.spec.ts",
	},
	PostGenerationCommands: []dotapi.Command{
		{Cmd: "pnpm install --dangerously-allow-all-builds"},
		{Cmd: "pnpm exec playwright install --with-deps chromium"},
	},
	TestCommands: []dotapi.Command{
		{Cmd: "pnpm exec playwright test"},
	},
	Validators: []dotapi.Validator{
		{
			Name: "playwright-setup-files",
			Checks: []dotapi.Check{
				{Type: dotapi.CheckFileExists, Path: "playwright.config.ts"},
				{Type: dotapi.CheckJSONKeyExists, Path: "package.json", Key: "devDependencies.@playwright/test"},
			},
		},
	},
}
