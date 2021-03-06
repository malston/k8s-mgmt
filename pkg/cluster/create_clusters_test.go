package cluster_test

import (
	"bytes"
	"errors"
	"strings"
	"sync"
	"testing"

	fakes "github.com/malston/k8s-mgmt/pkg/testing"
	v1 "k8s.io/api/core/v1"

	"github.com/malston/k8s-mgmt/pkg/cli"
	"github.com/malston/k8s-mgmt/pkg/config"
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
	stubClient := &stubPKSClient{}
	stubClient.wg.Add(5)
	conf := &cli.Config{
		Manager: newStubManager(
			[]*config.Cluster{
				{
					Name: "cluster-1",
				},
				{
					Name: "cluster-2",
				},
			}, nil, nil),
		Client:    fakes.NewKubeClient(),
		PKSClient: stubClient,
	}
	root := kmgmt.CreateRootCommand(conf)

	output := &bytes.Buffer{}
	root.SetOutput(output)
	conf.Stdout = output
	conf.Stderr = output
	root.SetArgs([]string{"create-clusters"})

	err := root.Execute()
	if err != nil {
		t.Fatalf("execute should not error: %s", err.Error())
	}

	if expected, actual := 5, stubClient.callCount(); expected != actual {
		t.Fatalf("expected %d, actual %d", expected, actual)
	}

	contents := output.String()
	if expected, actual := "Cluster cluster-1 created", contents; !strings.Contains(actual, expected) {
		t.Fatalf("expected %s, actual %s", expected, actual)
	}
	if expected, actual := "Cluster cluster-2 created\n", contents; !strings.Contains(actual, expected) {
		t.Fatalf("expected %s, actual %s", expected, actual)
	}
}

func TestCreateCluster_FailsWithError(t *testing.T) {
	stubClient := &stubPKSClient{
		stubCluster: []*stubCluster{
			{
				name: "cluster-1",
				err:  errors.New("failure"),
			},
			{
				name: "cluster-2",
			},
		}}
	stubClient.wg.Add(4)
	conf := &cli.Config{
		ConfigDir: "../config/testdata",
		Manager: newStubManager(
			[]*config.Cluster{
				{
					Name: "cluster-1",
				},
				{
					Name: "cluster-2",
				},
			}, nil, nil),
		Client:    fakes.NewKubeClient(),
		PKSClient: stubClient,
	}
	root := kmgmt.CreateRootCommand(conf)

	output := &bytes.Buffer{}
	root.SetOutput(output)
	conf.Stdout = output
	conf.Stderr = output
	root.SetArgs([]string{"create-clusters"})

	err := root.Execute()
	if err == nil {
		t.Fatalf("error should have occurred")
	}

	if expected, actual := 4, stubClient.callCount(); expected != actual {
		t.Fatalf("expected %d, actual %d", expected, actual)
	}

	if expected, actual := "Error: failed to create cluster cluster-1", output.String(); !strings.Contains(actual, expected) {
		t.Fatalf("expected %s, actual %s", expected, actual)
	}
}

type stubCluster struct {
	name string
	err  error
}

type stubPKSClient struct {
	called      int
	iterations  int
	stubCluster []*stubCluster
	wg          sync.WaitGroup
}

func (s *stubPKSClient) CreateCluster(cluster *config.Cluster) error {
	defer s.wg.Done()
	s.called++
	for _, c := range s.stubCluster {
		if cluster.Name == c.name {
			if c.err != nil {
				return c.err
			}
		}
	}
	return nil
}

func (s *stubPKSClient) ShowCluster(name string) (*config.Cluster, error) {
	defer s.wg.Done()
	defer func(i int) { s.iterations += i }(1)
	s.called++
	if s.iterations == 0 {
		return &config.Cluster{
			Name: name,
		}, nil
	}
	s.iterations = 0
	return &config.Cluster{
		Name:      name,
		IPAddress: "127.0.0.1",
	}, nil
}

func (s *stubPKSClient) callCount() int {
	s.wg.Wait()
	return s.called
}

type stubManager struct {
	clusters   []*config.Cluster
	namespaces []*config.Namespace
	err        error
}

func newStubManager(clusters []*config.Cluster, namespaces []*config.Namespace, err error) *stubManager {
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

func (m *stubManager) GetResourceQuota(cluster string, namespace string) (*v1.ResourceQuota, error) {
	if m.err != nil {
		return nil, m.err
	}
	return nil, nil
}
