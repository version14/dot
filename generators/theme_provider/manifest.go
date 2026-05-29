package themeprovider

import "github.com/version14/dot/pkg/dotapi"

var Manifest = dotapi.Manifest{
	Name:        "theme_provider",
	Version:     "0.1.0",
	Description: "Custom ThemeProvider with CSS variables for light/dark mode",
	DependsOn:   []string{"typescript_base"},
	Outputs: []string{
		"src/providers/ThemeProvider.tsx",
		"src/hooks/useTheme.ts",
		"src/styles/theme.css",
	},
	TestCommands: []dotapi.Command{
		{Cmd: "pnpm exec tsc --noEmit"},
	},
	Validators: []dotapi.Validator{
		{
			Name: "theme-provider-files",
			Checks: []dotapi.Check{
				{Type: dotapi.CheckFileExists, Path: "src/providers/ThemeProvider.tsx"},
				{Type: dotapi.CheckFileExists, Path: "src/styles/theme.css"},
			},
		},
	},
}
