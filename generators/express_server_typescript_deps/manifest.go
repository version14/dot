package expressservertypescriptdeps

import "github.com/version14/dot/pkg/dotapi"

var Manifest = dotapi.Manifest{
	Name:        "express_server_typescript_deps",
	Version:     "0.1.0",
	Description: "Express + CORS + dotenv npm dependencies and dev/build/start scripts for TypeScript projects",
	DependsOn:   []string{"express_server_entrypoint"},
	Outputs:     []string{},
	PostGenerationCommands: []dotapi.Command{
		{Cmd: "pnpm install"},
	},
	Validators: []dotapi.Validator{
		{
			Name: "express-typescript-deps",
			Checks: []dotapi.Check{
				{Type: dotapi.CheckJSONKeyExists, Path: "package.json", Key: "dependencies.express"},
				{Type: dotapi.CheckJSONKeyExists, Path: "package.json", Key: "scripts.dev"},
				{Type: dotapi.CheckJSONKeyExists, Path: "package.json", Key: "scripts.build"},
				{Type: dotapi.CheckJSONKeyExists, Path: "package.json", Key: "scripts.start"},
			},
		},
	},
}
