package backend_templates

import (
	"github.com/version14/dot/generators"
	clean_generator "github.com/version14/dot/generators/typescript/backend/architecture/clean-architecture"
	hexagonal_generator "github.com/version14/dot/generators/typescript/backend/architecture/hexagonal"
	mvc_generator "github.com/version14/dot/generators/typescript/backend/architecture/mvc"
	"github.com/version14/dot/internal/question"
)

// Databases

var dbSchemaQ = question.Text("Database schema", "db-schema").
	Description("Define tables, fields, constraints and relationships (reviewable later)").
	Q()

var databasesQ = question.Select("Databases", "databases").
	Multiple().
	Choice("PostgreSQL", "postgres").
	Choice("MySQL", "mysql").
	Choice("MongoDB", "mongodb").
	Choice("Redis", "redis").
	Choice("None", "none").
	Then(dbSchemaQ).
	Q()

// TypeScript tooling

var tsFormatterQ = question.Select("Formatter", "ts-formatter").
	// Choice("Prettier", "prettier", databasesQ).
	ChoiceWithGen("Biome", "biome", generators.Generators.Typescript.Linters.Biome.Func(), databasesQ).
	// Choice("None", "none", databasesQ).
	Q()

var tsLinterQ = question.Select("Linter", "ts-linter").
	Choice("ESLint", "eslint", tsFormatterQ).
	Choice("Biome", "biome", tsFormatterQ).
	Choice("None", "none", tsFormatterQ).
	Q()

var tsArchitectureQ = question.Select("Architecture pattern", "ts-architecture").
	ChoiceWithGen("Clean Architecture", "clean", clean_generator.CleanArchitectureTS.Func(), tsLinterQ).
	ChoiceWithGen("MVC", "mvc", mvc_generator.MvcTS.Func(), tsLinterQ).
	ChoiceWithGen("Hexagonal", "hexagonal", hexagonal_generator.HexagonalTS.Func(), tsLinterQ).
	Q()
