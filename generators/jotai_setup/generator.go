package jotaisetup

import (
	"github.com/version14/dot/internal/state"
	"github.com/version14/dot/pkg/dotapi"
)

type Generator struct{}

func New() *Generator { return &Generator{} }

func (g *Generator) Name() string    { return Manifest.Name }
func (g *Generator) Version() string { return Manifest.Version }

func (g *Generator) Generate(ctx *dotapi.Context) error {
	if err := ctx.State.UpdateJSON("package.json", func(d *state.JSONDoc) error {
		d.Merge(map[string]interface{}{
			"dependencies": map[string]interface{}{
				"jotai": "^2.0.0",
			},
		})
		return nil
	}); err != nil {
		return err
	}

	ctx.State.WriteFile("src/atoms/counter.atom.ts", []byte(counterAtom), state.ContentRaw)
	return nil
}

const counterAtom = `import { atom } from "jotai";

export const countAtom = atom(0);

export const incrementAtom = atom(null, (get, set) =>
  set(countAtom, get(countAtom) + 1),
);

export const decrementAtom = atom(null, (get, set) =>
  set(countAtom, get(countAtom) - 1),
);
`
