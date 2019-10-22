package pks

import (
	"encoding/json"
	"fmt"

	"github.com/malston/k8s-mgmt/pkg/config"
	"github.com/malston/k8s-mgmt/pkg/exec"
)

type Client interface {
	CreateCluster(cluster *config.Cluster) error
	ShowCluster(name string) (*config.Cluster, error)
}

func (c *pksClient) CreateCluster(cluster *config.Cluster) error {
	cmd := []string{"create-cluster", cluster.Name,
		"--plan", cluster.Plan,
		"--num-nodes", cluster.NumNodes,
		"--external-hostname", cluster.ExternalHostname}
	if cluster.NetworkProfile != "" {
		cmd = append(cmd, "--network-profile", cluster.NetworkProfile)
	}
	return c.Run("pks", cmd...)
}

func (c *pksClient) ShowCluster(name string) (*config.Cluster, error) {
	cmdOut, err := c.RunOutput("pks", "show-cluster", name, "--json")
	if err != nil {
		return nil, err
	}
	var cluster *config.Cluster
	err = json.Unmarshal(cmdOut, &cluster)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling json into cluster, %v", err)
	}
	return cluster, nil
}

func NewClient(clr exec.CommandLineRunner) Client {
	return &pksClient{clr}
}

type pksClient struct {
	exec.CommandLineRunner
}
