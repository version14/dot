package pluginreposkeleton

import "github.com/version14/dot/pkg/dotapi"

// Manifest declares plugin_repo_skeleton — the generator that scaffolds a
// brand-new git-publishable DOT plugin repository. The plugin-template flow
// invokes this exclusively (no DependsOn on base_project, because plugin
// repos want a plugin-shaped README, not the generic project README).
var Manifest = dotapi.Manifest{
	Name:        "plugin_repo_skeleton",
	Version:     "0.1.6",
	Description: "Scaffolds a publishable DOT plugin repository (go.mod + plugin.go + plugin.json + README + LICENSE)",
	Outputs: []string{
		"go.mod",
		"plugin.go",
		"plugin.json",
		"README.md",
		"LICENSE",
		".gitignore",
	},
	PostGenerationCommands: []dotapi.Command{
		{Cmd: "go mod tidy"},
		{Cmd: "git init"},
	},
	Validators: []dotapi.Validator{
		{
			Name: "plugin-skeleton",
			Checks: []dotapi.Check{
				{Type: dotapi.CheckFileExists, Path: "plugin.go"},
				{Type: dotapi.CheckFileExists, Path: "plugin.json"},
				{Type: dotapi.CheckFileExists, Path: "go.mod"},
				{Type: dotapi.CheckJSONKeyExists, Path: "plugin.json", Key: "id"},
				{Type: dotapi.CheckJSONKeyExists, Path: "plugin.json", Key: "version"},
			},
		},
	},
}
