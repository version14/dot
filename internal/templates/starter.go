package templates

import "github.com/version14/dot/internal/question"

// StarterQuestions is the root of the dot init survey.
var StarterQuestions = question.Text("Project name", "project-name").
	Then(
		question.Select("Monorepo?", "project-monorepo").
			Description("Does your project have more than one app in a single repository?").
			// Choice("Yes", "yes", question.Loop("How many apps?", "apps-count", AppConfigWithName)). //TODO: Make this work
			Choice("No", "no", AppTypeQuestions).
			Q(),
	).
	Q()
