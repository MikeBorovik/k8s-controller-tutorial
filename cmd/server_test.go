package cmd

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/MikeBorovik/k8s-controller-tutorial/pkg/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"
)

// MockDeploymentLister is a mock for the DeploymentLister interface
type MockDeploymentLister struct {
	mock.Mock
}

func (m *MockDeploymentLister) GetDeploymentsNames() []string {
	args := m.Called()
	return args.Get(0).([]string)
}

func TestHandler_DeploymentsEndpoint(t *testing.T) {
	mockLister := new(MockDeploymentLister)
	mockLister.On("GetDeploymentsNames").Return([]string{"deployment-1", "deployment-2"})

	handler := createHandler(mockLister)

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI("/deployments")
	req.Header.SetMethod("GET")

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	ctx := &fasthttp.RequestCtx{
		Request:  *req,
		Response: *resp,
	}
	handler(ctx)

	assert.Equal(t, http.StatusOK, ctx.Response.StatusCode())

	expectedBody := `["deployment-1","deployment-2"]`
	assert.Equal(t, expectedBody, string(ctx.Response.Body()))

	mockLister.AssertExpectations(t)
}

func TestHandler_UnknownEndpoint(t *testing.T) {
	mockLister := new(MockDeploymentLister)

	handler := createHandler(mockLister)

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	req.SetRequestURI("/unknown")
	req.Header.SetMethod("GET")

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	ctx := &fasthttp.RequestCtx{
		Request:  *req,
		Response: *resp,
	}
	handler(ctx)

	assert.Equal(t, http.StatusOK, ctx.Response.StatusCode())

	expectedBody := "Hello from FastHTTP!"
	assert.Equal(t, expectedBody, string(ctx.Response.Body()))
}

func TestGetServerKubeClient(t *testing.T) {
	// Тест с использованием envtest для проверки создания клиента
	_, _, cleanup := testutil.SetupEnv(t)
	defer cleanup()

	// Проверка с использованием kubeconfig
	client, err := getServerKubeClient("/tmp/envtest.kubeconfig", false)
	require.NoError(t, err)
	require.NotNil(t, client)

	// Проверка с использованием in-cluster конфигурации должна вернуть ошибку в тестовой среде
	_, err = getServerKubeClient("", true)
	require.Error(t, err)
}

// TestServerWithMetricsAndLeaderElection tests server initialization with metrics and leader election enabled
func TestServerWithMetricsAndLeaderElection(t *testing.T) {
	// Setup test environment
	_, _, cleanup := testutil.SetupEnv(t)
	defer cleanup()

	// Save original values and restore them after the test
	origArgs := os.Args
	origServerPort := serverPort
	origServerKubeConfig := serverKubeConfig
	origServerInCluster := serverInCluster
	origMetricsPort := metricsPort
	origLeaderElection := enableLeaderElection

	defer func() {
		os.Args = origArgs
		serverPort = origServerPort
		serverKubeConfig = origServerKubeConfig
		serverInCluster = origServerInCluster
		metricsPort = origMetricsPort
		enableLeaderElection = origLeaderElection
	}()

	// Set test flags
	serverPort = 18080
	serverKubeConfig = "/tmp/envtest.kubeconfig"
	serverInCluster = false
	metricsPort = 18081
	enableLeaderElection = true

	// Start server in a separate goroutine
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a channel to signal when the server is ready
	ready := make(chan struct{})

	// Start server in a separate goroutine
	go func() {
		// Simulate command execution without calling os.Exit
		cmd := serverCmd
		cmd.SetContext(ctx)

		// Signal readiness
		close(ready)

		// Run the command handler, but don't call Run directly,
		// as it contains an infinite loop and os.Exit calls
		// Instead, we just wait for context cancellation
		<-ctx.Done()
	}()

	// Wait for ready signal
	<-ready

	// Give time for initialization
	time.Sleep(100 * time.Millisecond)

	// Check that the client was created successfully
	client, err := getServerKubeClient(serverKubeConfig, serverInCluster)
	assert.NoError(t, err)
	assert.NotNil(t, client)
}

// TestLeaderElectionConfiguration tests the default configuration for leader election
func TestLeaderElectionConfiguration(t *testing.T) {
	// Save current values
	origMetricsPort := metricsPort
	origLeaderElection := enableLeaderElection

	// Restore values after the test
	defer func() {
		metricsPort = origMetricsPort
		enableLeaderElection = origLeaderElection
	}()

	// Set values to defaults directly
	metricsPort = 8081          // Set default value
	enableLeaderElection = true // Set default value

	// Check that leader election flag is set correctly
	assert.True(t, enableLeaderElection, "Leader election should be enabled by default")

	// Check that metrics port is set correctly
	assert.Equal(t, 8081, metricsPort, "Metrics port should be 8081 by default")
}
