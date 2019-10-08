package cluster_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/malston/k8s-mgmt/pkg/cli"
	"github.com/malston/k8s-mgmt/pkg/k8s"
	"github.com/malston/k8s-mgmt/pkg/kmgmt"
)

func TestCreateClusters_ErrorsWithArgs(t *testing.T) {
	output := &bytes.Buffer{}
	conf := cli.NewConfig("../config/testdata")
	conf.Client = k8s.NewClient("../k8s/testdata/.kube/config")
	root := kmgmt.CreateRootCommand(conf)
	root.SetOutput(output)
	root.SetArgs([]string{"create-clusters", "anything"})

	err := root.Execute()
	if err == nil {
		t.Fatalf("execute should error with args")
	}

	contents := output.String()
	if expected, actual := "Error: accepts 0 arg(s), received 1", contents; !strings.Contains(actual, expected) {
		t.Fatal("expected error to contain message: Error: accepts 0 arg(s), received 1")
	}
}

func TestCreateClusters(t *testing.T) {
	output := &bytes.Buffer{}
	conf := cli.NewConfig("../config/testdata")
	conf.Client = k8s.NewClient("../k8s/testdata/.kube/config")
	root := kmgmt.CreateRootCommand(conf)
	root.SetOutput(output)
	root.SetArgs([]string{"create-clusters"})

	err := root.Execute()
	if err != nil {
		t.Fatalf("error should not have occurred")
	}

}
