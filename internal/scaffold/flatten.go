package scaffold

import (
	"fmt"

	q "github.com/version14/dot/internal/question"
)

// conditionType represents one visibility condition: "key must equal value" (or not, if negate).
// A flatNode is shown only when ALL its conditions are satisfied.
type conditionType struct {
	key, value string
	negate     bool
}

// flatNode is a question tagged with the conditions for it to be visible.
type flatNode struct {
	question   *q.Question
	conditions []conditionType
}

// loopEntry defers a loop's body execution until after the main form.
type loopEntry struct {
	loop       *q.LoopAction
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

// flatten walks the question tree, accumulating visibility conditions.
// Returns an error if two questions with the same key appear on the same linear
// path (silent data loss). Keys that appear in mutually exclusive branches are
// allowed — only one branch runs, so no collision occurs at answer time.
// Plugin authors must namespace their keys: "plugin-name.question-key".
func flatten(question *q.Question, conditions []conditionType, nodes *[]flatNode, loops *[]loopEntry) error {
	return flattenNode(question, conditions, nodes, loops, make(map[string]struct{}))
}

func flattenNode(question *q.Question, conditions []conditionType, nodes *[]flatNode, loops *[]loopEntry, seen map[string]struct{}) error {
	if question == nil {
		return nil
	}

	if question.TextQuestion != nil {
		return flattenText(question, conditions, nodes, loops, seen)
	}
	if question.OptionQuestion != nil {
		return flattenOption(question, conditions, nodes, loops, seen)
	}
	if question.LoopQuestion != nil {
		return flattenLoop(question, conditions, nodes, loops, seen)
	}
	if question.IfQuestion != nil {
		return flattenIf(question, conditions, nodes, loops, seen)
	}
	return nil
}

func flattenText(question *q.Question, conditions []conditionType, nodes *[]flatNode, loops *[]loopEntry, seen map[string]struct{}) error {
	tq := question.TextQuestion
	if err := checkDuplicate(seen, tq.Value); err != nil {
		return err
	}
	*nodes = append(*nodes, flatNode{question: question, conditions: conditions})
	if tq.Next != nil {
		return flattenNode(tq.Next.Question, conditions, nodes, loops, seen)
	}
	return nil
}

func flattenOption(question *q.Question, conditions []conditionType, nodes *[]flatNode, loops *[]loopEntry, seen map[string]struct{}) error {
	oq := question.OptionQuestion
	if err := checkDuplicate(seen, oq.Value); err != nil {
		return err
	}
	*nodes = append(*nodes, flatNode{question: question, conditions: conditions})
	for _, option := range oq.Options {
		if option.Next != nil {
			optCond := withCond(conditions, conditionType{key: oq.Value, value: option.Value})
			// Each option branch gets its own copy of seen so that keys
			// reused across mutually exclusive branches don't collide.
			if err := flattenNode(option.Next.Question, optCond, nodes, loops, copyMap(seen)); err != nil {
				return err
			}
		}
	}
	if oq.Next != nil {
		return flattenNode(oq.Next.Question, conditions, nodes, loops, seen)
	}
	return nil
}

func flattenLoop(question *q.Question, conditions []conditionType, nodes *[]flatNode, loops *[]loopEntry, seen map[string]struct{}) error {
	lq := question.LoopQuestion
	if err := checkDuplicate(seen, lq.Value); err != nil {
		return err
	}
	// Insert a text node for the count prompt; defer the body until post-form.
	countQ := &q.Question{
		TextQuestion: &q.TextQuestion{
			Label: lq.Label,
			Value: lq.Value,
		},
	}
	*nodes = append(*nodes, flatNode{question: countQ, conditions: conditions})
	*loops = append(*loops, loopEntry{loop: lq, conditions: conditions})
	return nil
}

func flattenIf(question *q.Question, conditions []conditionType, nodes *[]flatNode, loops *[]loopEntry, seen map[string]struct{}) error {
	ia := question.IfQuestion
	// Then: visible only when key == value.
	thenCond := withCond(conditions, conditionType{key: ia.Key, value: ia.Value, negate: false})
	// Else: visible only when key != value.
	elseCond := withCond(conditions, conditionType{key: ia.Key, value: ia.Value, negate: true})
	// Each branch gets its own copy of seen — they are mutually exclusive.
	if err := flattenNode(ia.Then, thenCond, nodes, loops, copyMap(seen)); err != nil {
		return err
	}
	return flattenNode(ia.Else, elseCond, nodes, loops, copyMap(seen))
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
		return fmt.Errorf("duplicate question key %q — plugin authors must namespace keys (e.g. \"my-plugin.%s\")", key, key)
	}
	seen[key] = struct{}{}
	return nil
}
