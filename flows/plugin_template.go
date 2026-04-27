package flows

import (
	"strings"

	"github.com/version14/dot/internal/flow"
	"github.com/version14/dot/internal/spec"
)

// PluginTemplateFlow scaffolds a publishable DOT plugin repository. The user
// names the plugin, picks a module path, and chooses whether to include
// sample injection / generator code; the resolver wires up exactly one
// invocation of the plugin_repo_skeleton generator with all answers visible.
//
// Question shape:
//
//	project_name (text)              — also the plugin id
//	module_path  (text)              — defaults to github.com/<author>/<id>
//	plugin_description (text)
//	plugin_author (text)
//	plugin_year (text)
//	plugin_include_injection (confirm)
//	plugin_include_generator (confirm)
//	confirm_generate (confirm)
func PluginTemplateFlow() *FlowDef {
	confirmGenerate := &flow.ConfirmQuestion{
		QuestionBase: flow.QuestionBase{ID_: "confirm_generate"},
		Label:        "Scaffold the plugin repo now?",
		Default:      true,
		Then:         &flow.Next{End: true},
		Else:         &flow.Next{End: true},
	}

	includeGenerator := &flow.ConfirmQuestion{
		QuestionBase: flow.QuestionBase{ID_: "plugin_include_generator"},
		Label:        "Include a sample generator?",
		Default:      true,
		Then:         &flow.Next{Question: confirmGenerate},
		Else:         &flow.Next{Question: confirmGenerate},
	}

	includeInjection := &flow.ConfirmQuestion{
		QuestionBase: flow.QuestionBase{ID_: "plugin_include_injection"},
		Label:        "Include a sample InsertAfter injection?",
		Default:      true,
		Then:         &flow.Next{Question: includeGenerator},
		Else:         &flow.Next{Question: includeGenerator},
	}

	year := &flow.TextQuestion{
		QuestionBase: flow.QuestionBase{
			ID_:   "plugin_year",
			Next_: &flow.Next{Question: includeInjection},
		},
		Label:   "Copyright year",
		Default: "2026",
	}

	author := &flow.TextQuestion{
		QuestionBase: flow.QuestionBase{
			ID_:   "plugin_author",
			Next_: &flow.Next{Question: year},
		},
		Label:    "Author name (used in LICENSE)",
		Default:  "Anonymous",
		Validate: nonEmpty,
	}

	description := &flow.TextQuestion{
		QuestionBase: flow.QuestionBase{
			ID_:   "plugin_description",
			Next_: &flow.Next{Question: author},
		},
		Label:    "One-line description",
		Default:  "A DOT plugin",
		Validate: nonEmpty,
	}

	modulePath := &flow.TextQuestion{
		QuestionBase: flow.QuestionBase{
			ID_:   "module_path",
			Next_: &flow.Next{Question: description},
		},
		Label:       "Go module path",
		Description: "Used in go.mod and the install URL. Typically github.com/<you>/<plugin-id>.",
		Default:     "github.com/your-org/my-plugin",
		Validate:    validateModulePath,
	}

	projectName := &flow.TextQuestion{
		QuestionBase: flow.QuestionBase{
			ID_:   "project_name",
			Next_: &flow.Next{Question: modulePath},
		},
		Label:       "Plugin id (lowercase, no dots)",
		Description: "Used as the namespace prefix for every contributed ID.",
		Default:     "my-plugin",
		Validate:    validatePluginID,
	}

	return &FlowDef{
		ID:          "plugin-template",
		Title:       "Plugin Repository Template",
		Description: "Scaffold a publishable DOT plugin repo (go.mod + plugin.go + manifest + README + LICENSE).",
		Root:        projectName,
		Generators:  resolvePluginTemplateGenerators,
	}
}

func resolvePluginTemplateGenerators(_ *spec.ProjectSpec) []Invocation {
	return []Invocation{
		{Name: "plugin_repo_skeleton"},
	}
}

// validatePluginID enforces the plugin id naming rules: non-empty, no dots
// (because dots are reserved for the namespace separator), and ASCII-friendly
// to avoid filesystem surprises.
func validatePluginID(s string) error {
	s = strings.TrimSpace(s)
	if s == "" {
		return errEmpty
	}
	if strings.Contains(s, ".") {
		return errString("plugin id must not contain '.' (dot is reserved for namespacing)")
	}
	if strings.ContainsAny(s, " /\\") {
		return errString("plugin id must not contain spaces or path separators")
	}
	return nil
}

// validateModulePath does a light sanity check — Go itself is the source of
// truth for module path validity, but we catch obviously-wrong inputs early.
func validateModulePath(s string) error {
	s = strings.TrimSpace(s)
	if s == "" {
		return errEmpty
	}
	if !strings.Contains(s, "/") {
		return errString("module path must look like host/owner/repo")
	}
	return nil
}
