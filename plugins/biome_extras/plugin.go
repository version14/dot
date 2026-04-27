// Package biomeextras is the canonical in-tree demo plugin. It exists so
// every part of DOT's plugin contract — namespacing, InsertAfter injection,
// Generator registration, ResolveExtras gating — has a real-world example.
//
// The plugin contributes one extra question after `use_biome` ("Enable Biome
// strict mode?") and one generator (biome_extras.strict_writer) that runs
// when the user answers true. Selecting "no" leaves the project unchanged.
//
// Note: this plugin imports only pkg/dotapi and pkg/dotplugin — no internal/*
// reach-through. It is shape-identical to any plugin you would publish to a
// public git repo.
package biomeextras

import (
	"github.com/version14/dot/pkg/dotapi"
	"github.com/version14/dot/pkg/dotplugin"
)

// PluginID is the namespace prefix for everything this plugin contributes.
const PluginID dotplugin.PluginID = "biome_extras"

// init registers the plugin with the loader at program startup. The CLI's
// main.go imports this package for its side effect.
func init() {
	dotplugin.RegisterBuiltin(&Provider{})
}

// Provider is the plugin's surface for the DOT loader.
type Provider struct{}

func (Provider) ID() dotplugin.PluginID { return PluginID }

func (Provider) Generators() []dotplugin.Entry {
	return []dotplugin.Entry{
		{Manifest: strictWriterManifest, Generator: &strictWriter{}},
	}
}

// Injections lists the flow hooks this plugin installs. We add one
// InsertAfter hook on "use_biome" so the strict-mode question shows up
// immediately after the host flow's biome question.
func (Provider) Injections() []*dotplugin.Injection {
	strictMode := &dotplugin.ConfirmQuestion{
		QuestionBase: dotplugin.QuestionBase{ID_: "biome_extras.strict_mode"},
		Label:        "Enable Biome strict mode? (catches more issues but is noisier)",
		Default:      false,
		Then:         &dotplugin.Next{End: true},
		Else:         &dotplugin.Next{End: true},
	}

	return []*dotplugin.Injection{
		{
			Plugin:   PluginID,
			TargetID: "use_biome",
			Kind:     dotplugin.InjectInsertAfter,
			Question: strictMode,
		},
	}
}

// ResolveExtras returns the strict_writer invocation only when both the host
// flow's "use_biome" answer and our own "biome_extras.strict_mode" answer are
// true. Otherwise it returns nil so the generator is skipped entirely.
func (Provider) ResolveExtras(s *dotplugin.ProjectSpec) []dotplugin.Invocation {
	if s == nil {
		return nil
	}
	if useBiome, _ := s.Answers["use_biome"].(bool); !useBiome {
		return nil
	}
	if strict, _ := s.Answers["biome_extras.strict_mode"].(bool); !strict {
		return nil
	}
	return []dotplugin.Invocation{{Name: "biome_extras.strict_writer"}}
}

// ── Generator: biome_extras.strict_writer ──────────────────────────────────

var strictWriterManifest = dotapi.Manifest{
	Name:        "biome_extras.strict_writer",
	Version:     "0.1.0",
	Description: "Promotes selected Biome lint rules to errors when strict mode is on",
	DependsOn:   []string{"biome_config"},
	Outputs:     []string{"biome.json"},
	Validators: []dotapi.Validator{
		{
			Name: "biome-strict",
			Checks: []dotapi.Check{
				{Type: dotapi.CheckJSONKeyExists, Path: "biome.json", Key: "linter.rules.style"},
			},
		},
	},
}

type strictWriter struct{}

func (g *strictWriter) Name() string    { return strictWriterManifest.Name }
func (g *strictWriter) Version() string { return strictWriterManifest.Version }

func (g *strictWriter) Generate(ctx *dotapi.Context) error {
	return ctx.State.UpdateJSON("biome.json", func(d *dotplugin.JSONDoc) error {
		d.Merge(map[string]interface{}{
			"linter": map[string]interface{}{
				"rules": map[string]interface{}{
					"style": map[string]interface{}{
						"useImportType":      "error",
						"noNonNullAssertion": "error",
					},
					"suspicious": map[string]interface{}{
						"noExplicitAny": "error",
					},
				},
			},
		})
		return nil
	})
}
