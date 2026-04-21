package frontend_templates

import (
	frontend_react_generator "github.com/version14/dot/generators/typescript/frontend/react"
	"github.com/version14/dot/internal/question"
)

var frontendArchitectureQ = question.Select("Frontend architecture", "frontend-architecture").
	ChoiceWithGen("Feature-sliced Design", "feature-sliced", frontend_react_generator.ReactTS.Func()).
	ChoiceWithGen("Atomic Design", "atomic", frontend_react_generator.ReactTS.Func()).
	ChoiceWithGen("Container-Presentational", "container-presentational", frontend_react_generator.ReactTS.Func()).
	Q()

var FrontendQuestions = question.Select("Language", "frontend-language").
	Choice("TypeScript", "typescript",
		question.Select("Framework", "frontend-framework").
			Choice("React", "react", frontendArchitectureQ).
			// Choice("Next.js", "nextjs", frontendArchitectureQ). //TODO: MAKE THIS WORK
			Q(),
	).
	Q()
