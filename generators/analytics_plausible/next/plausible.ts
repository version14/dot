import Plausible from "plausible-tracker";

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
