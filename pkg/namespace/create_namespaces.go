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
package namespace

import (
	"fmt"
	"strings"

	"github.com/malston/k8s-mgmt/pkg/cli"
	"github.com/spf13/cobra"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
		Use:   "create-namespaces <cluster-name>",
		Short: "Creates namespaces",
		Long: strings.TrimSpace(`
Loops through files under the config directory, finds all the namespace folders, 
opens each namespace.yml file, and creates a new namespace based on contents of the file.
`),
		RunE: c.runE,
		Args: cobra.ExactArgs(1),
	}
}

func (c *create) runE(cmd *cobra.Command, args []string) error {
	client := c.c.Client
	m := c.c.Manager
	clusterName := args[0]
	if client.CurrentContext() != clusterName {
		err := client.SetContext(clusterName)
		//TODO wrap error
		if err != nil {
			return err
		}
	}
	namespaces, err := m.GetNamespaces(clusterName)
	if err != nil {
		return err
	}
	if len(namespaces) == 0 {
		return fmt.Errorf("no namespaces found for cluster %s", clusterName)
	}
	for _, ns := range namespaces {
		n, e := client.Core().Namespaces().Create(&v1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: ns.Name,
			},
		})
		if e != nil {
			return e
		}

		c.c.Printf("Namespace %s created\n", n.GetName())
	}

	return err
}
