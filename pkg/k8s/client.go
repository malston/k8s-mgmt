package k8s

import (
	"fmt"

	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

type Client interface {
	Core() corev1.CoreV1Interface
	CurrentContext() string
	SetContext(context string) error
}

func (c *client) Core() corev1.CoreV1Interface {
	return c.lazyLoadKubernetesClientsetOrDie().CoreV1()
}

func (c *client) CurrentContext() string {
	if c.context != "" {
		return c.context
	}
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
		panic(err)
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
				panic(err)
			}
			c.kubeClientset = kubernetes.NewForConfigOrDie(restConfig)
			c.context = context
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

type client struct {
	kubeClientset  *kubernetes.Clientset
	kubeConfigFile string
	kubeConfig     clientcmd.ClientConfig
	context        string
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
