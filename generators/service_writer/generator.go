package servicewriter

import (
	"fmt"

	"github.com/version14/dot/internal/render"
	"github.com/version14/dot/internal/state"
	"github.com/version14/dot/pkg/dotapi"
)

// Generator writes one service skeleton per invocation. The loop frame
// passes "name" (and optionally "port") via ctx.Answers because the resolver
// flattened them with FlattenScope before calling Generate.
type Generator struct{}

func New() *Generator { return &Generator{} }

func (g *Generator) Name() string    { return Manifest.Name }
func (g *Generator) Version() string { return Manifest.Version }

func (g *Generator) Generate(ctx *dotapi.Context) error {
	name, _ := ctx.Answers["name"].(string)
	if name == "" {
		return fmt.Errorf("service_writer: missing 'name' in scoped answers")
	}

	// Port defaults to 3000; loops typically vary it per iteration.
	port := 3000
	if p, ok := ctx.Answers["port"].(int); ok {
		port = p
	} else if pf, ok := ctx.Answers["port"].(float64); ok {
		port = int(pf)
	}

	data := map[string]interface{}{
		"Name": name,
		"Port": port,
	}

	mainTS, err := render.Render(mainTSTmpl, data)
	if err != nil {
		return fmt.Errorf("service_writer: render main.ts: %w", err)
	}
	pkgJSON, err := render.Render(packageJSONTmpl, data)
	if err != nil {
		return fmt.Errorf("service_writer: render package.json: %w", err)
	}

	ctx.State.WriteFile("services/"+name+"/src/main.ts", mainTS, state.ContentRaw)
	ctx.State.WriteFile("services/"+name+"/package.json", pkgJSON, state.ContentJSON)
	return nil
}

const mainTSTmpl = `// Service: {{.Name}}
import http from "node:http";

const PORT = {{.Port}};
const server = http.createServer((req, res) => {
  res.writeHead(200, { "Content-Type": "application/json" });
  res.end(JSON.stringify({ service: "{{.Name}}", ok: true }));
});
server.listen(PORT, () => {
  console.log("[{{.Name}}] listening on", PORT);
});
`

const packageJSONTmpl = `{
  "name": "{{.Name}}",
  "private": true,
  "type": "module",
  "scripts": {
    "dev": "node --watch src/main.ts",
    "start": "node src/main.ts"
  }
}
`
