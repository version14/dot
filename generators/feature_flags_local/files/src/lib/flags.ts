let flagsCache: Record<string, boolean> | null = null;

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
