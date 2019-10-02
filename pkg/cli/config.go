package cli

import (
	"os"
	"path/filepath"

	"github.com/malston/k8s-mgmt/pkg/config"
	"github.com/malston/k8s-mgmt/pkg/k8s"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

type Config struct {
	KubeConfigFile string
	k8s.Client
	config.Manager
}

func NewConfig() *Config {
	c := &Config{}

	// cobra.OnInitialize(c.initViperConfig)
	cobra.OnInitialize(c.initKubeConfig)
	cobra.OnInitialize(c.init)

	return c
}

func (c *Config) initKubeConfig() {
	if c.KubeConfigFile != "" {
		return
	}
	if kubeEnvConf, ok := os.LookupEnv("KUBECONFIG"); ok {
		c.KubeConfigFile = kubeEnvConf
	} else {
		home, err := homedir.Dir()
		if err != nil {
			// c.Errorf("%s\n", err)
			os.Exit(1)
		}
		c.KubeConfigFile = filepath.Join(home, ".kube", "config")
	}
}

func (c *Config) init() {
	if c.Client == nil {
		c.Client = k8s.NewClient(c.KubeConfigFile)
	}
	if c.Manager == nil {
		c.Manager, _ = config.NewManager("/Users/malston/workspace/k8s-mgmt/testdata")
	}
}
