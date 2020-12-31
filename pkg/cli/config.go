package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/malston/k8s-mgmt/pkg/config"
	"github.com/malston/k8s-mgmt/pkg/exec"
	"github.com/malston/k8s-mgmt/pkg/k8s"
	"github.com/malston/k8s-mgmt/pkg/pks"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Config struct {
	CompiledEnv
	KubeConfigFile  string
	ViperConfigFile string
	k8s.Client
	PKSClient pks.Client
	config.Manager
	ConfigDir string
	Name      string
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

func NewDefaultConfig() (*Config, error) {
	home, err := homedir.Dir()
	if err != nil {
		return nil, err
	}
	configDir := filepath.Join(home, ".k8s-mgmt", "config")

	if mgmtEnvConf, ok := os.LookupEnv("K8SMGMT_HOME"); ok {
		configDir = mgmtEnvConf
	}

	return NewConfig(configDir), nil
}

func NewConfig(configDir string) *Config {
	c := &Config{
		CompiledEnv: env,
		ConfigDir:   configDir,
		Name:        "k8s-mgmt",
		Stdin:       os.Stdin,
		Stdout:      os.Stdout,
		Stderr:      os.Stderr,
	}

	cobra.OnInitialize(c.initViperConfig)
	cobra.OnInitialize(c.initKubeConfig)
	cobra.OnInitialize(c.init)

	return c
}

// initViperConfig reads in config file and ENV variables if set.
func (c *Config) initViperConfig() {
	if c.ViperConfigFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(c.ViperConfigFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			// avoid color since we don't know if it should be enabled yet
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// Search config in home directory with name ".k8s-mgmt" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName("." + c.Name)
	}
	viper.SetEnvPrefix(c.Name)
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	if err == nil {
		c.Eprintf("Using config file: %s\n", viper.ConfigFileUsed())
	}
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
			os.Exit(1)
		}
		c.KubeConfigFile = filepath.Join(home, ".kube", "config")
	}
}

func (c *Config) init() {
	if c.Client == nil {
		c.Client = k8s.NewClient(c.KubeConfigFile)
	}
	if c.PKSClient == nil {
		clr := exec.NewCommandLineRunner(os.Stdout, os.Stderr)
		c.PKSClient = pks.NewClient(clr)
	}
	if c.Manager == nil {
		var err error
		c.Manager, err = config.NewManager(c.ConfigDir)
		if err != nil {
			os.Exit(1)
		}
	}
}
