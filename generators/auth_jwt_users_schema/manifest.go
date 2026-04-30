package authjwtusersschema

import "github.com/version14/dot/pkg/dotapi"

var Manifest = dotapi.Manifest{
	Name:        "auth_jwt_users_schema",
	Version:     "0.1.0",
	Description: "Drizzle users table schema for JWT authentication",
	DependsOn:   []string{"drizzle_postgres_adapter"},
	Outputs: []string{
		"src/db/schema/users.table.ts",
	},
	Validators: []dotapi.Validator{
		{
			Name: "auth_jwt_users_schema",
			Checks: []dotapi.Check{
				{Type: dotapi.CheckFileExists, Path: "src/db/schema/users.table.ts"},
			},
		},
	},
}
