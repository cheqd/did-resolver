package services

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	didTypes "github.com/cheqd/cheqd-node/api/v2/cheqd/did/v2"
	"github.com/cheqd/did-resolver/types"
	"github.com/rs/zerolog/log"
)

// Custom error types for better error handling
var (
	ErrNoHealthyEndpoints = fmt.Errorf("no healthy endpoints available")
	ErrEndpointUnavailable = fmt.Errorf("endpoint is currently unavailable")
	ErrEndpointTimeout = fmt.Errorf("endpoint health check timeout")
	ErrEndpointConnectionFailed = fmt.Errorf("failed to establish connection to endpoint")
)

// EndpointHealth represents the health status of an endpoint
type EndpointHealth struct {
	Network      types.Network
	Endpoint     types.Endpoint
	IsHealthy    bool
	LastCheck    time.Time
	FailureCount int
	Mutex        sync.RWMutex
}

// EndpointManager manages endpoint health and fallback logic
type EndpointManager struct {
	config              types.Config
	endpoints           map[string]*EndpointHealth
	mutex               sync.RWMutex
	healthCheckInterval time.Duration
	healthDataTTL      time.Duration
	healthTimeout       time.Duration
	stopChan           chan struct{}
	wg                 sync.WaitGroup
}

// NewEndpointManager creates a new endpoint manager
func NewEndpointManager(config types.Config) *EndpointManager {
	em := &EndpointManager{
		endpoints:           make(map[string]*EndpointHealth),
		config:              config,
		healthTimeout:       15 * time.Second,
		healthCheckInterval: 60 * time.Second,
		healthDataTTL:       120 * time.Second,
		stopChan:            make(chan struct{}),
	}

	em.initializeEndpoints()
	em.performStartupHealthCheck()

	em.startBackgroundHealthChecker()

	return em
}

// initializeEndpoints sets up endpoint health tracking based on configuration
func (em *EndpointManager) initializeEndpoints() {
	em.endpoints = make(map[string]*EndpointHealth)

	for _, network := range em.config.Networks {
		namespace := network.Namespace
		
		// Initialize endpoints using their role-based keys
		for _, endpoint := range network.Endpoints {
			key := fmt.Sprintf("%s-%s", namespace, endpoint.Role)
			em.endpoints[key] = &EndpointHealth{
				Network:      network,
				Endpoint:     endpoint,
				IsHealthy:    true, // Assume healthy initially
				LastCheck:    time.Now(),
				FailureCount: 0,
			}
		}
	}
}

// performStartupHealthCheck runs health checks on startup and validates system can start
func (em *EndpointManager) performStartupHealthCheck() {
	log.Info().Msg("Performing initial health check on all endpoints...")
	em.performInitialHealthCheck()
	
	// Verify at least one endpoint is healthy before allowing server to start
	if !em.hasAnyHealthyEndpoints() {
		log.Fatal().Msg("No healthy endpoints available - server cannot start. Check endpoint configuration and network connectivity.")
	}
	
	log.Info().Msg("Initial health check completed - server can start")
}

// GetHealthyEndpoint returns the best healthy endpoint for a given namespace
// Priority: Primary endpoint if healthy, otherwise fallback if healthy
func (em *EndpointManager) GetHealthyEndpoint(namespace string) (*types.Network, error) {
	em.mutex.RLock()
	defer em.mutex.RUnlock()

	var healthyPrimary *EndpointHealth
	var healthyFallback *EndpointHealth

	// Find healthy endpoints by role
	primaryKey := fmt.Sprintf("%s-%s", namespace, types.EndpointRolePrimary)
	fallbackKey := fmt.Sprintf("%s-%s", namespace, types.EndpointRoleFallback)
	
	if endpointHealth, exists := em.endpoints[primaryKey]; exists && em.isEndpointHealthy(endpointHealth) {
		healthyPrimary = endpointHealth
	}
	
	if endpointHealth, exists := em.endpoints[fallbackKey]; exists && em.isEndpointHealthy(endpointHealth) {
		healthyFallback = endpointHealth
	}

	// Priority 1: Use primary endpoint if healthy
	if healthyPrimary != nil {
		log.Debug().Msgf("Using primary endpoint %s for namespace %s", healthyPrimary.Endpoint.URL, healthyPrimary.Network.Namespace)
		return em.createNetworkWithEndpoint(healthyPrimary), nil
	}

	// Priority 2: Use fallback endpoint if primary is unhealthy
	if healthyFallback != nil {
		log.Debug().Msgf("Using fallback endpoint %s for namespace %s (primary unavailable)", healthyFallback.Endpoint.URL, healthyFallback.Network.Namespace)
		return em.createNetworkWithEndpoint(healthyFallback), nil
	}

	// No healthy endpoints found
	log.Warn().Msgf("No healthy endpoints found for namespace %s - background health checker will update status", namespace)
	return nil, ErrNoHealthyEndpoints
}

