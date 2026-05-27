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

	ctx.State.WriteFile("src/lib/ga4.ts", []byte(ga4TS), state.ContentRaw)
	ctx.State.WriteFile(".env.example", []byte(envExample), state.ContentRaw)

	return nil
}

const ga4TS = `import ReactGA from "react-ga4";

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

const envExample = `# Google Analytics 4
VITE_GA_MEASUREMENT_ID=G-XXXXXXXXXX
`
