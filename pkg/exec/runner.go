package exec

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"
)

type CommandLineRunner interface {
	Run(name string, arg ...string) error
	RunOutput(name string, arg ...string) ([]byte, error)
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

func WithContext(ctx context.Context) option {
	return func(r *clr) {
		r.ctx = ctx
	}
}

type clr struct {
	ctx    context.Context
	Stdout io.Writer
	Stderr io.Writer
}

func (c *clr) Run(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	if c.ctx != nil {
		cmd = exec.CommandContext(c.ctx, name, arg...)
	}

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("command '%s' failed '%v'", name, err)
	}

	outStr, _ := ioutil.ReadAll(stdout)
	fmt.Fprintf(c.Stdout, "%s", outStr)

	errStr, _ := ioutil.ReadAll(stderr)
	fmt.Fprintf(c.Stderr, "%s", errStr)

	return cmd.Wait()
}

func (c *clr) RunOutput(name string, arg ...string) ([]byte, error) {
	cmd := exec.Command(name, arg...)
	if c.ctx != nil {
		cmd = exec.CommandContext(c.ctx, name, arg...)
	}

	return cmd.CombinedOutput()
}
