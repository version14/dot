package templates

import (
	"github.com/version14/dot/internal/question"
	frontend_templates "github.com/version14/dot/internal/templates/frontend"
)

var AppTypeQuestions = question.Select("App type", "app-type").
	Choice("Frontend", "frontend", frontend_templates.FrontendQuestions).
	// Choice("Backend", "backend", backend_templates.BackendQuestions). //TODO: MAKE THIS WORK
	Q()

var AppConfigWithName = question.Text("App name", "app-name").
	Then(AppTypeQuestions).
	Q()
