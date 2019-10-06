package cli_test

import (
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
