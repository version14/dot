package cli

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/version14/dot/internal/plugin"
)

// runPlugin dispatches the `dot plugin <subcommand>` family.
//
// Subcommands:
//
//	list                  Show built-in providers + on-disk plugins
//	install <source>      Install from a remote git URL or shorthand
//	uninstall <id>        Remove an installed plugin
func runPlugin(ctx context.Context, args []string) int {
	if len(args) == 0 {
		printPluginUsage(os.Stdout)
		return 0
	}
	sub, rest := args[0], args[1:]
	switch sub {
	case "list", "ls":
		return runPluginList()
	case "install", "add":
		return runPluginInstall(ctx, rest)
	case "uninstall", "remove", "rm":
		return runPluginUninstall(rest)
	default:
		fmt.Fprintf(os.Stderr, "dot plugin: unknown subcommand %q\n\n", sub)
		printPluginUsage(os.Stderr)
		return 2
	}
}

func printPluginUsage(w io.Writer) {
	fmt.Fprintln(w, "Usage:")
	fmt.Fprintln(w, "  dot plugin list                                  List built-in + installed plugins")
	fmt.Fprintln(w, "  dot plugin install <source> [-ref REF]           Install from git remote")
	fmt.Fprintln(w, "  dot plugin install -from PATH                    (dev) install from a local copy")
	fmt.Fprintln(w, "  dot plugin uninstall <id>                        Remove an installed plugin")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Source forms:")
	fmt.Fprintln(w, "  github.com/owner/repo                  → https://github.com/owner/repo.git")
	fmt.Fprintln(w, "  github.com/owner/repo@v1.2.0           → clone, then checkout v1.2.0")
	fmt.Fprintln(w, "  https://example.com/path/repo.git      → clone direct")
	fmt.Fprintln(w, "  git@github.com:owner/repo.git          → ssh clone (ref via -ref)")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Examples:")
	fmt.Fprintln(w, "  dot plugin install github.com/version14/dot-plugin-biome-extras")
	fmt.Fprintln(w, "  dot plugin install github.com/me/my-plugin -ref v0.2.0")
	fmt.Fprintln(w, "  dot plugin install -from ./my-local-plugin")
}

// runPluginList enumerates built-in providers (compiled in) plus any plugins
// discovered under PluginsDir().
func runPluginList() int {
	rt, err := DefaultRuntime()
	if err != nil {
		PrintError(err.Error())
		return 1
	}

	PrintHeading("Built-in plugins")
	if len(rt.Plugins) == 0 {
		PrintInfo("  (none)")
	}
	for _, p := range rt.Plugins {
		fmt.Printf("  %-24s  %d generators · %d injections\n",
			p.ID(),
			len(p.Generators()),
			len(p.Injections()),
		)
	}

	installed, err := plugin.List()
	if err != nil {
		PrintError(err.Error())
		return 1
	}
	PrintHeading("Installed plugins (~/.dot/plugins)")
	if len(installed) == 0 {
		PrintInfo("  (none)")
		return 0
	}
	for _, p := range installed {
		desc := p.Description
		if desc == "" {
			desc = "(no description)"
		}
		fmt.Printf("  %-24s  v%-8s  %s\n", p.ID, p.Version, desc)
	}
	return 0
}

// runPluginInstall accepts either a remote source as the positional arg or
// `-from PATH` for development installs. The remote path is the primary
// flow; -from exists so plugin authors can iterate locally before pushing.
func runPluginInstall(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("plugin install", flag.ContinueOnError)
	from := fs.String("from", "", "(dev) local directory to copy into the plugin store")
	ref := fs.String("ref", "", "git ref (tag/branch/commit) to checkout after clone")
	overrideID := fs.String("id", "", "override the plugin id recorded in plugin.json")
	overrideVer := fs.String("version", "", "override the version recorded in plugin.json")
	if err := fs.Parse(args); err != nil {
		return 2
	}

	source := ""
	if fs.NArg() > 0 {
		source = fs.Arg(0)
	}

	if source == "" && *from == "" {
		PrintError("provide a source URL or -from PATH")
		printPluginUsage(os.Stderr)
		return 2
	}
	if source != "" && *from != "" {
		PrintError("pass either a source URL OR -from PATH, not both")
		return 2
	}

	if source != "" {
		PrintInfo(fmt.Sprintf("→ cloning %s%s", source, refSuffix(*ref)))
	} else {
		PrintInfo(fmt.Sprintf("→ copying from %s", *from))
	}

	installed, err := plugin.Install(ctx, plugin.InstallSpec{
		Source:          source,
		Ref:             *ref,
		LocalPath:       *from,
		OverrideID:      *overrideID,
		OverrideVersion: *overrideVer,
	})
	if err != nil {
		PrintError(err.Error())
		return 1
	}

	PrintSuccess(fmt.Sprintf("installed %s@%s → %s", installed.ID, installed.Version, installed.Dir))
	PrintInfo("rebuild `dot` (or restart it) so the plugin's init() runs and registers its hooks")
	return 0
}

func refSuffix(ref string) string {
	if ref == "" {
		return ""
	}
	return "@" + ref
}

// runPluginUninstall removes a plugin's directory. Idempotent.
func runPluginUninstall(args []string) int {
	if len(args) == 0 {
		PrintError("plugin id is required")
		return 2
	}
	id := args[0]
	if err := plugin.Uninstall(id); err != nil {
		PrintError(err.Error())
		return 1
	}
	PrintSuccess(fmt.Sprintf("uninstalled %s", id))
	return 0
}
