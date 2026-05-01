package expresserrormiddleware

import (
	"embed"
	"strings"

	"github.com/version14/dot/internal/render"
	"github.com/version14/dot/internal/state"
	"github.com/version14/dot/pkg/dotapi"
)

type Generator struct{}

func New() *Generator { return &Generator{} }

func (g *Generator) Name() string    { return Manifest.Name }
func (g *Generator) Version() string { return Manifest.Version }

//go:embed all:files
var fs embed.FS

const errorImport = "import { errorMiddleware } from './shared/middlewares/error.middleware';\n"
const errorUse = "app.use(errorMiddleware);\n"

func (g *Generator) Generate(ctx *dotapi.Context) error {
	renderer := render.NewLocalFolderRenderer(ctx.State)
	if err := renderer.Render(fs, nil); err != nil {
		return err
	}

	// Inject error middleware into app.ts
	f, ok := ctx.State.GetFile("src/app.ts")
	if !ok {
		return nil
	}
	content := string(f.Content)
	if strings.Contains(content, "errorMiddleware") {
		return nil
	}

	// Add import at the top (after last existing import line)
	content = errorImport + content

	// Add app.use(errorMiddleware) just before "export default app;"
	content = strings.Replace(content, "export default app;", errorUse+"\nexport default app;", 1)

	ctx.State.WriteFile("src/app.ts", []byte(content), state.ContentRaw)
	return nil
}
