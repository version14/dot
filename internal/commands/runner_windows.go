//go:build windows

package commands

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"
)

func (r *Runner) runBackground(ctx context.Context, c PlannedCommand, wd string, stdout, stderr fileOrBuf) error {
	delay := c.ReadyDelay
	if delay <= 0 {
		delay = defaultBackgroundReadyDelay
	}

	cmd := exec.Command("cmd.exe", "/C", c.Cmd)
	cmd.Dir = wd
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	cmd.Env = os.Environ()

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
		if err := cmd.Process.Kill(); err != nil {
			r.Logger.Errorf("failed to kill process: %v", err)
		}
		<-exited
		return nil
	}
}
