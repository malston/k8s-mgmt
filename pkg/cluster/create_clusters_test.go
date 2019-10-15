package cluster_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	fakes "github.com/malston/k8s-mgmt/pkg/testing"

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
	conf := &cli.Config{
		ConfigDir: "../config/testdata",
		Manager: newSpyManager(
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

	if expected, actual := 2, stubClient.called; expected != actual {
		t.Fatalf("expected %d, actual %d", expected, actual)
	}

	contents := output.String()
	if expected, actual := "Cluster cluster-1 created\nCluster cluster-2 created\n", contents; !strings.Contains(actual, expected) {
		t.Fatalf("expected %s, actual %s", expected, actual)
	}
}

func TestCreateCluster_FailsWithError(t *testing.T) {
	stubClient := &stubPKSClient{err: errors.New("failure")}
	conf := &cli.Config{
		ConfigDir: "../config/testdata",
		Manager: newSpyManager(
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

	if expected, actual := 2, stubClient.called; expected != actual {
		t.Fatalf("expected %d, actual %d", expected, actual)
	}

	if expected, actual := "Error: failed to create cluster cluster-1\nfailed to create cluster cluster-2", output.String(); !strings.Contains(actual, expected) {
		t.Fatalf("expected %s, actual %s", expected, actual)
	}
}

type stubPKSClient struct {
	called int
	err    error
}

func (m *stubPKSClient) CreateCluster(cluster *config.Cluster) error {
	m.called++
	if m.err != nil {
		return m.err
	}
	return nil
}

type spyManager struct {
	clusters   []*config.Cluster
	namespaces []*config.Namespace
	err        error
}

func newSpyManager(clusters []*config.Cluster, namespaces []*config.Namespace, err error) *spyManager {
	return &spyManager{clusters, namespaces, err}
}

func (m *spyManager) GetClusters() ([]*config.Cluster, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.clusters, nil
}

func (m *spyManager) GetNamespaces(cluster string) ([]*config.Namespace, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.namespaces, nil
}
