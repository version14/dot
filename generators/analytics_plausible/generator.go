package analyticsplausible

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
	plausibleContent := plausibleViteTS
	if framework == "next" {
		plausibleContent = plausibleNextTS
	}

	if err := ctx.State.UpdateJSON("package.json", func(d *state.JSONDoc) error {
		d.Merge(map[string]interface{}{
			"dependencies": map[string]interface{}{
				"plausible-tracker": "^0.3.0",
			},
		})
		return nil
	}); err != nil {
		return err
	}

	ctx.State.WriteFile("src/lib/plausible.ts", []byte(plausibleContent), state.ContentRaw)
	ctx.State.WriteFile(".env.example", []byte(envExample), state.ContentRaw)

	return nil
}

const plausibleViteTS = `import Plausible from "plausible-tracker";

const DOMAIN =
  import.meta.env.VITE_PLAUSIBLE_DOMAIN ?? window.location.hostname;
const API_HOST =
  import.meta.env.VITE_PLAUSIBLE_API_HOST ?? "https://plausible.io";

const plausible = Plausible({ domain: DOMAIN, apiHost: API_HOST });

export function initPlausible() {
  plausible.enableAutoPageviews();
}

export function trackEvent(
  name: string,
  props?: Record<string, string | number | boolean>,
) {
  plausible.trackEvent(name, { props });
}
`

const plausibleNextTS = `import Plausible from "plausible-tracker";

const DOMAIN =
  process.env.NEXT_PUBLIC_PLAUSIBLE_DOMAIN ??
  (typeof window !== "undefined" ? window.location.hostname : "");
const API_HOST =
  process.env.NEXT_PUBLIC_PLAUSIBLE_API_HOST ?? "https://plausible.io";

const plausible = Plausible({ domain: DOMAIN, apiHost: API_HOST });

export function initPlausible() {
  plausible.enableAutoPageviews();
}

export function trackEvent(
  name: string,
  props?: Record<string, string | number | boolean>,
) {
  plausible.trackEvent(name, { props });
}
`

const envExample = `# Plausible Analytics
VITE_PLAUSIBLE_DOMAIN=yourdomain.com
VITE_PLAUSIBLE_API_HOST=https://plausible.io
`
