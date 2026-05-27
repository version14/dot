package themeprovider

import (
	"github.com/version14/dot/internal/state"
	"github.com/version14/dot/pkg/dotapi"
)

type Generator struct{}

func New() *Generator { return &Generator{} }

func (g *Generator) Name() string    { return Manifest.Name }
func (g *Generator) Version() string { return Manifest.Version }

func (g *Generator) Generate(ctx *dotapi.Context) error {
	ctx.State.WriteFile("src/styles/theme.css", []byte(themeCSS), state.ContentRaw)
	ctx.State.WriteFile("src/providers/ThemeProvider.tsx", []byte(themeProviderTSX), state.ContentRaw)
	ctx.State.WriteFile("src/hooks/useTheme.ts", []byte(useThemeTS), state.ContentRaw)
	return nil
}

const themeCSS = `:root {
  --color-background: #ffffff;
  --color-foreground: #0a0a0a;
  --color-primary: #3b82f6;
  --color-primary-foreground: #ffffff;
  --color-muted: #f4f4f5;
  --color-muted-foreground: #71717a;
  --color-border: #e4e4e7;
}

[data-theme="dark"] {
  --color-background: #0a0a0a;
  --color-foreground: #f4f4f5;
  --color-primary: #60a5fa;
  --color-primary-foreground: #0a0a0a;
  --color-muted: #27272a;
  --color-muted-foreground: #a1a1aa;
  --color-border: #27272a;
}
`

const themeProviderTSX = `import React, { createContext, useContext, useEffect, useState } from "react";

type Theme = "light" | "dark" | "system";

interface ThemeContextValue {
  theme: Theme;
  setTheme: (theme: Theme) => void;
}

const ThemeContext = createContext<ThemeContextValue | null>(null);

export function ThemeProvider({
  children,
  defaultTheme = "system",
}: {
  children: React.ReactNode;
  defaultTheme?: Theme;
}) {
  const [theme, setThemeState] = useState<Theme>(
    () => (localStorage.getItem("theme") as Theme) ?? defaultTheme,
  );

  useEffect(() => {
    const root = document.documentElement;
    const resolved =
      theme === "system"
        ? window.matchMedia("(prefers-color-scheme: dark)").matches
          ? "dark"
          : "light"
        : theme;
    root.setAttribute("data-theme", resolved);
    localStorage.setItem("theme", theme);
  }, [theme]);

  return (
    <ThemeContext.Provider value={{ theme, setTheme: setThemeState }}>
      {children}
    </ThemeContext.Provider>
  );
}

export function useThemeContext() {
  const ctx = useContext(ThemeContext);
  if (!ctx) throw new Error("useThemeContext must be used within ThemeProvider");
  return ctx;
}
`

const useThemeTS = `export { useThemeContext as useTheme } from "../providers/ThemeProvider";
`
