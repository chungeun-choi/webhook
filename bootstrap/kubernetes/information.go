package kubernetes

import "os"

type K8sInformation struct {
	Namespace      string
	KubeConfigPath string
}

func LoadInformation() *K8sInformation {
	var k8sInformation K8sInformation
	// load envs from kubernetes
	k8sInformation.Namespace = os.Getenv("POD_NAMESPACE")
	k8sInformation.KubeConfigPath = os.Getenv("KUBECONFIG")
	return &k8sInformation
}
