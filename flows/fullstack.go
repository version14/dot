package flows

import (
	"github.com/version14/dot/internal/flow"
	"github.com/version14/dot/internal/spec"
)

// FullstackFlow demonstrates a richer flow with branching (UI library choice)
// and an IfQuestion that gates the backend section. It produces invocations
// that include react_app for any TypeScript frontend.
//
// Question shape:
//
//	project_name (text)
//	  → stack (select) [ts | polyglot]
//	    → ui_library (if stack=ts/polyglot) [react | none]
//	      → use_biome (confirm)
//	        → confirm_generate
func FullstackFlow() *FlowDef {
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

	uiLibrary := &flow.OptionQuestion{
		QuestionBase: flow.QuestionBase{ID_: "ui_library"},
		Label:        "UI library",
		Description:  "Choose the frontend framework (or none for headless).",
		Options: []*flow.Option{
			{Label: "React + Vite", Value: "react", Next: &flow.Next{Question: useBiome}},
			{Label: "None (just TypeScript)", Value: "none", Next: &flow.Next{Question: useBiome}},
		},
	}

	stack := &flow.OptionQuestion{
		QuestionBase: flow.QuestionBase{ID_: "stack"},
		Label:        "Stack",
		Options: []*flow.Option{
			{Label: "TypeScript only", Value: "typescript", Next: &flow.Next{Question: uiLibrary}},
			{Label: "Polyglot (TS frontend + Go backend)", Value: "polyglot", Next: &flow.Next{Question: uiLibrary}},
		},
	}

	projectName := &flow.TextQuestion{
		QuestionBase: flow.QuestionBase{
			ID_:   "project_name",
			Next_: &flow.Next{Question: stack},
		},
		Label:       "Project name",
		Description: "The repo / package name.",
		Default:     "fullstack-app",
		Validate:    nonEmpty,
	}

	return &FlowDef{
		ID:          "fullstack",
		Title:       "Fullstack Application",
		Description: "TypeScript (and optional Go backend) with a configurable UI library.",
		Root:        projectName,
		Generators:  resolveFullstackGenerators,
	}
}

func resolveFullstackGenerators(s *spec.ProjectSpec) []Invocation {
	if s == nil {
		return nil
	}

	out := []Invocation{
		{Name: "base_project"},
		{Name: "typescript_base"},
	}

	if ui, _ := s.Answers["ui_library"].(string); ui == "react" {
		out = append(out, Invocation{Name: "react_app"})
	}
	if useBiome, _ := s.Answers["use_biome"].(bool); useBiome {
		out = append(out, Invocation{Name: "biome_config"})
	}
	return out
}
