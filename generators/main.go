package generators

import (
	common_generators "github.com/version14/dot/generators/common"
	typescript_generators "github.com/version14/dot/generators/typescript"
)

type GeneratorsMapping struct {
	Common     common_generators.CommonGeneratorsMap
	Typescript typescript_generators.TypescriptGeneratorsMap
}

var Generators = GeneratorsMapping{
	Common:     common_generators.CommonGenerators,
	Typescript: typescript_generators.TypescriptGenerators,
}
