package featureflagslocal

import (
	"github.com/version14/dot/internal/state"
	"github.com/version14/dot/pkg/dotapi"
)

type Generator struct{}

func New() *Generator { return &Generator{} }

func (g *Generator) Name() string    { return Manifest.Name }
func (g *Generator) Version() string { return Manifest.Version }

func (g *Generator) Generate(ctx *dotapi.Context) error {
	ctx.State.WriteFile("public/flags.json", []byte(flagsJSON), state.ContentRaw)
	ctx.State.WriteFile("src/lib/flags.ts", []byte(flagsTS), state.ContentRaw)
	ctx.State.WriteFile("src/hooks/useFeatureFlag.ts", []byte(useFeatureFlagTS), state.ContentRaw)
	return nil
}

const flagsJSON = `{
  "new-dashboard": false,
  "beta-features": false
}
`

const flagsTS = `let flagsCache: Record<string, boolean> | null = null;

async function loadFlags(): Promise<Record<string, boolean>> {
  if (flagsCache) return flagsCache;
  const res = await fetch("/flags.json");
  flagsCache = await res.json();
  return flagsCache ?? {};

}

export async function getFlag(flag: string): Promise<boolean> {
  const flags = await loadFlags();
  return flags[flag] ?? false;
}
`

const useFeatureFlagTS = `import { useState, useEffect } from "react";
import { getFlag } from "../lib/flags";

export function useFeatureFlag(flag: string): boolean {
  const [enabled, setEnabled] = useState(false);
  useEffect(() => {
    getFlag(flag).then(setEnabled);
  }, [flag]);
  return enabled;
}
`
