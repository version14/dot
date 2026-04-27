package baseproject

import "github.com/version14/dot/pkg/dotapi"

// Manifest declares the base_project generator. It always runs first and
// creates the universal scaffolding (README, .gitignore, LICENSE) every
// project gets regardless of language or framework.
var Manifest = dotapi.Manifest{
	Name:        "base_project",
	Version:     "0.1.0",
	Description: "Universal project scaffolding (README, .gitignore, LICENSE)",
	Outputs: []string{
		"README.md",
		".gitignore",
		"LICENSE",
	},
	Validators: []dotapi.Validator{
		{
			Name: "base-files",
			Checks: []dotapi.Check{
				{Type: dotapi.CheckFileExists, Path: "README.md"},
				{Type: dotapi.CheckFileExists, Path: ".gitignore"},
				{Type: dotapi.CheckFileExists, Path: "LICENSE"},
			},
		},
	},
}
