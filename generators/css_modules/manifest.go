package cssmodules

import "github.com/version14/dot/pkg/dotapi"

var Manifest = dotapi.Manifest{
	Name:          "css_modules",
	Version:       "0.1.0",
	Description:   "CSS Modules setup with global styles and example component stylesheet",
	DependsOn:     []string{"typescript_base"},
	ConflictsWith: []string{"tailwind_v4", "panda_css", "shadcn_ui"},
	Outputs: []string{
		"src/styles/global.css",
		"src/styles/App.module.css",
	},
	TestCommands: []dotapi.Command{
		{Cmd: "pnpm exec tsc --noEmit"},
	},
	Validators: []dotapi.Validator{
		{
			Name: "css-modules-files",
			Checks: []dotapi.Check{
				{Type: dotapi.CheckFileExists, Path: "src/styles/global.css"},
			},
		},
	},
}
