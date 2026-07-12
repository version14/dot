package vitesttestinglibrary

import "github.com/version14/dot/pkg/dotapi"

var Manifest = dotapi.Manifest{
	Name:        "vitest_testing_library",
	Version:     "0.2.1",
	Description: "Vitest + React Testing Library for React+Vite projects",
	DependsOn:   []string{"react_app"},
	Outputs: []string{
		"vitest.config.ts",
		"src/test/setup.ts",
		"src/test/App.test.tsx",
	},
	PostGenerationCommands: []dotapi.Command{
		{Cmd: "pnpm install --dangerously-allow-all-builds"},
	},
	TestCommands: []dotapi.Command{
		{Cmd: "pnpm exec vitest run"},
	},
	Validators: []dotapi.Validator{
		{
			Name: "vitest-testing-library-files",
			Checks: []dotapi.Check{
				{Type: dotapi.CheckFileExists, Path: "vitest.config.ts"},
				{Type: dotapi.CheckJSONKeyExists, Path: "package.json", Key: "devDependencies.vitest"},
				{Type: dotapi.CheckJSONKeyExists, Path: "package.json", Key: "scripts.test"},
			},
		},
	},
}
