package prettierconfig

import "github.com/version14/dot/pkg/dotapi"

var Manifest = dotapi.Manifest{
	Name:        "prettier_config",
	Version:     "0.1.0",
	Description: "Base Prettier configuration: .prettierrc and .prettierignore",
	DependsOn:   []string{"typescript_base"},
	Outputs: []string{
		".prettierrc",
		".prettierignore",
	},
	Validators: []dotapi.Validator{
		{
			Name: "prettier-config",
			Checks: []dotapi.Check{
				{Type: dotapi.CheckFileExists, Path: ".prettierrc"},
			},
		},
	},
}
