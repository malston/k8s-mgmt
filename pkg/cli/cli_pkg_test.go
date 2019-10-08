package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

func TestInitViperConfig(t *testing.T) {
	defer viper.Reset()

	c := NewConfig("../config/testdata")
	output := &bytes.Buffer{}
	c.Stdout = output
	c.Stderr = output

	c.ViperConfigFile = "testdata/.k8s-mgmt.yaml"
	c.initViperConfig()

	expectedViperSettings := map[string]interface{}{}
	if diff := cmp.Diff(expectedViperSettings, viper.AllSettings()); diff != "" {
		t.Errorf("Unexpected viper settings (-expected, +actual): %s", diff)
	}
	if diff := cmp.Diff("Using config file: testdata/.k8s-mgmt.yaml", strings.TrimSpace(output.String())); diff != "" {
		t.Errorf("Unexpected output (-expected, +actual): %s", diff)
	}
}

func TestInitViperConfig_HomeDir(t *testing.T) {
	defer viper.Reset()

	home, homeisset := os.LookupEnv("HOME")
	defer func() {
		homedir.Reset()
		if homeisset {
			os.Setenv("HOME", home)
		} else {
			os.Unsetenv("HOME")
		}
	}()

	c := NewConfig("../config/testdata")
	output := &bytes.Buffer{}
	c.Stdout = output
	c.Stderr = output

	os.Setenv("HOME", "testdata")
	c.initViperConfig()

	expectedViperSettings := map[string]interface{}{}
	if diff := cmp.Diff(expectedViperSettings, viper.AllSettings()); diff != "" {
		t.Errorf("Unexpected viper settings (-expected, +actual): %s", diff)
	}
}

func TestInitKubeConfig_Flag(t *testing.T) {
	c := NewConfig("../config/testdata")
	output := &bytes.Buffer{}
	c.Stdout = output
	c.Stderr = output

	c.KubeConfigFile = "testdata/.kube/config"
	c.initKubeConfig()

	if expected, actual := "testdata/.kube/config", c.KubeConfigFile; expected != actual {
		t.Errorf("Expected kubeconfig path %q, actually %q", expected, actual)
	}
	if diff := cmp.Diff("", strings.TrimSpace(output.String())); diff != "" {
		t.Errorf("Unexpected output (-expected, +actual): %s", diff)
	}
}

func TestInitKubeConfig_EnvVar(t *testing.T) {
	kubeconfig, kubeconfigisset := os.LookupEnv("KUBECONFIG")
	defer func() {
		if kubeconfigisset {
			os.Setenv("KUBECONFIG", kubeconfig)
		} else {
			os.Unsetenv("KUBECONFIG")
		}
	}()

	c := NewConfig("../config/testdata")
	output := &bytes.Buffer{}
	c.Stdout = output
	c.Stderr = output

	os.Setenv("KUBECONFIG", "testdata/.kube/config")
	c.initKubeConfig()

	if expected, actual := "testdata/.kube/config", c.KubeConfigFile; expected != actual {
		t.Errorf("Expected kubeconfig path %q, actually %q", expected, actual)
	}
	if diff := cmp.Diff("", strings.TrimSpace(output.String())); diff != "" {
		t.Errorf("Unexpected output (-expected, +actual): %s", diff)
	}
}

func TestInitKubeConfig_HomeDir(t *testing.T) {
	home, homeisset := os.LookupEnv("HOME")
	defer func() {
		homedir.Reset()
		if homeisset {
			os.Setenv("HOME", home)
		} else {
			os.Unsetenv("HOME")
		}
	}()

	c := NewConfig("../config/testdata")
	output := &bytes.Buffer{}
	c.Stdout = output
	c.Stderr = output

	os.Setenv("HOME", "testdata")
	c.initKubeConfig()

	if expected, actual := filepath.FromSlash("testdata/.kube/config"), c.KubeConfigFile; expected != actual {
		t.Errorf("Expected kubeconfig path %q, actually %q", expected, actual)
	}
	if diff := cmp.Diff("", strings.TrimSpace(output.String())); diff != "" {
		t.Errorf("Unexpected output (-expected, +actual): %s", diff)
	}
}

func TestInit(t *testing.T) {
	c := NewConfig("../config/testdata")
	output := &bytes.Buffer{}
	c.Stdout = output
	c.Stderr = output

	c.KubeConfigFile = "testdata/.kube/config"
	c.init()

	if diff := cmp.Diff("", strings.TrimSpace(output.String())); diff != "" {
		t.Errorf("Unexpected output (-expected, +actual): %s", diff)
	}
	if c.Client == nil {
		t.Errorf("Expected c.Client tp be set, actually %v", c.Client)
	}
	if c.Manager == nil {
		t.Errorf("Expected c.Manager tp be set, actually %v", c.Manager)
	}
}
