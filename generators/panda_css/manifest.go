package pandacss

import "github.com/version14/dot/pkg/dotapi"

var Manifest = dotapi.Manifest{
	Name:          "panda_css",
	Version:       "0.1.0",
	Description:   "Panda CSS v1 with styled-system code generation",
	DependsOn:     []string{"typescript_base"},
	ConflictsWith: []string{"tailwind_v4", "shadcn_ui"},
	Outputs: []string{
		"panda.config.ts",
		"postcss.config.mjs",
	},
	PostGenerationCommands: []dotapi.Command{
		// panda codegen runs via the prepare script on pnpm install
		{Cmd: "pnpm install --dangerously-allow-all-builds"},
	},
	TestCommands: []dotapi.Command{
		{Cmd: "pnpm exec tsc --noEmit"},
	},
	Validators: []dotapi.Validator{
		{
			Name: "panda-css-files",
			Checks: []dotapi.Check{
				{Type: dotapi.CheckFileExists, Path: "panda.config.ts"},
				{Type: dotapi.CheckJSONKeyExists, Path: "package.json", Key: "devDependencies.@pandacss/dev"},
			},
		},
	},
}
