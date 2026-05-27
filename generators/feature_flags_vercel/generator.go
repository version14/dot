package featureflagsvercel

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
				"@vercel/edge-config": "^1.0.0",
				"@vercel/flags":       "^3.0.0",
			},
		})
		return nil
	}); err != nil {
		return err
	}

	ctx.State.WriteFile("src/lib/flags.ts", []byte(flagsTS), state.ContentRaw)
	ctx.State.WriteFile("src/hooks/useFeatureFlag.ts", []byte(useFeatureFlagTS), state.ContentRaw)
	ctx.State.WriteFile(".env.example", []byte(envExample), state.ContentRaw)

	return nil
}

const flagsTS = `// Vercel Edge Config feature flags
// Configure flags in your Vercel dashboard or edge-config.json

export type FeatureFlag = "new-dashboard" | "beta-features";

export async function getFlag(flag: FeatureFlag): Promise<boolean> {
  try {
    const { get } = await import("@vercel/edge-config");
    return (await get<boolean>(flag)) ?? false;
  } catch {
    return false;
  }
}
`

const useFeatureFlagTS = `import { useState, useEffect } from "react";
import { getFlag, type FeatureFlag } from "../lib/flags";

export function useFeatureFlag(flag: FeatureFlag): boolean {
  const [enabled, setEnabled] = useState(false);
  useEffect(() => {
    getFlag(flag).then(setEnabled);
  }, [flag]);
  return enabled;
}
`

const envExample = `# Vercel Edge Config
EDGE_CONFIG=your_edge_config_connection_string
`
