package kube_cluster

import (
	"github.com/chungeun-choi/webhook/service/endpoint_construct"
	"github.com/chungeun-choi/webhook/service/mutating"
)

type Cluster struct {
	Name            string
	KubeAPIClient   *ClientInterface
	mutatingManager *mutating.MutatingManager
	endpointManager *endpoint_construct.EndpointManager
}

func NewCluster(name string, kubeAPIClient *ClientInterface, mutatingManager *mutating.MutatingManager, endpointManager *endpoint_construct.EndpointManager) *Cluster {
	return &Cluster{
		Name:            name,
		KubeAPIClient:   kubeAPIClient,
		mutatingManager: mutatingManager,
		endpointManager: endpointManager,
	}
}

func (c *Cluster) GetClusterInfo() {}

func (c *Cluster) GetMutatingManager() *mutating.MutatingManager {
	return c.mutatingManager
}

func (c *Cluster) GetEndpointManager() *endpoint_construct.EndpointManager {
	return c.endpointManager
}

func (c *Cluster) UpdateMutatingManager(mutatingManager *mutating.MutatingManager) {
	c.mutatingManager = mutatingManager
}

func (c *Cluster) UpdateEndpointManager(endpointManager *endpoint_construct.EndpointManager) {
	c.endpointManager = endpointManager
}

func (c *Cluster) Delete() {

}
