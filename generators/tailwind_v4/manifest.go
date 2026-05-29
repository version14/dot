package tailwindv4

import "github.com/version14/dot/pkg/dotapi"

var Manifest = dotapi.Manifest{
	Name:          "tailwind_v4",
	Version:       "0.2.0",
	Description:   "Tailwind CSS v4 (CSS-first configuration)",
	DependsOn:     []string{"typescript_base"},
	ConflictsWith: []string{"panda_css", "shadcn_ui"},
	Outputs: []string{
		"src/styles/globals.css",
		"postcss.config.mjs",
	},
	PostGenerationCommands: []dotapi.Command{
		{Cmd: "pnpm install --dangerously-allow-all-builds"},
	},
	TestCommands: []dotapi.Command{
		{Cmd: "pnpm exec tsc --noEmit"},
	},
	Validators: []dotapi.Validator{
		{
			Name: "tailwind-v4-files",
			Checks: []dotapi.Check{
				{Type: dotapi.CheckFileExists, Path: "src/styles/globals.css"},
				{Type: dotapi.CheckJSONKeyExists, Path: "package.json", Key: "dependencies.tailwindcss"},
			},
		},
	},
}
