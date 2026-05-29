package seoreact

import "github.com/version14/dot/pkg/dotapi"

var Manifest = dotapi.Manifest{
	Name:          "seo_react",
	Version:       "0.2.0",
	Description:   "SEO setup for React+Vite using react-helmet-async",
	DependsOn:     []string{"react_app"},
	ConflictsWith: []string{"seo_next"},
	Outputs: []string{
		"src/providers/HelmetProvider.tsx",
		"src/components/SEO.tsx",
	},
	PostGenerationCommands: []dotapi.Command{
		{Cmd: "pnpm install --dangerously-allow-all-builds"},
	},
	TestCommands: []dotapi.Command{
		{Cmd: "pnpm exec tsc --noEmit"},
	},
	Validators: []dotapi.Validator{
		{
			Name: "seo-react-files",
			Checks: []dotapi.Check{
				{Type: dotapi.CheckFileExists, Path: "src/components/SEO.tsx"},
			},
		},
	},
}
