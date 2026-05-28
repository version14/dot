import posthog from "posthog-js";

export function initPostHog() {
  posthog.init(import.meta.env.VITE_POSTHOG_KEY ?? "", {
    api_host: import.meta.env.VITE_POSTHOG_HOST ?? "https://app.posthog.com",
    loaded: (ph) => {
      if (import.meta.env.DEV) ph.opt_out_capturing();
    },
  });
}

export { posthog };
