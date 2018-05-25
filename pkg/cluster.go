package gube

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
)

type Cluster interface {
	GetClientConfig() (*restclient.Config, error)
	GetClientset() (*kubernetes.Clientset, error)

	GetNodeCount() (int, error)
	GetNodes() (map[string]corev1.Node, error)
	KubeconfigProvider
}

type cluster struct {
	access    KubeconfigProvider
	config    *restclient.Config
	clientset *kubernetes.Clientset
}

func NewCluster(cfg KubeconfigProvider) Cluster {
	return (&cluster{}).new(cfg)
}

func (this *cluster) new(cfg KubeconfigProvider) *cluster {
	this.access = cfg
	return this
}

func (this *cluster) GetKubeconfig() ([]byte, error) {
	return this.access.GetKubeconfig()
}

func (this *cluster) GetClientConfig() (*restclient.Config, error) {
	if this.config == nil {
		bytes, err := this.GetKubeconfig()
		if err != nil {
			return nil, err
		}
		this.config, err = NewConfigFromBytes(bytes)
		if err != nil {
			return nil, err
		}
	}
	return this.config, nil
}

func (this *cluster) GetClientset() (*kubernetes.Clientset, error) {
	if this.clientset == nil {
		bytes, err := this.GetKubeconfig()
		if err != nil {
			return nil, err
		}
		this.clientset, err = NewClientFromBytes(bytes)
		if err != nil {
			return nil, err
		}
	}
	return this.clientset, nil
}

func (this *cluster) GetNodeCount() (int, error) {
	cs, err := this.GetClientset()
	if err != nil {
		return 0, err
	}
	list, err := cs.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		return 0, fmt.Errorf("failed to get node count: %s", err)
	}
	return len(list.Items), nil
}

func (this *cluster) GetNodes() (map[string]corev1.Node, error) {
	cs, err := this.GetClientset()
	if err != nil {
		return nil, err
	}
	list, err := cs.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get nodes: %s", err)
	}
	nodes := map[string]corev1.Node{}
	for _, n := range list.Items {
		nodes[n.GetName()] = n
	}
	return nodes, nil
}
