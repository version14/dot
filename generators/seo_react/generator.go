package seoreact

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
				"react-helmet-async": "^2.0.0",
			},
		})
		return nil
	}); err != nil {
		return err
	}

	ctx.State.WriteFile("src/providers/HelmetProvider.tsx", []byte(helmetProviderTSX), state.ContentRaw)
	ctx.State.WriteFile("src/components/SEO.tsx", []byte(seoTSX), state.ContentRaw)

	return nil
}

const helmetProviderTSX = `import { HelmetProvider as BaseHelmetProvider } from "react-helmet-async";

export function HelmetProvider({ children }: { children: React.ReactNode }) {
  return <BaseHelmetProvider>{children}</BaseHelmetProvider>;
}
`

const seoTSX = `import { Helmet } from "react-helmet-async";

interface SEOProps {
  title?: string;
  description?: string;
  canonical?: string;
}

export function SEO({ title, description, canonical }: SEOProps) {
  return (
    <Helmet>
      {title && <title>{title}</title>}
      {description && <meta name="description" content={description} />}
      {canonical && <link rel="canonical" href={canonical} />}
    </Helmet>
  );
}
`
