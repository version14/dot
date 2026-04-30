package authjwtmvcroute

import "github.com/version14/dot/pkg/dotapi"

var Manifest = dotapi.Manifest{
	Name:        "auth_jwt_mvc_route",
	Version:     "0.1.0",
	Description: "JWT auth route and controller for MVC architecture",
	DependsOn:   []string{"auth_jwt_vanilla"},
	Outputs: []string{
		"src/routes/auth.route.ts",
		"src/controllers/auth.controller.ts",
	},
	PostGenerationCommands: []dotapi.Command{
		{Cmd: "pnpm install"},
	},
	Validators: []dotapi.Validator{
		{
			Name: "auth_jwt_mvc_route",
			Checks: []dotapi.Check{
				{Type: dotapi.CheckFileExists, Path: "src/routes/auth.route.ts"},
				{Type: dotapi.CheckFileExists, Path: "src/controllers/auth.controller.ts"},
			},
		},
	},
}
