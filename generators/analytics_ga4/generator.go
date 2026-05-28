package analyticsga4

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

	ga4Content := ga4ViteTS
	envContent := envExampleVite
	if framework == "next" {
		ga4Content = ga4NextTS
		envContent = envExampleNext
	}

	if err := ctx.State.UpdateJSON("package.json", func(d *state.JSONDoc) error {
		d.Merge(map[string]interface{}{
			"dependencies": map[string]interface{}{
				"react-ga4": "^2.0.0",
			},
		})
		return nil
	}); err != nil {
		return err
	}

	ctx.State.WriteFile("src/lib/ga4.ts", []byte(ga4Content), state.ContentRaw)
	ctx.State.WriteFile(".env.example", []byte(envContent), state.ContentRaw)

	return nil
}

const ga4ViteTS = `import ReactGA from "react-ga4";

const GA_MEASUREMENT_ID = import.meta.env.VITE_GA_MEASUREMENT_ID ?? "";

export function initGA() {
  if (!GA_MEASUREMENT_ID) return;
  ReactGA.initialize(GA_MEASUREMENT_ID);
}

export function trackPageView(path: string) {
  ReactGA.send({ hitType: "pageview", page: path });
}

export function trackEvent(action: string, category: string, label?: string) {
  ReactGA.event({ action, category, label });
}
`

const ga4NextTS = `import ReactGA from "react-ga4";

const GA_MEASUREMENT_ID = process.env.NEXT_PUBLIC_GA_MEASUREMENT_ID ?? "";

export function initGA() {
  if (!GA_MEASUREMENT_ID) return;
  ReactGA.initialize(GA_MEASUREMENT_ID);
}

export function trackPageView(path: string) {
  ReactGA.send({ hitType: "pageview", page: path });
}

export function trackEvent(action: string, category: string, label?: string) {
  ReactGA.event({ action, category, label });
}
`

const envExampleVite = `# Google Analytics 4
VITE_GA_MEASUREMENT_ID=G-XXXXXXXXXX
`

const envExampleNext = `# Google Analytics 4
NEXT_PUBLIC_GA_MEASUREMENT_ID=G-XXXXXXXXXX
`
