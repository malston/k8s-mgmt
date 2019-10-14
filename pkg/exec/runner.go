package exec

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"
)

type CommandLineRunner interface {
	Run(name string, arg ...string) error
}

type option func(*clr)

func NewCommandLineRunner(w io.Writer, opts ...option) CommandLineRunner {
	clr := &clr{
		Stdout: w,
		Stderr: w,
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
	// cmd.Stdout = m.Stdout
	// cmd.Stderr = m.Stderr

	// err := cmd.Run()
	// if err != nil {
	// 	return err
	// }
	// var stdoutBuf, stderrBuf bytes.Buffer
	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()

	var errStdout, errStderr error
	stdout := io.MultiWriter(os.Stdout, m.Stdout)
	stderr := io.MultiWriter(os.Stderr, m.Stderr)
	err := cmd.Start()
	if err != nil {
		log.Fatalf("cmd.Start() failed with '%s'\n", err)
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		_, errStdout = io.Copy(stdout, stdoutIn)
		wg.Done()
	}()

	_, errStderr = io.Copy(stderr, stderrIn)
	wg.Wait()

	err = cmd.Wait()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	if errStdout != nil || errStderr != nil {
		log.Fatal("failed to capture stdout or stderr\n")
	}
	// outStr, errStr := string(stdoutBuf.Bytes()), string(stderrBuf.Bytes())
	fmt.Printf("\nout:\n%s\nerr:\n%s\n", m.Stdout, m.Stderr)
	return nil
}
