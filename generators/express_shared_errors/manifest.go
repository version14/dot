package expresssharederrors

import "github.com/version14/dot/pkg/dotapi"

var Manifest = dotapi.Manifest{
	Name:        "express_shared_errors",
	Version:     "0.1.0",
	Description: "Shared error classes: AppError base class with NotFoundError, ValidationError, UnauthorizedError, ForbiddenError, ConflictError",
	DependsOn:   []string{},
	Outputs: []string{
		"src/shared/errors/app.error.ts",
		"src/shared/errors/not-found.error.ts",
		"src/shared/errors/validation.error.ts",
		"src/shared/errors/unauthorized.error.ts",
		"src/shared/errors/forbidden.error.ts",
		"src/shared/errors/conflict.error.ts",
		"src/shared/errors/index.ts",
	},
	Validators: []dotapi.Validator{
		{
			Name: "express-shared-errors",
			Checks: []dotapi.Check{
				{Type: dotapi.CheckFileExists, Path: "src/shared/errors/app.error.ts"},
				{Type: dotapi.CheckFileExists, Path: "src/shared/errors/not-found.error.ts"},
				{Type: dotapi.CheckFileExists, Path: "src/shared/errors/validation.error.ts"},
				{Type: dotapi.CheckFileExists, Path: "src/shared/errors/unauthorized.error.ts"},
				{Type: dotapi.CheckFileExists, Path: "src/shared/errors/forbidden.error.ts"},
				{Type: dotapi.CheckFileExists, Path: "src/shared/errors/conflict.error.ts"},
				{Type: dotapi.CheckFileExists, Path: "src/shared/errors/index.ts"},
			},
		},
	},
}
