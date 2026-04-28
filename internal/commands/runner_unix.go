//go:build !windows

package commands

import (
	"context"
	"fmt"
	"os/exec"
	"syscall"
	"time"
)

// runBackground starts the command, waits ReadyDelay, verifies it's still
// running, then sends SIGTERM (with a SIGKILL fallback after 5s).
func (r *Runner) runBackground(ctx context.Context, c PlannedCommand, wd string, stdout, stderr fileOrBuf) error {
	delay := c.ReadyDelay
	if delay <= 0 {
		delay = defaultBackgroundReadyDelay
	}

	cmd := exec.Command("/bin/sh", "-c", c.Cmd)
	cmd.Dir = wd
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	cmd.Env = getEnviron()
	// Put the child in its own process group so we can kill its descendants.
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start: %w", err)
	}

	exited := make(chan error, 1)
	go func() { exited <- cmd.Wait() }()

	select {
	case err := <-exited:
		if err != nil {
			return fmt.Errorf("died before ready: %w", err)
		}
		return fmt.Errorf("exited before ready (delay %s)", delay)

	case <-time.After(delay):
		_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGTERM)
		select {
		case <-exited:
			return nil
		case <-time.After(5 * time.Second):
			_ = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
			<-exited
			return nil
		}
	}
}
