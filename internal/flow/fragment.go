package flow

// Fragment registry — named contextual resolvers that return the next node based on loaded plugins

type FlowFragment struct {
	ID      string
	Resolve func(ctx *FlowContext) *Next
}
