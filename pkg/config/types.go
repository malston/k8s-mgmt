package config

// ClusterConfig
type Config struct {
	Clusters []*Cluster `yaml:"clusters"`
	// Quota         []Quota
	// Limit         []Limit
	// NetworkPolicy []NetworkPolicy
}

// Cluster
type Cluster struct {
	Name             string       `yaml:"name"`
	Namespaces       []*Namespace `yaml:"namespaces"`
	Plan             string       `yaml:"plan"`
	ExternalHostname string       `yaml:"external-hostname"`
	NumNodes         string       `yaml:"num-nodes"`
	NetworkProfile   string       `yaml:"network-profile"`
	// ClusterRole
	// ClusterRoleBinding
}

// Namespace
type Namespace struct {
	Name string `yaml:"name"`
	// Role
	// RoleBinding
	// Quota         []Quota
	// Limit         []Limit
	// NetworkPolicy []NetworkPolicy
}

// type ClusterInfo struct {
// }

// type Limit struct {
// }

// type Quota struct {
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
