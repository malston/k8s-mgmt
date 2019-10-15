package pks

import (
	"github.com/malston/k8s-mgmt/pkg/config"
	"github.com/malston/k8s-mgmt/pkg/exec"
)

type Client interface {
	CreateCluster(cluster *config.Cluster) error
}

func (c *pksClient) CreateCluster(cluster *config.Cluster) error {
	return c.Run("pks", "create-cluster", cluster.Name,
		"--plan", cluster.Plan,
		"--num-nodes", cluster.NumNodes,
		"--external-hostname", cluster.ExternalHostname)
}

func NewClient(clr exec.CommandLineRunner) Client {
	return &pksClient{clr}
}

type pksClient struct {
	exec.CommandLineRunner
}
