package expresstestsetup

import "github.com/version14/dot/pkg/dotapi"

var Manifest = dotapi.Manifest{
	Name:        "express_test_setup",
	Version:     "0.1.0",
	Description: "Vitest configuration and testing dependencies (vitest, supertest) for Express TypeScript projects",
	DependsOn:   []string{"express_server_entrypoint"},
	Outputs:     []string{"vitest.config.ts"},
	TestCommands: []dotapi.Command{
		{Cmd: "pnpm exec vitest run unit"},
		{Cmd: "pnpm exec vitest run feature"},
		{Cmd: "test -f .env.example && cp -n .env.example .env; pnpm exec vitest run e2e"},
	},
	Validators: []dotapi.Validator{
		{
			Name: "express-test-setup",
			Checks: []dotapi.Check{
				{Type: dotapi.CheckFileExists, Path: "vitest.config.ts"},
				{Type: dotapi.CheckJSONKeyExists, Path: "package.json", Key: "devDependencies.vitest"},
			},
		},
	},
}
