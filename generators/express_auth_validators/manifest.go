package expressauthvalidators

import "github.com/version14/dot/pkg/dotapi"

var Manifest = dotapi.Manifest{
	Name:        "express_auth_validators",
	Version:     "0.1.0",
	Description: "Zod schemas for auth endpoint input validation: register, login, refresh",
	DependsOn:   []string{"express_server_entrypoint"},
	Outputs:     []string{"src/shared/validators/auth.validators.ts"},
	Validators: []dotapi.Validator{
		{
			Name: "express-auth-validators",
			Checks: []dotapi.Check{
				{Type: dotapi.CheckFileExists, Path: "src/shared/validators/auth.validators.ts"},
				{Type: dotapi.CheckJSONKeyExists, Path: "package.json", Key: "dependencies.zod"},
			},
		},
	},
}
