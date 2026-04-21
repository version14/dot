package generator

import (
	"github.com/version14/dot/internal/generator"
	"github.com/version14/dot/internal/new-test/registry/questions"
	"github.com/version14/dot/internal/spec"
)

// Generator is the concrete scaffold unit.
// ApplyFunction produces the files; PostApplyFunction produces the shell
// commands to run after (e.g. pnpm install). extends hooks let community
// plugins inject extra FileOps into an existing generator without forking it.
type Generator struct {
	Name              string
	Language          string
	ApplyFunction     func(s spec.Spec) ([]generator.FileOp, error)
	PostApplyFunction func(s spec.Spec) []generator.PostOp // optional
	extends           []func(s spec.Spec, ops []generator.FileOp)
}

// Extend registers a hook that receives the FileOps produced by ApplyFunction
// and may append to them. Used by community plugins to extend official generators.
func (g *Generator) Extend(f func(s spec.Spec, ops []generator.FileOp)) {
	g.extends = append(g.extends, f)
}

// Apply runs ApplyFunction then all Extend hooks.
func (g *Generator) Apply(s spec.Spec) ([]generator.FileOp, error) {
	ops, err := g.ApplyFunction(s)
	if err != nil {
		return nil, err
	}
	for _, f := range g.extends {
		f(s, ops)
	}
	return ops, nil
}

// PostApply returns the post-generation shell commands, or nil if none.
func (g *Generator) PostApply(s spec.Spec) []generator.PostOp {
	if g.PostApplyFunction == nil {
		return nil
	}
	return g.PostApplyFunction(s)
}

// Func returns a questions.GeneratorFunc that wraps this Generator.
// Use this to attach the generator to a question's Next field:
//
//	questions.Select(...).ChoiceWithGen("React", "react", ReactTS.Func())
func (g *Generator) Func() questions.GeneratorFunc {
	return func(s spec.Spec) ([]generator.FileOp, []generator.PostOp, error) {
		ops, err := g.Apply(s)
		if err != nil {
			return nil, nil, err
		}
		return ops, g.PostApply(s), nil
	}
}
