// k8s/client.go
package k8s

import (
	"context"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// Client wraps the Kubernetes client
type Client struct {
	clientset *kubernetes.Clientset
	namespace string
}

// creates k8s client
func NewClient() (*Client, error) {
	var kubeconfig string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &Client{
		clientset: clientset,
		namespace: "default",
	}, nil
}

func (c *Client) SetNamespace(namespace string) {
	c.namespace = namespace
}

func (c *Client) GetNamespace() string {
	return c.namespace
}

// TestConnection tests the connection to the Kubernetes cluster
func (c *Client) TestConnection(ctx context.Context) error {
	_, err := c.clientset.Discovery().ServerVersion()
	return err
}
