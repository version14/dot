import { HelmetProvider as BaseHelmetProvider } from "react-helmet-async";

export function HelmetProvider({ children }: { children: React.ReactNode }) {
  return <BaseHelmetProvider>{children}</BaseHelmetProvider>;
}
