// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package wire

import (
	"context"

	kubernetes_core "github.com/solo-io/service-mesh-hub/pkg/clients/kubernetes/core"
	"github.com/solo-io/service-mesh-hub/pkg/security/certgen"
	"github.com/solo-io/service-mesh-hub/pkg/wire_providers"
	mc_wire "github.com/solo-io/service-mesh-hub/services/common/multicluster/wire"
	csr_generator "github.com/solo-io/service-mesh-hub/services/csr-agent/pkg/csr-generator"
)

// Injectors from wire.go:

func InitializeCsrAgent(ctx context.Context) (CsrAgentContext, error) {
	config, err := mc_wire.LocalKubeConfigProvider()
	if err != nil {
		return CsrAgentContext{}, err
	}
	asyncManager, err := mc_wire.LocalManagerProvider(ctx, config)
	if err != nil {
		return CsrAgentContext{}, err
	}
	virtualMeshCertificateSigningRequestEventWatcher := csr_generator.CsrControllerProviderLocal(asyncManager)
	virtualMeshCSRDataSourceFactory := csr_generator.NewVirtualMeshCSRDataSourceFactory()
	clientset, err := wire_providers.NewSecurityClients(config)
	if err != nil {
		return CsrAgentContext{}, err
	}
	virtualMeshCertificateSigningRequestClient := wire_providers.NewVirtualMeshCertificateSigningRequestClient(clientset)
	client := mc_wire.DynamicClientProvider(asyncManager)
	secretClient := kubernetes_core.NewSecretClient(client)
	signer := certgen.NewSigner()
	privateKeyGenerator := csr_generator.NewPrivateKeyGenerator()
	certClient := csr_generator.NewCertClient(secretClient, signer, privateKeyGenerator)
	istioCSRGenerator := csr_generator.NewIstioCSRGenerator(virtualMeshCertificateSigningRequestClient, secretClient, certClient, signer)
	virtualMeshCSRProcessor := csr_generator.NewCsrAgentIstioProcessor(istioCSRGenerator)
	csrAgentContext := CsrAgentContextProvider(ctx, asyncManager, virtualMeshCertificateSigningRequestEventWatcher, virtualMeshCSRDataSourceFactory, virtualMeshCSRProcessor, virtualMeshCertificateSigningRequestClient)
	return csrAgentContext, nil
}
