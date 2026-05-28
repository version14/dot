package authclerkfrontend

import (
	"github.com/version14/dot/internal/state"
	"github.com/version14/dot/pkg/dotapi"
)

type Generator struct{}

func New() *Generator { return &Generator{} }

func (g *Generator) Name() string    { return Manifest.Name }
func (g *Generator) Version() string { return Manifest.Version }

func (g *Generator) Generate(ctx *dotapi.Context) error {
	framework, _ := ctx.Answers["framework"].(string)

	var pkg, providerContent, hookContent string
	if framework == "next" {
		pkg = "@clerk/nextjs"
		providerContent = clerkProviderNext
		hookContent = clerkHookNext
	} else {
		pkg = "@clerk/clerk-react"
		providerContent = clerkProviderReact
		hookContent = clerkHookReact
	}

	if err := ctx.State.UpdateJSON("package.json", func(d *state.JSONDoc) error {
		d.Merge(map[string]interface{}{
			"dependencies": map[string]interface{}{
				pkg: "^5.0.0",
			},
		})
		return nil
	}); err != nil {
		return err
	}

	envContent := envExampleVite
	if framework == "next" {
		envContent = envExampleNext
	}

	ctx.State.WriteFile("src/providers/ClerkProvider.tsx", []byte(providerContent), state.ContentRaw)
	ctx.State.WriteFile("src/hooks/useAuth.ts", []byte(hookContent), state.ContentRaw)
	ctx.State.WriteFile(".env.example", []byte(envContent), state.ContentRaw)

	return nil
}

const clerkProviderNext = `// For Next.js: wrap your root layout with ClerkProvider from @clerk/nextjs
// See: https://clerk.com/docs/quickstarts/nextjs
export { ClerkProvider } from "@clerk/nextjs";
`

const clerkHookNext = `export { useAuth, useUser, useClerk } from "@clerk/nextjs";
`

const clerkProviderReact = `import { ClerkProvider as BaseClerkProvider } from "@clerk/clerk-react";

const PUBLISHABLE_KEY = import.meta.env.VITE_CLERK_PUBLISHABLE_KEY;
if (!PUBLISHABLE_KEY) throw new Error("Missing VITE_CLERK_PUBLISHABLE_KEY");

export function ClerkProvider({ children }: { children: React.ReactNode }) {
  return (
    <BaseClerkProvider publishableKey={PUBLISHABLE_KEY}>
      {children}
    </BaseClerkProvider>
  );
}
`

const clerkHookReact = `export {
  useAuth,
  useUser,
  useClerk,
  SignedIn,
  SignedOut,
  SignInButton,
} from "@clerk/clerk-react";
`

const envExampleVite = `# Clerk
VITE_CLERK_PUBLISHABLE_KEY=pk_test_your_publishable_key_here
`

const envExampleNext = `# Clerk
NEXT_PUBLIC_CLERK_PUBLISHABLE_KEY=pk_test_your_publishable_key_here
`
