package cli

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/version14/dot/flows"
	"github.com/version14/dot/internal/commands"
	"github.com/version14/dot/internal/generator"
)

// Dispatch routes os.Args-style arguments to the matching subcommand.
// args MUST NOT include the program name (i.e. caller passes os.Args[1:]).
//
// Returns the desired process exit code so main can `os.Exit(code)`.
func Dispatch(ctx context.Context, args []string, toolVersion string) int {
	if len(args) == 0 {
		printUsage(os.Stdout, toolVersion)
		return 0
	}

	cmd, rest := args[0], args[1:]
	switch cmd {
	case "version", "--version", "-v":
		fmt.Printf("dot %s\n", toolVersion)
		return 0

	case "help", "--help", "-h":
		printUsage(os.Stdout, toolVersion)
		return 0

	case "flows":
		return runListFlows(os.Stdout)

	case "generators":
		return runListGenerators(os.Stdout)

	case "scaffold":
		return runScaffold(ctx, rest, toolVersion)

	case "update":
		return runUpdate(ctx, rest, toolVersion)

	case "self-update":
		return runSelfUpdate()

	case "doctor":
		return runDoctor(ctx, rest)

	case "plugin", "plugins":
		return runPlugin(ctx, rest)

	default:
		fmt.Fprintf(os.Stderr, "dot: unknown command %q\n\n", cmd)
		printUsage(os.Stderr, toolVersion)
		return 2
	}
}

// printUsage writes the top-level help text.
func printUsage(w io.Writer, version string) {
	PrintBanner()
	fmt.Fprintf(w, "dot %s — generative project scaffolding\n\n", version)
	fmt.Fprintln(w, "Usage:")
	fmt.Fprintln(w, "  dot scaffold [flow-id] [-out DIR]   Run an interactive scaffold flow")
	fmt.Fprintln(w, "  dot update [PATH]                   Re-run generators against an existing project")
	fmt.Fprintln(w, "  dot self-update                     Update dot to the latest release")
	fmt.Fprintln(w, "  dot doctor [PATH]                   Diagnose drift between spec and current tools")
	fmt.Fprintln(w, "  dot plugin <list|install|uninstall> Manage installable plugins")
	fmt.Fprintln(w, "  dot flows                           List available flows")
	fmt.Fprintln(w, "  dot generators                      List registered generators")
	fmt.Fprintln(w, "  dot version                         Print the tool version")
	fmt.Fprintln(w, "  dot help                            Show this message")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Examples:")
	fmt.Fprintln(w, "  dot scaffold                        # pick flow interactively")
	fmt.Fprintln(w, "  dot scaffold monorepo               # use the monorepo flow")
	fmt.Fprintln(w, "  dot scaffold fullstack -out /tmp")
	fmt.Fprintln(w, "  dot update ./my-project             # re-run generators after an upgrade")
	fmt.Fprintln(w, "  dot doctor ./my-project             # check spec ↔ installed drift")
}

// runListFlows prints every registered flow ID and title.
func runListFlows(w io.Writer) int {
	reg := flows.Default()
	all := reg.All()

	PrintHeading("Flows")
	for _, f := range all {
		fmt.Fprintf(w, "  %-12s  %s\n", f.ID, f.Title)
		if f.Description != "" {
			fmt.Fprintf(w, "  %-12s  %s\n", "", InfoText(f.Description))
		}
	}
	return 0
}

// runListGenerators prints every registered generator and a short summary.
func runListGenerators(w io.Writer) int {
	reg, err := DefaultGeneratorRegistry()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	PrintHeading("Generators")
	for _, e := range reg.All() {
		fmt.Fprintf(w, "  %-18s  v%-7s  %s\n", e.Manifest.Name, e.Manifest.Version, e.Manifest.Description)
		if len(e.Manifest.DependsOn) > 0 {
			fmt.Fprintf(w, "  %-18s  %s\n", "", InfoText("depends on: "+joinNames(e.Manifest.DependsOn)))
		}
	}
	return 0
}

