package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/version14/dot/internal/cli"

	// Built-in plugins. Imported for their side-effect: each calls
	// plugin.RegisterBuiltin(...) in its init() so DefaultRuntime sees them.
	_ "github.com/version14/dot/plugins/biome_extras"
)

var toolVersion string

func main() {
	// Cancel ctx on Ctrl-C so long-running scaffolds + post-gen commands
	// stop promptly without leaving dangling child processes.
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	os.Exit(cli.Dispatch(ctx, os.Args[1:], toolVersion))
}
