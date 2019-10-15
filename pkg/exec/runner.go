package exec

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"
	"time"
)

type CommandLineRunner interface {
	Run(name string, arg ...string) error
}

type option func(*clr)

func NewCommandLineRunner(stdOut io.Writer, stdErr io.Writer, opts ...option) CommandLineRunner {
	clr := &clr{
		Stdout: stdOut,
		Stderr: stdErr,
	}

	for _, o := range opts {
		o(clr)
	}

	return clr
}

func WithTimeout(duration time.Duration) option {
	return func(r *clr) {
		r.timeout = duration
	}
}

type clr struct {
	timeout time.Duration
	Stdout  io.Writer
	Stderr  io.Writer
}

func (m *clr) Run(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	if m.timeout != 0 {
		ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
		defer cancel()
		cmd = exec.CommandContext(ctx, name, arg...)
	}

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("command '%s' failed '%v'", name, err)
	}

	outStr, _ := ioutil.ReadAll(stdout)
	fmt.Fprintf(m.Stdout, "%s", outStr)

	errStr, _ := ioutil.ReadAll(stderr)
	fmt.Fprintf(m.Stderr, "%s", errStr)

	return cmd.Wait()
}
