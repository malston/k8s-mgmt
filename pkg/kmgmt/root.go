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
package kmgmt

import (
	"github.com/malston/k8s-mgmt/pkg/cli"
	"github.com/malston/k8s-mgmt/pkg/k8s"
	"github.com/malston/k8s-mgmt/pkg/namespace"

	"github.com/spf13/cobra"
)

func CreateRootCommand(k k8s.Client, config *cli.Config) *cobra.Command {

	rootCmd := &cobra.Command{
		Use:   "kmgmt",
		Short: "A brief description of your application",
		Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		//	Run: func(cmd *cobra.Command, args []string) { },
	}

	namespaceCmd := namespace.NewCommand(config, k)
	rootCmd.AddCommand(namespaceCmd)
	// rootCmd.AddCommand(cluster.NewCommand())

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	return rootCmd
}
