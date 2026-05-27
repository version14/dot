package storybooksetup

import "github.com/version14/dot/pkg/dotapi"

var Manifest = dotapi.Manifest{
	Name:        "storybook_setup",
	Version:     "0.1.0",
	Description: "Storybook v8 component development environment",
	DependsOn:   []string{"typescript_base"},
	Outputs: []string{
		".storybook/main.ts",
		".storybook/preview.ts",
		"src/stories/Button.stories.tsx",
	},
	PostGenerationCommands: []dotapi.Command{
		{Cmd: "pnpm install --dangerously-allow-all-builds"},
	},
	TestCommands: []dotapi.Command{
		{Cmd: "pnpm exec storybook build"},
	},
	Validators: []dotapi.Validator{
		{
			Name: "storybook-setup-files",
			Checks: []dotapi.Check{
				{Type: dotapi.CheckFileExists, Path: ".storybook/main.ts"},
				{Type: dotapi.CheckJSONKeyExists, Path: "package.json", Key: "devDependencies.storybook"},
			},
		},
	},
}
