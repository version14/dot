// Vercel Edge Config feature flags
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
