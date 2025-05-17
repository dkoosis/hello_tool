// file: cmd/hello-tool-base/main_test.go
package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os" // Keep this if you use it, or remove if not
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	// Adjust import path to your actual internal packages
	"github.com/dkoosis/hello-tool-base/internal/buildinfo"
	"github.com/dkoosis/hello-tool-base/internal/config"
	"github.com/dkoosis/hello-tool-base/internal/logging"
)

var (
	log logging.Logger // Declare log variable
)

// TestMain is executed before any other tests in this package.
func TestMain(m *testing.M) {
	logging.SetupDefaultLogger("error")
	log = logging.GetLogger("main_test")
	cfg = config.DefaultConfig()
	os.Exit(m.Run())
}

// TestHelloHandler_ReturnsGreeting_When_NameParameterProvided (ADR-008 Naming)
func TestHelloHandler_ReturnsGreeting_When_NameParameterProvided(t *testing.T) {
	// Arrange
	req := httptest.NewRequest("GET", "/hello?name=TestUser", nil)
	rr := httptest.NewRecorder()

	// Act
	helloHandler(rr, req)

	// Assert
	resp := rr.Result()
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			t.Logf("Warning: error closing response body: %v", err) // Log if you care, or just assign to _
		}
	}() // Corrected: Check error from resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Status code should be OK")

	var response GreetingResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	require.NoError(t, err, "Should be no error decoding JSON response")

	assert.Contains(t, response.Message, "Hello, TestUser", "Response message should contain the provided name")
}

// TestHelloHandler_ReturnsBadRequest_When_NameParameterMissing (ADR-008 Naming)
func TestHelloHandler_ReturnsBadRequest_When_NameParameterMissing(t *testing.T) {
	// Arrange
	req := httptest.NewRequest("GET", "/hello", nil)
	rr := httptest.NewRecorder()

	// Act
	helloHandler(rr, req)

	// Assert
	resp := rr.Result()
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			t.Logf("Warning: error closing response body: %v", err)
		}
	}() // Corrected: Check error from resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, "Status code should be Bad Request")

	var errorResponse ClientErrorResponse
	err := json.NewDecoder(resp.Body).Decode(&errorResponse)
	require.NoError(t, err, "Should be no error decoding JSON error response")

	assert.Contains(t, errorResponse.Error, "Invalid Request Parameter", "Error message should indicate invalid parameter")
	assert.Contains(t, errorResponse.Details, "The 'name' query parameter is required.", "Error details should specify missing 'name'")
}

// TestRootHandler_ReturnsBuildInfo_When_PathIsRoot (ADR-008 Naming)
func TestRootHandler_ReturnsBuildInfo_When_PathIsRoot(t *testing.T) {
	// Arrange
	originalVersion := buildinfo.Version
	originalCommit := buildinfo.CommitHash
	originalBuildDate := buildinfo.BuildDate
	buildinfo.Version = "test-v1.0"
	buildinfo.CommitHash = "testcommit123"
	buildinfo.BuildDate = "2024-01-01T12:00:00Z"
	defer func() {
		buildinfo.Version = originalVersion
		buildinfo.CommitHash = originalCommit
		buildinfo.BuildDate = originalBuildDate
	}()
	if cfg == nil {
		cfg = config.DefaultConfig()
	}
	cfg.Server.Name = "TestService"

	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	// Act
	rootHandler(rr, req)

	// Assert
	resp := rr.Result()
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			t.Logf("Warning: error closing response body: %v", err)
		}
	}() // Corrected: Check error from resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Status code should be OK")
	assert.Equal(t, "text/plain; charset=utf-8", resp.Header.Get("Content-Type"), "Content-Type should be text/plain")

	body := rr.Body.String()
	assert.Contains(t, body, "Hello World Root!", "Body should contain root greeting")
	assert.Contains(t, body, "Service: TestService", "Body should contain the configured service name")
	assert.Contains(t, body, "Version: test-v1.0", "Body should contain test version")
	assert.Contains(t, body, "Commit:  testcommit123", "Body should contain test commit hash")
	assert.Contains(t, body, "Built:   2024-01-01T12:00:00Z", "Body should contain test build date")
}

// TestRootHandler_ReturnsNotFound_When_PathIsNotRoot (ADR-008 Naming)
func TestRootHandler_ReturnsNotFound_When_PathIsNotRoot(t *testing.T) {
	// Arrange
	req := httptest.NewRequest("GET", "/notfound", nil)
	rr := httptest.NewRecorder()

	// Act
	rootHandler(rr, req)

	// Assert
	resp := rr.Result()
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			t.Logf("Warning: error closing response body: %v", err)
		}
	}() // Corrected: Check error from resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode, "Status code should be Not Found for incorrect path")
}

// TestHealthHandler_ReturnsHealthyStatus_When_Called (ADR-008 Naming)
func TestHealthHandler_ReturnsHealthyStatus_When_Called(t *testing.T) {
	// Arrange
	originalVersion := buildinfo.Version
	originalCommit := buildinfo.CommitHash
	originalBuildDate := buildinfo.BuildDate
	buildinfo.Version = "health-v1.1"
	buildinfo.CommitHash = "healthSHA"
	buildinfo.BuildDate = "2024-02-02T10:00:00Z"
	defer func() {
		buildinfo.Version = originalVersion
		buildinfo.CommitHash = originalCommit
		buildinfo.BuildDate = originalBuildDate
	}()

	req := httptest.NewRequest("GET", "/health", nil)
	rr := httptest.NewRecorder()

	// Act
	healthHandler(rr, req)

	// Assert
	resp := rr.Result()
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			t.Logf("Warning: error closing response body: %v", err)
		}
	}() // Corrected: Check error from resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Status code should be OK")

	var healthStatus struct {
		Status    string `json:"status"`
		Version   string `json:"version"`
		Commit    string `json:"commit"`
		BuildDate string `json:"buildDate"`
	}
	err := json.NewDecoder(resp.Body).Decode(&healthStatus)
	require.NoError(t, err, "Should be no error decoding JSON health response")

	assert.Equal(t, "OK", healthStatus.Status, "Health status should be OK")
	assert.Equal(t, "health-v1.1", healthStatus.Version, "Version should match build info")
	assert.Equal(t, "healthSHA", healthStatus.Commit, "Commit should match build info")
	assert.Equal(t, "2024-02-02T10:00:00Z", healthStatus.BuildDate, "BuildDate should match build info")
}
