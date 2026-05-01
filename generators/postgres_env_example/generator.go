package postgresenvexample

import (
	"fmt"

	"github.com/version14/dot/internal/state"
	"github.com/version14/dot/pkg/dotapi"
)

type Generator struct{}

func New() *Generator { return &Generator{} }

func (g *Generator) Name() string    { return Manifest.Name }
func (g *Generator) Version() string { return Manifest.Version }

func (g *Generator) Generate(ctx *dotapi.Context) error {
	projectName := ctx.Spec.Metadata.ProjectName
	if projectName == "" {
		projectName = "app"
	}

	existing := ""
	if f, ok := ctx.State.GetFile(".env.example"); ok {
		existing = string(f.Content)
	}

	dbURL := fmt.Sprintf("postgresql://postgres:postgres@localhost:5433/%s", projectName)
	updated := existing + fmt.Sprintf("\n# PostgreSQL\nDATABASE_URL=%s\n", dbURL)

	ctx.State.WriteFile(".env.example", []byte(updated), state.ContentRaw)
	return nil
}
