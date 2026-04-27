// Package pluginreposkeleton scaffolds a publishable DOT plugin repository.
//
// The output is a complete Go module (go.mod + plugin.go + plugin.json +
// README + LICENSE + .gitignore) that someone else can install with
//
//	dot plugin install github.com/<your-author>/<plugin-id>
//
// Optional sample injection / sample generator code is included based on
// the user's answers, so the generated repo is either a minimal skeleton
// or a fully working example.
package pluginreposkeleton

import (
	"fmt"
	"strings"

	"github.com/version14/dot/internal/render"
	"github.com/version14/dot/internal/state"
	"github.com/version14/dot/pkg/dotapi"
)

// dotModulePath is the import path plugin authors use to depend on dot's
// public API. Pinned to a recent default; users can edit go.mod afterward.
const dotModulePath = "github.com/version14/dot"

// dotModuleVersion is the semver pinned in the generated go.mod. Authors will
// typically `go get -u github.com/version14/dot@latest` after scaffolding.
const dotModuleVersion = "v0.1.0"

type Generator struct{}

func New() *Generator { return &Generator{} }

func (g *Generator) Name() string    { return Manifest.Name }
func (g *Generator) Version() string { return Manifest.Version }

// Generate builds the plugin repo skeleton from the populated answers.
//
// Required answers:
//
//	project_name (string)            — also used as the plugin id
//	module_path  (string)            — go module path (typically github.com/<author>/<id>)
//	plugin_description (string)      — human-readable summary
//	plugin_author (string)           — author name (for LICENSE)
//	plugin_include_injection (bool)  — emit a sample InsertAfter injection?
//	plugin_include_generator (bool)  — emit a sample generator?
func (g *Generator) Generate(ctx *dotapi.Context) error {
	id := stringAnswer(ctx.Answers, "project_name", ctx.Spec.Metadata.ProjectName)
	if id == "" {
		return fmt.Errorf("plugin_repo_skeleton: missing project_name")
	}
	if strings.Contains(id, ".") {
		return fmt.Errorf("plugin_repo_skeleton: plugin id %q must not contain '.'", id)
	}

	modulePath := stringAnswer(ctx.Answers, "module_path", "github.com/your-org/"+id)
	desc := stringAnswer(ctx.Answers, "plugin_description", "A DOT plugin")
	author := stringAnswer(ctx.Answers, "plugin_author", "Anonymous")
	year := stringAnswer(ctx.Answers, "plugin_year", "2026")

	includeInjection, _ := ctx.Answers["plugin_include_injection"].(bool)
	includeGenerator, _ := ctx.Answers["plugin_include_generator"].(bool)

	data := map[string]interface{}{
		"PluginID":         id,
		"PackageName":      packageNameFor(id),
		"ModulePath":       modulePath,
		"Description":      desc,
		"Author":           author,
		"Year":             year,
		"IncludeInjection": includeInjection,
		"IncludeGenerator": includeGenerator,
		"DotModule":        dotModulePath,
		"DotVersion":       dotModuleVersion,
	}

	if err := writeRendered(ctx.State, "go.mod", goModTmpl, data, state.ContentRaw); err != nil {
		return err
	}
	if err := writeRendered(ctx.State, "plugin.json", pluginJSONTmpl, data, state.ContentJSON); err != nil {
		return err
	}
	if err := writeRendered(ctx.State, "plugin.go", pluginGoTmpl, data, state.ContentRaw); err != nil {
		return err
	}
	if err := writeRendered(ctx.State, "README.md", readmeTmpl, data, state.ContentRaw); err != nil {
		return err
	}
	if err := writeRendered(ctx.State, "LICENSE", licenseTmpl, data, state.ContentRaw); err != nil {
		return err
	}
	ctx.State.WriteFile(".gitignore", []byte(gitignoreContent), state.ContentRaw)
	return nil
}

func writeRendered(s *state.VirtualProjectState, path, tmpl string, data interface{}, ct state.ContentType) error {
	out, err := render.Render(tmpl, data)
	if err != nil {
		return fmt.Errorf("plugin_repo_skeleton: render %s: %w", path, err)
	}
	s.WriteFile(path, out, ct)
	return nil
}

// packageNameFor turns a plugin id ("biome_extras") into a valid Go package
// identifier ("biome_extras" → "biomeextras"). Hyphens are also stripped.
func packageNameFor(id string) string {
	out := strings.ReplaceAll(id, "-", "")
	out = strings.ReplaceAll(out, "_", "")
	return strings.ToLower(out)
}

func stringAnswer(answers map[string]interface{}, key, fallback string) string {
	if v, ok := answers[key]; ok {
		if s, ok := v.(string); ok && s != "" {
			return s
		}
	}
	return fallback
}

// ── Templates ─────────────────────────────────────────────────────────────

const goModTmpl = `module {{.ModulePath}}

go 1.26

require {{.DotModule}} {{.DotVersion}}
`

