package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/version14/dot/internal/generator"
	"github.com/version14/dot/internal/pipeline"
	"github.com/version14/dot/internal/scaffold"
	"github.com/version14/dot/internal/spec"
	templates "github.com/version14/dot/internal/templates"
)

func cmdInit() error {
	runner := scaffold.Runner{Flow: templates.StarterQuestions}
	if err := runner.Run(); err != nil {
		return fmt.Errorf("survey: %w", err)
	}

	result := runner.Result

	var confirmed bool
	confirmForm := huh.NewForm(
		huh.NewGroup(
			huh.NewNote().
				Title("Your choices").
				Description(buildRecap(result)),
			huh.NewConfirm().
				Title("Generate the project?").
				Affirmative("Yes, let's go").
				Negative("Cancel").
				Value(&confirmed),
		),
	).WithTheme(themeDot())
	if err := confirmForm.Run(); err != nil {
		return fmt.Errorf("confirm: %w", err)
	}
	if !confirmed {
		fmt.Println(mutedStyle.Render("cancelled."))
		return nil
	}

	s := spec.Spec{
		Project: spec.ProjectSpec{
			Name:     result.Get("project-name"),
			Language: result.Get("frontend-language"),
			Type:     "frontend",
		},
		Extensions: map[string]any{
			"framework":    result.Get("frontend-framework"),
			"architecture": result.Get("frontend-architecture"),
		},
	}

	activations := scaffold.Collect(templates.StarterQuestions, result)
	if len(activations) == 0 {
		fmt.Println(mutedStyle.Render("no generators activated — nothing to scaffold"))
		return nil
	}

	var fileOps []generator.FileOp
	var postOps []generator.PostOp

	for _, a := range activations {
		fops, pops, err := a.Fn(s)
		if err != nil {
			return fmt.Errorf("generator [%s=%s]: %w", a.QuestionKey, a.AnswerValue, err)
		}
		fileOps = append(fileOps, fops...)
		postOps = append(postOps, pops...)
	}

	if err := pipeline.Run(fileOps); err != nil {
		return fmt.Errorf("pipeline: %w", err)
	}
	fmt.Printf("%s  wrote %d file(s)\n", successStyle.Render("✓"), len(fileOps))

	for _, pop := range postOps {
		fmt.Printf("%s  running: %s %s\n", mutedStyle.Render("→"), pop.Command, strings.Join(pop.Args, " "))
		cmd := exec.Command(pop.Command, pop.Args...)
		if pop.Dir != "" && pop.Dir != "." {
			cmd.Dir = pop.Dir
		}
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("%s %v: %w", pop.Command, pop.Args, err)
		}
	}

	fmt.Println(successStyle.Render("✓  done."))
	return nil
}

func buildRecap(result *scaffold.Result) string {
	var lines []string
	for _, e := range result.Entries {
		label := keyToLabel(e.Key)
		switch {
		case len(e.Iterations) > 0:
			lines = append(lines, fmt.Sprintf("%-22s %d iteration(s)", label+":", len(e.Iterations)))
		case len(e.Multi) > 0:
			lines = append(lines, fmt.Sprintf("%-22s %s", label+":", strings.Join(e.Multi, ", ")))
		case e.Value != "":
			lines = append(lines, fmt.Sprintf("%-22s %s", label+":", e.Value))
		}
	}
	return strings.Join(lines, "\n")
}

func keyToLabel(key string) string {
	parts := strings.Split(key, "-")
	for i, p := range parts {
		if len(p) > 0 {
			parts[i] = strings.ToUpper(p[:1]) + p[1:]
		}
	}
	return strings.Join(parts, " ")
}
