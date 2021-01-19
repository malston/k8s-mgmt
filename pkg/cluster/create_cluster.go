package cluster

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/malston/k8s-mgmt/pkg/cli"
	"github.com/malston/k8s-mgmt/pkg/config"
	"github.com/spf13/cobra"
)

func NewCommand(conf *cli.Config) *cobra.Command {
	c := &create{
		c: conf,
	}
	cmd := c.command()

	return cmd
}

type create struct {
	*cobra.Command
	c *cli.Config
}

func (c *create) command() *cobra.Command {
	return &cobra.Command{
		Use:   "create-clusters",
		Short: "Creates clusters",
		Long: strings.TrimSpace(`
Loops through files under the config directory, finds all the cluster folders, 
opens each cluster.yml file, and creates a new cluster based on contents of the file.
`),
		RunE: c.runE,
		Args: cobra.ExactArgs(0),
	}
}

func (c *create) runE(cmd *cobra.Command, args []string) error {
	m := c.c.Manager
	errors := make(chan error)
	started := make(chan *config.Cluster)
	var errs []string
	var startedWG sync.WaitGroup

	clusters, _ := m.GetClusters()
	for _, cl := range clusters {
		startedWG.Add(1)
		go c.createCluster(cl, &startedWG, started, errors)
	}

	go func() {
		startedWG.Wait()
		close(started)
		close(errors)
	}()

	go func() {
		if err := <-errors; err != nil {
			errs = append(errs, err.Error())
		}
	}()

	var finishedWG sync.WaitGroup
	finished := make(chan *config.Cluster)
	for cl := range started {
		finishedWG.Add(1)
		c.c.Printf("Waiting for cluster %s to complete\n", cl.Name)
		go c.waitForClusterCompletion(cl.Name, &finishedWG, finished)
	}

	go func() {
		finishedWG.Wait()
		close(finished)
	}()

	for cl := range finished {
		c.c.Printf("Cluster %s created\n", cl.Name)
	}

	if len(errs) > 0 {
		return fmt.Errorf(strings.Join(errs, "\n"))
	}

	return nil
}

func (c *create) createCluster(cluster *config.Cluster, wg *sync.WaitGroup, results chan<- *config.Cluster, errors chan<- error) {
	defer wg.Done()
	err := c.c.CreateCluster(cluster)
	if err != nil {
		errors <- fmt.Errorf("failed to create cluster %s, %v", cluster.Name, err)
		return
	}
	results <- cluster
}

func (c *create) waitForClusterCompletion(name string, wg *sync.WaitGroup, results chan<- *config.Cluster) {
	defer wg.Done()
	for {
		cluster, _ := c.c.ShowCluster(name)
		if len(cluster.IPAddress) > 0 {
			results <- cluster
			return
		}
		c.c.Printf("%s", ".")
		time.Sleep(2 * time.Second)
	}
}
