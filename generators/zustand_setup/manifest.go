package zustandsetup

import "github.com/version14/dot/pkg/dotapi"

var Manifest = dotapi.Manifest{
	Name:          "zustand_setup",
	Version:       "0.1.0",
	Description:   "Zustand global state management with example counter store",
	DependsOn:     []string{"typescript_base"},
	ConflictsWith: []string{"jotai_setup"},
	Outputs: []string{
		"src/stores/counter.store.ts",
	},
	PostGenerationCommands: []dotapi.Command{
		{Cmd: "pnpm install --dangerously-allow-all-builds"},
	},
	TestCommands: []dotapi.Command{
		{Cmd: "pnpm exec tsc --noEmit"},
	},
	Validators: []dotapi.Validator{
		{
			Name: "zustand-setup-files",
			Checks: []dotapi.Check{
				{Type: dotapi.CheckFileExists, Path: "src/stores/counter.store.ts"},
				{Type: dotapi.CheckJSONKeyExists, Path: "package.json", Key: "dependencies.zustand"},
			},
		},
	},
}
