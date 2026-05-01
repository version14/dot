package postgresenvexample

import "github.com/version14/dot/pkg/dotapi"

var Manifest = dotapi.Manifest{
	Name:        "postgres_env_example",
	Version:     "0.1.1",
	Description: "Appends PostgreSQL DATABASE_URL to .env.example",
	DependsOn:   []string{"express_server_entrypoint"},
	Outputs:     []string{},
	Validators:  []dotapi.Validator{},
	PostGenerationCommands: []dotapi.Command{
		{Cmd: "cp .env.example .env"},
	},
}
