package cmd

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
