package version14ui

import "github.com/version14/dot/pkg/dotapi"

var Manifest = dotapi.Manifest{
	Name:          "version14_ui",
	Version:       "0.4.0",
	Description:   "@version14/ui component library",
	DependsOn:     []string{"typescript_base"},
	ConflictsWith: []string{"shadcn_ui", "ark_ui"},
	Outputs: []string{
		"src/components/V14Example.tsx",
	},
	PostGenerationCommands: []dotapi.Command{
		{Cmd: "pnpm install --dangerously-allow-all-builds"},
	},
	TestCommands: []dotapi.Command{
		{Cmd: "pnpm exec tsc --noEmit"},
	},
	Validators: []dotapi.Validator{
		{
			Name: "version14-ui-files",
			Checks: []dotapi.Check{
				{Type: dotapi.CheckJSONKeyExists, Path: "package.json", Key: "dependencies.@version14/ui"},
				{Type: dotapi.CheckJSONKeyExists, Path: "package.json", Key: "dependencies.@ark-ui/react"},
			},
		},
	},
}
