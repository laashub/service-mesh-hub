package config_lookup

import (
	"context"

	"github.com/rotisserie/eris"
	discovery_v1alpha1 "github.com/solo-io/service-mesh-hub/pkg/api/discovery.zephyr.solo.io/v1alpha1"
	kubernetes_core "github.com/solo-io/service-mesh-hub/pkg/clients/kubernetes/core"
	zephyr_discovery "github.com/solo-io/service-mesh-hub/pkg/clients/zephyr/discovery"
	"github.com/solo-io/service-mesh-hub/pkg/kubeconfig"
)

var (
	FailedToFindKubeConfigSecret = func(err error, clusterName string) error {
		return eris.Wrapf(err, "Failed to find kube config secret for cluster %s", clusterName)
	}
	FailedToConvertSecretToKubeConfig = func(err error, clusterName string) error {
		return eris.Wrapf(err, "Failed to convert kube config secret for cluster %s to REST config", clusterName)
	}
	ClusterNotFound = func(clusterName string) error {
		return eris.Errorf("Cluster '%s' was not found", clusterName)
	}
)

func NewKubeConfigLookup(
	kubeClusterClient zephyr_discovery.KubernetesClusterClient,
	secrestClient kubernetes_core.SecretsClient,
	secretToConfigConverter kubeconfig.SecretToConfigConverter,
) KubeConfigLookup {
	return &kubeConfigLookup{
		kubeClusterClient:       kubeClusterClient,
		secretsClient:           secrestClient,
		secretToConfigConverter: secretToConfigConverter,
	}
}

type kubeConfigLookup struct {
	secretsClient           kubernetes_core.SecretsClient
	secretToConfigConverter kubeconfig.SecretToConfigConverter
	kubeClusterClient       zephyr_discovery.KubernetesClusterClient
}

func (k *kubeConfigLookup) FromCluster(ctx context.Context, clusterName string) (config *kubeconfig.Config, err error) {
	var kubeCluster *discovery_v1alpha1.KubernetesCluster
	allClusters, err := k.kubeClusterClient.List(ctx)
	if err != nil {
		return nil, err
	}
	for _, foundCluster := range allClusters.Items {
		if foundCluster.GetName() == clusterName {
			kubeCluster = &foundCluster
			break
		}
	}

	if kubeCluster == nil {
		return nil, ClusterNotFound(clusterName)
	}

	cfgSecretRef := kubeCluster.Spec.GetSecretRef()
	secret, err := k.secretsClient.Get(ctx, cfgSecretRef.GetName(), cfgSecretRef.GetNamespace())
	if err != nil {
		return nil, FailedToFindKubeConfigSecret(err, kubeCluster.GetName())
	}

	clusterName, config, err = k.secretToConfigConverter(secret)
	if err != nil {
		return nil, FailedToConvertSecretToKubeConfig(err, kubeCluster.GetName())
	}

	return config, nil
}
