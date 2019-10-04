package namespace_test

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/malston/k8s-mgmt/pkg/k8s"
	"github.com/malston/k8s-mgmt/pkg/kmgmt"
	fakes "github.com/malston/k8s-mgmt/pkg/testing"
	"github.com/onsi/gomega/gbytes"
)

func TestCreateNamespacesErrorsWithoutArgs(t *testing.T) {
	buffer := gbytes.NewBuffer()
	// k := fakes.NewClient()
	kubeConfigFile := "../k8s/testdata/.kube/config"
	c := k8s.NewClient(kubeConfigFile)
	root := kmgmt.CreateRootCommand(c, "../config/testdata")
	root.SetOutput(buffer)
	root.SetArgs([]string{"create-namespaces"})

	err := root.Execute()
	if err == nil {
		t.Fatalf("execute should error without args")
	}

	contents := string(buffer.Contents())
	if !strings.Contains(contents, "Error: accepts 1 arg(s), received 0\n") {
		t.Fatal("expected error message: Error: accepts 1 arg(s), received 0")
	}
}

func TestCreateNamespacesInvalidCluster(t *testing.T) {
	buffer := gbytes.NewBuffer()
	// k := fakes.NewClient()
	kubeConfigFile := "../k8s/testdata/.kube/config"
	c := k8s.NewClient(kubeConfigFile)
	root := kmgmt.CreateRootCommand(c, "../config/testdata")
	root.SetOutput(buffer)
	root.SetArgs([]string{"create-namespaces", "cluster-noexiste"})

	err := root.Execute()
	if err == nil {
		t.Fatalf("execute should error without args")
	}

	contents := string(buffer.Contents())
	if !strings.Contains(contents, "context 'cluster-noexiste' not found\n") {
		t.Fatal("expected error message: context 'cluster-noexiste' not found")
	}
}

func TestCreateNamespacesValidCluster(t *testing.T) {
	c := fakes.NewClient()
	root := kmgmt.CreateRootCommand(c, "../config/testdata")

	output := &bytes.Buffer{}
	root.SetOutput(output)
	c.Stdout = output
	c.Stderr = output
	root.SetArgs([]string{"create-namespaces", "cluster-1"})

	err := root.Execute()
	if err != nil {
		fmt.Printf("error: %s", err.Error())
		t.Fatalf("execute should not error: %s", err.Error())
	}

	outputString := output.String()
	fmt.Println(outputString)
	if !strings.Contains(outputString, "Namespace namespace-1 created\nNamespace namespace-2 created\n") {
		t.Fatal("expected namespaces to be created")
	}
}
