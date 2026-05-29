package tanstackrouter

import "github.com/version14/dot/pkg/dotapi"

var Manifest = dotapi.Manifest{
	Name:          "tanstack_router",
	Version:       "0.2.0",
	Description:   "TanStack Router type-safe file-based routing for React+Vite",
	DependsOn:     []string{"react_app"},
	ConflictsWith: []string{"react_router_v7"},
	Outputs: []string{
		"src/routes/__root.tsx",
		"src/routes/index.tsx",
	},
	PostGenerationCommands: []dotapi.Command{
		{Cmd: "pnpm install --dangerously-allow-all-builds"},
		{Cmd: "pnpm exec vite build"},
	},
	TestCommands: []dotapi.Command{
		{Cmd: "pnpm exec tsc --noEmit"},
		{Cmd: "pnpm exec vite build"},
	},
	Validators: []dotapi.Validator{
		{
			Name: "tanstack-router-files",
			Checks: []dotapi.Check{
				{Type: dotapi.CheckFileExists, Path: "src/routes/__root.tsx"},
				{Type: dotapi.CheckJSONKeyExists, Path: "package.json", Key: "dependencies.@tanstack/react-router"},
			},
		},
	},
}
