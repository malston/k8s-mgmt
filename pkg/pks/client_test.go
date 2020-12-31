package pks_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"testing"

	"github.com/malston/k8s-mgmt/pkg/config"
	"github.com/malston/k8s-mgmt/pkg/pks"
)

func TestCreateCluster(t *testing.T) {
	clr := &stubCommandLineRunner{}
	pks := pks.NewPKSClient(clr)
	err := pks.CreateCluster(&config.Cluster{Name: "some-cluster", NetworkProfile: "some-nw-profile"})
	if err != nil {
		t.Fatalf("error should not have occured %v", err)
	}
	if expected, actual := 1, clr.runCalled; expected != actual {
		t.Fatalf("expected %d, actual %d", expected, actual)
	}
}

func TestCreateCluster_InvalidCluster(t *testing.T) {
	clr := &stubCommandLineRunner{err: errors.New("blah")}
	pks := pks.NewPKSClient(clr)
	err := pks.CreateCluster(&config.Cluster{Name: "some-cluster"})
	if err == nil {
		t.Fatal("error should have occured")
	}
	if expected, actual := 1, clr.runCalled; expected != actual {
		t.Fatalf("expected %d, actual %d", expected, actual)
	}
}

func TestShowCluster(t *testing.T) {
	clr := &stubCommandLineRunner{}
	resp := new(bytes.Buffer)
	json.NewEncoder(resp).Encode(&config.Cluster{Name: "some-cluster"})
	clr.output = resp.Bytes()
	pks := pks.NewPKSClient(clr)
	cluster, err := pks.ShowCluster("some-cluster")
	if err != nil {
		t.Fatalf("error should not have occured %v", err)
	}
	if cluster == nil {
		t.Fatal("cluster should not be nil")
	}
	if expected, actual := 1, clr.runCalled; expected != actual {
		t.Fatalf("expected %d, actual %d", expected, actual)
	}
}

func TestShowCluster_InvalidCluster(t *testing.T) {
	clr := &stubCommandLineRunner{err: errors.New("blah")}
	pks := pks.NewPKSClient(clr)
	_, err := pks.ShowCluster("some-cluster")
	if err == nil {
		t.Fatal("error should have occured")
	}
	if expected, actual := 1, clr.runCalled; expected != actual {
		t.Fatalf("expected %d, actual %d", expected, actual)
	}
}

func TestShowCluster_InvalidCommandOutput(t *testing.T) {
	clr := &stubCommandLineRunner{}
	clr.output = nil
	pks := pks.NewPKSClient(clr)
	_, err := pks.ShowCluster("some-cluster")
	if err == nil {
		t.Fatal("error should have occured")
	}
	if expected, actual := 1, clr.runCalled; expected != actual {
		t.Fatalf("expected %d, actual %d", expected, actual)
	}
}

type stubCommandLineRunner struct {
	runCalled int
	output    []byte
	err       error
}

func (m *stubCommandLineRunner) Run(name string, arg ...string) error {
	m.runCalled++
	if m.err != nil {
		return m.err
	}
	return nil
}

func (m *stubCommandLineRunner) RunOutput(name string, arg ...string) ([]byte, error) {
	m.runCalled++
	if m.err != nil {
		return m.output, m.err
	}
	return m.output, nil
}
