package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/version14/dot/internal/generator"
	"github.com/version14/dot/internal/new-test/core"
	templates "github.com/version14/dot/internal/new-test/question-tempaltes"
	"github.com/version14/dot/internal/pipeline"
	"github.com/version14/dot/internal/spec"
)

func main() {
	// 1. Run the survey
	runner := core.Runner{Flow: templates.StarterQuestions}
	if err := runner.Run(); err != nil {
		fatal("survey", err)
	}

	result := runner.Result

	// 2. Recap + confirm
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
	)
	if err := confirmForm.Run(); err != nil {
		fatal("confirm", err)
	}
	if !confirmed {
		fmt.Println("Cancelled.")
		return
	}

	// 3. Build the Spec from survey answers
	s := spec.Spec{
		Project: spec.ProjectSpec{
			Name:     result.Get("project-name"),
			Language: result.Get("frontend-language"),
			Type:     spec.ProjectTypeFrontend,
		},
		Extensions: map[string]any{
			"framework":    result.Get("frontend-framework"),
			"architecture": result.Get("frontend-architecture"),
		},
	}

	// 4. Collect activated generators
	activations := core.Collect(templates.StarterQuestions, result)
	if len(activations) == 0 {
		fmt.Println("no generators activated — nothing to scaffold")
		return
	}

	// 5. Run generators → collect FileOps and PostOps
	var fileOps []generator.FileOp
	var postOps []generator.PostOp

	for _, a := range activations {
		fops, pops, err := a.Fn(s)
		if err != nil {
			fatal(fmt.Sprintf("generator [%s=%s]", a.QuestionKey, a.AnswerValue), err)
		}
		fileOps = append(fileOps, fops...)
		postOps = append(postOps, pops...)
	}

	// 6. Write files
	if err := pipeline.Run(fileOps); err != nil {
		fatal("pipeline", err)
	}
	fmt.Printf("wrote %d file(s)\n", len(fileOps))

	// 7. Run post-ops (pnpm install, etc.)
	for _, pop := range postOps {
		fmt.Printf("running: %s %s\n", pop.Command, strings.Join(pop.Args, " "))
		cmd := exec.Command(pop.Command, pop.Args...)
		if pop.Dir != "" && pop.Dir != "." {
			cmd.Dir = pop.Dir
		}
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fatal(fmt.Sprintf("%s %v", pop.Command, pop.Args), err)
		}
	}

	fmt.Println("done.")
}

func buildRecap(result *core.Result) string {
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

func fatal(context string, err error) {
	fmt.Fprintf(os.Stderr, "error [%s]: %v\n", context, err)
	os.Exit(1)
}
