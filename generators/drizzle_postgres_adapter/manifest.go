package drizzlepostgresadapter

import "github.com/version14/dot/pkg/dotapi"

var Manifest = dotapi.Manifest{
	Name:        "drizzle_postgres_adapter",
	Version:     "0.1.0",
	Description: "PostgreSQL driver and database connection file for Drizzle ORM (src/db/index.ts)",
	DependsOn:   []string{"drizzle_typescript_deps"},
	Outputs: []string{
		"src/db/index.ts",
	},
	Validators: []dotapi.Validator{
		{
			Name: "drizzle-postgres-adapter",
			Checks: []dotapi.Check{
				{Type: dotapi.CheckFileExists, Path: "src/db/index.ts"},
				{Type: dotapi.CheckJSONKeyExists, Path: "package.json", Key: "dependencies.postgres"},
			},
		},
	},
}
