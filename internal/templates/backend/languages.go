package backend_templates

import "github.com/version14/dot/internal/question"

// tsFrameworkQ routes TypeScript services to their framework-specific flow.
// Express gets architecture selection; NestJS skips it (opinionated by design).
var tsFrameworkQ = question.Select("Framework", "ts-framework").
	Choice("Express", "express", tsArchitectureQ).
	Choice("NestJS", "nestjs", tsLinterQ).
	Q()
