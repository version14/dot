package drizzletypescriptdeps

import "github.com/version14/dot/pkg/dotapi"

var Manifest = dotapi.Manifest{
	Name:        "drizzle_typescript_deps",
	Version:     "0.1.0",
	Description: "Adds drizzle-orm and drizzle-kit to package.json and db:* scripts",
	DependsOn:   []string{"drizzle_config_base"},
	Outputs:     []string{},
	PostGenerationCommands: []dotapi.Command{
		{Cmd: "pnpm db:generate"},
	},
	Validators: []dotapi.Validator{
		{
			Name: "drizzle-typescript-deps",
			Checks: []dotapi.Check{
				{Type: dotapi.CheckJSONKeyExists, Path: "package.json", Key: "dependencies.drizzle-orm"},
				{Type: dotapi.CheckJSONKeyExists, Path: "package.json", Key: "devDependencies.drizzle-kit"},
				{Type: dotapi.CheckJSONKeyExists, Path: "package.json", Key: "scripts.db:push"},
			},
		},
	},
}
