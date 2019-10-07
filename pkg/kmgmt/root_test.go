package kmgmt_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/malston/k8s-mgmt/pkg/cli"
	"github.com/malston/k8s-mgmt/pkg/k8s"
	"github.com/malston/k8s-mgmt/pkg/kmgmt"
)

func TestRootCommandWithHelpFlag(t *testing.T) {
	kubeConfigFile := "."
	k := k8s.NewClient(kubeConfigFile)
	c := cli.NewConfig("../config/testdata")

	root := kmgmt.CreateRootCommand(k, c)

	output := &bytes.Buffer{}
	root.SetOutput(output)
	root.SetArgs([]string{"--help"})

	err := root.Execute()

	if err != nil {
		t.Fatalf("execute should not error, %s", err.Error())
	}

	contents := output.String()
	if !strings.Contains(contents, "Help message for toggle\n") {
		t.Fatal("expected help message to be given")
	}
}
