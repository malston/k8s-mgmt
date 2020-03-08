package resourcequota

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

import (
	"github.com/malston/k8s-mgmt/pkg/cli"
	"github.com/spf13/cobra"
)

// NewCommand Configure cobra cli sub function
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
		Use:   "create-resourcequota <cluster> <namespace>",
		Short: "Create resourcequota for a given cluster and namespace",
		RunE:  c.runE,
		Args:  cobra.ExactArgs(2),
	}
}

func (c *create) runE(cmd *cobra.Command, args []string) error {
	client := c.c.Client
	m := c.c.Manager
	cluster := args[0]
	namespace := args[1]
	rq, err := m.GetResourceQuota(cluster, namespace)
	if err != nil {
		return err
	}
	r, e := client.Core().ResourceQuotas(namespace).Create(rq)
	if e != nil {
		return e
	}

	c.c.Printf("resourcequota/%s created\n", r.GetName())

	return err
}
