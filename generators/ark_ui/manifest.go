package arkui

import "github.com/version14/dot/pkg/dotapi"

var Manifest = dotapi.Manifest{
	Name:          "ark_ui",
	Version:       "0.2.0",
	Description:   "Ark UI headless component library",
	DependsOn:     []string{"typescript_base"},
	ConflictsWith: []string{"shadcn_ui"},
	Outputs: []string{
		"src/components/ui/button.tsx",
		"src/components/ui/index.ts",
	},
	PostGenerationCommands: []dotapi.Command{
		{Cmd: "pnpm install --dangerously-allow-all-builds"},
	},
	TestCommands: []dotapi.Command{
		{Cmd: "pnpm exec tsc --noEmit"},
	},
	Validators: []dotapi.Validator{
		{
			Name: "ark-ui-files",
			Checks: []dotapi.Check{
				{Type: dotapi.CheckFileExists, Path: "src/components/ui/button.tsx"},
				{Type: dotapi.CheckJSONKeyExists, Path: "package.json", Key: "dependencies.@ark-ui/react"},
			},
		},
	},
}
