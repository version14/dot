package authvanillafrontend

import (
	"github.com/version14/dot/internal/state"
	"github.com/version14/dot/pkg/dotapi"
)

type Generator struct{}

func New() *Generator { return &Generator{} }

func (g *Generator) Name() string    { return Manifest.Name }
func (g *Generator) Version() string { return Manifest.Version }

func (g *Generator) Generate(ctx *dotapi.Context) error {
	ctx.State.WriteFile("src/lib/auth.ts", []byte(authTS), state.ContentRaw)
	ctx.State.WriteFile("src/providers/AuthProvider.tsx", []byte(authProviderTSX), state.ContentRaw)
	ctx.State.WriteFile("src/hooks/useAuth.ts", []byte(useAuthTS), state.ContentRaw)
	return nil
}

const authTS = `const API_URL = import.meta.env.VITE_API_URL ?? "http://localhost:3000";

let accessToken: string | null = null;

export async function login(email: string, password: string): Promise<void> {
  const res = await fetch(` + "`" + `${API_URL}/auth/login` + "`" + `, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password }),
  });
  if (!res.ok) throw new Error("Login failed");
  const data = await res.json();
  accessToken = data.accessToken;
}

export async function logout(): Promise<void> {
  await fetch(` + "`" + `${API_URL}/auth/logout` + "`" + `, { method: "POST" });
  accessToken = null;
}

export function getAccessToken(): string | null {
  return accessToken;
}

export function isAuthenticated(): boolean {
  return accessToken !== null;
}
`

const authProviderTSX = `import React, { createContext, useContext, useState, useCallback } from "react";
import { login as apiLogin, logout as apiLogout, isAuthenticated } from "../lib/auth";

interface AuthContextValue {
  isAuth: boolean;
  login: (email: string, password: string) => Promise<void>;
  logout: () => Promise<void>;
}

const AuthContext = createContext<AuthContextValue | null>(null);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [isAuth, setIsAuth] = useState(isAuthenticated());

  const login = useCallback(async (email: string, password: string) => {
    await apiLogin(email, password);
    setIsAuth(true);
  }, []);

  const logout = useCallback(async () => {
    await apiLogout();
    setIsAuth(false);
  }, []);

  return (
    <AuthContext.Provider value={{ isAuth, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuthContext() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("useAuthContext must be used within AuthProvider");
  return ctx;
}
`

const useAuthTS = `export { useAuthContext as useAuth } from "../providers/AuthProvider";
`
