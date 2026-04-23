package generators

import (
	clean_arch "github.com/version14/dot/generators/typescript/backend/architecture/clean-architecture"
	hex_arch "github.com/version14/dot/generators/typescript/backend/architecture/hexagonal"
	mvc_arch "github.com/version14/dot/generators/typescript/backend/architecture/mvc"
	typescript_express_generator "github.com/version14/dot/generators/typescript/backend/frameworks/express"
	typescript_frontend_react_generator "github.com/version14/dot/generators/typescript/frontend/react"
	"github.com/version14/dot/internal/scaffold"
)

type backendFrameworksMap struct {
	Express scaffold.Generator
}

type backendArchitecturesMap struct {
	CleanArchitecture scaffold.Generator
	Hexagonal         scaffold.Generator
	MVC               scaffold.Generator
}

type backendGeneratorsMap struct {
	Framework    backendFrameworksMap
	Architecture backendArchitecturesMap
}

type frontendGeneratorsMap struct {
	React scaffold.Generator
}

type TypescriptGeneratorsMap struct {
	Backend  backendGeneratorsMap
	Frontend frontendGeneratorsMap
}

var TypescriptGenerators = TypescriptGeneratorsMap{
	Backend: backendGeneratorsMap{
		Framework: backendFrameworksMap{
			Express: *typescript_express_generator.Generator,
		},
		Architecture: backendArchitecturesMap{
			CleanArchitecture: *clean_arch.CleanArchitectureTS,
			Hexagonal:         *hex_arch.HexagonalTS,
			MVC:               *mvc_arch.MvcTS,
		},
	},
	Frontend: frontendGeneratorsMap{
		React: *typescript_frontend_react_generator.Generator,
	},
}
