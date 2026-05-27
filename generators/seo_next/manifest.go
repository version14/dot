package seonext

import "github.com/version14/dot/pkg/dotapi"

var Manifest = dotapi.Manifest{
	Name:          "seo_next",
	Version:       "0.1.0",
	Description:   "SEO setup for Next.js using next-seo",
	DependsOn:     []string{"nextjs_base"},
	ConflictsWith: []string{"seo_react"},
	Outputs: []string{
		"src/lib/seo.config.ts",
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
			Name: "seo-next-files",
			Checks: []dotapi.Check{
				{Type: dotapi.CheckFileExists, Path: "src/lib/seo.config.ts"},
			},
		},
	},
}
