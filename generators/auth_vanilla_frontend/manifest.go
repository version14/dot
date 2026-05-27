package authvanillafrontend

import "github.com/version14/dot/pkg/dotapi"

var Manifest = dotapi.Manifest{
	Name:          "auth_vanilla_frontend",
	Version:       "0.1.0",
	Description:   "Custom JWT authentication context and useAuth hook (from scratch)",
	DependsOn:     []string{"typescript_base"},
	ConflictsWith: []string{"auth_clerk_frontend", "auth_better_auth_frontend"},
	Outputs: []string{
		"src/lib/auth.ts",
		"src/providers/AuthProvider.tsx",
		"src/hooks/useAuth.ts",
	},
	TestCommands: []dotapi.Command{
		{Cmd: "pnpm exec tsc --noEmit"},
	},
	Validators: []dotapi.Validator{
		{
			Name: "auth-vanilla-frontend-files",
			Checks: []dotapi.Check{
				{Type: dotapi.CheckFileExists, Path: "src/lib/auth.ts"},
			},
		},
	},
}