// createNetworkWithEndpoint creates a Network with only the specified healthy endpoint
func (em *EndpointManager) createNetworkWithEndpoint(endpointHealth *EndpointHealth) *types.Network {
	network := endpointHealth.Network
	network.Endpoints = []types.Endpoint{endpointHealth.Endpoint}
	return &network
}

// isEndpointHealthy checks if an endpoint is healthy, considering stale data
func (em *EndpointManager) isEndpointHealthy(endpointHealth *EndpointHealth) bool {
	endpointHealth.Mutex.RLock()
	defer endpointHealth.Mutex.RUnlock()

	isHealthy := endpointHealth.IsHealthy
	lastCheck := endpointHealth.LastCheck

	// Check if health data is stale
	if time.Since(lastCheck) > em.healthDataTTL {
		log.Debug().Msgf("Health data for endpoint %s is stale, marking as potentially unhealthy", endpointHealth.Endpoint.URL)
		return false
	}

	return isHealthy
}

// MarkEndpointUnhealthy marks an endpoint as unhealthy
func (em *EndpointManager) MarkEndpointUnhealthy(network types.Network) {
	em.updateEndpointHealth(network, false)
}

// MarkEndpointHealthy marks an endpoint as healthy
func (em *EndpointManager) MarkEndpointHealthy(network types.Network) {
	em.updateEndpointHealth(network, true)
}

// updateEndpointHealth updates the health status of an endpoint
func (em *EndpointManager) updateEndpointHealth(network types.Network, isHealthy bool) {
	em.mutex.RLock()
	defer em.mutex.RUnlock()

	// Update health for each endpoint in the network using role-based keys
	for _, endpoint := range network.Endpoints {
		key := fmt.Sprintf("%s-%s", network.Namespace, endpoint.Role)
		if endpointHealth, exists := em.endpoints[key]; exists {
			endpointHealth.Mutex.Lock()
			if isHealthy {
				if !endpointHealth.IsHealthy {
					log.Info().Msgf("Marked endpoint %s as healthy again", endpoint.URL)
				}
				endpointHealth.FailureCount = 0
			} else {
				endpointHealth.FailureCount++
				log.Warn().Msgf("Marked endpoint %s as unhealthy (failure count: %d)", endpoint.URL, endpointHealth.FailureCount)
			}
			endpointHealth.IsHealthy = isHealthy
			endpointHealth.LastCheck = time.Now()
			endpointHealth.Mutex.Unlock()
		}
	}
}

// startBackgroundHealthChecker starts the background health checker goroutine
func (em *EndpointManager) startBackgroundHealthChecker() {
	em.wg.Add(1)
	go func() {
		defer em.wg.Done()
		ticker := time.NewTicker(em.healthCheckInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				em.performPeriodicHealthChecks()
			case <-em.stopChan:
				return
			}
		}
	}()
}

// performHealthChecks performs health checks on all endpoints
func (em *EndpointManager) performHealthChecks(logMessage string) {
	em.mutex.RLock()
	namespaces := make([]string, 0, len(em.endpoints))
	for key := range em.endpoints {
		namespace := strings.Split(key, "-")[0]
		namespaces = append(namespaces, namespace)
	}
	em.mutex.RUnlock()

	for _, namespace := range namespaces {
		em.checkAllEndpointsHealth(namespace)
	}
	
	log.Info().Msg(logMessage)
}

// performInitialHealthCheck performs health checks on all endpoints on startup
func (em *EndpointManager) performInitialHealthCheck() {
	em.performHealthChecks("Initial health check completed")
}

// performPeriodicHealthChecks performs health checks on all endpoints
func (em *EndpointManager) performPeriodicHealthChecks() {
	em.performHealthChecks("Periodic health check completed")
}

