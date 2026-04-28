package flows

import (
	"github.com/version14/dot/internal/flow"
	"github.com/version14/dot/internal/spec"
)

// MonorepoFlow is the default DOT scaffolding flow. It walks the user through
// project name → monorepo structure → language stack → linting choices.
//
// Question IDs are kept stable: re-runs of `dot scaffold` reuse the persisted
// answers from .dot/spec.json keyed by these IDs.
func MonorepoFlow() *FlowDef {
	confirmGenerate := &flow.ConfirmQuestion{
		QuestionBase: flow.QuestionBase{ID_: "confirm-generate"},
		Label:        "Generate the project now?",
		Default:      true,
		Then:         &flow.Next{End: true},
		Else:         &flow.Next{End: true},
	}

	linter := &flow.OptionQuestion{
		QuestionBase: flow.QuestionBase{ID_: "ts-backend-linter"},
		Label:        "Choose a linter.",
		Options: []*flow.Option{
			{Label: "Biome", Value: "biome", Next: &flow.Next{Question: confirmGenerate}},
			{Label: "Prettier", Value: "prettier", Next: &flow.Next{Question: confirmGenerate}},
		},
	}

	formatter := &flow.OptionQuestion{
		QuestionBase: flow.QuestionBase{ID_: "ts-backend-formatter"},
		Label:        "Choose a formatter.",
		Options: []*flow.Option{
			{Label: "Biome", Value: "biome", Next: &flow.Next{Question: linter}},
			{Label: "Prettier", Value: "prettier", Next: &flow.Next{Question: linter}},
		},
	}

	architecture := &flow.OptionQuestion{
		QuestionBase: flow.QuestionBase{ID_: "ts-backend-architecture"},
		Label:        "Choose your architecture.",
		Options: []*flow.Option{
			{Label: "Clean Architecture", Value: "clean-architecture", Next: &flow.Next{Question: formatter}},
			{Label: "MVC", Value: "mvc-architecture", Next: &flow.Next{Question: formatter}},
			// {Label: "Hexagonal", Value: "hexagonal-architecture", Next: &flow.Next{Question: formatter}},
		},
	}

	framework := &flow.OptionQuestion{
		QuestionBase: flow.QuestionBase{ID_: "ts-backend-framework"},
		Label:        "Library / Framework",
		Description:  "Choose a library or framework to scaffold your backend.",
		Options: []*flow.Option{
			{Label: "Express", Value: "express", Next: &flow.Next{Question: architecture}},
		},
	}

	stack := &flow.OptionQuestion{
		QuestionBase: flow.QuestionBase{ID_: "stack"},
		Label:        "Primary language stack",
		Description:  "DOT will scaffold the matching toolchain.",
		Options: []*flow.Option{
			{Label: "TypeScript", Value: "typescript", Next: &flow.Next{Question: framework}},
			// {Label: "Go", Value: "go", Next: &flow.Next{Question: confirmGenerate}},
		},
	}

	monorepoType := &flow.OptionQuestion{
		QuestionBase: flow.QuestionBase{ID_: "monorepo_type"},
		Label:        "Monorepo style",
		Options: []*flow.Option{
			{Label: "Single app (no monorepo)", Value: "single", Next: &flow.Next{Question: stack}},
			// {Label: "Turborepo", Value: "turborepo", Next: &flow.Next{Question: stack}},
		},
	}

	projectName := &flow.TextQuestion{
		QuestionBase: flow.QuestionBase{
			ID_:   "project_name",
			Next_: &flow.Next{Question: monorepoType},
		},
		Label:       "Project name",
		Description: "Used as the package name and root directory.",
		Default:     "my-project",
		Validate:    nonEmpty,
	}

	return &FlowDef{
		ID:          "monorepo",
		Title:       "Monorepo / Project Wizard",
		Description: "Scaffold a new project with optional monorepo, language, and tooling.",
		Root:        projectName,
		Generators:  resolveMonorepoGenerators,
	}
}

// resolveMonorepoGenerators maps the populated spec to the ordered generator
// invocations. Order is significant: dependents come after their deps.
func resolveMonorepoGenerators(s *spec.ProjectSpec) []Invocation {
	if s == nil {
		return nil
	}

	out := []Invocation{
		{Name: "base_project"},
	}

	if stack, _ := s.Answers["stack"].(string); stack == "typescript" {
		out = append(out, Invocation{Name: "typescript_base"})
	}
	if architecture, _ := s.Answers["ts-backend-architecture"].(string); architecture == "clean-architecture" {
		out = append(out, Invocation{Name: "backend_architecture_clean_architecture"})
	} else if architecture == "mvc-architecture" {
		out = append(out, Invocation{Name: "backend_architecture_mvc"})
	}

	return out
}

func nonEmpty(s string) error {
	if s == "" {
		return errEmpty
	}
	return nil
}

// errEmpty isj reused so we don't allocate per validate call.
var errEmpty = errString("required")

type errString string

func (e errString) Error() string { return string(e) }
