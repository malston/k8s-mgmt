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
	ResourceQuota []ResourceQuota
	// Limit         []Limit
	// NetworkPolicy []NetworkPolicy
}

// type ClusterInfo struct {
// }

// type Limit struct {
// }

// ResourceQuota defines quotas for a kubernetes namespace
type ResourceQuota struct {
	Name           string `yaml:"name"`
	RequestsCPU    string `yaml:"requests.cpu"`
	RequestsMemory string `yaml:"reqeusts.memory"`
	LimitsCPU      string `yaml:"limits.cpu"`
	LimitsMemory   string `yaml:"limits.memory"`
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
