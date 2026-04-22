package backend_templates

import (
	typescript_base_generator "github.com/version14/dot/generators/common/typescript/base"
	"github.com/version14/dot/internal/question"
)

// ServiceDetailFlow is the shared core: language → framework → arch →
// linter → formatter → databases → db-schema.
// Used by both monolith (directly) and each microservice (inside a loop).
var ServiceDetailFlow = question.Select("Language", "service-language").
	// Choice("Go", "go", goArchitectureQ). // TODO: Need to implement
	ChoiceWithGen("TypeScript", "typescript", typescript_base_generator.BaseTypescriptTS.Func(), tsFrameworkQ).
	Q()

// perServiceHostingQ is asked only when the user picked "per-service" hosting
// strategy earlier. After selection it continues into ServiceDetailFlow.
var perServiceHostingQ = question.Select("Hosting for this service", "service-hosting").
	Choice("GCP", "gcp").
	Choice("AWS", "aws").
	Choice("Docker", "docker").
	Choice("Other", "other").
	Then(ServiceDetailFlow).
	Q()

// ServiceFlow adds a service name + conditional per-service hosting before
// ServiceDetailFlow. Used inside the microservices service loop.
var ServiceFlow = question.Text("Service name", "service-name").
	Then(
		question.If("ms-hosting-strategy", question.ComparisonEqual, "per-service").
			Then(perServiceHostingQ).
			Else(ServiceDetailFlow).
			Q(),
	).
	Q()
