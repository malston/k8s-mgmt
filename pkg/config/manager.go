package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

// Manager provides all the configuration mgmt functions
type Manager interface {
	GetClusters() ([]*Cluster, error)
	GetNamespaces(cluster string) ([]*Namespace, error)
}

type configmanager struct {
	homeDir string
	config  *Config
}

// NewManager returns a new Manager
func NewManager(cfgHome string) (Manager, error) {
	_, err := ioutil.ReadDir(cfgHome)
	if err != nil {
		return nil, err
	}

	m := &configmanager{
		homeDir: cfgHome,
	}
	return m, nil
}

func (m *configmanager) GetClusters() ([]*Cluster, error) {
	err := m.lazyLoadConfig()
	if err != nil {
		return nil, err
	}
	return m.config.Clusters, nil
}

func (m *configmanager) GetNamespaces(cluster string) ([]*Namespace, error) {
	err := m.lazyLoadConfig()
	if err != nil {
		return nil, err
	}

	if len(m.config.Clusters) == 0 {
		return nil, fmt.Errorf("cluster %s does not exist in config folder", cluster)
	}

	var n []*Namespace
	for _, c := range m.config.Clusters {
		if c.Name == cluster {
			return append(n, c.Namespaces...), nil
		}
	}

	return n, nil
}

func (m *configmanager) lazyLoadConfig() error {
	if m.config == nil {
		err := m.loadConfig()
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *configmanager) loadConfig() error {
	c := &Config{}
	cfgHome := m.homeDir
	var cluster *Cluster
	err := filepath.Walk(cfgHome, func(path string, f os.FileInfo, err error) error {
		filePath := path

		if s, ok := directoryContains(filePath, "cluster.yml"); ok {
			cluster = &Cluster{}
			cluster.Name = s[len(s)-2]
			// fmt.Printf("setting cluster name to %s\n", cluster.Name)

			err := loadClusterConfig(filePath, cluster)
			if err != nil {
				return err
			}

			c.Clusters = append(c.Clusters, cluster)
			cluster.Namespaces = make([]*Namespace, 0)
		}

		if s, ok := directoryContains(filePath, "namespace.yml"); ok {
			namespace := &Namespace{}
			namespace.Name = s[len(s)-2]
			// fmt.Printf("adding namespace '%s' to cluster '%s'\n", namespace.Name, cluster.Name)

			err := loadNamespaceConfig(filePath, namespace)
			if err != nil {
				return err
			}

			cluster.Namespaces = append(cluster.Namespaces, namespace)
		}
		return nil
	})
	if err != nil {
		return err
	}
	m.config = c
	return nil
}

func directoryContains(path string, e string) ([]string, bool) {
	s := strings.Split(path, "/")

	for _, a := range s {
		if a == e {
			return s, true
		}
	}
	return s, false
}

func loadClusterConfig(file string, cluster *Cluster) error {
	bytes, err := readFile(file)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(bytes, cluster)
	if err != nil {
		return fmt.Errorf("error loading cluster from: %s: %v", file, err)
	}

	return nil
}

func loadNamespaceConfig(file string, namespace *Namespace) error {
	bytes, err := readFile(file)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(bytes, namespace)
	if err != nil {
		return fmt.Errorf("error loading namespace from %s: %v", file, err)
	}

	return nil
}

func readFile(file string) ([]byte, error) {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %s: %v", file, err)
	}
	return bytes, nil
}
