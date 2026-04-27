package flows

import (
	"github.com/version14/dot/internal/flow"
	"github.com/version14/dot/internal/spec"
)

// MicroservicesFlow scaffolds a project containing N service folders. The
// flow is the canonical example of a LoopQuestion in DOT: the user answers
// "how many services?" and then names each one, and the resolver expands
// those iterations into one Invocation per service via LoopFrames.
//
// Question shape:
//
//	project_name (text)
//	  → services (LoopQuestion)
//	      Body:
//	        name (text, "Service name")
//	        port (text, "Port", default 3000)
//	  → confirm_generate
func MicroservicesFlow() *FlowDef {
	confirmGenerate := &flow.ConfirmQuestion{
		QuestionBase: flow.QuestionBase{ID_: "confirm_generate"},
		Label:        "Generate the project now?",
		Default:      true,
		Then:         &flow.Next{End: true},
		Else:         &flow.Next{End: true},
	}

	servicesLoop := &flow.LoopQuestion{
		QuestionBase: flow.QuestionBase{
			ID_:   "services",
			Next_: &flow.Next{Question: confirmGenerate},
		},
		Label: "services",
		Body: []flow.Question{
			&flow.TextQuestion{
				QuestionBase: flow.QuestionBase{ID_: "name"},
				Label:        "Service name",
				Default:      "svc",
				Validate:     nonEmpty,
			},
			&flow.TextQuestion{
				QuestionBase: flow.QuestionBase{ID_: "port"},
				Label:        "Port",
				Default:      "3000",
			},
		},
		Continue: &flow.Next{Question: confirmGenerate},
	}

	projectName := &flow.TextQuestion{
		QuestionBase: flow.QuestionBase{
			ID_:   "project_name",
			Next_: &flow.Next{Question: servicesLoop},
		},
		Label:    "Project name",
		Default:  "platform",
		Validate: nonEmpty,
	}

	return &FlowDef{
		ID:          "microservices",
		Title:       "Microservices Platform",
		Description: "Scaffold a project with N service folders driven by a LoopQuestion.",
		Root:        projectName,
		Generators:  resolveMicroservicesGenerators,
	}
}

// resolveMicroservicesGenerators emits one base_project + one service_writer
// invocation per recorded service iteration. Each service_writer carries a
// LoopStack frame so FlattenScope makes the iteration's name/port visible to
// the generator's Context.Answers.
func resolveMicroservicesGenerators(s *spec.ProjectSpec) []Invocation {
	if s == nil {
		return nil
	}
	out := []Invocation{
		{Name: "base_project"},
	}

	rawIters, ok := s.Answers["services"].([]map[string]flow.AnswerNode)
	if !ok {
		// JSON / other adapters yield []interface{} of map[string]interface{}.
		// Coerce via the more general path.
		coerced := coerceLoopAnswers(s.Answers["services"])
		if coerced == nil {
			return out
		}
		rawIters = coerced
	}

	for i, iter := range rawIters {
		out = append(out, Invocation{
			Name: "service_writer",
			LoopStack: []flow.LoopFrame{{
				QuestionID: "services",
				Index:      i,
				Answers:    iter,
			}},
		})
	}
	return out
}

// coerceLoopAnswers normalizes whatever AnswerNode the loop accumulator left
// behind into the engine's preferred shape ([]map[string]flow.AnswerNode).
//
// Why the coercion: when fixtures load via JSON, arrays unmarshal as
// []interface{} and inner objects as map[string]interface{} — the type system
// can't reconcile that with the engine's typed slice without a manual hop.
func coerceLoopAnswers(raw flow.AnswerNode) []map[string]flow.AnswerNode {
	asTyped, ok := raw.([]map[string]flow.AnswerNode)
	if ok {
		return asTyped
	}
	asAny, ok := raw.([]interface{})
	if !ok {
		return nil
	}
	// flow.AnswerNode is `interface{}`, so map[string]interface{} and
	// map[string]flow.AnswerNode are the same Go type — single case suffices.
	out := make([]map[string]flow.AnswerNode, 0, len(asAny))
	for _, e := range asAny {
		if typed, ok := e.(map[string]interface{}); ok {
			out = append(out, typed)
		}
	}
	return out
}
