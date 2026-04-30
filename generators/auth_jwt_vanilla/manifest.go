package authjwtvanilla

import "github.com/version14/dot/pkg/dotapi"

var Manifest = dotapi.Manifest{
	Name:        "auth_jwt_vanilla",
	Version:     "0.1.0",
	Description: "Vanilla JWT authentication: src/lib/jwt.ts utility and src/middleware/auth.middleware.ts",
	DependsOn:   []string{"express_server_entrypoint"},
	Outputs: []string{
		"src/lib/jwt.ts",
		"src/middleware/auth.middleware.ts",
	},
	PostGenerationCommands: []dotapi.Command{
		{Cmd: "pnpm install"},
	},
	Validators: []dotapi.Validator{
		{
			Name: "auth-jwt-vanilla",
			Checks: []dotapi.Check{
				{Type: dotapi.CheckFileExists, Path: "src/lib/jwt.ts"},
				{Type: dotapi.CheckFileExists, Path: "src/middleware/auth.middleware.ts"},
				{Type: dotapi.CheckJSONKeyExists, Path: "package.json", Key: "dependencies.jsonwebtoken"},
			},
		},
	},
}
