package sentryfrontend

import (
	"github.com/version14/dot/internal/state"
	"github.com/version14/dot/pkg/dotapi"
)

type Generator struct{}

func New() *Generator { return &Generator{} }

func (g *Generator) Name() string    { return Manifest.Name }
func (g *Generator) Version() string { return Manifest.Version }

func (g *Generator) Generate(ctx *dotapi.Context) error {
	framework, _ := ctx.Answers["framework"].(string)

	var sentryPkg, sentryContent string
	if framework == "next" {
		sentryPkg = "@sentry/nextjs"
		sentryContent = sentryNext
		ctx.State.WriteFile("sentry.client.config.ts", []byte(sentryClientConfig), state.ContentRaw)
	} else {
		sentryPkg = "@sentry/react"
		sentryContent = sentryReact
	}

	if err := ctx.State.UpdateJSON("package.json", func(d *state.JSONDoc) error {
		d.Merge(map[string]interface{}{
			"dependencies": map[string]interface{}{
				sentryPkg: "^8.0.0",
			},
		})
		return nil
	}); err != nil {
		return err
	}

	envContent := envExampleVite
	if framework == "next" {
		envContent = envExampleNext
	}

	ctx.State.WriteFile("src/lib/sentry.ts", []byte(sentryContent), state.ContentRaw)
	ctx.State.WriteFile(".env.example", []byte(envContent), state.ContentRaw)

	return nil
}

const sentryNext = `// For Next.js, configure Sentry in sentry.client.config.ts and sentry.server.config.ts
// See: https://docs.sentry.io/platforms/javascript/guides/nextjs/
export * from "@sentry/nextjs";
`

const sentryClientConfig = `import * as Sentry from "@sentry/nextjs";

Sentry.init({
  dsn: process.env.NEXT_PUBLIC_SENTRY_DSN,
  tracesSampleRate: 0.1,
});
`

const sentryReact = `import * as Sentry from "@sentry/react";

export function initSentry() {
  Sentry.init({
    dsn: import.meta.env.VITE_SENTRY_DSN,
    tracesSampleRate: 0.1,
    environment: import.meta.env.MODE,
  });
}

export { Sentry };
`

const envExampleVite = `# Sentry
VITE_SENTRY_DSN=https://your_dsn@sentry.io/your_project_id
`

const envExampleNext = `# Sentry
NEXT_PUBLIC_SENTRY_DSN=https://your_dsn@sentry.io/your_project_id
`
