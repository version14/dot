package authclerkfrontend

import "github.com/version14/dot/pkg/dotapi"

var Manifest = dotapi.Manifest{
	Name:          "auth_clerk_frontend",
	Version:       "0.1.0",
	Description:   "Clerk authentication provider and useAuth hook",
	DependsOn:     []string{"typescript_base"},
	ConflictsWith: []string{"auth_better_auth_frontend", "auth_vanilla_frontend"},
	Outputs: []string{
		"src/providers/ClerkProvider.tsx",
		"src/hooks/useAuth.ts",
	},
	PostGenerationCommands: []dotapi.Command{
		{Cmd: "pnpm install --dangerously-allow-all-builds"},
	},
	TestCommands: []dotapi.Command{
		{Cmd: "pnpm exec tsc --noEmit"},
	},
	Validators: []dotapi.Validator{
		{
			Name: "auth-clerk-frontend-files",
			Checks: []dotapi.Check{
				{Type: dotapi.CheckFileExists, Path: "src/providers/ClerkProvider.tsx"},
			},
		},
	},
}