// runSelfUpdate checks for and installs the latest dot release.
func runSelfUpdate() int {
	if err := cmdSelfUpdate(); err != nil {
		PrintError(err.Error())
		return 1
	}
	return 0
}

// runScaffold parses scaffold-specific flags and delegates to Scaffold.
// After scaffolding succeeds it also runs the deduplicated PostGenerationCommands
// with a Docker-style spinner UX (output captured, surfaced only on failure).
func runScaffold(ctx context.Context, args []string, toolVersion string) int {
	fs := flag.NewFlagSet("scaffold", flag.ContinueOnError)
	out := fs.String("out", ".", "parent directory the project will be created in")
	skipPost := fs.Bool("skip-post", false, "skip post-generation commands (e.g. pnpm install)")
	if err := fs.Parse(args); err != nil {
		return 2
	}

	flowID := ""
	if fs.NArg() > 0 {
		flowID = fs.Arg(0)
	}

	PrintBanner()

	flowReg := flows.Default()
	def, err := pickFlow(flowReg, flowID)
	if err != nil {
		fmt.Fprintln(os.Stderr, "dot scaffold:", err)
		return 1
	}

	rt, err := DefaultRuntime()
	if err != nil {
		fmt.Fprintln(os.Stderr, "dot scaffold:", err)
		return 1
	}

	absOut, err := filepath.Abs(*out)
	if err != nil {
		fmt.Fprintln(os.Stderr, "dot scaffold:", err)
		return 1
	}

	logger := NewStepLogger()
	res, err := Scaffold(ctx, ScaffoldOptions{
		Flow:        def,
		Registry:    rt.Generators,
		Hooks:       rt.Hooks,
		Fragments:   rt.Fragments,
		Plugins:     rt.Plugins,
		OutputDir:   absOut,
		ToolVersion: toolVersion,
		Logger:      logger,
	})
	if err != nil {
		if errors.Is(err, ErrAborted) {
			PrintInfo("aborted")
			return 1
		}
		PrintError(err.Error())
		return 1
	}

	// Post-generation commands.
	if !*skipPost {
		plan := PlanPostGenCommands(res.Spec, res.Manifests)
		if len(plan) > 0 {
			PrintHeading(fmt.Sprintf("post-gen commands (%d)", len(plan)))
			runner := commands.NewRunner(res.ProjectRoot, logger)
			if err := RunCommandsQuiet(ctx, runner, plan, 2); err != nil {
				PrintError("post-gen failed: " + err.Error())
				return 1
			}
		}
	}

	PrintSuccess(fmt.Sprintf("scaffolded %s in %s", res.Spec.Metadata.ProjectName, res.ProjectRoot))
	return 0
}

// pickFlow resolves a flow ID, listing options when ambiguous or empty.
func pickFlow(reg *flows.Registry, id string) (*flows.FlowDef, error) {
	if id != "" {
		def, ok := reg.Get(id)
		if !ok {
			return nil, fmt.Errorf("unknown flow %q (try `dot flows`)", id)
		}
		return def, nil
	}

	all := reg.All()
	if len(all) == 1 {
		return all[0], nil
	}
	// Multiple flows + no explicit pick: show a list and let the user choose
	// by re-running with an explicit flow-id. (Keeps the CLI deterministic.)
	fmt.Fprintln(os.Stderr, "Multiple flows available — re-run with one of:")
	for _, f := range all {
		fmt.Fprintf(os.Stderr, "  dot scaffold %s\n", f.ID)
	}
	return nil, fmt.Errorf("flow not specified")
}

func joinNames(names []string) string {
	out := ""
	for i, n := range names {
		if i > 0 {
			out += ", "
		}
		out += n
	}
	return out
}

// Compile-time guard so we don't accidentally drop generator imports while
// we add features here. The reference is harmless at runtime.
var _ = generator.NewRegistry
