package jotaisetup

import "github.com/version14/dot/pkg/dotapi"

var Manifest = dotapi.Manifest{
	Name:          "jotai_setup",
	Version:       "0.1.0",
	Description:   "Jotai atomic state management with example counter atom",
	DependsOn:     []string{"typescript_base"},
	ConflictsWith: []string{"zustand_setup"},
	Outputs: []string{
		"src/atoms/counter.atom.ts",
	},
	PostGenerationCommands: []dotapi.Command{
		{Cmd: "pnpm install --dangerously-allow-all-builds"},
	},
	TestCommands: []dotapi.Command{
		{Cmd: "pnpm exec tsc --noEmit"},
	},
	Validators: []dotapi.Validator{
		{
			Name: "jotai-setup-files",
			Checks: []dotapi.Check{
				{Type: dotapi.CheckFileExists, Path: "src/atoms/counter.atom.ts"},
				{Type: dotapi.CheckJSONKeyExists, Path: "package.json", Key: "dependencies.jotai"},
			},
		},
	},
}
