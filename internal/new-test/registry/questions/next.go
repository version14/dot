package questions

import (
	"github.com/version14/dot/internal/generator"
	"github.com/version14/dot/internal/spec"
)

// GeneratorFunc is the type for generator functions stored in Next.
// Returns file operations and post-generation shell commands (e.g. pnpm install).
type GeneratorFunc func(s spec.Spec) ([]generator.FileOp, []generator.PostOp, error)

// Next is the continuation after a question or option.
// Generator is optional — set it to trigger a generator when this path is taken.
type Next struct {
	Generator GeneratorFunc
	Question  *Question
}
