package pks

import (
	"os"
	"time"

	"github.com/malston/k8s-mgmt/pkg/config"
	"github.com/malston/k8s-mgmt/pkg/exec"
)

type Client interface {
	CreateCluster(cluster *config.Cluster) error
}

func (m *pksClient) CreateCluster(cluster *config.Cluster) error {
	return m.Run("pks", "create-cluster", cluster.Name, "-p", cluster.Plan, "-n", cluster.NumNodes)
}

func NewClient() Client {
	clr := exec.NewCommandLineRunner(os.Stdout, exec.WithTimeout(100*time.Millisecond))
	return &pksClient{clr}
}

type pksClient struct {
	exec.CommandLineRunner
}
