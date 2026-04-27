package commands

import (
	"fmt"
	"time"
)

// PlannedCommand is a fully-resolved command ready to run: tokens already
// substituted, working directory anchored to a project-root-relative path.
//
// Background = true is for long-running commands (dev servers, watchers).
// The runner starts them, waits ReadyDelay, then sends SIGTERM. If the
// process dies before ReadyDelay, that counts as a failure.
type PlannedCommand struct {
	Cmd        string        // post-token-substitution shell line
	WorkDir    string        // project-relative; "" means root
	Source     string        // generator name that contributed this command
	Background bool          // long-running: start, wait, then kill
	ReadyDelay time.Duration // wait this long before considering ready
}

// Key returns the dedup signature: same Cmd + WorkDir + Background = same.
// Source is intentionally excluded so equivalent commands from different
// generators (e.g. two TS generators both running `pnpm install`) collapse.
// Background IS part of the key — a background dev-server and a foreground
// install of the same string are distinct intents.
func (c PlannedCommand) Key() string {
	return fmt.Sprintf("%s::%s::%v", c.WorkDir, c.Cmd, c.Background)
}

// Dedup returns cmds with duplicates removed, preserving the order of first
// appearance. The Source of the first occurrence is kept.
func Dedup(cmds []PlannedCommand) []PlannedCommand {
	seen := make(map[string]bool, len(cmds))
	out := make([]PlannedCommand, 0, len(cmds))
	for _, c := range cmds {
		k := c.Key()
		if seen[k] {
			continue
		}
		seen[k] = true
		out = append(out, c)
	}
	return out
}
