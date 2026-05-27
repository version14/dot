package authbetterauthfrontend

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
	authClientContent := authClientViteTS
	if framework == "next" {
		authClientContent = authClientNextTS
	}

	if err := ctx.State.UpdateJSON("package.json", func(d *state.JSONDoc) error {
		d.Merge(map[string]interface{}{
			"dependencies": map[string]interface{}{
				"better-auth": "^1.0.0",
			},
		})
		return nil
	}); err != nil {
		return err
	}

	ctx.State.WriteFile("src/lib/authClient.ts", []byte(authClientContent), state.ContentRaw)
	ctx.State.WriteFile("src/providers/AuthProvider.tsx", []byte(authProviderTSX), state.ContentRaw)
	ctx.State.WriteFile("src/hooks/useAuth.ts", []byte(useAuthTS), state.ContentRaw)

	return nil
}

const authClientViteTS = `import { createAuthClient } from "better-auth/react";

export const authClient = createAuthClient({
  baseURL: import.meta.env.VITE_API_URL ?? "http://localhost:3000",
});

export const { signIn, signOut, signUp, useSession } = authClient;
`

const authClientNextTS = `import { createAuthClient } from "better-auth/react";

export const authClient = createAuthClient({
  baseURL: process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:3000",
});

export const { signIn, signOut, signUp, useSession } = authClient;
`

const authProviderTSX = `import React, { createContext, useContext } from "react";
import { useSession } from "../lib/authClient";

interface AuthContextValue {
  session: ReturnType<typeof useSession>["data"];
}

const AuthContext = createContext<AuthContextValue | null>(null);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const { data: session } = useSession();
  return (
    <AuthContext.Provider value={{ session }}>{children}</AuthContext.Provider>
  );
}

export function useAuthContext() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("useAuthContext must be used within AuthProvider");
  return ctx;
}
`

const useAuthTS = `export { useSession } from "../lib/authClient";
export { useAuthContext } from "../providers/AuthProvider";
`
