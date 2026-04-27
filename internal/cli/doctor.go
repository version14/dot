package cli

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/version14/dot/flows"
	"github.com/version14/dot/internal/dotdir"
	"github.com/version14/dot/internal/generator"
	"github.com/version14/dot/internal/versioning"
	"github.com/version14/dot/pkg/dotapi"
)

// DoctorReport summarizes the state of an existing scaffolded project.
type DoctorReport struct {
	ProjectRoot string

	// FlowOK is false when the spec's FlowID is not in the current registry.
	FlowOK bool

	// MissingGenerators is the set of names referenced by the manifest but
	// absent from the current generator registry.
	MissingGenerators []string

	// VersionDrift records generators whose currently-installed version no
	// longer satisfies the constraint recorded in the spec.
	VersionDrift []DriftEntry

	// ValidationFailures is the result of running every Validator's Checks
	// against the on-disk project.
	ValidationFailures []generator.ValidationFailure
}

// DriftEntry is one generator with mismatched constraint vs installed.
type DriftEntry struct {
	Name             string
	Constraint       string // from spec
	InstalledVersion string // from registry
}

// OK reports whether the project is healthy (no missing gens, no drift, no
// failed validators, flow still known).
func (r *DoctorReport) OK() bool {
	return r.FlowOK &&
		len(r.MissingGenerators) == 0 &&
		len(r.VersionDrift) == 0 &&
		len(r.ValidationFailures) == 0
}

// Doctor inspects an existing project and returns a DoctorReport describing
// every drift between what the spec recorded and what the running tool can
// honour today.
func Doctor(ctx context.Context, projectRoot string, registry *generator.Registry) (*DoctorReport, error) {
	_ = ctx
	if registry == nil {
		return nil, fmt.Errorf("cli: doctor: nil registry")
	}

	abs, err := filepath.Abs(projectRoot)
	if err != nil {
		return nil, fmt.Errorf("cli: doctor: %w", err)
	}

	s, err := dotdir.LoadSpec(abs)
	if err != nil {
		return nil, fmt.Errorf("cli: doctor: load spec: %w", err)
	}

	report := &DoctorReport{ProjectRoot: abs}

	// 1. Is the originating flow still registered?
	flowReg := flows.Default()
	def, ok := flowReg.Get(s.FlowID)
	report.FlowOK = ok

	// 2. Are all generators referenced by the flow's resolver still available?
	if def != nil && def.Generators != nil {
		flowInvs := def.Generators(s)
		mans := make([]dotapi.Manifest, 0, len(flowInvs))
		for _, fi := range flowInvs {
			entry, ok := registry.Get(fi.Name)
			if !ok {
				report.MissingGenerators = append(report.MissingGenerators, fi.Name)
				continue
			}
			mans = append(mans, entry.Manifest)

			// Drift check (only if a constraint was recorded in the spec).
			cnstStr, hasCnst := s.GeneratorConstraints[fi.Name]
			if !hasCnst {
				continue
			}
			cnst, err := versioning.ParseConstraint(cnstStr)
			if err != nil {
				continue // malformed constraint = treat as no constraint
			}
			installed, err := versioning.Parse(entry.Manifest.Version)
			if err != nil {
				continue
			}
			if !cnst.Allows(installed) {
				report.VersionDrift = append(report.VersionDrift, DriftEntry{
					Name:             fi.Name,
					Constraint:       cnstStr,
					InstalledVersion: installed.String(),
				})
			}
		}

		// 3. Run validators against the on-disk tree.
		failures, err := generator.RunValidators(abs, mans)
		if err != nil {
			return nil, fmt.Errorf("cli: doctor: validators: %w", err)
		}
		report.ValidationFailures = failures
	}

	return report, nil
}

// runDoctor is the CLI dispatch for `dot doctor [path]`.
func runDoctor(ctx context.Context, args []string) int {
	root := "."
	if len(args) > 0 {
		root = args[0]
	}

	rt, err := DefaultRuntime()
	if err != nil {
		PrintError(err.Error())
		return 1
	}

	report, err := Doctor(ctx, root, rt.Generators)
	if err != nil {
		PrintError(err.Error())
		return 1
	}

	PrintHeading(fmt.Sprintf("Project health: %s", report.ProjectRoot))

	if !report.FlowOK {
		PrintError("originating flow is no longer registered (re-installs may be required)")
	}
	if len(report.MissingGenerators) > 0 {
		PrintError("missing generators:")
		for _, n := range report.MissingGenerators {
			fmt.Printf("    - %s\n", n)
		}
	}
	if len(report.VersionDrift) > 0 {
		PrintWarn("version drift:")
		for _, d := range report.VersionDrift {
			fmt.Printf("    - %s: installed %s, constraint %s\n", d.Name, d.InstalledVersion, d.Constraint)
		}
	}
	if len(report.ValidationFailures) > 0 {
		PrintError("validator failures:")
		for _, f := range report.ValidationFailures {
			fmt.Printf("    - %s\n", f.String())
		}
	}

	if report.OK() {
		PrintSuccess("project is healthy")
		return 0
	}
	return 1
}
