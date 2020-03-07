package config

import v1 "k8s.io/api/core/v1"

// The Config defines a global kubernetes cluster config
type Config struct {
	Clusters []*Cluster `yaml:"clusters"`
	// Quota         []Quota
	// Limit         []Limit
	// NetworkPolicy []NetworkPolicy
}

// The Cluster struct defins a kubernetes cluster
type Cluster struct {
	Name             string       `yaml:"name"`
	IPAddress        string       `yaml:"kubernetes_master_ips"`
	Namespaces       []*Namespace `yaml:"namespaces"`
	Plan             string       `yaml:"plan"`
	ExternalHostname string       `yaml:"external-hostname"`
	NumNodes         string       `yaml:"num-nodes"`
	NetworkProfile   string       `yaml:"network-profile"`
	// ClusterRole
	// ClusterRoleBinding
}

// The Namespace defins namespaces for kubernetes
type Namespace struct {
	Name string `yaml:"name"`
	// Role
	// RoleBinding
	ResourceQuota *v1.ResourceQuota

	// Limit         []Limit
	// NetworkPolicy []NetworkPolicy
}

// type ClusterInfo struct {
// }

// type Limit struct {
// }

// type NetworkPolicy struct {
// }

// type ClusterRole struct {
// }

// type ClusterRoleBinding struct {
// }

// type Role struct {
// }

// type RoleBinding struct {
// }
