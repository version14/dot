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
	"embed"
	"fmt"
	"strings"

	"github.com/version14/dot/internal/render"
	"github.com/version14/dot/pkg/dotapi"
)

// dotModulePath is the import path plugin authors use to depend on dot's
// public API. Pinned to a recent default; users can edit go.mod afterward.
const dotModulePath = "github.com/version14/dot"

// dotModuleVersion is the semver pinned in the generated go.mod. Authors will
// typically `go get -u github.com/version14/dot@latest` after scaffolding.
const dotModuleVersion = "v0.1.6"

type Generator struct{}

func New() *Generator { return &Generator{} }

func (g *Generator) Name() string    { return Manifest.Name }
func (g *Generator) Version() string { return Manifest.Version }

//go:embed all:files
var fs embed.FS

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

	renderer := render.NewLocalFolderRenderer(ctx.State)
	return renderer.Render(fs, data)
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
