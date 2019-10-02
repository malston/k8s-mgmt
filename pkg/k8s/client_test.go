package k8s_test

import (
	"testing"

	"github.com/malston/k8s-mgmt/pkg/k8s"
)

func TestNewClient(t *testing.T) {
	client := k8s.NewClient("testdata/.kube/config")

	if client.Core() == nil {
		t.Errorf("Expected Core client to not be nil")
	}
}
