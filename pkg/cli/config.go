package cli

import (
	"fmt"
	"io"
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
	ConfigDir string
	Stdin     io.Reader
	Stdout    io.Writer
	Stderr    io.Writer
}

func (c *Config) Printf(format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(c.Stdout, format, a...)
}

func (c *Config) Eprintf(format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(c.Stderr, format, a...)
}

func NewConfig(configDir string) *Config {
	c := &Config{
		ConfigDir: configDir,
		Stdin:     os.Stdin,
		Stdout:    os.Stdout,
		Stderr:    os.Stderr,
	}

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
		var err error
		c.Manager, err = config.NewManager(c.ConfigDir)
		if err != nil {
			fmt.Print(err.Error() + "\n")
			os.Exit(1)
		}
	}
}
