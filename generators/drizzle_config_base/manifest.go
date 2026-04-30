package drizzleconfigbase

import "github.com/version14/dot/pkg/dotapi"

var Manifest = dotapi.Manifest{
	Name:        "drizzle_config_base",
	Version:     "0.1.0",
	Description: "Base Drizzle ORM configuration: drizzle.config.ts and src/db/schema directory",
	DependsOn:   []string{"typescript_base"},
	Outputs: []string{
		"drizzle.config.ts",
		"src/db/schema/index.ts",
	},
	Validators: []dotapi.Validator{
		{
			Name: "drizzle-config-base",
			Checks: []dotapi.Check{
				{Type: dotapi.CheckFileExists, Path: "drizzle.config.ts"},
				{Type: dotapi.CheckFileExists, Path: "src/db/schema/index.ts"},
			},
		},
	},
}
