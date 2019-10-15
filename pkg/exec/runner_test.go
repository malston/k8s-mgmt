package exec_test

import (
	"bytes"
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
			timeout:        100 * time.Millisecond,
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
			clr := exec.NewCommandLineRunner(&stdOut, &stdErr, exec.WithTimeout(test.timeout))
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
