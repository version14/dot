package biomeconfig

import "github.com/version14/dot/pkg/dotapi"

const biomeConfigFileName = "biome.json"

// Manifest declares biome_config — a linter+formatter setup using Biome.
// Depends on typescript_base since it modifies package.json scripts.
var Manifest = dotapi.Manifest{
	Name:        "biome_config",
	Version:     "0.1.0",
	Description: "Biome lint + format configuration",
	DependsOn:   []string{"typescript_base"},
	Outputs: []string{
		biomeConfigFileName,
	},
	PostGenerationCommands: []dotapi.Command{
		{Cmd: "pnpm install"},
	},
	TestCommands: []dotapi.Command{
		{Cmd: "pnpm install"},
		{Cmd: "pnpm exec biome check ."},
	},
	Validators: []dotapi.Validator{
		{
			Name: "biome-config",
			Checks: []dotapi.Check{
				{Type: dotapi.CheckFileExists, Path: biomeConfigFileName},
				{Type: dotapi.CheckJSONKeyExists, Path: biomeConfigFileName, Key: "linter.enabled"},
				{Type: dotapi.CheckJSONKeyExists, Path: "package.json", Key: "scripts.lint"},
			},
		},
	},
}
