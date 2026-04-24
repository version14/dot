// Package typescript_express_shared provides shared Express template files
// (error middleware, env config) that are architecture-agnostic.
// Callers pass a destPrefix to control where the files land in the output tree:
//   - "src"        → used by MVC (flat src/ layout)
//   - "src/shared" → used by Clean Architecture and Hexagonal (shared/ layer)
package typescript_express_shared

import (
	"embed"
	"fmt"

	"github.com/version14/dot/internal/generator"
	"github.com/version14/dot/internal/generator/genfs"
)

//go:embed all:files
var sharedFiles embed.FS

// SharedOps returns FileOps for error middleware and env config placed under
// destPrefix. The sub-paths (middlewares/ and config/) are preserved.
func SharedOps(destPrefix, generatorName string) ([]generator.FileOp, error) {
	ops, err := genfs.RenderDir(sharedFiles, "files", generatorName, nil)
	if err != nil {
		return nil, fmt.Errorf("express-shared: %w", err)
	}
	for i := range ops {
		ops[i].Path = destPrefix + "/" + ops[i].Path
	}
	return ops, nil
}
