package k8s

import (
	"fmt"
	"io"

	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

type Client interface {
	Core() corev1.CoreV1Interface
	CurrentContext() string
	SetContext(context string) error
	Printf(format string, a ...interface{}) (n int, err error)
	Eprintf(format string, a ...interface{}) (n int, err error)
}

func (c *client) Core() corev1.CoreV1Interface {
	return c.lazyLoadKubernetesClientsetOrDie().CoreV1()
}

func (c *client) CurrentContext() string {
	rc, err := c.lazyLoadKubeConfig().RawConfig()
	if err != nil {
		panic(err)
	}

	return rc.CurrentContext
}

func (c *client) SetContext(context string) error {
	if context == c.CurrentContext() {
		return nil
	}
	rc, err := c.lazyLoadKubeConfig().RawConfig()
	if err != nil {
		return err
	}

	for name, cluster := range rc.Clusters {
		if context == name {
			c.kubeConfig = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
				&clientcmd.ClientConfigLoadingRules{ExplicitPath: c.kubeConfigFile},
				&clientcmd.ConfigOverrides{ClusterInfo: clientcmdapi.Cluster{Server: cluster.Server},
					CurrentContext: context},
			)
			restConfig, err := c.kubeConfig.ClientConfig()
			if err != nil {
				return err
			}
			c.kubeClientset = kubernetes.NewForConfigOrDie(restConfig)
			return nil
		}
	}
	return fmt.Errorf("context '%s' not found", context)
}

func NewClient(kubeConfigFile string) Client {
	return &client{
		kubeConfigFile: kubeConfigFile,
	}
}

func (c *client) Printf(format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(c.Stdout, format, a...)
}

func (c *client) Eprintf(format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(c.Stderr, format, a...)
}

type client struct {
	kubeClientset  *kubernetes.Clientset
	kubeConfigFile string
	kubeConfig     clientcmd.ClientConfig
	restConfig     *rest.Config
	Stdin          io.Reader
	Stdout         io.Writer
	Stderr         io.Writer
}

func (c *client) lazyLoadKubeConfig() clientcmd.ClientConfig {
	if c.kubeConfig == nil {
		c.kubeConfig = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			&clientcmd.ClientConfigLoadingRules{ExplicitPath: c.kubeConfigFile},
			&clientcmd.ConfigOverrides{ClusterInfo: clientcmdapi.Cluster{Server: ""}},
		)
	}
	return c.kubeConfig
}

func (c *client) lazyLoadKubernetesClientsetOrDie() *kubernetes.Clientset {
	if c.kubeClientset == nil {
		kubeConfig := c.lazyLoadKubeConfig()
		restConfig, err := kubeConfig.ClientConfig()
		if err != nil {
			panic(err)
		}
		c.kubeClientset = kubernetes.NewForConfigOrDie(restConfig)
	}
	return c.kubeClientset
}
