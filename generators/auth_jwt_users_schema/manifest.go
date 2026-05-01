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
	TestCommands: []dotapi.Command{
		{Cmd: "test -f .env || cp .env.example .env"},
		{Cmd: "docker compose down -v 2>/dev/null || true"},
		{Cmd: "docker compose up -d && sleep 5"},
		{Cmd: "pnpm exec drizzle-kit push --force"},
		{Cmd: "bash -c 'pnpm exec vitest run db; EXIT_CODE=$?; docker compose down -v; exit $EXIT_CODE'"},
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
