package backendArchitectureMVC

import "github.com/version14/dot/pkg/dotapi"

// Manifest declares backend-architecture-mvc — the generator that scaffolds an
// MVC-based backend structure.
var Manifest = dotapi.Manifest{
	Name:        "backend_architecture_mvc",
	Version:     "1.0.0",
	Description: "Scaffolds a base structure for a backend architecture using MVC",
	Outputs: []string{
		"src/controllers/.gitkeep",
		"src/views/.gitkeep",
		"src/models/.gitkeep",
		"src/routes/.gitkeep",

		"src/shared/validators/.gitkeep",
		"src/shared/errors/.gitkeep",
		"src/shared/libs/.gitkeep",
		"src/shared/middlewares/.gitkeep",
		"src/shared/services/.gitkeep",
	},
	Validators: []dotapi.Validator{
		{
			Name: "backend_architecture_mvc",
			Checks: []dotapi.Check{
				{Type: dotapi.CheckFileExists, Path: "src/controllers/.gitkeep"},
				{Type: dotapi.CheckFileExists, Path: "src/views/.gitkeep"},
				{Type: dotapi.CheckFileExists, Path: "src/models/.gitkeep"},
				{Type: dotapi.CheckFileExists, Path: "src/routes/.gitkeep"},

				{Type: dotapi.CheckFileExists, Path: "src/shared/validators/.gitkeep"},
				{Type: dotapi.CheckFileExists, Path: "src/shared/errors/.gitkeep"},
				{Type: dotapi.CheckFileExists, Path: "src/shared/libs/.gitkeep"},
				{Type: dotapi.CheckFileExists, Path: "src/shared/middlewares/.gitkeep"},
				{Type: dotapi.CheckFileExists, Path: "src/shared/services/.gitkeep"},
			},
		},
	},
}