// checkAllEndpointsHealth performs health checks on all endpoints for a namespace
func (em *EndpointManager) checkAllEndpointsHealth(namespace string) {
	em.mutex.RLock()
	defer em.mutex.RUnlock()

	// Find endpoints by role-based keys
	primaryKey := fmt.Sprintf("%s-%s", namespace, types.EndpointRolePrimary)
	fallbackKey := fmt.Sprintf("%s-%s", namespace, types.EndpointRoleFallback)
	
	primaryEndpoint, _ := em.endpoints[primaryKey]
	fallbackEndpoint, _ := em.endpoints[fallbackKey]

	// Check primary endpoint first (priority order)
	if primaryEndpoint != nil {
		log.Debug().Msgf("Checking health for endpoint primary: %s", primaryEndpoint.Endpoint.URL)
		em.checkEndpointHealth(primaryEndpoint)
		time.Sleep(100 * time.Millisecond) // Small delay between checks
	}

	// Then check fallback endpoint
	if fallbackEndpoint != nil {
		log.Debug().Msgf("Checking health for endpoint fallback: %s", fallbackEndpoint.Endpoint.URL)
		em.checkEndpointHealth(fallbackEndpoint)
		time.Sleep(100 * time.Millisecond) // Small delay between checks
	}
}

// checkEndpointHealth performs a health check on a single endpoint
func (em *EndpointManager) checkEndpointHealth(endpointHealth *EndpointHealth) {
	if em.performSingleHealthCheck(&endpointHealth.Endpoint) {
		em.markEndpointHealthy(endpointHealth)
	} else {
		em.markEndpointUnhealthy(endpointHealth)
	}
}

// markEndpointHealthy marks an endpoint as healthy
func (em *EndpointManager) markEndpointHealthy(endpointHealth *EndpointHealth) {
	endpointHealth.Mutex.Lock()
	endpointHealth.IsHealthy = true
	endpointHealth.FailureCount = 0
	endpointHealth.LastCheck = time.Now()
	endpointHealth.Mutex.Unlock()
}

// markEndpointUnhealthy marks an endpoint as unhealthy
func (em *EndpointManager) markEndpointUnhealthy(endpointHealth *EndpointHealth) {
	endpointHealth.Mutex.Lock()
	endpointHealth.IsHealthy = false
	endpointHealth.FailureCount++
	endpointHealth.LastCheck = time.Now()
	endpointHealth.Mutex.Unlock()
}

// performSingleHealthCheck performs a simple, fast health check
func (em *EndpointManager) performSingleHealthCheck(endpoint *types.Endpoint) bool {
	conn, err := openGRPCConnectionWithTimeout(endpoint.URL, endpoint.UseTls, em.healthTimeout)
	if err != nil {
		log.Debug().Err(err).Msgf("Health check failed for endpoint %s: connection failed", endpoint.URL)
		return false
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), em.healthTimeout)
	defer cancel()

	client := didTypes.NewQueryClient(conn)
	
	// Try to query metadata for a non-existent DID - this should return NotFound, not connection error
	_, err = client.AllDidDocVersionsMetadata(ctx, &didTypes.QueryAllDidDocVersionsMetadataRequest{
		Id: "did:cheqd:testnet:healthcheck",
	})
	
	if err != nil {
		// Check if it's a gRPC status error (like NotFound) vs connection error
		if grpcStatus, ok := status.FromError(err); ok {
			// gRPC status errors mean the service is working but returned an error
			// Connection/transport errors mean the service is unhealthy
			if grpcStatus.Code() == codes.NotFound || grpcStatus.Code() == codes.InvalidArgument {
				log.Debug().Msgf("Health check passed for endpoint %s: service responded with %s", endpoint.URL, grpcStatus.Code())
				return true
			}
		}
		
		log.Debug().Err(err).Msgf("Health check failed for endpoint %s: service error", endpoint.URL)
		return false
	}
	
	log.Debug().Msgf("Health check passed for endpoint %s: service responded successfully", endpoint.URL)
	return true
}

// hasAnyHealthyEndpoints checks if there are any healthy endpoints in the manager
func (em *EndpointManager) hasAnyHealthyEndpoints() bool {
	em.mutex.RLock()
	defer em.mutex.RUnlock()

	for _, endpointHealth := range em.endpoints {
		if em.isEndpointHealthy(endpointHealth) {
			return true
		}
	}
	return false
}
 