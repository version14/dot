package sentryfrontend

import "github.com/version14/dot/pkg/dotapi"

var Manifest = dotapi.Manifest{
	Name:        "sentry_frontend",
	Version:     "0.3.0",
	Description: "Sentry error tracking and performance monitoring",
	DependsOn:   []string{"typescript_base"},
	Outputs: []string{
		"src/lib/sentry.ts",
	},
	PostGenerationCommands: []dotapi.Command{
		{Cmd: "pnpm install --dangerously-allow-all-builds"},
	},
	TestCommands: []dotapi.Command{
		{Cmd: "pnpm exec tsc --noEmit"},
	},
	Validators: []dotapi.Validator{
		{
			Name: "sentry-frontend-files",
			Checks: []dotapi.Check{
				{Type: dotapi.CheckFileExists, Path: "src/lib/sentry.ts"},
			},
		},
	},
}
