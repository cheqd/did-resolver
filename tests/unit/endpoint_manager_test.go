package unit

import (
	"testing"
	"time"

	"github.com/cheqd/did-resolver/services"
	"github.com/cheqd/did-resolver/types"
)

func TestEndpointManager_BasicCreation(t *testing.T) {
	// Test basic creation without fallback endpoints
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
				UseTls:   true,
				Timeout:  5 * time.Second,
			},
		},
	}

	em := services.NewEndpointManager(config)
	if em == nil {
		t.Error("expected EndpointManager to be created, got nil")
	}
}

func TestEndpointManager_FallbackEnabled(t *testing.T) {
	// Test creation with fallback endpoints enabled
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
				UseTls:   true,
				Timeout:  5 * time.Second,
			},
		},
	}

	em := services.NewEndpointManager(config)
	if em == nil {
		t.Error("expected EndpointManager to be created, got nil")
	}
}

func TestEndpointManager_GetHealthyEndpoint_ValidEndpoint(t *testing.T) {
	// Test that we can get a healthy endpoint
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
				UseTls:   true,
				Timeout:  5 * time.Second,
			},
		},
	}

	em := services.NewEndpointManager(config)
	if em == nil {
		t.Error("expected EndpointManager to be created, got nil")
	}
	
	// Wait for initial health check to complete
	time.Sleep(3 * time.Second)
	
	// Test that we can get a healthy endpoint
	network, err := em.GetHealthyEndpoint("mainnet")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if network == nil {
		t.Error("expected healthy endpoint, got nil")
	}
	
	t.Logf("Got healthy endpoint: %s", network.Endpoints[0].URL)
	if network.Endpoints[0].URL != "grpc.cheqd.network:443" {
		t.Errorf("expected endpoint %s, got %s", "grpc.cheqd.network:443", network.Endpoints[0].URL)
	}
}

func TestEndpointManager_HealthCheckInvalidURL(t *testing.T) {
	// Configure invalid endpoint for mainnet and valid endpoint for testnet
	config := types.Config{
		EnableFallbackEndpoints: false,
		Networks: []types.Network{
			{
				Namespace: "mainnet",
				Endpoints: []types.Endpoint{
					{
						URL:     "invalid-grpc.cheqd.net:443", // Invalid endpoint
						UseTls:  true,
						Timeout: 5 * time.Second,
						Role:    types.EndpointRolePrimary,
					},
				},
				UseTls:   true,
				Timeout:  5 * time.Second,
			},
			{
				Namespace: "testnet",
				Endpoints: []types.Endpoint{
					{
						URL:     "grpc.cheqd.network:443", // Valid endpoint
						UseTls:  true,
						Timeout: 5 * time.Second,
						Role:    types.EndpointRolePrimary,
					},
				},
				UseTls:   true,
				Timeout:  5 * time.Second,
			},
		},
	}

	em := services.NewEndpointManager(config)
	if em == nil {
		t.Fatal("expected EndpointManager to be created, got nil")
	}

	// Wait for initial health check to complete
	time.Sleep(3 * time.Second)

	// Request healthy endpoint for mainnet (invalid) and expect an error
	_, err := em.GetHealthyEndpoint("mainnet")
	if err == nil {
		t.Fatal("expected error for invalid mainnet endpoint, got nil")
	}

	// And for testnet (valid) we should get a healthy endpoint
	network, err := em.GetHealthyEndpoint("testnet")
	if err != nil || network == nil {
		t.Fatalf("expected healthy endpoint for testnet, got err=%v network=%v", err, network)
	}
}

func TestEndpointManager_GetHealthyEndpoint_NonExistentNamespace(t *testing.T) {
	// Test that requesting a non-existent namespace returns an error
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
				UseTls:   true,
				Timeout:  5 * time.Second,
			},
		},
	}

	em := services.NewEndpointManager(config)
	
	// Test getting endpoint for non-existent namespace
	_, err := em.GetHealthyEndpoint("nonexistent")
	if err == nil {
		t.Error("expected error for non-existent namespace, got nil")
	}
	
	// Verify it's the correct error type
	if err.Error() != "no healthy endpoints available" {
		t.Errorf("expected 'no healthy endpoints available' error, got: %v", err)
	}
}

func TestEndpointManager_FallbackBehavior(t *testing.T) {
	// Test fallback behavior with one valid and one invalid endpoint
	config := types.Config{
		EnableFallbackEndpoints: true,
		Networks: []types.Network{
			{
				Namespace: "mainnet",
				Endpoints: []types.Endpoint{
					{
						URL:     "invalid-grpc.cheqd.network:443", // Primary (invalid)
						UseTls:  true,
						Timeout: 5 * time.Second,
						Role:    types.EndpointRolePrimary,
					},
					{
						URL:     "grpc.cheqd.network:443", // Fallback (valid)
						UseTls:  true,
						Timeout: 5 * time.Second,
						Role:    types.EndpointRoleFallback,
					},
				},
				UseTls:   true,
				Timeout:  5 * time.Second,
			},
		},
	}

	em := services.NewEndpointManager(config)
	
	// Wait for initial health check to complete
	time.Sleep(3 * time.Second)
	
	// Test that we get a healthy endpoint (should be fallback)
	network, err := em.GetHealthyEndpoint("mainnet")
	if err != nil {
		t.Errorf("expected healthy endpoint, got error: %v", err)
	}
	if network == nil {
		t.Error("expected healthy endpoint, got nil")
	}
	
	// Verify it's using the fallback endpoint
	if network.Endpoints[0].URL != "grpc.cheqd.network:443" {
		t.Errorf("expected fallback endpoint 'grpc.cheqd.network:443', got '%s'", network.Endpoints[0].URL)
	}
	
	t.Logf("Successfully using fallback endpoint: %s", network.Endpoints[0].URL)
} 