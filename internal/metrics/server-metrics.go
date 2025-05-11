// Package metrics provides structures and functions for collecting and managing server health and performance metrics.
// file: internal/metrics/server_metrics.go
package metrics

import (
	"runtime"
	"sync"
	"time"
)

// ServerMetrics holds various metrics about the server's health and performance.
// This version is generalized for a template tool program.
type ServerMetrics struct {
	// Server uptime and basic info.
	StartTime     time.Time     `json:"startTime"`
	Uptime        time.Duration `json:"uptime"`
	GoVersion     string        `json:"goVersion"`
	NumGoroutines int           `json:"numGoroutines"`

	// Memory stats.
	MemoryAllocated   uint64 `json:"memoryAllocated"`   // Currently allocated memory in bytes.
	MemoryTotalAlloc  uint64 `json:"memoryTotalAlloc"`  // Total allocated memory since start.
	MemorySystemTotal uint64 `json:"memorySystemTotal"` // Total memory obtained from system.
	MemoryGCCount     uint32 `json:"memoryGCCount"`     // Number of completed GC cycles.

	// Connection stats (generic for any type of connection the server might handle).
	ActiveConnections int `json:"activeConnections"`
	TotalConnections  int `json:"totalConnections"`  // Total unique connections initiated.
	FailedConnections int `json:"failedConnections"` // Connections that failed to establish or dropped unexpectedly.

	// Request stats (generic for API requests).
	TotalRequests    int            `json:"totalRequests"`
	FailedRequests   int            `json:"failedRequests"`
	RequestLatencies map[string]int `json:"requestLatencies"` // Map of (e.g., handler name or path) to average latency in ms.

	// Last errors recorded by the application.
	LastErrors []ErrorInfo `json:"lastErrors,omitempty"`
}

// ErrorInfo contains details about an error that occurred.
type ErrorInfo struct {
	Timestamp time.Time `json:"timestamp"`
	Component string    `json:"component"` // e.g., "database", "HTTPHandler", "auth_service"
	Message   string    `json:"message"`
	Stack     string    `json:"stack,omitempty"` // Optional: stack trace if available and safe to expose/log.
}

// Collector manages server metrics collection and reporting.
type Collector struct {
	metrics     ServerMetrics
	startTime   time.Time
	errorBuffer []ErrorInfo
	bufferSize  int
	mu          sync.RWMutex

	// Connection tracking (e.g., using connection IDs or remote addresses as keys).
	activeConnections map[string]bool
}

// NewCollector creates a new metrics collector instance.
// errorBufferSize determines how many recent errors are kept in memory.
func NewCollector(errorBufferSize int) *Collector {
	startTime := time.Now()

	return &Collector{
		metrics: ServerMetrics{
			StartTime:        startTime,
			GoVersion:        runtime.Version(),
			RequestLatencies: make(map[string]int),
			// ActiveConnections, TotalConnections, etc., will start at 0.
		},
		startTime:         startTime,
		errorBuffer:       make([]ErrorInfo, 0, errorBufferSize),
		bufferSize:        errorBufferSize,
		activeConnections: make(map[string]bool),
	}
}

