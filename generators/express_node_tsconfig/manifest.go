package expressnodetsconfig

import "github.com/version14/dot/pkg/dotapi"

var Manifest = dotapi.Manifest{
	Name:        "express_node_tsconfig",
	Version:     "0.1.0",
	Description: "Overrides tsconfig.json compiler options for Node.js/CommonJS Express backend",
	DependsOn:   []string{"typescript_base"},
	Outputs:     []string{},
	Validators: []dotapi.Validator{
		{
			Name: "express-node-tsconfig",
			Checks: []dotapi.Check{
				{Type: dotapi.CheckJSONKeyExists, Path: "tsconfig.json", Key: "compilerOptions.module"},
			},
		},
	},
}
