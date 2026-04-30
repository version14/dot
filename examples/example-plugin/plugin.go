// Example external plugin for DOT.
//
// This file is the minimal-but-complete blueprint for an external plugin.
// Copy it (or use `dot scaffold plugin-template` to generate one), adjust
// the PluginID + injections + generators, and either:
//
//   - Vendor it into your DOT fork and import for side-effect from cmd/dot,
//   - Or publish it as a git repo and `dot plugin install github.com/you/it`
//     once dynamic loading lands.
//
// Conventions enforced by the engine:
//
//  1. PluginID must NOT contain '.'.
//  2. Every contributed ID (generator names, question IDs, option values)
//     must start with "<PluginID>.".
//  3. Provider methods must be safe to call multiple times — the loader
//     instantiates each plugin exactly once but calls Generators(),
//     Injections(), and ResolveExtras() during normal flow runs.
package exampleplugin

import (
	"github.com/version14/dot/pkg/dotapi"
	"github.com/version14/dot/pkg/dotplugin"
)

// PluginID is the namespace prefix for everything this plugin contributes.
const PluginID dotplugin.PluginID = "example"

func init() {
	dotplugin.RegisterBuiltin(&Provider{})
}

// Provider is the loader entry point.
type Provider struct{}

func (Provider) ID() dotplugin.PluginID { return PluginID }

// Generators returns the (Manifest, Generator) pairs to register.
func (Provider) Generators() []dotplugin.Entry {
	return []dotplugin.Entry{
		{Manifest: editorConfigManifest, Generator: &editorConfigWriter{}},
	}
}

// Injections demonstrates two of the three injection kinds:
//
//   - InsertAfter: adds an "Add .editorconfig?" confirm question after the
//     base flow's "use_biome" question (chosen because it appears in every
//     in-tree flow that touches TypeScript).
//   - AddOption: adds an "Editor: VSCode-only" option to the "stack"
//     OptionQuestion. The option's Next is End so picking it short-circuits
//     the rest of the flow (illustrates how plugins can carve sub-paths).
func (Provider) Injections() []*dotplugin.Injection {
	addEditorConfig := &dotplugin.ConfirmQuestion{
		QuestionBase: dotplugin.QuestionBase{ID_: "example.add_editorconfig"},
		Label:        "Add .editorconfig for cross-IDE consistency?",
		Default:      true,
		Then:         &dotplugin.Next{End: true},
		Else:         &dotplugin.Next{End: true},
	}

	return []*dotplugin.Injection{
		{
			Plugin:   PluginID,
			TargetID: "use_biome",
			Kind:     dotplugin.InjectInsertAfter,
			Question: addEditorConfig,
		},
		{
			Plugin:   PluginID,
			TargetID: "stack",
			Kind:     dotplugin.InjectAddOption,
			Option: &dotplugin.Option{
				Label: "VSCode workspace only",
				Value: "example.vscode_only",
				Next:  &dotplugin.Next{End: true},
			},
		},
	}
}

func (Provider) ResolveExtras(s *dotplugin.ProjectSpec) []dotplugin.Invocation {
	if s == nil {
		return nil
	}
	if want, _ := s.Answers["example.add_editorconfig"].(bool); !want {
		return nil
	}
	return []dotplugin.Invocation{{Name: "example.editorconfig_writer"}}
}

// ── Generator: example.editorconfig_writer ─────────────────────────────────

const editorConfigFileName = ".editorconfig"

var editorConfigManifest = dotapi.Manifest{
	Name:        "example.editorconfig_writer",
	Version:     "0.1.0",
	Description: "Writes a sensible .editorconfig at project root",
	DependsOn:   []string{"base_project"},
	Outputs:     []string{editorConfigFileName},
	Validators: []dotapi.Validator{
		{
			Name: "editorconfig-present",
			Checks: []dotapi.Check{
				{Type: dotapi.CheckFileExists, Path: editorConfigFileName},
			},
		},
	},
}

type editorConfigWriter struct{}

func (g *editorConfigWriter) Name() string    { return editorConfigManifest.Name }
func (g *editorConfigWriter) Version() string { return editorConfigManifest.Version }

func (g *editorConfigWriter) Generate(ctx *dotapi.Context) error {
	const editorConfig = `# https://editorconfig.org/
root = true

[*]
end_of_line = lf
insert_final_newline = true
charset = utf-8
indent_style = space
indent_size = 2
trim_trailing_whitespace = true

[*.md]
trim_trailing_whitespace = false
`
	ctx.State.WriteFile(".editorconfig", []byte(editorConfig), dotplugin.ContentRaw)
	return nil
}
