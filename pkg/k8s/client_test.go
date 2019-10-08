package k8s_test

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/malston/k8s-mgmt/pkg/cli"
	"github.com/malston/k8s-mgmt/pkg/k8s"
)

func TestNewClient(t *testing.T) {
	cli.NewConfig("../config/testdata")
	client := k8s.NewClient("testdata/.kube/config")

	if client.Core() == nil {
		t.Errorf("Expected Core client to not be nil")
	}

	if client.CurrentContext() == "" {
		t.Errorf("Expected current context to not be empty")
	}

	if client.SetContext("my-context") != nil {
		t.Errorf("Expected SetContext not to return error")
	}
}

func TestContextSettings(t *testing.T) {
	client := k8s.NewClient("testdata/.kube/config")
	if diff := cmp.Diff(client.CurrentContext(), "my-context"); diff != "" {
		t.Errorf("Unexpected context (-expected, +actual): %s", diff)
	}
	client.SetContext("minikube")
	if diff := cmp.Diff(client.CurrentContext(), "minikube"); diff != "" {
		t.Errorf("Unexpected context (-expected, +actual): %s", diff)
	}
	err := client.SetContext("no-existe")
	if err == nil {
		t.Fatal("should return error when given a context name that doesn't exist")
	}
	if err.Error() != fmt.Sprintf("context '%s' not found", "no-existe") {
		t.Fatal("should return correct error message")
	}
}
