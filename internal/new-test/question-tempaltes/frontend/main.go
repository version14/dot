package frontend_templates

import (
	frontend_react_generator "github.com/version14/dot/internal/new-test/generators/typescript/frontend/react"
	"github.com/version14/dot/internal/new-test/registry/questions"
)

var frontendArchitectureQ = questions.Select("Frontend architecture", "frontend-architecture").
	ChoiceWithGen("Feature-sliced Design", "feature-sliced", frontend_react_generator.ReactTS.Func()).
	ChoiceWithGen("Atomic Design", "atomic", frontend_react_generator.ReactTS.Func()).
	ChoiceWithGen("Container-Presentational", "container-presentational", frontend_react_generator.ReactTS.Func()).
	Q()

var FrontendQuestions = questions.Select("Language", "frontend-language").
	Choice("TypeScript", "typescript",
		questions.Select("Framework", "frontend-framework").
			Choice("React", "react", frontendArchitectureQ).
			Choice("Next.js", "nextjs", frontendArchitectureQ).
			Q(),
	).
	Q()
