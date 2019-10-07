package cli_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/malston/k8s-mgmt/pkg/cli"
)

func TestNewConfig(t *testing.T) {
	conf := cli.NewConfig("../config/testdata")

	if expected, actual := os.Stdin, conf.Stdin; expected != actual {
		t.Errorf("Expected stdin to be %v, actually %v", expected, actual)
	}
	if expected, actual := os.Stdout, conf.Stdout; expected != actual {
		t.Errorf("Expected stdout to be %v, actually %v", expected, actual)
	}
	if expected, actual := os.Stderr, conf.Stderr; expected != actual {
		t.Errorf("Expected stderr to be %v, actually %v", expected, actual)
	}
}

func TestConfig_Print(t *testing.T) {
	config := cli.NewConfig("../config/testdata")

	tests := []struct {
		name    string
		format  string
		args    []interface{}
		printer func(format string, a ...interface{}) (n int, err error)
		stdout  string
		stderr  string
	}{{
		name:    "Printf",
		format:  "%s",
		args:    []interface{}{"hello"},
		printer: config.Printf,
		stdout:  "hello",
	}, {
		name:    "Eprintf",
		format:  "%s",
		args:    []interface{}{"hello"},
		printer: config.Eprintf,
		stderr:  "hello",
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			stdout := &bytes.Buffer{}
			stderr := &bytes.Buffer{}
			config.Stdout = stdout
			config.Stderr = stderr

			_, err := test.printer(test.format, test.args...)

			if err != nil {
				t.Errorf("Expected no error, actually %q", err)
			}
			if expected, actual := test.stdout, stdout.String(); expected != actual {
				t.Errorf("Expected stdout to be %q, actually %q", expected, actual)
			}
			if expected, actual := test.stderr, stderr.String(); expected != actual {
				t.Errorf("Expected stderr to be %q, actually %q", expected, actual)
			}
		})
	}
}
