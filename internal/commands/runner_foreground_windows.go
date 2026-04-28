//go:build windows

package commands

import (
	"context"
	"os/exec"
)

func (r *Runner) runForeground(ctx context.Context, c PlannedCommand, wd string, stdout, stderr fileOrBuf) error {
	cmd := exec.CommandContext(ctx, "cmd.exe", "/C", c.Cmd)
	cmd.Dir = wd
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	cmd.Env = getEnviron()
	return cmd.Run()
}
