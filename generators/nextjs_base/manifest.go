package nextjsbase

import "github.com/version14/dot/pkg/dotapi"

var Manifest = dotapi.Manifest{
	Name:        "nextjs_base",
	Version:     "0.1.0",
	Description: "Next.js 15 App Router project with TypeScript",
	DependsOn:   []string{"typescript_base"},
	Outputs: []string{
		"next.config.ts",
		"src/app/layout.tsx",
		"src/app/page.tsx",
		"src/app/globals.css",
	},
	PostGenerationCommands: []dotapi.Command{
		{Cmd: "pnpm install --dangerously-allow-all-builds"},
	},
	TestCommands: []dotapi.Command{
		{Cmd: "pnpm exec tsc --noEmit"},
		{Cmd: "pnpm exec next build"},
	},
	Validators: []dotapi.Validator{
		{
			Name: "nextjs-base-files",
			Checks: []dotapi.Check{
				{Type: dotapi.CheckFileExists, Path: "next.config.ts"},
				{Type: dotapi.CheckFileExists, Path: "src/app/layout.tsx"},
				{Type: dotapi.CheckFileExists, Path: "src/app/page.tsx"},
				{Type: dotapi.CheckJSONKeyExists, Path: "package.json", Key: "dependencies.next"},
			},
		},
	},
}
