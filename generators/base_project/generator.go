package baseproject

import (
	"context"

	"github.com/version14/dot/internal/render"
	"github.com/version14/dot/pkg/dotapi"
)

// Generator implements dotapi.Generator for base_project. It writes universal
// project files from a GitHub template repository. It must run first because
// later generators may modify these (e.g. typescript_base appends node_modules
// to .gitignore).
type Generator struct{}

// New constructs a Generator instance for registration.
func New() *Generator { return &Generator{} }

func (g *Generator) Name() string    { return Manifest.Name }
func (g *Generator) Version() string { return Manifest.Version }

func (g *Generator) Generate(ctx *dotapi.Context) error {
	snapshot, err := getBaseFilesSnapshot()
	if err != nil {
		return err
	}

	if err := render.PopulateStateFromSnapshot(ctx.State, snapshot); err != nil {
		return err
	}

	return nil
}

func getBaseFilesSnapshot() (*render.RepoSnapshot, error) {
	fetcher := render.NewGitHubArchiveFetcher()
	snapshot, err := fetcher.FetchRepo(
		context.TODO(),
		"https://github.com/mathieusouflis/github-template.git",
		render.FetchOptions{},
	)
	if err != nil {
		return nil, err
	}

	return snapshot, nil
}
