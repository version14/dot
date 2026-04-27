package dotapi

// Generator is the contract every generator implements. It is invoked once per
// matching flow path; loop bodies cause repeated invocations with different
// scoped Answers in Context.
type Generator interface {
	Name() string
	Version() string
	Generate(ctx *Context) error
}
