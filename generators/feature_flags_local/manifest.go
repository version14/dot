package featureflagslocal

import "github.com/version14/dot/pkg/dotapi"

var Manifest = dotapi.Manifest{
	Name:        "feature_flags_local",
	Version:     "0.1.0",
	Description: "Local JSON file-based feature flags with typed hook",
	DependsOn:   []string{"typescript_base"},
	Outputs: []string{
		"public/flags.json",
		"src/lib/flags.ts",
		"src/hooks/useFeatureFlag.ts",
	},
	TestCommands: []dotapi.Command{
		{Cmd: "pnpm exec tsc --noEmit"},
	},
	Validators: []dotapi.Validator{
		{
			Name: "feature-flags-local-files",
			Checks: []dotapi.Check{
				{Type: dotapi.CheckFileExists, Path: "public/flags.json"},
				{Type: dotapi.CheckFileExists, Path: "src/lib/flags.ts"},
			},
		},
	},
}
