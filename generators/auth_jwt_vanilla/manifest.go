package authjwtvanilla

import "github.com/version14/dot/pkg/dotapi"

var Manifest = dotapi.Manifest{
	Name:        "auth_jwt_vanilla",
	Version:     "0.1.0",
	Description: "Vanilla JWT authentication: src/shared/services/jwt.ts utility and src/shared/middlewares/auth.middleware.ts",
	DependsOn:   []string{"express_server_entrypoint"},
	Outputs: []string{
		"src/shared/services/jwt.ts",
		"src/shared/middlewares/auth.middleware.ts",
	},
	TestCommands: []dotapi.Command{},
	Validators: []dotapi.Validator{
		{
			Name: "auth-jwt-vanilla",
			Checks: []dotapi.Check{
				{Type: dotapi.CheckFileExists, Path: "src/shared/services/jwt.ts"},
				{Type: dotapi.CheckFileExists, Path: "src/shared/middlewares/auth.middleware.ts"},
				{Type: dotapi.CheckJSONKeyExists, Path: "package.json", Key: "dependencies.jsonwebtoken"},
			},
		},
	},
}
