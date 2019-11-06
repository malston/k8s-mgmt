package config

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
	Resourcequota []Resourcequota
	// Limit         []Limit
	// NetworkPolicy []NetworkPolicy
}

// type ClusterInfo struct {
// }

// type Limit struct {
// }

// The Resourcequota struct defines quotas for a kubernetes namespace
type Resourcequota struct {
	Name           string  `yaml:"name"`
	Requestscpu    float64 `yaml:"requests.cpu"`
	Requestsmemory string  `yaml:"reqeusts.memory"`
	Limitscpu      float64 `yaml:"limits.cpu"`
	Limitsmemory   string  `yaml:"limits.memory"`
}

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
