package authbetterauthschema

import "github.com/version14/dot/pkg/dotapi"

var Manifest = dotapi.Manifest{
	Name:        "auth_better_auth_schema",
	Version:     "0.1.0",
	Description: "Drizzle schema tables required by BetterAuth (user, session, account, verification)",
	DependsOn:   []string{"auth_better_auth"},
	Outputs: []string{
		"src/db/schema/auth.schema.ts",
	},
	Validators: []dotapi.Validator{
		{
			Name: "auth_better_auth_schema",
			Checks: []dotapi.Check{
				{Type: dotapi.CheckFileExists, Path: "src/db/schema/auth.schema.ts"},
			},
		},
	},
}