const pluginJSONTmpl = `{
  "id": "{{.PluginID}}",
  "version": "0.1.0",
  "description": "{{.Description}}",
  "entry_point": "plugin.go"
}
`

const pluginGoTmpl = `// Package {{.PackageName}} implements the "{{.PluginID}}" DOT plugin.
//
// Install with:
//
//	dot plugin install {{.ModulePath}}
//
// See https://github.com/version14/dot for the plugin authoring guide.
package {{.PackageName}}

import (
	"github.com/version14/dot/pkg/dotapi"
	"github.com/version14/dot/pkg/dotplugin"
)

// PluginID is the namespace prefix for everything this plugin contributes.
// Per dot's convention, every contributed ID (generator names, question IDs,
// option values) MUST start with "{{.PluginID}}.".
const PluginID dotplugin.PluginID = "{{.PluginID}}"

func init() {
	dotplugin.RegisterBuiltin(&Provider{})
}

// Provider is this plugin's loader entry point.
type Provider struct{}

func (Provider) ID() dotplugin.PluginID { return PluginID }

func (Provider) Generators() []dotplugin.Entry {
{{- if .IncludeGenerator }}
	return []dotplugin.Entry{
		{Manifest: sampleManifest, Generator: &sampleGen{}},
	}
{{- else }}
	return nil
{{- end }}
}

func (Provider) Injections() []*dotplugin.Injection {
{{- if .IncludeInjection }}
	// Sample InsertAfter: adds an extra confirm question after "use_biome".
	// Replace this with whatever hooks your plugin actually wants.
	confirm := &dotplugin.ConfirmQuestion{
		QuestionBase: dotplugin.QuestionBase{ID_: "{{.PluginID}}.enabled"},
		Label:        "Enable {{.PluginID}}?",
		Default:      true,
		Then:         &dotplugin.Next{End: true},
		Else:         &dotplugin.Next{End: true},
	}
	return []*dotplugin.Injection{
		{
			Plugin:   PluginID,
			TargetID: "use_biome",
			Kind:     dotplugin.InjectInsertAfter,
			Question: confirm,
		},
	}
{{- else }}
	return nil
{{- end }}
}

func (Provider) ResolveExtras(s *dotplugin.ProjectSpec) []dotplugin.Invocation {
{{- if and .IncludeInjection .IncludeGenerator }}
	if s == nil {
		return nil
	}
	if enabled, _ := s.Answers["{{.PluginID}}.enabled"].(bool); !enabled {
		return nil
	}
	return []dotplugin.Invocation{{"{"}}{Name: "{{.PluginID}}.sample"}}
{{- else }}
	return nil
{{- end }}
}

{{- if .IncludeGenerator }}

// ── Generator: {{.PluginID}}.sample ─────────────────────────────────────

var sampleManifest = dotapi.Manifest{
	Name:        "{{.PluginID}}.sample",
	Version:     "0.1.0",
	Description: "Sample generator scaffolded by {{.PluginID}}",
	DependsOn:   []string{"base_project"},
	Outputs:     []string{"{{.PluginID}}.txt"},
}

type sampleGen struct{}

func (g *sampleGen) Name() string    { return sampleManifest.Name }
func (g *sampleGen) Version() string { return sampleManifest.Version }

func (g *sampleGen) Generate(ctx *dotapi.Context) error {
	ctx.State.WriteFile("{{.PluginID}}.txt", []byte("hello from {{.PluginID}}\n"), dotplugin.ContentRaw)
	return nil
}
{{- end }}
`

const readmeTmpl = `# {{.PluginID}}

> {{.Description}}

A [DOT](https://github.com/version14/dot) plugin.

## Installation

` + "```" + `bash
dot plugin install {{.ModulePath}}
` + "```" + `

To pin a specific version:

` + "```" + `bash
dot plugin install {{.ModulePath}}@v0.1.0
` + "```" + `

After install, rebuild ` + "`dot`" + ` (or restart it) so the plugin's ` + "`init()`" + ` runs.

## Development

` + "```" + `bash
git clone {{.ModulePath}}
cd $(basename {{.ModulePath}})
go mod tidy
` + "```" + `

Iterate locally with:

` + "```" + `bash
dot plugin install -from .
` + "```" + `

## What this plugin does

{{ if .IncludeInjection }}- Adds a question after the host flow's ` + "`use_biome`" + ` step (sample InsertAfter injection).
{{ end }}{{ if .IncludeGenerator }}- Registers a generator that writes ` + "`{{.PluginID}}.txt`" + ` to the project root.
{{ end }}{{ if not (or .IncludeInjection .IncludeGenerator) }}- Replace this list with what your plugin actually does.
{{ end }}

## License

MIT — see [LICENSE](./LICENSE).
`

const licenseTmpl = `MIT License

Copyright (c) {{.Year}} {{.Author}}

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND.
`

const gitignoreContent = `# Editor / OS
.DS_Store
*.swp
.idea/
.vscode/

# Build artifacts
dist/
*.log
`
