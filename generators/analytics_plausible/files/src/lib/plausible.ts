import Plausible from "plausible-tracker";

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
