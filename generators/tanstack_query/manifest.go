package tanstackquery

import "github.com/version14/dot/pkg/dotapi"

var Manifest = dotapi.Manifest{
	Name:        "tanstack_query",
	Version:     "0.1.0",
	Description: "TanStack Query v5 server-state management with QueryProvider and devtools",
	DependsOn:   []string{"react_app"},
	Outputs: []string{
		"src/providers/query-provider.tsx",
	},
	Validators: []dotapi.Validator{
		{
			Name: "tanstack-query-files",
			Checks: []dotapi.Check{
				{Type: dotapi.CheckFileExists, Path: "src/providers/query-provider.tsx"},
				{Type: dotapi.CheckJSONKeyExists, Path: "package.json", Key: "dependencies.@tanstack/react-query"},
			},
		},
	},
	PostGenerationCommands: []dotapi.Command{
		{Cmd: "pnpm install --dangerously-allow-all-builds"},
	},
	TestCommands: []dotapi.Command{
		{Cmd: "pnpm exec tsc --noEmit"},
		{Cmd: "pnpm exec vite build"},
	},
}
