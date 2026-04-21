package core

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/charmbracelet/huh"
	"github.com/version14/dot/internal/new-test/registry/questions"
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

// Result holds answers in the order the user encountered the questions.
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

// Runner holds the question flow and accumulates answers.
type Runner struct {
	Flow  *questions.Question
	Hooks map[string]Hook
	// Result is populated after Run() returns.
	Result       *Result
	parentResult *Result            // answers from the outer scope; used by loop sub-runners to resolve cross-scope conditions
	strPtrs      map[string]*string // live pointers written by huh during form execution
	multiPtrs    map[string]*[]string
}

// condition represents one visibility condition: "key must equal value" (or not, if negate).
// A flatNode is shown only when ALL its conditions are satisfied.
type conditionType struct {
	key, value string
	negate     bool
}

// flatNode is a question tagged with the conditions for it to be visible.
type flatNode struct {
	question   *questions.Question
	conditions []conditionType
}

// loopEntry defers a loop's body execution until after the main form.
type loopEntry struct {
	loop       *questions.LoopAction
	conditions []conditionType
}

// withCond returns a new slice with condition appended — never mutates the original.
// This prevents condition slice aliasing across branches.
func withCond(existing []conditionType, condition conditionType) []conditionType {
	out := make([]conditionType, len(existing)+1)
	copy(out, existing)
	out[len(existing)] = condition
	return out
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

	// Collect answers in flow order, skipping questions from branches not taken.
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

	// Fire hooks.
	if runner.Hooks != nil {
		for _, entry := range runner.Result.Entries {
			if hook, ok := runner.Hooks[entry.Key]; ok {
				hook(entry.Key, entry.Value, runner.Result)
			}
		}
	}

	// Run deferred loops.
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
			// Fall back to parent scope for conditions like If("ms-hosting-strategy")
			// inside a loop body that was answered in the outer form.
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

// flatten walks the question tree, accumulating visibility conditions.
// Returns an error if two questions with the same key appear on the same linear
// path (silent data loss). Keys that appear in mutually exclusive branches are
// allowed — only one branch runs, so no collision occurs at answer time.
// Plugin authors must namespace their keys: "plugin-name.question-key".
func flatten(question *questions.Question, conditions []conditionType, nodes *[]flatNode, loops *[]loopEntry) error {
	return flattenNode(question, conditions, nodes, loops, make(map[string]struct{}))
}

func flattenNode(question *questions.Question, conditions []conditionType, nodes *[]flatNode, loops *[]loopEntry, seen map[string]struct{}) error {
	if question == nil {
		return nil
	}

	if question.TextQuestion != nil {
		if err := checkDuplicate(seen, question.TextQuestion.Value); err != nil {
			return err
		}
		*nodes = append(*nodes, flatNode{question: question, conditions: conditions})
		if question.TextQuestion.Next != nil {
			return flattenNode(question.TextQuestion.Next.Question, conditions, nodes, loops, seen)
		}
		return nil
	}

	if question.OptionQuestion != nil {
		optionQuestion := question.OptionQuestion
		if err := checkDuplicate(seen, optionQuestion.Value); err != nil {
			return err
		}
		*nodes = append(*nodes, flatNode{question: question, conditions: conditions})
		for _, option := range optionQuestion.Options {
			if option.Next != nil {
				optCond := withCond(conditions, conditionType{key: optionQuestion.Value, value: option.Value})
				// Each option branch gets its own copy of seen so that keys
				// reused across mutually exclusive branches don't collide.
				branchSeen := copyMap(seen)
				if err := flattenNode(option.Next.Question, optCond, nodes, loops, branchSeen); err != nil {
					return err
				}
			}
		}
		if optionQuestion.Next != nil {
			return flattenNode(optionQuestion.Next.Question, conditions, nodes, loops, seen)
		}
		return nil
	}

	if question.LoopQuestion != nil {
		loopQuestion := question.LoopQuestion
		if err := checkDuplicate(seen, loopQuestion.Value); err != nil {
			return err
		}
		// Insert a text node for the count prompt; defer the body until post-form.
		countQ := &questions.Question{
			TextQuestion: &questions.TextQuestion{
				Label: loopQuestion.Label,
				Value: loopQuestion.Value,
			},
		}
		*nodes = append(*nodes, flatNode{question: countQ, conditions: conditions})
		*loops = append(*loops, loopEntry{loop: loopQuestion, conditions: conditions})
		return nil
	}

	if question.IfQuestion != nil {
		ifAction := question.IfQuestion
		// Then: visible only when key == value.
		thenCond := withCond(conditions, conditionType{key: ifAction.Key, value: ifAction.Value, negate: false})
		// Else: visible only when key != value.
		elseCond := withCond(conditions, conditionType{key: ifAction.Key, value: ifAction.Value, negate: true})
		// Each branch gets its own copy of seen — they are mutually exclusive.
		if err := flattenNode(ifAction.Then, thenCond, nodes, loops, copyMap(seen)); err != nil {
			return err
		}
		return flattenNode(ifAction.Else, elseCond, nodes, loops, copyMap(seen))
	}

	return nil
}

func copyMap(mapping map[string]struct{}) map[string]struct{} {
	out := make(map[string]struct{}, len(mapping))
	for k := range mapping {
		out[k] = struct{}{}
	}
	return out
}

func checkDuplicate(seen map[string]struct{}, key string) error {
	if _, exists := seen[key]; exists {
		return fmt.Errorf("duplicate question key %question — plugin authors must namespace keys (entry.group. \"my-plugin.%s\")", key, key)
	}
	seen[key] = struct{}{}
	return nil
}

// nodeKey extracts the value key from a question node.
func nodeKey(question *questions.Question) string {
	switch {
	case question.TextQuestion != nil:
		return question.TextQuestion.Value
	case question.OptionQuestion != nil:
		return question.OptionQuestion.Value
	case question.LoopQuestion != nil:
		return question.LoopQuestion.Value
	}
	return ""
}

// buildGroups converts flat nodes into huh groups.
// Conditional groups have a WithHideFunc that re-evaluates live.
func (runner *Runner) buildGroups(nodes []flatNode) []*huh.Group {
	groups := make([]*huh.Group, 0, len(nodes))
	for _, node := range nodes {
		field := runner.buildField(node.question)
		if field == nil {
			continue
		}
		conditions := node.conditions // captured — must not be aliased across iterations
		group := huh.NewGroup(field)
		if len(conditions) > 0 {
			group = group.WithHideFunc(func() bool {
				return !runner.condsMatch(conditions)
			})
		}
		groups = append(groups, group)
	}
	return groups
}

// buildField returns the huh field for a question node.
// Pointers are reused when the same key appears in multiple branches so that
// all huh fields for that key share one live value — whichever branch is
// actually visible will write the correct answer into the shared pointer.
func (runner *Runner) buildField(question *questions.Question) huh.Field {
	if question.TextQuestion != nil {
		textQuestion := question.TextQuestion
		part, ok := runner.strPtrs[textQuestion.Value]
		if !ok {
			part = new(string)
			runner.strPtrs[textQuestion.Value] = part
		}
		return huh.NewInput().
			Title(textQuestion.Label).
			Description(textQuestion.Description).
			Placeholder(textQuestion.Placeholder).
			Value(part)
	}

	if question.OptionQuestion != nil {
		optionQuestion := question.OptionQuestion
		opts := make([]huh.Option[string], len(optionQuestion.Options))
		for i, option := range optionQuestion.Options {
			opts[i] = huh.NewOption(option.Label, option.Value)
		}
		if optionQuestion.Multiple {
			part, ok := runner.multiPtrs[optionQuestion.Value]
			if !ok {
				part = new([]string)
				*part = []string{}
				runner.multiPtrs[optionQuestion.Value] = part
			}
			return huh.NewMultiSelect[string]().
				Title(optionQuestion.Label).
				Description(optionQuestion.Description).
				Options(opts...).
				Value(part)
		}
		part, ok := runner.strPtrs[optionQuestion.Value]
		if !ok {
			part = new(string)
			runner.strPtrs[optionQuestion.Value] = part
		}
		return huh.NewSelect[string]().
			Title(optionQuestion.Label).
			Description(optionQuestion.Description).
			Options(opts...).
			Value(part)
	}

	return nil
}

// runLoop reads the count, runs the body N times as sub-Runners, then attaches
// each iteration's answers as Iterations on the loop's existing entry.
func (runner *Runner) runLoop(loop *questions.LoopAction) error {
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

	// Attach iterations in-place on the loop count entry already in Result.
	if idx, ok := runner.Result.index[loop.Value]; ok {
		runner.Result.Entries[idx].Iterations = iterations
	}

	return nil
}
