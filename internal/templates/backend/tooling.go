package backend_templates

import "github.com/version14/dot/internal/question"

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

// Go tooling

var goFormatterQ = question.Select("Formatter", "go-formatter").
	Choice("gofmt", "gofmt", databasesQ).
	Choice("None", "none", databasesQ).
	Q()

var goLinterQ = question.Select("Linter", "go-linter").
	Choice("golangci-lint", "golangci-lint", goFormatterQ).
	Choice("None", "none", goFormatterQ).
	Q()

var goArchitectureQ = question.Select("Architecture pattern", "go-architecture").
	Choice("Clean Architecture", "clean", goLinterQ).
	Choice("MVC", "mvc", goLinterQ).
	Choice("Hexagonal", "hexagonal", goLinterQ).
	Q()

// TypeScript tooling

var tsFormatterQ = question.Select("Formatter", "ts-formatter").
	Choice("Prettier", "prettier", databasesQ).
	Choice("Biome", "biome", databasesQ).
	Choice("None", "none", databasesQ).
	Q()

var tsLinterQ = question.Select("Linter", "ts-linter").
	Choice("ESLint", "eslint", tsFormatterQ).
	Choice("Biome", "biome", tsFormatterQ).
	Choice("None", "none", tsFormatterQ).
	Q()

var tsArchitectureQ = question.Select("Architecture pattern", "ts-architecture").
	Choice("Clean Architecture", "clean", tsLinterQ).
	Choice("MVC", "mvc", tsLinterQ).
	Choice("Hexagonal", "hexagonal", tsLinterQ).
	Q()
