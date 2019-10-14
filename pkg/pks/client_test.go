package pks_test

import (
	"errors"
	"testing"

	"github.com/malston/k8s-mgmt/pkg/config"
	"github.com/malston/k8s-mgmt/pkg/pks"
)

func TestCreateCluster(t *testing.T) {
	clr := &stubCommandLineRunner{}
	client := pks.NewClient(clr)
	err := client.CreateCluster(&config.Cluster{Name: "some-cluster"})
	if err != nil {
		t.Fatalf("error should not have occured %v", err)
	}
	if expected, actual := 1, clr.runCalled; expected != actual {
		t.Fatalf("expected %d, actual %d", expected, actual)
	}
}

func TestCreateCluster_InvalidCluster(t *testing.T) {
	clr := &stubCommandLineRunner{err: errors.New("blah")}
	client := pks.NewClient(clr)
	err := client.CreateCluster(&config.Cluster{Name: "some-cluster"})
	if err == nil {
		t.Fatal("error should have occured")
	}
	if expected, actual := 1, clr.runCalled; expected != actual {
		t.Fatalf("expected %d, actual %d", expected, actual)
	}
}

type stubCommandLineRunner struct {
	runCalled int
	err       error
}

func (m *stubCommandLineRunner) Run(name string, arg ...string) error {
	m.runCalled++
	if m.err != nil {
		return m.err
	}
	return nil
}
