package authjwtcleanarchmodule

import "github.com/version14/dot/pkg/dotapi"

var Manifest = dotapi.Manifest{
	Name:        "auth_jwt_clean_arch_module",
	Version:     "0.1.0",
	Description: "JWT auth module for Clean Architecture (use-cases, controller, repository, domain)",
	DependsOn:   []string{"auth_jwt_vanilla", "auth_jwt_users_schema"},
	Outputs: []string{
		"src/modules/auth/domain/entities/user.entity.ts",
		"src/modules/auth/domain/interfaces/user.repository.interface.ts",
		"src/modules/auth/domain/interfaces/refresh-token.repository.interface.ts",
		"src/modules/auth/application/use-cases/login.use-case.ts",
		"src/modules/auth/application/use-cases/register.use-case.ts",
		"src/modules/auth/application/use-cases/refresh.use-case.ts",
		"src/modules/auth/application/use-cases/logout.use-case.ts",
		"src/modules/auth/application/controllers/auth.controller.ts",
		"src/routes/auth.route.ts",
		"src/modules/auth/infrastructure/database/repositories/user.repository.ts",
		"src/modules/auth/infrastructure/database/repositories/refresh-token.repository.ts",
	},
	PostGenerationCommands: []dotapi.Command{
		{Cmd: "pnpm install"},
	},
	Validators: []dotapi.Validator{
		{
			Name: "auth_jwt_clean_arch_module",
			Checks: []dotapi.Check{
				{Type: dotapi.CheckFileExists, Path: "src/modules/auth/domain/entities/user.entity.ts"},
				{Type: dotapi.CheckFileExists, Path: "src/modules/auth/domain/interfaces/user.repository.interface.ts"},
				{Type: dotapi.CheckFileExists, Path: "src/modules/auth/domain/interfaces/refresh-token.repository.interface.ts"},
				{Type: dotapi.CheckFileExists, Path: "src/modules/auth/application/use-cases/login.use-case.ts"},
				{Type: dotapi.CheckFileExists, Path: "src/modules/auth/application/use-cases/register.use-case.ts"},
				{Type: dotapi.CheckFileExists, Path: "src/modules/auth/application/use-cases/refresh.use-case.ts"},
				{Type: dotapi.CheckFileExists, Path: "src/modules/auth/application/use-cases/logout.use-case.ts"},
				{Type: dotapi.CheckFileExists, Path: "src/modules/auth/application/controllers/auth.controller.ts"},
				{Type: dotapi.CheckFileExists, Path: "src/routes/auth.route.ts"},
				{Type: dotapi.CheckFileExists, Path: "src/modules/auth/infrastructure/database/repositories/user.repository.ts"},
				{Type: dotapi.CheckFileExists, Path: "src/modules/auth/infrastructure/database/repositories/refresh-token.repository.ts"},
			},
		},
	},
}
