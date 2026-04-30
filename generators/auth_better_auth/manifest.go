package authbetterauth

import "github.com/version14/dot/pkg/dotapi"

var Manifest = dotapi.Manifest{
	Name:        "auth_better_auth",
	Version:     "0.1.0",
	Description: "BetterAuth setup with Drizzle adapter: src/lib/auth.ts and auth route handler",
	DependsOn:   []string{"drizzle_postgres_adapter"},
	Outputs: []string{
		"src/lib/auth.ts",
		"src/routes/auth.route.ts",
	},
	PostGenerationCommands: []dotapi.Command{
		{Cmd: "pnpm install"},
	},
	Validators: []dotapi.Validator{
		{
			Name: "auth-better-auth",
			Checks: []dotapi.Check{
				{Type: dotapi.CheckFileExists, Path: "src/lib/auth.ts"},
				{Type: dotapi.CheckJSONKeyExists, Path: "package.json", Key: "dependencies.better-auth"},
			},
		},
	},
}
