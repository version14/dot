package prettierexpressrules

import (
	"encoding/json"

	"github.com/version14/dot/internal/state"
	"github.com/version14/dot/pkg/dotapi"
)

type Generator struct{}

func New() *Generator { return &Generator{} }

func (g *Generator) Name() string    { return Manifest.Name }
func (g *Generator) Version() string { return Manifest.Version }

func (g *Generator) Generate(ctx *dotapi.Context) error {
	// Read existing .prettierrc and merge backend-specific rules
	var existing map[string]interface{}
	if f, ok := ctx.State.GetFile(".prettierrc"); ok {
		if err := json.Unmarshal(f.Content, &existing); err != nil {
			existing = map[string]interface{}{}
		}
	} else {
		existing = map[string]interface{}{}
	}

	existing["printWidth"] = 100
	existing["endOfLine"] = "lf"
	existing["bracketSpacing"] = true
	existing["arrowParens"] = "always"

	out, err := json.MarshalIndent(existing, "", "  ")
	if err != nil {
		return err
	}
	out = append(out, '\n')

	ctx.State.WriteFile(".prettierrc", out, state.ContentRaw)
	return nil
}
