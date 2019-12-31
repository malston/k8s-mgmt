package resourcequota_test

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/malston/k8s-mgmt/pkg/cli"
	"github.com/malston/k8s-mgmt/pkg/config"
	"github.com/malston/k8s-mgmt/pkg/kmgmt"
	fakes "github.com/malston/k8s-mgmt/pkg/testing"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	k8stesting "k8s.io/client-go/testing"
)

func TestCreateResourceQuota_ErrorsWithoutArgs(t *testing.T) {
	output := &bytes.Buffer{}
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
			&config.ResourceQuota{},
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
	if !strings.Contains(contents, "Error: accepts 1 arg(s), received 0\n") {
		t.Fatalf("expected error message: Error: accepts 1 arg(s), received 0; got %s", contents)
	}
}

func TestCreateResourceQuota(t *testing.T) {
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
			&config.ResourceQuota{
				Name:           "default-mem-cpu-quotas",
				RequestsCPU:    "1",
				RequestsMemory: "1Gi",
				LimitsCPU:      "2",
				LimitsMemory:   "2Gi",
			},
			nil),
	}
	conf.Client = fakes.NewKubeClient()
	root := kmgmt.CreateRootCommand(conf)

	output := &bytes.Buffer{}
	root.SetOutput(output)
	conf.Stdout = output
	conf.Stderr = output
	root.SetArgs([]string{"create-resourcequota", "namespace-1"})

	err := root.Execute()
	if err != nil {
		t.Fatalf("execute should not error: %s", err.Error())
	}

	contents := output.String()
	if !strings.Contains(contents, "resourcequota/default-mem-cpu-quotas created\n") {
		t.Fatalf("expected ResourceQuota to be created, got %s", contents)
	}
}

func TestCreateResourceQuota_NamespacesNotFound(t *testing.T) {
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
			fmt.Errorf("namespace not found")),
	}
	conf.Client = fakes.NewKubeClient()
	root := kmgmt.CreateRootCommand(conf)

	output := &bytes.Buffer{}
	root.SetOutput(output)
	conf.Stdout = output
	conf.Stderr = output
	root.SetArgs([]string{"create-resourcequota", "namespace-noexiste"})

	err := root.Execute()
	if err == nil {
		t.Fatalf("execute should return error namespace not found")
	}

	contents := output.String()
	if !strings.Contains(contents, "namespace not found\n") {
		t.Fatalf("expected error message: 'namespace not found', got %s", contents)
	}
}

func TestCreateResourceQuota_InvalidResourceQuotaConfig(t *testing.T) {
	c := fakes.NewKubeClient()
	c.FakeKubeClientset.Fake.PrependReactor("create", "resourcequotas", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, &v1.ResourceQuota{}, errors.New("error creating resource quota")
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
			&config.ResourceQuota{
				Name:           "default-mem-cpu-quotas",
				RequestsCPU:    "1",
				RequestsMemory: "1Gi",
				LimitsCPU:      "2",
				LimitsMemory:   "2Gi",
			},
			nil),
	}
	conf.Client = c
	root := kmgmt.CreateRootCommand(conf)

	output := &bytes.Buffer{}
	root.SetOutput(output)
	conf.Stdout = output
	conf.Stderr = output
	root.SetArgs([]string{"create-resourcequota", "namespace-1"})

	err := root.Execute()
	if err == nil {
		t.Fatalf("execute should return error")
	}

	contents := output.String()
	if !strings.Contains(contents, "Error: error creating resource quota") {
		t.Fatalf("expected error message: 'error creating resource quota', got %s", contents)
	}
}

type stubManager struct {
	clusters      []*config.Cluster
	namespaces    []*config.Namespace
	resourcequota *config.ResourceQuota
	err           error
}

func newManager(clusters []*config.Cluster, namespaces []*config.Namespace, resourcequota *config.ResourceQuota, err error) *stubManager {
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

func (m *stubManager) GetResourceQuota(namespace string) (*config.ResourceQuota, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.resourcequota, nil
}
