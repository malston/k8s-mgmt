/*
Copyright © 2019 Mark Alston <marktalston@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
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
	clusters, _ := m.GetClusters()

	errors := make(chan error)
	results := make(chan *config.Cluster)
	var errs []string
	var wg1 sync.WaitGroup
	for _, cl := range clusters {
		wg1.Add(1)
		go c.createCluster(cl, &wg1, results, errors)
	}

	go func() {
		wg1.Wait()
		close(results)
		close(errors)
	}()

	go func() {
		if err := <-errors; err != nil {
			errs = append(errs, err.Error())
		}
	}()

	var wg2 sync.WaitGroup
	results2 := make(chan *config.Cluster)
	for res := range results {
		wg2.Add(1)
		c.c.Printf("Waiting for cluster %s to complete\n", res.Name)
		go c.waitForClusterCompletion(res.Name, &wg2, results2)
	}

	go func() {
		wg2.Wait()
		close(results2)
	}()

	for res := range results2 {
		c.c.Printf("Cluster %s created\n", res.Name)
	}

	if len(errs) > 0 {
		return fmt.Errorf(strings.Join(errs, "\n"))
	}

	return nil
}

func (c *create) createCluster(cluster *config.Cluster, wg *sync.WaitGroup, results chan<- *config.Cluster, errors chan<- error) {
	defer wg.Done()
	err := c.c.PKSClient.CreateCluster(cluster)
	if err != nil {
		errors <- fmt.Errorf("failed to create cluster %s, %v", cluster.Name, err)
		return
	}
	results <- cluster
}

func (c *create) waitForClusterCompletion(name string, wg *sync.WaitGroup, results chan<- *config.Cluster) {
	defer wg.Done()
	for {
		cluster, _ := c.c.PKSClient.ShowCluster(name)
		if len(cluster.IPAddress) > 0 {
			results <- cluster
			return
		}
		c.c.Printf("%s", ".")
		time.Sleep(2 * time.Second)
	}
}
