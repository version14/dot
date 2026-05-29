package shadcnui

import "github.com/version14/dot/pkg/dotapi"

var Manifest = dotapi.Manifest{
	Name:          "shadcn_ui",
	Version:       "0.2.0",
	Description:   "shadcn/ui component library with Tailwind CSS v4",
	DependsOn:     []string{"typescript_base"},
	ConflictsWith: []string{"tailwind_v4", "panda_css", "css_modules"},
	Outputs: []string{
		"components.json",
		"src/lib/utils.ts",
		"src/components/ui/button.tsx",
	},
	PostGenerationCommands: []dotapi.Command{
		{Cmd: "pnpm install --dangerously-allow-all-builds"},
	},
	TestCommands: []dotapi.Command{
		{Cmd: "pnpm exec tsc --noEmit"},
	},
	Validators: []dotapi.Validator{
		{
			Name: "shadcn-ui-files",
			Checks: []dotapi.Check{
				{Type: dotapi.CheckFileExists, Path: "components.json"},
				{Type: dotapi.CheckFileExists, Path: "src/components/ui/button.tsx"},
				{Type: dotapi.CheckJSONKeyExists, Path: "package.json", Key: "dependencies.tailwindcss"},
			},
		},
	},
}
