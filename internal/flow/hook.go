package flow

import (
	"fmt"
	"strings"
)

// PluginID identifies a plugin. It is the namespace prefix for every ID the
// plugin contributes to the flow (question IDs, option values). PluginIDs
// must be unique across all loaded plugins.
type PluginID string

// InjectionKind enumerates the three ways a plugin may modify a flow:
//
//  1. Replace      — fully swap the targeted question for one the plugin defines.
//     The replacement provides its own Next chain; the original
//     flow beyond the target is discarded for this branch.
//
//  2. AddOption    — append an Option to an existing OptionQuestion. The
//     option's Next decides where selecting it leads.
//
//  3. InsertAfter  — splice a new question (or chain of questions) after the
//     targeted question. The targeted question's normal Next
//     resumes once the inserted chain reaches an End edge.
//
// Plugins target a question by its ID; no named "hook constants" are required.
type InjectionKind int

const (
	InjectReplace InjectionKind = iota
	InjectAddOption
	InjectInsertAfter
)

// Injection is one plugin contribution targeting one question ID.
//
// Exactly one of Replacement, Option, or Question must be set, matching Kind.
// All IDs the plugin contributes (Replacement.ID(), Option.Value, Question.ID())
// must be prefixed with "<Plugin>." — enforced at registration time.
type Injection struct {
	Plugin   PluginID
	TargetID string
	Kind     InjectionKind

	Replacement Question // when Kind == InjectReplace
	Option      *Option  // when Kind == InjectAddOption
	Question    Question // when Kind == InjectInsertAfter
}

// HookRegistry holds plugin injections keyed by the question ID they target.
// The engine consults it during traversal to discover plugin contributions.
type HookRegistry struct {
	byTarget map[string][]*Injection
	plugins  map[PluginID]bool
}

func NewHookRegistry() *HookRegistry {
	return &HookRegistry{
		byTarget: map[string][]*Injection{},
		plugins:  map[PluginID]bool{},
	}
}

// RegisterPlugin records a plugin ID. Must be called before any Inject for
// that plugin; the registry uses it to enforce the prefix rule.
func (r *HookRegistry) RegisterPlugin(id PluginID) error {
	if id == "" {
		return fmt.Errorf("flow: plugin ID required")
	}
	if strings.Contains(string(id), ".") {
		return fmt.Errorf("flow: plugin ID %q must not contain '.'", id)
	}
	if r.plugins[id] {
		return fmt.Errorf("flow: plugin %q already registered", id)
	}
	r.plugins[id] = true
	return nil
}

// Inject registers a plugin's injection. Returns an error if the plugin is
// not registered, the kind is unknown, the required field for that kind is
// nil, or any contributed ID lacks the plugin's prefix.
func (r *HookRegistry) Inject(inj *Injection) error {
	if inj == nil {
		return fmt.Errorf("flow: nil injection")
	}
	if !r.plugins[inj.Plugin] {
		return fmt.Errorf("flow: plugin %q not registered", inj.Plugin)
	}
	if inj.TargetID == "" {
		return fmt.Errorf("flow: injection requires TargetID")
	}
	if err := validateInjection(inj); err != nil {
		return err
	}
	r.byTarget[inj.TargetID] = append(r.byTarget[inj.TargetID], inj)
	return nil
}

func validateInjection(inj *Injection) error {
	prefix := string(inj.Plugin) + "."
	switch inj.Kind {
	case InjectReplace:
		if inj.Replacement == nil {
			return fmt.Errorf("flow: InjectReplace requires Replacement")
		}
		if !strings.HasPrefix(inj.Replacement.ID(), prefix) {
			return fmt.Errorf("flow: replacement ID %q must start with %q", inj.Replacement.ID(), prefix)
		}
	case InjectAddOption:
		if inj.Option == nil {
			return fmt.Errorf("flow: InjectAddOption requires Option")
		}
		if !strings.HasPrefix(inj.Option.Value, prefix) {
			return fmt.Errorf("flow: option value %q must start with %q", inj.Option.Value, prefix)
		}
	case InjectInsertAfter:
		if inj.Question == nil {
			return fmt.Errorf("flow: InjectInsertAfter requires Question")
		}
		if !strings.HasPrefix(inj.Question.ID(), prefix) {
			return fmt.Errorf("flow: inserted question ID %q must start with %q", inj.Question.ID(), prefix)
		}
	default:
		return fmt.Errorf("flow: unknown injection kind %d", inj.Kind)
	}
	return nil
}

// For returns every injection targeting targetID, in registration order.
func (r *HookRegistry) For(targetID string) []*Injection {
	return r.byTarget[targetID]
}

// Plugins returns the IDs of every registered plugin.
func (r *HookRegistry) Plugins() []PluginID {
	out := make([]PluginID, 0, len(r.plugins))
	for k := range r.plugins {
		out = append(out, k)
	}
	return out
}

// ForKind partitions a target's injections by kind.
// It is the public form of byKind, available to packages outside flow (e.g. cli).
func (r *HookRegistry) ForKind(targetID string) (replace []Question, options []*Option, after []Question) {
	return r.byKind(targetID)
}

// byKind partitions a target's injections by kind for the engine.
func (r *HookRegistry) byKind(targetID string) (replace []Question, options []*Option, after []Question) {
	for _, inj := range r.byTarget[targetID] {
		switch inj.Kind {
		case InjectReplace:
			replace = append(replace, inj.Replacement)
		case InjectAddOption:
			options = append(options, inj.Option)
		case InjectInsertAfter:
			after = append(after, inj.Question)
		}
	}
	return
}
