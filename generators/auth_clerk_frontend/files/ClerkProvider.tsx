import { ClerkProvider as BaseClerkProvider } from "@clerk/clerk-react";

const PUBLISHABLE_KEY = import.meta.env.VITE_CLERK_PUBLISHABLE_KEY;
if (!PUBLISHABLE_KEY) throw new Error("Missing VITE_CLERK_PUBLISHABLE_KEY");

export function ClerkProvider({ children }: { children: React.ReactNode }) {
  return (
    <BaseClerkProvider publishableKey={PUBLISHABLE_KEY}>
      {children}
    </BaseClerkProvider>
  );
}
