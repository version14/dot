package registry

import "github.com/version14/dot/internal/new-test/registry/questions"

// PluginKind describes how a plugin attaches to a flow.
type PluginKind string

const (
	// PluginKindSubflow adds an option to an existing OptionQuestion and
	// routes it into the plugin's Flow when selected. Use this when your
	// plugin extends an existing decision point (e.g. adding "Rust" to the
	// service-language question).
	PluginKindSubflow PluginKind = "subflow"

	// PluginKindNewFlow adds a brand-new top-level variant at the app-type
	// question. Use this when your plugin introduces a new app category
	// (e.g. adding "Mobile" alongside Frontend/Backend).
	PluginKindNewFlow PluginKind = "new-flow"
)

// Plugin is a declarative description of how to modify a Template.
// It is pure data: applying it is the Template's responsibility.
type Plugin struct {
	Name        string
	Kind        PluginKind
	Language    string // optional scope filter (e.g. "go", "typescript")
	AttachTo    string // question key to attach to (ignored for NewFlow — always "app-type")
	OptionLabel string
	OptionValue string
	Flow        *questions.Question
}

// AttachmentKey resolves the question Value the plugin targets.
func (plugin *Plugin) AttachmentKey() string {
	if plugin.Kind == PluginKindNewFlow {
		return "app-type"
	}
	return plugin.AttachTo
}

// ApplyPlugin attaches the plugin's flow to the given question tree.
// For shared package-level vars, call questions.DeepCopy first.
func ApplyPlugin(flow *questions.Question, plugin *Plugin) *questions.Question {
	return flow.AttachFlow(plugin.AttachmentKey(), plugin.OptionLabel, plugin.OptionValue, plugin.Flow)
}

// ApplyRegistry applies every plugin in the registry in order.
func ApplyRegistry(flow *questions.Question, registry *Registry) *questions.Question {
	for i := range registry.Plugins {
		flow = ApplyPlugin(flow, &registry.Plugins[i])
	}
	return flow
}
