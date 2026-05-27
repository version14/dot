package featureflagsposthog

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
				"posthog-js": "^1.0.0",
			},
		})
		return nil
	}); err != nil {
		return err
	}

	ctx.State.WriteFile("src/lib/posthog.ts", []byte(posthogTS), state.ContentRaw)
	ctx.State.WriteFile("src/providers/PostHogProvider.tsx", []byte(posthogProviderTSX), state.ContentRaw)
	ctx.State.WriteFile(".env.example", []byte(envExample), state.ContentRaw)

	return nil
}

const posthogTS = `import posthog from "posthog-js";

export function initPostHog() {
  posthog.init(import.meta.env.VITE_POSTHOG_KEY ?? "", {
    api_host: import.meta.env.VITE_POSTHOG_HOST ?? "https://app.posthog.com",
    loaded: (ph) => {
      if (import.meta.env.DEV) ph.opt_out_capturing();
    },
  });
}

export { posthog };
`

const posthogProviderTSX = `import React, { useEffect } from "react";
import { initPostHog } from "../lib/posthog";

export function PostHogProvider({ children }: { children: React.ReactNode }) {
  useEffect(() => {
    initPostHog();
  }, []);
  return <>{children}</>;
}
`

const envExample = `# PostHog
VITE_POSTHOG_KEY=phc_your_project_api_key
VITE_POSTHOG_HOST=https://app.posthog.com
`
