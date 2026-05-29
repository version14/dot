package featureflagsvercel

import "github.com/version14/dot/pkg/dotapi"

var Manifest = dotapi.Manifest{
	Name:        "feature_flags_vercel",
	Version:     "0.2.0",
	Description: "Vercel Edge Config feature flags with typed hook",
	DependsOn:   []string{"typescript_base"},
	Outputs: []string{
		"src/lib/flags.ts",
		"src/hooks/useFeatureFlag.ts",
	},
	PostGenerationCommands: []dotapi.Command{
		{Cmd: "pnpm install --dangerously-allow-all-builds"},
	},
	TestCommands: []dotapi.Command{
		{Cmd: "pnpm exec tsc --noEmit"},
	},
	Validators: []dotapi.Validator{
		{
			Name: "feature-flags-vercel-files",
			Checks: []dotapi.Check{
				{Type: dotapi.CheckFileExists, Path: "src/lib/flags.ts"},
				{Type: dotapi.CheckJSONKeyExists, Path: "package.json", Key: "dependencies.@vercel/edge-config"},
			},
		},
	},
}
