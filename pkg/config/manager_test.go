package config_test

import (
	"fmt"
	"testing"

	"github.com/malston/k8s-mgmt/pkg/config"
)

func TestInvalidValidConfigDirectory(t *testing.T) {
	_, err := config.NewManager("")
	if err == nil {
		t.Fatal("should return error when given an invalid directory")
	}
}

func TestYamlMarshalledIntoClusterConfig(t *testing.T) {
	m, err := config.NewManager("./testdata")
	if err != nil {
		t.Errorf("error parsing config %s", err)
	}
	clusters, err := m.GetClusters()
	if len(clusters) == 0 {
		t.Errorf("error populating clusters: %s", err)
	}
	if len(clusters) != 2 {
		t.Errorf("expected 2, got %d clusters", len(clusters))
	}
	if len(clusters[0].Namespaces) != 2 {
		t.Errorf("expected 2, got %d namespaces", len(clusters[0].Namespaces))
	}
	if len(clusters[1].Namespaces) != 1 {
		t.Errorf("expected 1, got %d namespaces", len(clusters[1].Namespaces))
	}
}

func TestClusterIsOverriddenInYaml(t *testing.T) {
	m, _ := config.NewManager("./testdata")
	clusters, _ := m.GetClusters()
	if clusters[0].Name != "cluster-1" {
		t.Errorf("expected cluster-1, got %s clusters", clusters[0].Name)
	}
	if clusters[1].Name != "my-cluster" {
		t.Errorf("expected my-cluster, got %s clusters", clusters[1].Name)
	}
}

func TestYamlMarshalledIntoNamespaceConfig(t *testing.T) {
	m, _ := config.NewManager("./testdata")
	namespaces, err := m.GetNamespaces("cluster-1")
	if err != nil {
		t.Errorf("error parsing config %s", err)
	}
	if namespaces[0].Name != "namespace-1" {
		t.Errorf("expected namespace-1, got %s namespaces", namespaces[0].Name)
	}
	if namespaces[1].Name != "namespace-2" {
		t.Errorf("expected namespace-2, got %s namespaces", namespaces[1].Name)
	}
}

func TestNamespaceIsOverriddenInYaml(t *testing.T) {
	m, _ := config.NewManager("./testdata")
	namespaces, err := m.GetNamespaces("my-cluster")
	if err != nil {
		t.Errorf("error parsing config %s", err)
	}
	if namespaces[0].Name != "overridden-namespace-1" {
		t.Errorf("expected namespace-1, got %s namespaces", namespaces[0].Name)
	}
}

func TestClusterDoesNotExisat(t *testing.T) {
	m, _ := config.NewManager("./testdata")
	name := "cluster-noexiste"
	_, err := m.GetNamespaces(name)
	if err == nil {
		t.Fatal("should return error when given a cluster name that doesn't exist")
	}
	if err.Error() != fmt.Sprintf("cluster %s does not exist in config folder", name) {
		t.Fatal("should return correct error message")
	}
}
