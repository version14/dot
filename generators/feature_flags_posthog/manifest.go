package featureflagsposthog

import "github.com/version14/dot/pkg/dotapi"

var Manifest = dotapi.Manifest{
	Name:        "feature_flags_posthog",
	Version:     "0.3.0",
	Description: "PostHog feature flags and analytics provider",
	DependsOn:   []string{"typescript_base"},
	Outputs: []string{
		"src/lib/posthog.ts",
		"src/providers/PostHogProvider.tsx",
	},
	PostGenerationCommands: []dotapi.Command{
		{Cmd: "pnpm install --dangerously-allow-all-builds"},
	},
	TestCommands: []dotapi.Command{
		{Cmd: "pnpm exec tsc --noEmit"},
	},
	Validators: []dotapi.Validator{
		{
			Name: "feature-flags-posthog-files",
			Checks: []dotapi.Check{
				{Type: dotapi.CheckFileExists, Path: "src/lib/posthog.ts"},
				{Type: dotapi.CheckJSONKeyExists, Path: "package.json", Key: "dependencies.posthog-js"},
			},
		},
	},
}
