package resourcequota_test

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

func TestCreateResourcequota_ErrorsWithoutArgs(t *testing.T) {
	output := &bytes.Buffer{}
	// conf := cli.NewConfig("../config/testdata")
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
			*config.Resourcequota{},
			fmt.Errorf("cluster doesn't exist")),
	}

	conf.Client = fakes.NewKubeClient()
	root := kmgmt.CreateRootCommand(conf)
	root.SetOutput(output)
	root.SetArgs([]string{"create-resourcequota"})

	err := root.Execute()
	if err == nil {
		t.Fatalf("execute should error without args")
	}

	contents := output.String()
	if !strings.Contains(contents, "Error: accepts 4 arg(s), received 0\n") {
		t.Fatal("expected error message: Error: accepts 4 arg(s), received 0")
	}
}

func TestCreateResourcequota_ValidCluster(t *testing.T) {
	conf := cli.NewConfig("../config/testdata")
	conf.Client = fakes.NewKubeClient()
	root := kmgmt.CreateRootCommand(conf)

	output := &bytes.Buffer{}
	root.SetOutput(output)
	conf.Stdout = output
	conf.Stderr = output
	root.SetArgs([]string{"create-resourcequota", "namespace-1", "kind"})

	err := root.Execute()
	if err != nil {
		t.Fatalf("execute should not error: %s", err.Error())
	}

	contents := output.String()
	if !strings.Contains(contents, "resourcequota/default-mem-cpu-quotas created\n") {
		t.Fatal("expected Resourcequota to be created")
	}
}

func TestCreateResourcequota_InvalidCluster(t *testing.T) {
	output := &bytes.Buffer{}
	// k := fakes.NewKubeClient()
	conf := cli.NewConfig("../config/testdata")
	conf.Client = k8s.NewClient("../k8s/testdata/.kube/config")
	root := kmgmt.CreateRootCommand(conf)
	root.SetOutput(output)
	root.SetArgs([]string{"create-resourcequota", "cluster-noexiste", "namespace-1"})

	err := root.Execute()
	if err == nil {
		t.Fatalf("execute should error without args")
	}

	contents := output.String()
	if !strings.Contains(contents, "context 'cluster-noexiste' not found\n") {
		t.Fatal("expected error message: context 'cluster-noexiste' not found")
	}
}

func TestCreateResourcequota_NamespacesNotFound(t *testing.T) {
	conf := cli.NewConfig("../config/testdata")
	conf.Client = fakes.NewKubeClient()
	root := kmgmt.CreateRootCommand(conf)

	output := &bytes.Buffer{}
	root.SetOutput(output)
	conf.Stdout = output
	conf.Stderr = output
	root.SetArgs([]string{"create-resourcequota", "namespace-3", "kind"})

	err := root.Execute()
	if err == nil {
		t.Fatalf("execute should return error namespace not found")
	}

	contents := output.String()
	if !strings.Contains(contents, "no namespaces found for cluster kind\n") {
		t.Fatal("expected namespace to exist to set quotas")
	}
}

func TestCreateResourcequota_ClusterDoesNotExist(t *testing.T) {
	conf := &cli.Config{
		ConfigDir: "../config/testdata",
		Manager: newManager(
			[]*config.Cluster{
				{
					Name: "cluster-1",
				},
			},
			nil,
			nil,
			fmt.Errorf("cluster doesn't exist")),
	}
	conf.Client = fakes.NewKubeClient()
	root := kmgmt.CreateRootCommand(conf)

	output := &bytes.Buffer{}
	root.SetOutput(output)
	conf.Stdout = output
	conf.Stderr = output
	root.SetArgs([]string{"create-resourcequota", "cluster-1"})

	err := root.Execute()
	if err == nil {
		t.Fatalf("execute should return error")
	}

	contents := output.String()
	if !strings.Contains(contents, "Error: cluster doesn't exist\n") {
		t.Fatal("expected cluster to exist")
	}
}

func TestCreateResourcequota_InvalidQuotaConfig(t *testing.T) {
	c := fakes.NewKubeClient()
	c.FakeKubeClientset.Fake.PrependReactor("create", "quotas", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, &v1.Namespace{}, errors.New("error creating quotas")
	})
	conf := &cli.Config{
		ConfigDir: "../config/testdata",
		Manager: newManager(
			[]*config.Cluster{
				{
					Name: "cluster-1",
				},
			},
			[]*config.Namespace{},
			[]*config.Resourcequota{},
			nil),
	}
	conf.Client = c
	root := kmgmt.CreateRootCommand(conf)

	output := &bytes.Buffer{}
	root.SetOutput(output)
	conf.Stdout = output
	conf.Stderr = output
	root.SetArgs([]string{"create-resourcequota", "kind"})

	err := root.Execute()
	if err == nil {
		t.Fatalf("execute should return error")
	}

	contents := output.String()
	if !strings.Contains(contents, "Error: error creating quota") {
		t.Fatal("expected quota to be created")
	}
}

type stubManager struct {
	clusters      []*config.Cluster
	namespaces    []*config.Namespace
	resourcequota *config.Resourcequota
	err           error
}

func newManager(clusters []*config.Cluster, namespaces []*config.Namespace, resourcequota *config.Resourcequota, err error) *stubManager {
	return &stubManager{clusters, namespaces, resourcequota, err}
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

func (m *stubManager) GetResourcequota(namespace string) (*config.Resourcequota, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.resourcequota, nil
}
