package authbetterauthfrontend

import "github.com/version14/dot/pkg/dotapi"

var Manifest = dotapi.Manifest{
	Name:          "auth_better_auth_frontend",
	Version:       "0.2.1",
	Description:   "Better Auth client setup with AuthProvider and useAuth hook",
	DependsOn:     []string{"typescript_base"},
	ConflictsWith: []string{"auth_clerk_frontend", "auth_vanilla_frontend"},
	Outputs: []string{
		"src/lib/authClient.ts",
		"src/providers/AuthProvider.tsx",
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
			Name: "auth-better-auth-frontend-files",
			Checks: []dotapi.Check{
				{Type: dotapi.CheckFileExists, Path: "src/lib/authClient.ts"},
				{Type: dotapi.CheckJSONKeyExists, Path: "package.json", Key: "dependencies.better-auth"},
			},
		},
	},
}
