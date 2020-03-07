package namespace_test

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/malston/k8s-mgmt/pkg/cli"
	"github.com/malston/k8s-mgmt/pkg/config"
	"github.com/malston/k8s-mgmt/pkg/k8s"
	"github.com/malston/k8s-mgmt/pkg/kmgmt"
	fakes "github.com/malston/k8s-mgmt/pkg/testing"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	k8stesting "k8s.io/client-go/testing"
)

func TestCreateNamespaces_ErrorsWithoutArgs(t *testing.T) {
	output := &bytes.Buffer{}
	// k := fakes.NewKubeClient()
	conf := cli.NewConfig("../config/testdata")
	conf.Client = k8s.NewClient("../k8s/testdata/.kube/config")
	root := kmgmt.CreateRootCommand(conf)
	root.SetOutput(output)
	root.SetArgs([]string{"create-namespaces"})

	err := root.Execute()
	if err == nil {
		t.Fatalf("execute should error without args")
	}

	contents := output.String()
	if !strings.Contains(contents, "Error: accepts 1 arg(s), received 0\n") {
		t.Fatal("expected error message: Error: accepts 1 arg(s), received 0")
	}
}

func TestCreateNamespaces_ValidCluster(t *testing.T) {
	conf := cli.NewConfig("../config/testdata")
	conf.Client = fakes.NewKubeClient()
	root := kmgmt.CreateRootCommand(conf)

	output := &bytes.Buffer{}
	root.SetOutput(output)
	conf.Stdout = output
	conf.Stderr = output
	root.SetArgs([]string{"create-namespaces", "cluster-1"})

	err := root.Execute()
	if err != nil {
		t.Fatalf("execute should not error: %s", err.Error())
	}

	contents := output.String()
	if !strings.Contains(contents, "Namespace namespace-1 created\nNamespace namespace-2 created\n") {
		t.Fatal("expected namespaces to be created")
	}
}

func TestCreateNamespaces_InvalidCluster(t *testing.T) {
	output := &bytes.Buffer{}
	// k := fakes.NewKubeClient()
	conf := cli.NewConfig("../config/testdata")
	conf.Client = k8s.NewClient("../k8s/testdata/.kube/config")
	root := kmgmt.CreateRootCommand(conf)
	root.SetOutput(output)
	root.SetArgs([]string{"create-namespaces", "cluster-noexiste"})

	err := root.Execute()
	if err == nil {
		t.Fatalf("execute should error without args")
	}

	contents := output.String()
	if !strings.Contains(contents, "context 'cluster-noexiste' not found\n") {
		t.Fatal("expected error message: context 'cluster-noexiste' not found")
	}
}

func TestCreateNamespaces_NamespacesNotFound(t *testing.T) {
	conf := cli.NewConfig("../config/testdata")
	conf.Client = fakes.NewKubeClient()
	root := kmgmt.CreateRootCommand(conf)

	output := &bytes.Buffer{}
	root.SetOutput(output)
	conf.Stdout = output
	conf.Stderr = output
	root.SetArgs([]string{"create-namespaces", "cluster-3"})

	err := root.Execute()
	if err == nil {
		t.Fatalf("execute should return error namespace not found")
	}

	contents := output.String()
	if !strings.Contains(contents, "no namespaces found for cluster cluster-3\n") {
		t.Fatal("expected namespaces to be created")
	}
}

func TestCreateNamespaces_ClusterDoesNotExist(t *testing.T) {
	conf := &cli.Config{
		ConfigDir: "../config/testdata",
		Manager: newManager(
			[]*config.Cluster{
				{
					Name: "cluster-1",
				},
			}, nil, fmt.Errorf("cluster doesn't exist")),
	}
	conf.Client = fakes.NewKubeClient()
	root := kmgmt.CreateRootCommand(conf)

	output := &bytes.Buffer{}
	root.SetOutput(output)
	conf.Stdout = output
	conf.Stderr = output
	root.SetArgs([]string{"create-namespaces", "cluster-1"})

	err := root.Execute()
	if err == nil {
		t.Fatalf("execute should return error")
	}

	contents := output.String()
	if !strings.Contains(contents, "Error: cluster doesn't exist\n") {
		t.Fatal("expected namespaces to be created")
	}
}

func TestCreateNamespaces_InvalidNamespace(t *testing.T) {
	c := fakes.NewKubeClient()
	c.FakeKubeClientset.Fake.PrependReactor("create", "namespaces", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, &v1.Namespace{}, errors.New("error creating namespace")
	})
	conf := &cli.Config{
		ConfigDir: "../config/testdata",
		Manager: newManager(
			[]*config.Cluster{
				{
					Name: "cluster-1",
				},
			},
			[]*config.Namespace{
				{
					Name: "namespace-1",
				},
			},
			nil),
	}
	conf.Client = c
	root := kmgmt.CreateRootCommand(conf)

	output := &bytes.Buffer{}
	root.SetOutput(output)
	conf.Stdout = output
	conf.Stderr = output
	root.SetArgs([]string{"create-namespaces", "cluster-1"})

	err := root.Execute()
	if err == nil {
		t.Fatalf("execute should return error")
	}

	contents := output.String()
	if !strings.Contains(contents, "Error: error creating namespace") {
		t.Fatalf("expected error message: 'error creating namespace', got %s", contents)
	}
}

type stubManager struct {
	clusters   []*config.Cluster
	namespaces []*config.Namespace
	err        error
}

func newManager(clusters []*config.Cluster, namespaces []*config.Namespace, err error) *stubManager {
	return &stubManager{clusters, namespaces, err}
}

func (m *stubManager) GetClusters() ([]*config.Cluster, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.clusters, nil
}

func (m *stubManager) GetNamespaces(cluster string) ([]*config.Namespace, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.namespaces, nil
}

func (m *stubManager) GetResourceQuota(cluster string) (*v1.ResourceQuota, error) {
	if m.err != nil {
		return nil, m.err
	}
	return nil, nil
}
