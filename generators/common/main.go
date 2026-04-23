package generators

import (
	frontend_architecture_generator "github.com/version14/dot/generators/common/frontend/architecture"
	typescript_base_generator "github.com/version14/dot/generators/common/typescript/base"
	"github.com/version14/dot/internal/scaffold"
)

type CommonGeneratorsMap struct {
	TypescriptBase       *scaffold.Generator
	FrontendArchitecture *scaffold.Generator
}

var CommonGenerators = CommonGeneratorsMap{
	TypescriptBase:       typescript_base_generator.BaseTypescriptTS,
	FrontendArchitecture: frontend_architecture_generator.ArchitectureTS,
}
