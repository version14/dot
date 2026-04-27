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
	// ── Tail nodes built first so we can wire Next pointers upward ──────────

	confirmGenerate := &flow.ConfirmQuestion{
		QuestionBase: flow.QuestionBase{ID_: "confirm_generate"},
		Label:        "Generate the project now?",
		Default:      true,
		Then:         &flow.Next{End: true},
		Else:         &flow.Next{End: true},
	}

	useBiome := &flow.ConfirmQuestion{
		QuestionBase: flow.QuestionBase{ID_: "use_biome"},
		Label:        "Add Biome (lint + format)?",
		Default:      true,
		Then:         &flow.Next{Question: confirmGenerate},
		Else:         &flow.Next{Question: confirmGenerate},
	}

	useReact := &flow.ConfirmQuestion{
		QuestionBase: flow.QuestionBase{ID_: "use_react"},
		Label:        "Set up a React app?",
		Default:      true,
		Then:         &flow.Next{Question: useBiome},
		Else:         &flow.Next{Question: useBiome},
	}

	stack := &flow.OptionQuestion{
		QuestionBase: flow.QuestionBase{ID_: "stack"},
		Label:        "Primary language stack",
		Description:  "DOT will scaffold the matching toolchain.",
		Options: []*flow.Option{
			{Label: "TypeScript", Value: "typescript", Next: &flow.Next{Question: useReact}},
			{Label: "Go", Value: "go", Next: &flow.Next{Question: confirmGenerate}},
			{Label: "Polyglot (TS + Go)", Value: "polyglot", Next: &flow.Next{Question: useReact}},
		},
	}

	monorepoType := &flow.OptionQuestion{
		QuestionBase: flow.QuestionBase{ID_: "monorepo_type"},
		Label:        "Monorepo style",
		Options: []*flow.Option{
			{Label: "Single app (no monorepo)", Value: "single", Next: &flow.Next{Question: stack}},
			{Label: "Turborepo", Value: "turborepo", Next: &flow.Next{Question: stack}},
			{Label: "Nx", Value: "nx", Next: &flow.Next{Question: stack}},
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

	stack, _ := s.Answers["stack"].(string)
	wantsTS := stack == "typescript" || stack == "polyglot"

	if wantsTS {
		out = append(out, Invocation{Name: "typescript_base"})

		if useReact, _ := s.Answers["use_react"].(bool); useReact {
			out = append(out, Invocation{Name: "react_app"})
		}
		if useBiome, _ := s.Answers["use_biome"].(bool); useBiome {
			out = append(out, Invocation{Name: "biome_config"})
		}
	}

	return out
}

func nonEmpty(s string) error {
	if s == "" {
		return errEmpty
	}
	return nil
}

// errEmpty is reused so we don't allocate per validate call.
var errEmpty = errString("required")

type errString string

func (e errString) Error() string { return string(e) }
