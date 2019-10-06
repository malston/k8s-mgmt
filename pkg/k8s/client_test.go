package k8s_test

import (
	"testing"

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
