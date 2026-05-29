import ReactGA from "react-ga4";

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
