package config_test

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/malston/k8s-mgmt/pkg/config"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	if len(clusters) != 3 {
		t.Errorf("expected 3, got %d clusters", len(clusters))
	}
	if len(clusters[0].Namespaces) != 2 {
		t.Errorf("expected 2, got %d namespaces", len(clusters[0].Namespaces))
		if clusters[0].Namespaces[0].ResourceQuota == nil {
			t.Errorf("expected resource quota not to be nil for namespace %s", clusters[0].Namespaces[0].Name)
		}
		if clusters[0].Namespaces[1].ResourceQuota == nil {
			t.Errorf("expected resource quota not to be nil for namespace %s", clusters[0].Namespaces[0].Name)
		}
	}
	if len(clusters[1].Namespaces) != 1 {
		t.Errorf("expected 1, got %d namespaces", len(clusters[1].Namespaces))
		if clusters[0].Namespaces[0].ResourceQuota == nil {
			t.Errorf("expected resource quota not to be nil for namespace %s", clusters[0].Namespaces[0].Name)
		}
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
	if namespaces[0].Name != "my-namespace" {
		t.Errorf("expected namespace-1, got %s namespaces", namespaces[0].Name)
	}
}

func TestNamespacesAreEmpty(t *testing.T) {
	m, _ := config.NewManager("./testdata")
	name := "cluster-noexiste"
	ns, _ := m.GetNamespaces(name)
	if len(ns) != 0 {
		t.Fatal("should return empty set of namespaces")
	}
}

func TestClusterDoesNotExist(t *testing.T) {
	m, err := config.NewManager("./testdata")
	if err != nil {
		t.Fatal("error should not occur")
	}
	name := "cluster-noexiste"
	ns, err := m.GetNamespaces(name)
	if len(ns) != 0 {
		t.Fatal("should return error when given a cluster name that doesn't exist")
	}
}

func TestClusterInFolderDoesNotExist(t *testing.T) {
	m, _ := config.NewManager("./testdata/noclusters")
	name := "cluster-noexiste"
	_, err := m.GetNamespaces(name)
	if err == nil {
		t.Fatal("should return error when given a cluster name that doesn't exist")
	}
	if err.Error() != fmt.Sprintf("cluster %s does not exist in config folder", name) {
		t.Fatal("should return correct error message")
	}
}

func TestResourceQuotaIsCreatedForAllNamespaces(t *testing.T) {
	m, _ := config.NewManager("./testdata")
	want := &v1.ResourceQuota{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ResourceQuota",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "default-mem-cpu-quota-for-cluster-1",
		},
		Spec: v1.ResourceQuotaSpec{
			Hard: v1.ResourceList{
				v1.ResourceRequestsCPU:    resource.MustParse("1"),
				v1.ResourceRequestsMemory: resource.MustParse("1Gi"),
				v1.ResourceLimitsCPU:      resource.MustParse("2"),
				v1.ResourceLimitsMemory:   resource.MustParse("2Gi"),
			},
		},
	}
	got, err := m.GetResourceQuota("cluster-1", "namespace-1")
	if err != nil {
		t.Errorf("error creating resource quota %s", err)
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("GetResourceQuota() mismatch (-want +got):\n%s", diff)
	}
	got, err = m.GetResourceQuota("cluster-1", "namespace-2")
	if err != nil {
		t.Errorf("error creating resource quota %s", err)
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("GetResourceQuota() mismatch (-want +got):\n%s", diff)
	}
	want = &v1.ResourceQuota{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ResourceQuota",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "default-mem-cpu-quota-for-cluster-2",
		},
		Spec: v1.ResourceQuotaSpec{
			Hard: v1.ResourceList{
				v1.ResourceRequestsCPU:    resource.MustParse("2"),
				v1.ResourceRequestsMemory: resource.MustParse("2Gi"),
				v1.ResourceLimitsCPU:      resource.MustParse("4"),
				v1.ResourceLimitsMemory:   resource.MustParse("4Gi"),
			},
		},
	}
	got, err = m.GetResourceQuota("my-cluster", "my-namespace")
	if err != nil {
		t.Errorf("error creating resource quota %s", err)
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("GetResourceQuota() mismatch (-want +got):\n%s", diff)
	}
}
