// Code generated by protoc-gen-solo-kit. DO NOT EDIT.

package v1

import (
	"context"
	"time"

	gloo_solo_io "github.com/solo-io/supergloo/pkg/api/external/gloo/v1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/factory"
	"github.com/solo-io/solo-kit/pkg/api/v1/clients/memory"
)

var _ = Describe("TranslatorEventLoop", func() {
	var (
		namespace string
		emitter   TranslatorEmitter
		err       error
	)

	BeforeEach(func() {

		meshClientFactory := &factory.MemoryResourceClientFactory{
			Cache: memory.NewInMemoryResourceCache(),
		}
		meshClient, err := NewMeshClient(meshClientFactory)
		Expect(err).NotTo(HaveOccurred())

		routingRuleClientFactory := &factory.MemoryResourceClientFactory{
			Cache: memory.NewInMemoryResourceCache(),
		}
		routingRuleClient, err := NewRoutingRuleClient(routingRuleClientFactory)
		Expect(err).NotTo(HaveOccurred())

		upstreamClientFactory := &factory.MemoryResourceClientFactory{
			Cache: memory.NewInMemoryResourceCache(),
		}
		upstreamClient, err := gloo_solo_io.NewUpstreamClient(upstreamClientFactory)
		Expect(err).NotTo(HaveOccurred())

		secretClientFactory := &factory.MemoryResourceClientFactory{
			Cache: memory.NewInMemoryResourceCache(),
		}
		secretClient, err := gloo_solo_io.NewSecretClient(secretClientFactory)
		Expect(err).NotTo(HaveOccurred())

		emitter = NewTranslatorEmitter(meshClient, routingRuleClient, upstreamClient, secretClient)
	})
	It("runs sync function on a new snapshot", func() {
		_, err = emitter.Mesh().Write(NewMesh(namespace, "jerry"), clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())
		_, err = emitter.RoutingRule().Write(NewRoutingRule(namespace, "jerry"), clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())
		_, err = emitter.Upstream().Write(gloo_solo_io.NewUpstream(namespace, "jerry"), clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())
		_, err = emitter.Secret().Write(gloo_solo_io.NewSecret(namespace, "jerry"), clients.WriteOpts{})
		Expect(err).NotTo(HaveOccurred())
		sync := &mockTranslatorSyncer{}
		el := NewTranslatorEventLoop(emitter, sync)
		_, err := el.Run([]string{namespace}, clients.WatchOpts{})
		Expect(err).NotTo(HaveOccurred())
		Eventually(func() bool { return sync.synced }, time.Second).Should(BeTrue())
	})
})

type mockTranslatorSyncer struct {
	synced bool
}

func (s *mockTranslatorSyncer) Sync(ctx context.Context, snap *TranslatorSnapshot) error {
	s.synced = true
	return nil
}
