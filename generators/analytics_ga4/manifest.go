package analyticsga4

import "github.com/version14/dot/pkg/dotapi"

var Manifest = dotapi.Manifest{
	Name:          "analytics_ga4",
	Version:       "0.1.0",
	Description:   "Google Analytics 4 (GA4) with pageview tracking",
	DependsOn:     []string{"typescript_base"},
	ConflictsWith: []string{"analytics_plausible"},
	Outputs: []string{
		"src/lib/ga4.ts",
	},
	PostGenerationCommands: []dotapi.Command{
		{Cmd: "pnpm install --dangerously-allow-all-builds"},
	},
	TestCommands: []dotapi.Command{
		{Cmd: "pnpm exec tsc --noEmit"},
	},
	Validators: []dotapi.Validator{
		{
			Name: "analytics-ga4-files",
			Checks: []dotapi.Check{
				{Type: dotapi.CheckFileExists, Path: "src/lib/ga4.ts"},
			},
		},
	},
}
