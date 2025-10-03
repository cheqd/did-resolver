package unit

import (
	"time"

	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("EndpointManager", func() {
	Describe("creation", func() {
		It("creates without fallback endpoints", func() {
			config := types.Config{
				EnableFallbackEndpoints: false,
				Networks: []types.Network{
					{
						Namespace: "mainnet",
						Endpoints: []types.Endpoint{
							{
								URL:     "grpc.cheqd.network:443",
								UseTls:  true,
								Timeout: 5 * time.Second,
								Role:    types.EndpointRolePrimary,
							},
						},
						UseTls:  true,
						Timeout: 5 * time.Second,
					},
				},
			}

			em := services.NewEndpointManager(config)
			Expect(em).ToNot(BeNil())
		})

		It("creates with fallback endpoints enabled", func() {
			config := types.Config{
				EnableFallbackEndpoints: true,
				Networks: []types.Network{
					{
						Namespace: "mainnet",
						Endpoints: []types.Endpoint{
							{
								URL:     "grpc.cheqd.network:443",
								UseTls:  true,
								Timeout: 5 * time.Second,
								Role:    types.EndpointRolePrimary,
							},
							{
								URL:     "archive-grpc.cheqd.net:443",
								UseTls:  true,
								Timeout: 5 * time.Second,
								Role:    types.EndpointRoleFallback,
							},
						},
						UseTls:  true,
						Timeout: 5 * time.Second,
					},
				},
			}

			em := services.NewEndpointManager(config)
			Expect(em).ToNot(BeNil())
		})
	})

	Describe("GetHealthyEndpoint", func() {
		It("returns a healthy endpoint when available", func() {
			config := types.Config{
				EnableFallbackEndpoints: false,
				Networks: []types.Network{
					{
						Namespace: "mainnet",
						Endpoints: []types.Endpoint{
							{
								URL:     "grpc.cheqd.network:443",
								UseTls:  true,
								Timeout: 5 * time.Second,
								Role:    types.EndpointRolePrimary,
							},
						},
						UseTls:  true,
						Timeout: 5 * time.Second,
					},
				},
			}

			em := services.NewEndpointManager(config)
			Expect(em).ToNot(BeNil())

			// Wait for initial health check to complete
			time.Sleep(3 * time.Second)

			network, err := em.GetHealthyEndpoint("mainnet")
			Expect(err).ToNot(HaveOccurred())
			Expect(network).ToNot(BeNil())
			Expect(network.Endpoints).To(HaveLen(1))
			Expect(network.Endpoints[0].URL).To(Equal("grpc.cheqd.network:443"))
		})

		It("returns error for invalid endpoint and success for valid namespace", func() {
			config := types.Config{
				EnableFallbackEndpoints: false,
				Networks: []types.Network{
					{
						Namespace: "mainnet",
						Endpoints: []types.Endpoint{
							{
								URL:     "invalid-grpc.cheqd.net:443",
								UseTls:  true,
								Timeout: 5 * time.Second,
								Role:    types.EndpointRolePrimary,
							},
						},
						UseTls:  true,
						Timeout: 5 * time.Second,
					},
					{
						Namespace: "testnet",
						Endpoints: []types.Endpoint{
							{
								URL:     "grpc.cheqd.network:443",
								UseTls:  true,
								Timeout: 5 * time.Second,
								Role:    types.EndpointRolePrimary,
							},
						},
						UseTls:  true,
						Timeout: 5 * time.Second,
					},
				},
			}

			em := services.NewEndpointManager(config)
			Expect(em).ToNot(BeNil())

			// Wait for initial health check to complete
			time.Sleep(3 * time.Second)

			_, err := em.GetHealthyEndpoint("mainnet")
			Expect(err).To(HaveOccurred())

			network, err := em.GetHealthyEndpoint("testnet")
			Expect(err).ToNot(HaveOccurred())
			Expect(network).ToNot(BeNil())
		})

		It("returns error for non-existent namespace", func() {
			config := types.Config{
				EnableFallbackEndpoints: false,
				Networks: []types.Network{
					{
						Namespace: "mainnet",
						Endpoints: []types.Endpoint{
							{
								URL:     "grpc.cheqd.network:443",
								UseTls:  true,
								Timeout: 5 * time.Second,
								Role:    types.EndpointRolePrimary,
							},
						},
						UseTls:  true,
						Timeout: 5 * time.Second,
					},
				},
			}

			em := services.NewEndpointManager(config)

			_, err := em.GetHealthyEndpoint("nonexistent")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("no healthy endpoints available"))
		})
	})

	Describe("fallback behavior", func() {
		It("uses fallback when primary is unhealthy", func() {
			config := types.Config{
				EnableFallbackEndpoints: true,
				Networks: []types.Network{
					{
						Namespace: "mainnet",
						Endpoints: []types.Endpoint{
							{
								URL:     "invalid-grpc.cheqd.network:443",
								UseTls:  true,
								Timeout: 5 * time.Second,
								Role:    types.EndpointRolePrimary,
							},
							{
								URL:     "grpc.cheqd.network:443",
								UseTls:  true,
								Timeout: 5 * time.Second,
								Role:    types.EndpointRoleFallback,
							},
						},
						UseTls:  true,
						Timeout: 5 * time.Second,
					},
				},
			}

			em := services.NewEndpointManager(config)

			// Wait for initial health check to complete
			time.Sleep(3 * time.Second)

			network, err := em.GetHealthyEndpoint("mainnet")
			Expect(err).ToNot(HaveOccurred())
			Expect(network).ToNot(BeNil())
			Expect(network.Endpoints[0].URL).To(Equal("grpc.cheqd.network:443"))
		})
	})
})
