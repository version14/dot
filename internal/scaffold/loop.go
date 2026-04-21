package scaffold

import (
	"strconv"

	q "github.com/version14/dot/internal/question"
)

// runLoop reads the count, runs the body N times as sub-Runners, then attaches
// each iteration's answers as Iterations on the loop's existing entry.
func (runner *Runner) runLoop(loop *q.LoopAction) error {
	count, err := strconv.Atoi(runner.Result.Get(loop.Value))
	if err != nil || count <= 0 {
		return nil
	}

	iterations := make([][]AnswerEntry, 0, count)
	for i := 0; i < count; i++ {
		sub := &Runner{Hooks: runner.Hooks, Flow: loop.Question, parentResult: runner.Result}
		if err := sub.Run(); err != nil {
			return err
		}
		iterations = append(iterations, sub.Result.Entries)
	}

	if idx, ok := runner.Result.index[loop.Value]; ok {
		runner.Result.Entries[idx].Iterations = iterations
	}

	return nil
}
