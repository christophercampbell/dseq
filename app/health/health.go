package health

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Status represents the health status
type Status string

const (
	StatusUp   Status = "UP"
	StatusDown Status = "DOWN"
)

// Check represents a health check
type Check struct {
	Name     string
	Check    func(context.Context) error
	Timeout  time.Duration
	Interval time.Duration
}

// Result represents a health check result
type Result struct {
	Status    Status                 `json:"status"`
	Timestamp time.Time              `json:"timestamp"`
	Details   map[string]CheckResult `json:"details"`
}

// CheckResult represents the result of a single check
type CheckResult struct {
	Status    Status    `json:"status"`
	Error     string    `json:"error,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

// Service represents a health check service
type Service struct {
	checks  map[string]*Check
	results map[string]CheckResult
	mu      sync.RWMutex
	server  *http.Server
}

// New creates a new health check service
func New(port int) *Service {
	s := &Service{
		checks:  make(map[string]*Check),
		results: make(map[string]CheckResult),
	}

	// Create HTTP server
	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.handleHealth)
	mux.HandleFunc("/ready", s.handleReady)

	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	return s
}

// AddCheck adds a health check
func (s *Service) AddCheck(check *Check) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.checks[check.Name] = check
	go s.runCheck(check)
}

// Start starts the health check service
func (s *Service) Start() error {
	return s.server.ListenAndServe()
}

// Stop stops the health check service
func (s *Service) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

// runCheck runs a health check periodically
func (s *Service) runCheck(check *Check) {
	ticker := time.NewTicker(check.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(context.Background(), check.Timeout)
			err := check.Check(ctx)
			cancel()

			s.mu.Lock()
			s.results[check.Name] = CheckResult{
				Status:    StatusUp,
				Timestamp: time.Now(),
			}
			if err != nil {
				s.results[check.Name] = CheckResult{
					Status:    StatusDown,
					Error:     err.Error(),
					Timestamp: time.Now(),
				}
			}
			s.mu.Unlock()
		}
	}
}

// handleHealth handles health check requests
func (s *Service) handleHealth(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := Result{
		Status:    StatusUp,
		Timestamp: time.Now(),
		Details:   make(map[string]CheckResult),
	}

	for name, checkResult := range s.results {
		result.Details[name] = checkResult
		if checkResult.Status == StatusDown {
			result.Status = StatusDown
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if result.Status == StatusDown {
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	fmt.Fprintf(w, `{"status":"%s","timestamp":"%s","details":%v}`,
		result.Status,
		result.Timestamp.Format(time.RFC3339),
		result.Details)
}

// handleReady handles readiness check requests
func (s *Service) handleReady(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Check if all required checks are up
	for _, checkResult := range s.results {
		if checkResult.Status == StatusDown {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
