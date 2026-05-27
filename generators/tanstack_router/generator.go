package tanstackrouter

import (
	"github.com/version14/dot/internal/state"
	"github.com/version14/dot/pkg/dotapi"
)

type Generator struct{}

func New() *Generator { return &Generator{} }

func (g *Generator) Name() string    { return Manifest.Name }
func (g *Generator) Version() string { return Manifest.Version }

func (g *Generator) Generate(ctx *dotapi.Context) error {
	if err := ctx.State.UpdateJSON("package.json", func(d *state.JSONDoc) error {
		d.Merge(map[string]interface{}{
			"dependencies": map[string]interface{}{
				"@tanstack/react-router": "^1.0.0",
			},
			"devDependencies": map[string]interface{}{
				"@tanstack/router-plugin": "^1.0.0",
			},
		})
		return nil
	}); err != nil {
		return err
	}

	ctx.State.WriteFile("vite.config.ts", []byte(viteConfig), state.ContentRaw)
	ctx.State.WriteFile("src/main.tsx", []byte(mainTSX), state.ContentRaw)
	ctx.State.WriteFile("src/routes/__root.tsx", []byte(rootTSX), state.ContentRaw)
	ctx.State.WriteFile("src/routes/index.tsx", []byte(indexTSX), state.ContentRaw)

	return nil
}

const viteConfig = `import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import { TanStackRouterVite } from "@tanstack/router-plugin/vite";

export default defineConfig({
  plugins: [TanStackRouterVite({ autoCodeSplitting: true }), react()],
});
`

const mainTSX = `import React from "react";
import ReactDOM from "react-dom/client";
import { RouterProvider, createRouter } from "@tanstack/react-router";
import { routeTree } from "./routeTree.gen";

const router = createRouter({ routeTree });

declare module "@tanstack/react-router" {
  interface Register {
    router: typeof router;
  }
}

const rootElement = document.getElementById("root");

if (!rootElement) {
  throw new Error("Root element #root not found");
}

ReactDOM.createRoot(rootElement).render(
  <React.StrictMode>
    <RouterProvider router={router} />
  </React.StrictMode>,
);
`

const rootTSX = `import { createRootRoute, Outlet } from "@tanstack/react-router";

export const Route = createRootRoute({
  component: () => <Outlet />,
});
`

const indexTSX = `import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/")({
  component: () => <h1>Home</h1>,
});
`
