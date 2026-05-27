package reactrouterv7

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
				"react-router":     "^7.0.0",
				"react-router-dom": "^7.0.0",
			},
		})
		return nil
	}); err != nil {
		return err
	}

	ctx.State.WriteFile("src/main.tsx", []byte(mainTSX), state.ContentRaw)
	ctx.State.WriteFile("src/router.tsx", []byte(routerTSX), state.ContentRaw)
	ctx.State.WriteFile("src/pages/Home.tsx", []byte(homeTSX), state.ContentRaw)

	return nil
}

const mainTSX = `import React from "react";
import ReactDOM from "react-dom/client";
import { RouterProvider } from "react-router-dom";
import { router } from "./router";

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

const routerTSX = `import { createBrowserRouter } from "react-router-dom";
import Home from "./pages/Home";

export const router = createBrowserRouter([
  { path: "/", element: <Home /> },
]);
`

const homeTSX = `export default function Home() {
  return <h1>Home</h1>;
}
`
