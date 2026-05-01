package expresserrormiddleware

import "github.com/version14/dot/pkg/dotapi"

var Manifest = dotapi.Manifest{
	Name:        "express_error_middleware",
	Version:     "0.1.0",
	Description: "Global Express error-handling middleware that catches AppError instances and formats JSON responses",
	DependsOn:   []string{"express_server_entrypoint", "express_shared_errors"},
	Outputs:     []string{"src/shared/middlewares/error.middleware.ts"},
	Validators: []dotapi.Validator{
		{
			Name: "express-error-middleware",
			Checks: []dotapi.Check{
				{Type: dotapi.CheckFileExists, Path: "src/shared/middlewares/error.middleware.ts"},
			},
		},
	},
}
