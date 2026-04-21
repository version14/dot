package templates

import (
	backend_templates "github.com/version14/dot/internal/new-test/question-tempaltes/backend"
	"github.com/version14/dot/internal/new-test/registry/questions"
)

var AppTypeQuestions = questions.Select("App type", "app-type").
	// Choice("Frontend", "frontend", frontend_templates.FrontendQuestions).
	Choice("Backend", "backend", backend_templates.BackendQuestions).
	Q()

var AppConfigWithName = questions.Text("App name", "app-name").
	Then(AppTypeQuestions).
	Q()
