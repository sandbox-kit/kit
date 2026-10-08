package commands

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
)

type step struct {
	dir  string
	name string
	args []string
}

type runner interface {
	run(context.Context, step, io.Writer, io.Writer) error
}

type execRunner struct{}

func (execRunner) run(ctx context.Context, task step, stdout, stderr io.Writer) error {
	cmd := exec.CommandContext(ctx, task.name, task.args...)
	cmd.Dir = task.dir
	cmd.Env = append(os.Environ(), "GOWORK=off")
	cmd.Stdout, cmd.Stderr = stdout, stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s in %s: %w", task.name, task.dir, err)
	}
	return nil
}
