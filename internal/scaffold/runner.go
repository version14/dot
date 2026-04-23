package scaffold

import (
	"encoding/json"
	"fmt"

	"github.com/charmbracelet/huh"
	q "github.com/version14/dot/internal/question"
)

// Hook is called after the form completes, once per answered question.
type Hook func(key, value string, result *Result)

// AnswerEntry is one answered question in flow order.
// Loop entries carry Iterations instead of a flat Value.
type AnswerEntry struct {
	Key        string          `json:"key"`
	Value      string          `json:"value,omitempty"`
	Multi      []string        `json:"multi,omitempty"`
	Iterations [][]AnswerEntry `json:"iterations,omitempty"`
}

// Result holds answers in the order the user encountered the question.
type Result struct {
	Entries []AnswerEntry
	index   map[string]int // key → position in Entries for O(1) lookup
}

// Get returns the string answer for key, or "" if not answered.
func (result *Result) Get(key string) string {
	if result.index == nil {
		return ""
	}
	if i, ok := result.index[key]; ok {
		return result.Entries[i].Value
	}
	return ""
}

// GetMulti returns the multi-select answer for key, or nil if not answered.
func (result *Result) GetMulti(key string) []string {
	if result.index == nil {
		return nil
	}
	if i, ok := result.index[key]; ok {
		return result.Entries[i].Multi
	}
	return nil
}

// ToJSON serializes the ordered answer list.
func (result *Result) ToJSON() ([]byte, error) {
	return json.MarshalIndent(result.Entries, "", "  ")
}

func (result *Result) add(entry AnswerEntry) {
	if result.index == nil {
		result.index = make(map[string]int)
	}
	result.index[entry.Key] = len(result.Entries)
	result.Entries = append(result.Entries, entry)
}

// Add appends an answer entry and indexes it by key. Exported so callers that
// synthesize a Result (tests, alternative input layers) can mirror what the
// survey runner produces.
func (result *Result) Add(entry AnswerEntry) {
	result.add(entry)
}

// Runner holds the question flow and accumulates answers.
type Runner struct {
	Flow  *q.Question
	Hooks map[string]Hook
	// Result is populated after Run() returns.
	Result       *Result
	parentResult *Result            // answers from the outer scope; used by loop sub-runners to resolve cross-scope conditions
	strPtrs      map[string]*string // live pointers written by huh during form execution
	multiPtrs    map[string]*[]string
}

// Run executes the survey and populates runner.Result.
//
// Algorithm:
//  1. flatten()     — walk the tree once; each node gets the conditions under
//     which it should be visible; loops are deferred.
//  2. buildGroups() — one huh.Group per node; conditional groups use
//     WithHideFunc reading live *string pointers.
//  3. form.Run()    — single combined form; branches hidden until conditions are met.
//  4. Collect       — iterate nodes in flow order, skip branches not taken, populate Result.
//  5. Hooks         — fire per-key hooks.
//  6. Loops         — run each loop body N times as sub-Runners; merge with prefix.
func (runner *Runner) Run() error {
	runner.Result = &Result{}
	runner.strPtrs = make(map[string]*string)
	runner.multiPtrs = make(map[string]*[]string)

	var nodes []flatNode
	var loops []loopEntry
	if err := flatten(runner.Flow, nil, &nodes, &loops); err != nil {
		return err
	}

	groups := runner.buildGroups(nodes)
	if len(groups) > 0 {
		form := huh.NewForm(groups...)
		if err := form.Run(); err != nil {
			return fmt.Errorf("form: %w", err)
		}
	}

	for _, node := range nodes {
		if !runner.condsMatch(node.conditions) {
			continue
		}
		key := nodeKey(node.question)
		if part, ok := runner.strPtrs[key]; ok {
			runner.Result.add(AnswerEntry{Key: key, Value: *part})
		} else if part, ok := runner.multiPtrs[key]; ok {
			runner.Result.add(AnswerEntry{Key: key, Multi: *part})
		}
	}

	if runner.Hooks != nil {
		for _, entry := range runner.Result.Entries {
			if hook, ok := runner.Hooks[entry.Key]; ok {
				hook(entry.Key, entry.Value, runner.Result)
			}
		}
	}

	for _, le := range loops {
		if !runner.condsMatch(le.conditions) {
			continue
		}
		if err := runner.runLoop(le.loop); err != nil {
			return err
		}
	}

	return nil
}

// condsMatch returns true when every condition is satisfied.
// Reads live strPtrs during form execution; falls back to Result after.
func (runner *Runner) condsMatch(conditions []conditionType) bool {
	for _, condition := range conditions {
		var actual string
		if part, ok := runner.strPtrs[condition.key]; ok {
			actual = *part // live — written by huh as the user interacts
		} else {
			actual = runner.Result.Get(condition.key)
			// Fall back to parent scope for conditions inside a loop body
			// that were answered in the outer form.
			if actual == "" && runner.parentResult != nil {
				actual = runner.parentResult.Get(condition.key)
			}
		}
		matches := actual == condition.value
		if condition.negate && matches {
			return false
		}
		if !condition.negate && !matches {
			return false
		}
	}
	return true
}
