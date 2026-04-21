package backend_templates

import "github.com/version14/dot/internal/question"

var serviceLoopQ = question.Loop("How many services?", "services-count", ServiceFlow)

// Local gateway

var localGatewayGoFrameworkQ = question.Select("Gateway framework", "gateway-go-framework").
	Choice("net/http", "nethttp", serviceLoopQ).
	Choice("Gin", "gin", serviceLoopQ).
	Q()

var localGatewayTsFrameworkQ = question.Select("Gateway framework", "gateway-ts-framework").
	Choice("Express", "express", serviceLoopQ).
	Choice("Fastify", "fastify", serviceLoopQ).
	Q()

var localGatewayQ = question.Select("Gateway language", "gateway-language").
	Choice("Go", "go", localGatewayGoFrameworkQ).
	Choice("TypeScript", "typescript", localGatewayTsFrameworkQ).
	Q()

// Microservices flow

var msGatewayQ = question.Select("Gateway type", "ms-gateway-type").
	Choice("Local", "local", localGatewayQ).
	Choice("Cloud", "cloud", serviceLoopQ).
	Q()

var msEventManagerQ = question.Select("Event manager", "ms-event-manager").
	Choice("None", "none", msGatewayQ).
	Choice("NATS", "nats", msGatewayQ).
	Choice("gRPC", "grpc", msGatewayQ).
	Q()

// MicroservicesQuestions is the entry point for the microservices backend flow.
var MicroservicesQuestions = question.Select("Hosting strategy", "ms-hosting-strategy").
	Choice("Global", "global", msEventManagerQ).
	Choice("Per service", "per-service", msEventManagerQ).
	Q()
