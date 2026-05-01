package expressratelimit

import "github.com/version14/dot/pkg/dotapi"

var Manifest = dotapi.Manifest{
	Name:        "express_rate_limit",
	Version:     "0.1.0",
	Description: "Adds express-rate-limit with a default 100 req/15min window applied globally",
	DependsOn:   []string{"express_server_entrypoint"},
	Outputs:     []string{"src/shared/middlewares/rate-limit.middleware.ts"},
	Validators: []dotapi.Validator{
		{
			Name: "express-rate-limit",
			Checks: []dotapi.Check{
				{Type: dotapi.CheckFileExists, Path: "src/shared/middlewares/rate-limit.middleware.ts"},
				{Type: dotapi.CheckJSONKeyExists, Path: "package.json", Key: "dependencies.express-rate-limit"},
			},
		},
	},
}
