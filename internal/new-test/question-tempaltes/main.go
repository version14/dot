package templates

import "github.com/version14/dot/internal/new-test/registry/questions"

// StarterQuestions is the root of the dot init survey.
var StarterQuestions = questions.Text("Project name", "project-name").
	Then(
		questions.Select("Monorepo?", "project-monorepo").
			Description("Does your project have more than one app in a single repository?").
			Choice("Yes", "yes", questions.Loop("How many apps?", "apps-count", AppConfigWithName)).
			Choice("No", "no", AppTypeQuestions).
			Q(),
	).
	Q()
