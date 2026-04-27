package flow

// FlowRunner executes a flow from a root question and returns the collected
// answers. Both FlowEngine (sequential per-question adapter model) and
// HuhFormRunner (single Huh form with reactive hide functions) implement this
// interface, so callers can swap rendering strategies without touching the
// engine or generator layers.
type FlowRunner interface {
	Run(root Question) (*FlowContext, error)
}
