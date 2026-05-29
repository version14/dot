import posthog from "posthog-js";

export function initPostHog() {
  posthog.init(process.env.NEXT_PUBLIC_POSTHOG_KEY ?? "", {
    api_host: process.env.NEXT_PUBLIC_POSTHOG_HOST ?? "https://app.posthog.com",
    loaded: (ph) => {
      if (process.env.NODE_ENV === "development") ph.opt_out_capturing();
    },
  });
}

export { posthog };
