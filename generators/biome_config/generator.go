package biomeconfig

import (
	"github.com/version14/dot/internal/deps"
	"github.com/version14/dot/internal/state"
	"github.com/version14/dot/pkg/dotapi"
)

type Generator struct{}

func New() *Generator { return &Generator{} }

func (g *Generator) Name() string    { return Manifest.Name }
func (g *Generator) Version() string { return Manifest.Version }

func (g *Generator) Generate(ctx *dotapi.Context) error {
	if err := ctx.State.UpdateJSON("biome.json", func(d *state.JSONDoc) error {
		d.Merge(map[string]interface{}{
			"$schema": "https://biomejs.dev/schemas/2.4.16/schema.json",
			"linter": map[string]interface{}{
				"enabled": true,
				"rules":   map[string]interface{}{"recommended": true},
			},
			"formatter": map[string]interface{}{
				"enabled":     true,
				"indentStyle": "space",
				"indentWidth": 2,
			},
			"files": map[string]interface{}{},
		})

		// Biome 2.x: files.ignore and experimentalScannerIgnores were removed.
		// Use files.includes negation patterns. Per the "Valid Folder Ignore
		// Pattern" docs, directory exclusions must use the **/<dir> form — a
		// bare name like !dist does not match the directory's contents.
		// !! is a hard-exclude (skip scanner) for the dot-owned meta dir.
		// "**" MUST be position 0; collect any entries added by earlier
		// generators (e.g. panda_css's !**/styled-system) and append them after.
		standard := []string{
			"**",
			"!!**/.dot",
			"!**/.next",
			"!**/coverage",
			"!**/dist",
			"!**/playwright-report",
			"!**/storybook-static",
			"!src/routeTree.gen.ts",
		}
		seen := make(map[string]struct{}, len(standard))
		final := make([]interface{}, 0, len(standard))
		for _, s := range standard {
			seen[s] = struct{}{}
			final = append(final, s)
		}
		if raw, ok := d.GetNested("files.includes"); ok {
			if arr, ok := raw.([]interface{}); ok {
				for _, v := range arr {
					if s, ok := v.(string); ok {
						if _, dup := seen[s]; !dup {
							seen[s] = struct{}{}
							final = append(final, s)
						}
					}
				}
			}
		}
		return d.SetNested("files.includes", final)
	}); err != nil {
		return err
	}

	return ctx.State.UpdateJSON("package.json", func(d *state.JSONDoc) error {
		d.Merge(map[string]interface{}{
			"scripts": map[string]interface{}{
				"lint":   "biome check .",
				"format": "biome format --write .",
			},
			"devDependencies": deps.NPM("@biomejs/biome"),
		})
		return nil
	})
}