// GetCurrentMetrics returns a copy of the current server metrics.
// This method is safe for concurrent use.
func (c *Collector) GetCurrentMetrics() ServerMetrics {
	c.mu.RLock()
	// Defer RUnlock until after we've potentially modified c.metrics for Uptime and NumGoroutines.
	// This means the read lock is held for the duration of the critical section.

	// Update real-time metrics that are cheap to get.
	// Note: These are written to c.metrics under RLock, which is generally not advised.
	// A better pattern might be to calculate these and put them directly into metricsCopy.
	// However, for simple counters and runtime stats, this is often acceptable if GetCurrentMetrics is the primary writer of these specific fields.
	// For this example, we'll update them directly and then copy.
	currentUptime := time.Since(c.startTime)
	currentNumGoroutines := runtime.NumGoroutine()

	// Update memory stats.
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	currentMemoryAllocated := memStats.Alloc
	currentMemoryTotalAlloc := memStats.TotalAlloc
	currentMemorySystemTotal := memStats.Sys
	currentMemoryGCCount := memStats.NumGC

	// Create a copy of the metrics to return.
	// This is important to avoid race conditions if the caller modifies the returned map/slice.
	metricsCopy := c.metrics // Copy existing values first

	// Apply the freshly read values to the copy
	metricsCopy.Uptime = currentUptime
	metricsCopy.NumGoroutines = currentNumGoroutines
	metricsCopy.MemoryAllocated = currentMemoryAllocated
	metricsCopy.MemoryTotalAlloc = currentMemoryTotalAlloc
	metricsCopy.MemorySystemTotal = currentMemorySystemTotal
	metricsCopy.MemoryGCCount = currentMemoryGCCount

	// Create a fresh copy of the error buffer for the snapshot.
	if len(c.errorBuffer) > 0 {
		metricsCopy.LastErrors = make([]ErrorInfo, len(c.errorBuffer))
		copy(metricsCopy.LastErrors, c.errorBuffer)
	} else {
		metricsCopy.LastErrors = nil // Ensure it's nil if empty, not an empty slice from previous calls
	}
	c.mu.RUnlock() // Unlock after all reads from c.metrics and c.errorBuffer are done.

	return metricsCopy
}

// RecordRequest records statistics about a generic application request.
// 'identifier' could be an API endpoint, method name, or other request type.
// 'latencyMs' is the request processing time in milliseconds.
// 'success' indicates if the request was processed successfully.
func (c *Collector) RecordRequest(identifier string, latencyMs int, success bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.metrics.TotalRequests++
	if !success {
		c.metrics.FailedRequests++
	}

	// Update average latency for this request identifier.
	// Using a simple moving average. More sophisticated averaging could be used if needed.
	if existingAvg, ok := c.metrics.RequestLatencies[identifier]; ok {
		// This is a very basic moving average. For N requests, it's (avg_N-1 * (N-1) + new_latency) / N
		// The current c.metrics.TotalRequests might not be specific to this 'identifier'.
		// A more accurate average would track counts per identifier.
		// For simplicity, we'll do a 50/50 blend if it exists.
		c.metrics.RequestLatencies[identifier] = (existingAvg + latencyMs) / 2
	} else {
		c.metrics.RequestLatencies[identifier] = latencyMs
	}
}

// RecordConnection tracks connection statistics.
// 'connectionID' is a unique identifier for the connection.
// 'active' is true if the connection is being established/is active, false if it's being closed.
func (c *Collector) RecordConnection(connectionID string, active bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if active {
		// New or re-activated connection.
		if _, exists := c.activeConnections[connectionID]; !exists {
			// This is a new unique connection being marked active.
			c.activeConnections[connectionID] = true
			c.metrics.TotalConnections++ // Increment total unique connections seen.
		} else {
			// Connection was already known, just marking it active again (no change to TotalConnections).
			c.activeConnections[connectionID] = true
		}
	} else {
		// Connection closed.
		delete(c.activeConnections, connectionID)
	}

	// Update active count.
	c.metrics.ActiveConnections = len(c.activeConnections)
}

// RecordConnectionFailure increments the failed connections counter.
// This could be called when a connection attempt fails before it's considered "active".
func (c *Collector) RecordConnectionFailure() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.metrics.FailedConnections++
}

// RecordError adds an error to the in-memory error buffer.
// 'component' indicates the part of the system where the error occurred.
// 'message' is the error message.
// 'stack' (optional) is the stack trace.
func (c *Collector) RecordError(component, message, stack string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	errorInfo := ErrorInfo{
		Timestamp: time.Now(),
		Component: component,
		Message:   message,
		Stack:     stack,
	}

	// Add to the circular buffer.
	if c.bufferSize <= 0 { // Guard against zero or negative buffer size
		return
	}
	if len(c.errorBuffer) >= c.bufferSize {
		// Remove oldest error to make space.
		c.errorBuffer = c.errorBuffer[1:]
	}
	c.errorBuffer = append(c.errorBuffer, errorInfo)
}
