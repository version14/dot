package reactrouterv7

import "github.com/version14/dot/pkg/dotapi"

var Manifest = dotapi.Manifest{
	Name:          "react_router_v7",
	Version:       "0.1.0",
	Description:   "React Router v7 client-side routing for React+Vite",
	DependsOn:     []string{"react_app"},
	ConflictsWith: []string{"tanstack_router"},
	Outputs: []string{
		"src/router.tsx",
		"src/pages/Home.tsx",
	},
	PostGenerationCommands: []dotapi.Command{
		{Cmd: "pnpm install --dangerously-allow-all-builds"},
	},
	TestCommands: []dotapi.Command{
		{Cmd: "pnpm exec tsc --noEmit"},
		{Cmd: "pnpm exec vite build"},
	},
	Validators: []dotapi.Validator{
		{
			Name: "react-router-v7-files",
			Checks: []dotapi.Check{
				{Type: dotapi.CheckFileExists, Path: "src/router.tsx"},
				{Type: dotapi.CheckJSONKeyExists, Path: "package.json", Key: "dependencies.react-router"},
			},
		},
	},
}
