package backend_templates

import (
	express_generator "github.com/version14/dot/generators/typescript/backend/frameworks/express"
	"github.com/version14/dot/internal/question"
)

// tsFrameworkQ routes TypeScript services to their framework-specific flow.
// Express gets architecture selection; NestJS skips it (opinionated by design).
var tsFrameworkQ = question.Select("Framework", "ts-framework").
	ChoiceWithGen("Express", "express", express_generator.ExpressTS.Func(), tsArchitectureQ).
	// ChoiceWithGen("NestJS", "nestjs", ..., tsLinterQ). // TODO: Need to implement
	Q()
