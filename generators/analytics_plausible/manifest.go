package analyticsplausible

import "github.com/version14/dot/pkg/dotapi"

var Manifest = dotapi.Manifest{
	Name:          "analytics_plausible",
	Version:       "0.1.0",
	Description:   "Plausible Analytics privacy-first tracking",
	DependsOn:     []string{"typescript_base"},
	ConflictsWith: []string{"analytics_ga4"},
	Outputs: []string{
		"src/lib/plausible.ts",
	},
	PostGenerationCommands: []dotapi.Command{
		{Cmd: "pnpm install --dangerously-allow-all-builds"},
	},
	TestCommands: []dotapi.Command{
		{Cmd: "pnpm exec tsc --noEmit"},
	},
	Validators: []dotapi.Validator{
		{
			Name: "analytics-plausible-files",
			Checks: []dotapi.Check{
				{Type: dotapi.CheckFileExists, Path: "src/lib/plausible.ts"},
			},
		},
	},
}
