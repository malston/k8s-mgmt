package exec_test

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/malston/k8s-mgmt/pkg/exec"
)

func TestCommandLineRunner(t *testing.T) {
	tests := []struct {
		name           string
		args           string
		timeout        time.Duration
		expectedOutput string
		expectError    bool
	}{
		{
			name:           "pks",
			args:           "--version",
			expectedOutput: "",
		},
		{
			name:        "blah",
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var stdOut, stdErr bytes.Buffer
			clr := exec.NewCommandLineRunner(&stdOut, &stdErr)
			err := clr.Run(test.name, test.args)
			if err != nil && !test.expectError {
				t.Fatalf("error should not have occurred: %s", err.Error())
			}
			actualStdOut := stdOut.String()
			if !strings.Contains(actualStdOut, test.expectedOutput) {
				t.Errorf("Unexpected output: %s", actualStdOut)
			}
		})
	}
}

func TestCommandLineRunnerWithOutput(t *testing.T) {
	tests := []struct {
		name           string
		args           string
		timeout        time.Duration
		expectedOutput string
		expectError    bool
	}{
		{
			name:           "pks",
			args:           "--version",
			expectedOutput: "PKS CLI version",
		},
		{
			name:           "pks",
			args:           "--version",
			expectedOutput: "",
		},
		{
			name:        "blah",
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var stdOut, stdErr bytes.Buffer
			clr := exec.NewCommandLineRunner(&stdOut, &stdErr)
			output, err := clr.RunOutput(test.name, test.args)
			if err != nil && !test.expectError {
				t.Fatalf("error should not have occurred: %s", err.Error())
			}
			if !strings.Contains(string(output), test.expectedOutput) {
				t.Errorf("Unexpected output: %s", output)
			}
		})
	}
}

func TestCommandLineRunnerWithTimeoutContext(t *testing.T) {
	var stdOut, stdErr bytes.Buffer
	ctx, cancel := context.WithTimeout(context.Background(), 0)
	defer cancel()
	clr := exec.NewCommandLineRunner(&stdOut, &stdErr, exec.WithContext(ctx))
	err := clr.Run("pks", "--version")
	if err == nil {
		t.Fatalf("error should have occurred")
	}
	actualStdOut := stdOut.String()
	if actualStdOut != "" {
		t.Errorf("Unexpected output: %s", actualStdOut)
	}
}
