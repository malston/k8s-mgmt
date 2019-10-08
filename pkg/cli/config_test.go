package cli_test

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/malston/k8s-mgmt/pkg/cli"
	homedir "github.com/mitchellh/go-homedir"
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

func TestNewDefaultConfig(t *testing.T) {
	conf, err := cli.NewDefaultConfig()
	if err != nil {
		t.Error("Expected error not to have occurred")
	}

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

func TestDefaultConfig_HomeDir(t *testing.T) {
	home, homeisset := os.LookupEnv("K8SMGMT_HOME")
	defer func() {
		homedir.Reset()
		if homeisset {
			os.Setenv("K8SMGMT_HOME", home)
		} else {
			os.Unsetenv("K8SMGMT_HOME")
		}
	}()

	os.Setenv("K8SMGMT_HOME", "testdata")

	c, err := cli.NewDefaultConfig()
	if err != nil {
		t.Error("Expected error not to have occurred")
	}

	output := &bytes.Buffer{}
	c.Stdout = output
	c.Stderr = output

	if expected, actual := "testdata", c.ConfigDir; expected != actual {
		t.Errorf("Expected config dir %s, actually %s", expected, actual)
	}
}

func TestDefaultConfig_HomeDirNotSet(t *testing.T) {
	c, err := cli.NewDefaultConfig()
	if err != nil {
		t.Error("Expected error not to have occurred")
	}

	output := &bytes.Buffer{}
	c.Stdout = output
	c.Stderr = output

	os.Setenv("HOME", "")

	if c.ConfigDir == "" {
		t.Errorf("Expected config directory to be set to default location")
	}
	if diff := cmp.Diff("", strings.TrimSpace(output.String())); diff != "" {
		t.Errorf("Unexpected output (-expected, +actual): %s", diff)
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
