package backend_templates

import "github.com/version14/dot/internal/question"

// BackendQuestions is the entry point for the backend app configuration flow.
var BackendQuestions = question.Select("Architecture type", "backend-architecture-type").
	Choice("Monolith", "monolith", ServiceDetailFlow).
	Choice("Microservices", "microservices", MicroservicesQuestions).
	Q()
