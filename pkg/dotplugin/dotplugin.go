// Package dotplugin is the public author-facing API for DOT plugins.
//
// External plugin repositories — the kind you publish on GitHub and let
// other people install via `dot plugin install github.com/you/your-plugin` —
// import this package together with pkg/dotapi. Every type, constant, and
// helper a plugin author needs is re-exported here, so no plugin ever has
// to reach into github.com/version14/dot/internal/* (which the Go toolchain
// would refuse outside the dot module).
//
// In-tree plugins (plugins/* in this repo) historically reached straight
// into internal/* because they could; the migration to dotplugin removed
// that special case so the in-tree plugins are now identical in shape to
// any plugin you'd find on GitHub.
package dotplugin

import (
	"github.com/version14/dot/internal/flow"
	"github.com/version14/dot/internal/generator"
	"github.com/version14/dot/internal/plugin"
	"github.com/version14/dot/internal/spec"
	"github.com/version14/dot/internal/state"
)

// ── Plugin identity & registration ────────────────────────────────────────

// PluginID is the namespace prefix for everything a plugin contributes.
// MUST NOT contain '.'. Every contributed ID (generator names, question IDs,
// option values) MUST start with "<PluginID>.".
type PluginID = flow.PluginID

// Provider is the interface a plugin's entry point implements.
type Provider = plugin.Provider

// RegisterBuiltin registers a Provider with the loader. Call this from your
// plugin package's init() function; the dot binary picks it up at startup
// when the package is imported (typically via blank import in cmd/dot/main.go
// when the plugin is vendored, or via dynamic loading once that lands).
func RegisterBuiltin(p Provider) { plugin.RegisterBuiltin(p) }

// ── Question DSL ──────────────────────────────────────────────────────────

type (
	Question        = flow.Question
	QuestionBase    = flow.QuestionBase
	OptionQuestion  = flow.OptionQuestion
	TextQuestion    = flow.TextQuestion
	ConfirmQuestion = flow.ConfirmQuestion
	LoopQuestion    = flow.LoopQuestion
	IfQuestion      = flow.IfQuestion
	Next            = flow.Next
	Option          = flow.Option
	Answer          = flow.Answer
	AnswerNode      = flow.AnswerNode
	LoopFrame       = flow.LoopFrame
	FlowContext     = flow.FlowContext
)

// ── Injection DSL ─────────────────────────────────────────────────────────

type (
	Injection     = flow.Injection
	InjectionKind = flow.InjectionKind
)

const (
	InjectReplace     = flow.InjectReplace
	InjectAddOption   = flow.InjectAddOption
	InjectInsertAfter = flow.InjectInsertAfter
)

// ── Generator wiring ──────────────────────────────────────────────────────

// Entry pairs a Manifest with a Generator implementation. Plugins return a
// slice of these from Provider.Generators().
type Entry = generator.Entry

// Invocation is one (generator-name, loop-stack) tuple — what
// Provider.ResolveExtras() returns.
type Invocation = generator.Invocation

// ── Spec ──────────────────────────────────────────────────────────────────

type (
	ProjectSpec     = spec.ProjectSpec
	ProjectMetadata = spec.ProjectMetadata
)

// ── Virtual filesystem ────────────────────────────────────────────────────

type (
	VirtualProjectState = state.VirtualProjectState
	FileNode            = state.FileNode
	ContentType         = state.ContentType
	JSONDoc             = state.JSONDoc
	YAMLDoc             = state.YAMLDoc
	GoMod               = state.GoMod
)

const (
	ContentRaw   = state.ContentRaw
	ContentJSON  = state.ContentJSON
	ContentYAML  = state.ContentYAML
	ContentGoMod = state.ContentGoMod
)
