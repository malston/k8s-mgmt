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
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/malston/k8s-mgmt/pkg/cli"
	"github.com/malston/k8s-mgmt/pkg/kmgmt"
	"github.com/mitchellh/go-homedir"
)

func main() {

	conf := newConfigClient()
	root := kmgmt.CreateRootCommand(conf)
	fmt.Println() // Print a blank line before output for readability

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
	fmt.Println() // Print a blank line after output for readability
}

func newConfigClient() *cli.Config {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Print(err.Error() + "\n")
		os.Exit(1)
	}
	configDir := filepath.Join(home, ".k8s-mgmt", "config")

	if mgmtEnvConf, ok := os.LookupEnv("K8SMGMT_HOME"); ok {
		configDir = mgmtEnvConf
	}

	return cli.NewConfig(configDir)
}
